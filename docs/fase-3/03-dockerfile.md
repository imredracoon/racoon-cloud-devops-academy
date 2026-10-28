# 3.3 Dockerfile

!!! obiettivi "Obiettivi"
    - Scrivere un Dockerfile e costruire un'immagine con `docker build`.
    - Conoscere le istruzioni principali e la differenza tra `CMD` ed `ENTRYPOINT`.
    - Sfruttare la cache degli strati e il file `.dockerignore`.
    - Usare le build multi-stage per ottenere immagini piccole.
    - Applicare le buone pratiche di sicurezza: utente non privilegiato, immagini minime, nessun segreto.
    - Containerizzare la Tasks API.

## La ricetta di un'immagine

Un **Dockerfile** è un file di testo con le istruzioni per costruire un'immagine, eseguite dall'alto in basso.
Podman e Buildah usano lo stesso formato, che chiamano anche **Containerfile**.

```dockerfile title="Dockerfile"
FROM python:3.12-slim
WORKDIR /app
COPY saluta.py .
CMD ["python", "saluta.py"]
```

```python title="saluta.py"
import os
print(f"Ciao, {os.environ.get('NOME', 'mondo')}!")
```

```bash
docker build -t saluta:1.0 .         # il punto è il "contesto": la cartella inviata al builder
docker run --rm saluta:1.0           # Ciao, mondo!
docker run --rm -e NOME=Vincenzo saluta:1.0
```

## Le istruzioni

| Istruzione | Funzione |
|---|---|
| `FROM immagine` | l'immagine di partenza. È sempre la prima istruzione. |
| `WORKDIR /percorso` | la cartella di lavoro per le istruzioni seguenti e per il container (la crea se manca) |
| `COPY origine destinazione` | copia file dal contesto di build nell'immagine |
| `RUN comando` | esegue un comando **durante la build** e salva il risultato in un nuovo strato |
| `ENV NOME=valore` | imposta una variabile d'ambiente, presente anche nel container |
| `ARG NOME=valore` | una variabile disponibile **solo durante la build** (`--build-arg NOME=...`) |
| `EXPOSE 8000` | documenta la porta su cui ascolta l'applicazione (non la pubblica: serve `-p`) |
| `USER utente` | l'utente con cui eseguire le istruzioni seguenti e il container |
| `CMD [...]` | il comando predefinito all'avvio del container |
| `ENTRYPOINT [...]` | l'eseguibile fisso del container |
| `HEALTHCHECK CMD ...` | un comando che verifica periodicamente se l'applicazione è sana |
| `LABEL chiave=valore` | metadati: autore, versione, repository del codice |

Esiste anche `ADD`, simile a `COPY` ma capace di scaricare URL ed estrarre archivi: usa `COPY` salvo necessità
precise, è più prevedibile.

### CMD ed ENTRYPOINT

- `CMD` è il comando **predefinito**: se passi un comando a `docker run immagine ALTRO`, `CMD` viene sostituito.
- `ENTRYPOINT` è l'eseguibile **fisso**: ciò che passi a `docker run` diventa un suo argomento; `CMD`, se presente,
  fornisce gli argomenti predefiniti.

```dockerfile
ENTRYPOINT ["ping"]
CMD ["-c", "3", "localhost"]
```

`docker run immagine` esegue `ping -c 3 localhost`; `docker run immagine -c 1 example.com` esegue
`ping -c 1 example.com`.

!!! tip "Usa sempre la forma exec, con le parentesi quadre"
    `CMD ["gunicorn", "app:app"]` (forma *exec*) avvia il programma direttamente come PID 1. La forma *shell*
    `CMD gunicorn app:app` lo avvia attraverso `/bin/sh -c`: la shell diventa PID 1 e **non inoltra** il segnale
    SIGTERM all'applicazione, che non si chiude in modo ordinato e viene uccisa dopo 10 secondi a ogni `docker stop`.

## La cache degli strati

Ogni istruzione che modifica il filesystem (`RUN`, `COPY`, `ADD`) crea uno **strato**. Docker conserva gli strati
in **cache**: se un'istruzione e tutto ciò che la precede non sono cambiati, riusa lo strato senza rieseguirla.
Appena un'istruzione cambia, **tutte quelle successive** vengono rieseguite.

Da qui la regola: **ciò che cambia raramente in alto, ciò che cambia spesso in basso**.

=== "Sbagliato"

    ```dockerfile
    FROM python:3.12-slim
    WORKDIR /app
    COPY . .
    RUN pip install -r requirements.txt     # rieseguito a OGNI modifica del codice
    CMD ["python", "app.py"]
    ```

=== "Corretto"

    ```dockerfile
    FROM python:3.12-slim
    WORKDIR /app
    COPY requirements.txt .
    RUN pip install --no-cache-dir -r requirements.txt   # rieseguito solo se cambiano le dipendenze
    COPY . .
    CMD ["python", "app.py"]
    ```

Altre regole che riducono dimensione e strati:

- unisci i comandi correlati in un solo `RUN` e pulisci **nello stesso** `RUN`: un file cancellato in uno strato
  successivo occupa comunque spazio nell'immagine;

    ```dockerfile
    RUN apt-get update \
        && apt-get install -y --no-install-recommends curl \
        && rm -rf /var/lib/apt/lists/*
    ```

- `pip install --no-cache-dir` evita di salvare nell'immagine la cache di pip.

### .dockerignore

Il **contesto di build** (la cartella indicata a `docker build`) viene inviato tutto al builder. Il file
`.dockerignore` esclude ciò che non serve, rendendo la build più veloce ed evitando di copiare per errore file
sensibili con `COPY . .`:

```text title=".dockerignore"
.git/
.venv/
__pycache__/
*.pyc
instance/
*.db
.env
tests/
```

## Build multi-stage

Spesso per **costruire** un'applicazione servono strumenti (compilatori, header, gestori di pacchetti) che non
servono per **eseguirla**. Una build **multi-stage** usa più `FROM` nello stesso Dockerfile: si costruisce in una
fase e si copia nell'immagine finale **solo il risultato**.

L'esempio classico è un programma Go, che produce un unico eseguibile senza dipendenze:

```dockerfile title="Dockerfile"
# Fase 1: compilazione, con tutto il toolchain di Go
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY main.go .
RUN CGO_ENABLED=0 go build -o /out/server main.go

# Fase 2: l'immagine finale contiene solo l'eseguibile
FROM scratch
COPY --from=build /out/server /server
USER 65534:65534
EXPOSE 8080
ENTRYPOINT ["/server"]
```

L'immagine `golang:1.25-alpine` supera i 300 MB; quella finale, basata su `scratch` (un'immagine vuota), pesa circa
**13 MB** e non contiene nemmeno una shell: niente da sfruttare per un attaccante.

Per un'applicazione Python come la nostra il guadagno è minore, perché l'interprete serve comunque. La build
multi-stage diventa utile quando alcune dipendenze vanno **compilate** (servono `gcc` e gli header nella prima fase,
ma non nella seconda).

## Le buone pratiche di sicurezza

1. **Non eseguire come root.** Crea un utente e passa a lui con `USER`. Se un attaccante compromette l'applicazione,
   si trova con i privilegi minimi.
2. **Immagini di base piccole e ufficiali**, con tag precisi: meno pacchetti, meno vulnerabilità.
3. **Nessun segreto nell'immagine.** Né in `COPY`, né in `ENV`, né in `ARG`: restano negli strati e chiunque abbia
   l'immagine può leggerli con `docker history` o estraendo gli strati. I segreti si passano **all'avvio**
   (variabili d'ambiente, file montati, gestori di segreti nel cloud).
4. **Ricostruisci spesso** le immagini, per includere gli aggiornamenti di sicurezza dell'immagine di base.
5. **Un processo per container.** Nginx e l'API vanno in due container separati, che si scalano e si aggiornano in
   modo indipendente.

## La Tasks API in un container

Ecco il Dockerfile della Tasks API (versione di riferimento in `labs/tasks-api/Dockerfile`):

```dockerfile title="Dockerfile"
# Immagine della Tasks API (Fase 3)
FROM python:3.12-slim

# Python: niente file .pyc e output non bufferizzato (i log arrivano subito)
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

# Utente senza privilegi con cui girerà l'applicazione
RUN useradd --system --uid 10001 --home-dir /app app \
    && install -d -o app -g app /app/instance

# Prima le dipendenze: questo strato resta in cache finché requirements.txt non cambia
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

USER app
EXPOSE 8000

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD ["python", "-c", "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8000/health', timeout=2)"]

CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "2", "--preload", "--access-logfile", "-", "app:create_app()"]
```

Alcuni dettagli:

- Gunicorn ascolta su `0.0.0.0` e non su `127.0.0.1`: nel container, `127.0.0.1` è raggiungibile solo dall'interno
  del container stesso, e la porta pubblicata non funzionerebbe.
- `--access-logfile -` scrive i log delle richieste su standard output, dove `docker logs` li trova.
- La cartella `instance/` appartiene all'utente `app`, così il database SQLite predefinito funziona anche senza
  PostgreSQL.
- `HEALTHCHECK` usa Python, già presente, invece di `curl`, che non c'è nell'immagine slim: lo stato compare in
  `docker ps` come `healthy` o `unhealthy`.
- `--preload` fa creare l'applicazione (e le tabelle del database) una sola volta, prima di avviare i worker. Così
  due worker non provano a creare le stesse tabelle nello stesso istante. Per questo, nella versione di riferimento,
  `create_app()` chiude le connessioni al database subito dopo `create_all()`: i worker ne apriranno di proprie.

Aggiungi anche il driver PostgreSQL a `requirements.txt`, ti servirà nella lezione 3.4:

```text title="requirements.txt"
Flask==3.1.2
Flask-SQLAlchemy==3.1.1
gunicorn==23.0.0
psycopg[binary]==3.3.6
```

```bash
docker build -t tasks-api:dev .
docker run -d --name api -p 8000:8000 tasks-api:dev
curl -s localhost:8000/health
docker ps              # dopo qualche secondo: (healthy)
docker logs api
```

## Esercizi

!!! esercizio "Esercizio 1 · Containerizza la Tasks API"
    Nel tuo repository `tasks-api` scrivi `Dockerfile` e `.dockerignore`, costruisci l'immagine, avviala e verifica con
    `curl` che `POST /tasks` e `GET /tasks` funzionino. Controlla con `docker exec api id` che il processo non giri
    come root.

??? success "Soluzione"
    Usa il Dockerfile e il `.dockerignore` della lezione, poi:
    ```bash
    docker build -t tasks-api:dev .
    docker run -d --name api -p 8000:8000 tasks-api:dev
    curl -s -X POST localhost:8000/tasks -H 'Content-Type: application/json' -d '{"title":"Nel container"}'
    curl -s localhost:8000/tasks
    docker exec api id            # uid=10001(app)
    ```
    **Spiegazione.** Se `id` mostra `uid=0(root)`, manca l'istruzione `USER` o si trova prima di istruzioni che
    richiedono root (come `RUN useradd`): l'ordine conta.

!!! esercizio "Esercizio 2 · Misura la cache"
    Modifica solo `app.py` (per esempio un messaggio) e ricostruisci, misurando il tempo con `time docker build ...`.
    Poi modifica `requirements.txt` aggiungendo un commento e ricostruisci di nuovo. Quali passaggi vengono rieseguiti
    nei due casi? Come lo vedi nell'output?

??? success "Soluzione"
    Nel primo caso l'output mostra `CACHED` per tutti i passi fino a `RUN pip install`, e solo `COPY app.py` e i
    successivi vengono rieseguiti: pochi secondi. Nel secondo caso cambia `COPY requirements.txt`, quindi vengono
    rieseguiti anche `pip install` e tutto ciò che segue.

    **Spiegazione.** Docker confronta le istruzioni e, per `COPY`, il contenuto dei file copiati. Basta un commento
    in `requirements.txt` per invalidare la cache da quel punto in poi.

!!! esercizio "Esercizio 3 · Trova gli errori"
    Questo Dockerfile funziona, ma ha almeno sei problemi. Trovali e correggili.

    ```dockerfile
    FROM python:latest
    COPY . /app
    WORKDIR /app
    RUN pip install -r requirements.txt
    ENV DATABASE_URL=postgresql+psycopg://admin:Password123@db/tasks
    RUN apt-get update
    RUN apt-get install -y curl
    EXPOSE 8000
    CMD gunicorn --bind 0.0.0.0:8000 "app:create_app()"
    ```

??? success "Soluzione"
    1. `python:latest`: tag non preciso e immagine completa, molto grande. Usa `python:3.12-slim`.
    2. `COPY . /app` prima di `pip install`: ogni modifica al codice invalida la cache delle dipendenze. E senza
       `.dockerignore` copia anche `.git`, `.venv` ed eventuali `.env`.
    3. `pip install` senza `--no-cache-dir`.
    4. Una password nell'immagine con `ENV`: è un segreto, va passato all'avvio.
    5. `apt-get update` e `apt-get install` in due `RUN` separati, senza pulizia: la cache del primo può restare
       vecchia e le liste dei pacchetti restano nell'immagine. Inoltre `curl` probabilmente non serve.
    6. Nessun `USER`: il container gira come root.
    7. `CMD` in forma shell: Gunicorn non riceve SIGTERM.

    ```dockerfile
    FROM python:3.12-slim
    ENV PYTHONUNBUFFERED=1
    WORKDIR /app
    RUN useradd --system --uid 10001 app
    COPY requirements.txt .
    RUN pip install --no-cache-dir -r requirements.txt
    COPY app.py .
    USER app
    EXPOSE 8000
    CMD ["gunicorn", "--bind", "0.0.0.0:8000", "app:create_app()"]
    ```

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale istruzione viene eseguita durante la build e crea un nuovo strato?","options":["CMD","RUN","EXPOSE","ENTRYPOINT"],"answer":[1],"explanation":"CMD ed ENTRYPOINT definiscono cosa eseguire all'avvio del container."},
 {"q":"Perché si copia requirements.txt e si installano le dipendenze prima di copiare il codice?","options":["Per sicurezza","Per sfruttare la cache: le dipendenze cambiano meno spesso del codice","Perché pip lo richiede","Per ridurre il numero di file"],"answer":[1],"explanation":"Una modifica al codice non obbliga a reinstallare le dipendenze."},
 {"q":"Cosa fa `EXPOSE 8000`?","options":["Pubblica la porta sull'host","Documenta la porta su cui ascolta l'applicazione","Apre il firewall","Avvia il server"],"answer":[1],"explanation":"Per pubblicarla serve -p al momento di docker run."},
 {"q":"Perché è preferibile `CMD [\"gunicorn\", ...]` a `CMD gunicorn ...`?","options":["È più corto","Con la forma exec l'applicazione è PID 1 e riceve i segnali come SIGTERM","La forma shell non è supportata","Per usare le variabili d'ambiente"],"answer":[1],"explanation":"Con la forma shell il PID 1 è /bin/sh, che non inoltra SIGTERM."},
 {"q":"Quali sono buone pratiche per un Dockerfile? (scegli 2)","options":["Eseguire l'applicazione con un utente non root","Usare tag precisi per l'immagine di base","Salvare le password con ENV","Usare sempre :latest","Installare strumenti di debug in produzione"],"answer":[0,1],"explanation":"I segreti non vanno mai nell'immagine."},
 {"q":"A cosa serve una build multi-stage?","options":["A eseguire più container","A costruire con tutti gli strumenti necessari e copiare nell'immagine finale solo il risultato","A creare più tag","A dividere i log"],"answer":[1],"explanation":"Immagini finali più piccole e con meno superficie d'attacco."}
]}
</script>

## Fonti

- [Docker Docs: Dockerfile reference](https://docs.docker.com/reference/dockerfile/)
- [Docker Docs: Building best practices](https://docs.docker.com/build/building/best-practices/)
- [Docker Docs: Build cache](https://docs.docker.com/build/cache/)
- [Docker Docs: Multi-stage builds](https://docs.docker.com/build/building/multi-stage/)
- [Docker Docs: .dockerignore](https://docs.docker.com/build/concepts/context/#dockerignore-files)
- [Gunicorn: Settings (preload, bind, accesslog)](https://docs.gunicorn.org/en/stable/settings.html)
- [OWASP: Docker Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Docker_Security_Cheat_Sheet.html)
