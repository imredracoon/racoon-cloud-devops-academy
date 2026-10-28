# 3.5 Docker Compose

!!! obiettivi "Obiettivi"
    - Descrivere un'applicazione a più container in un file `compose.yaml`.
    - Gestire l'intero stack con pochi comandi: avvio, log, aggiornamento, arresto.
    - Gestire configurazione e segreti con le variabili e il file `.env`.
    - Ordinare l'avvio dei servizi con gli healthcheck.
    - Avviare lo stack completo della Tasks API: Nginx, API, PostgreSQL.

## Il problema

Nella lezione precedente hai avviato due container con comandi lunghi, nell'ordine giusto, con reti e volumi creati
a mano. Con tre o quattro servizi diventa ingestibile e soggetto a errori. **Docker Compose** descrive l'intera
applicazione in un file YAML **dichiarativo**: dici *cosa* vuoi, e Compose crea reti, volumi e container nell'ordine
giusto.

È anche il primo esempio di un'idea che tornerà con Terraform e Kubernetes: **l'infrastruttura descritta in un file
versionato in Git**.

Compose si usa con `docker compose` (con lo spazio: è un plugin di Docker). Il vecchio comando `docker-compose`, con
il trattino, è la versione 1, non più mantenuta.

=== "Debian/Ubuntu"

    Il plugin si installa con il pacchetto `docker-compose-plugin` ([lezione 3.2](02-immagini-container.md)) ed è già presente nel Dev Container:

    ```bash
    docker compose version
    ```

=== "RHEL (Rocky/AlmaLinux)"

    `podman compose` è un comando di Podman che delega a un programma esterno: `docker-compose` oppure
    **podman-compose**, disponibile nel repository EPEL:

    ```bash
    sudo dnf install -y epel-release
    sudo dnf install -y podman-compose
    podman compose version
    ```

    I file `compose.yaml` sono gli stessi. Per i servizi da tenere attivi in produzione, su RHEL è preferibile
    usare i **Quadlet** ([lezione 3.2](02-immagini-container.md)), che integrano i container in systemd.

## Lo stack della Tasks API

```yaml title="compose.yaml"
# Stack completo della Tasks API: proxy Nginx, API, PostgreSQL (Fase 3)
# Avvio: cp .env.example .env && docker compose up -d --build
services:
  db:
    image: postgres:17
    environment:
      POSTGRES_DB: tasks
      POSTGRES_USER: tasks
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?imposta POSTGRES_PASSWORD nel file .env}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U tasks -d tasks"]
      interval: 5s
      timeout: 3s
      retries: 10
    networks: [backend]
    restart: unless-stopped

  api:
    build: .
    image: tasks-api:dev
    environment:
      DATABASE_URL: postgresql+psycopg://tasks:${POSTGRES_PASSWORD}@db:5432/tasks
    depends_on:
      db:
        condition: service_healthy
    networks: [backend, frontend]
    restart: unless-stopped

  proxy:
    image: nginx:1.30-alpine
    ports:
      - "8080:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      api:
        condition: service_healthy
    networks: [frontend]
    restart: unless-stopped

networks:
  frontend:
  backend:

volumes:
  pgdata:
```

```nginx title="nginx.conf"
# Reverse proxy verso la Tasks API (Fase 3)
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://api:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```text title=".env.example"
POSTGRES_PASSWORD=cambiami
```

Leggiamolo sezione per sezione.

### services

Ogni voce sotto `services` diventa uno o più container. Il nome del servizio (`db`, `api`, `proxy`) è anche il
**nome DNS** con cui gli altri servizi lo raggiungono: per questo `DATABASE_URL` punta a `db` e Nginx a `api`.

- `image`: l'immagine da usare. Con `build: .` Compose costruisce l'immagine dal Dockerfile nella cartella
  indicata, e `image` diventa il nome da dare al risultato.
- `environment`: le variabili d'ambiente del container.
- `ports`: le porte pubblicate, `"HOST:CONTAINER"` (sempre tra virgolette: YAML può interpretare alcuni valori con i
  due punti come numeri in base 60).
- `volumes`: volumi con nome (`pgdata:...`) o bind mount (`./nginx.conf:...`, con percorso relativo al file).
- `restart`: la politica di riavvio ([lezione 3.2](02-immagini-container.md)).

### depends_on e healthcheck

`depends_on` stabilisce l'**ordine di avvio**. Da solo garantisce solo che il container `db` sia **avviato**, non che
PostgreSQL sia **pronto** ad accettare connessioni (ci mette qualche secondo). Con `condition: service_healthy`,
Compose aspetta che l'**healthcheck** del servizio risulti sano:

- per `db` l'healthcheck è definito nel file, con `pg_isready`, lo strumento di PostgreSQL che verifica se il server
  accetta connessioni;
- per `api` è quello definito con `HEALTHCHECK` nel Dockerfile ([lezione 3.3](03-dockerfile.md)).

La catena è quindi: `db` sano → avvio di `api` → `api` sana → avvio di `proxy`.

### networks e volumes

In fondo al file si dichiarano reti e volumi con nome. Compose li crea con un **prefisso**, il nome del progetto
(per impostazione predefinita il nome della cartella): `tasks-api_backend`, `tasks-api_pgdata`. Senza la sezione
`networks`, Compose crea comunque una rete predefinita per tutti i servizi del file.

### Variabili e segreti: il file .env

Compose legge automaticamente il file `.env` nella stessa cartella e sostituisce le variabili `${NOME}` nel
`compose.yaml`:

- `${POSTGRES_PASSWORD}` viene sostituita con il valore del file `.env` (o dell'ambiente della shell);
- `${POSTGRES_PASSWORD:?messaggio}` fa **fallire** Compose con quel messaggio se la variabile manca: meglio un errore
  chiaro che un database con password vuota;
- `${PORTA:-8080}` usa 8080 se la variabile non è impostata.

Il file `.env` contiene segreti e **non va in Git**: si versiona un `.env.example` con valori finti, e ognuno crea il
proprio `.env`.

```bash
docker compose config        # mostra il file con le variabili già sostituite: utile per il debug
```

!!! info "I segreti in produzione"
    Le variabili d'ambiente sono comode ma non perfette: sono visibili con `docker inspect` a chi ha accesso al demone.
    Compose supporta anche i `secrets`, montati come file in `/run/secrets/`. Nel cloud e in Kubernetes userai
    gestori di segreti dedicati (Fasi 5 e 8).

## I comandi

```bash
cp .env.example .env              # e cambia la password
docker compose up -d --build      # costruisce, crea reti e volumi, avvia tutto in background
docker compose ps                 # stato dei servizi, con lo stato di salute
docker compose logs -f api        # i log di un servizio (senza nome: tutti)
docker compose exec db psql -U tasks -d tasks     # un comando in un servizio
docker compose restart api
docker compose stop               # ferma i container, senza eliminarli
docker compose down               # ferma ed elimina container e reti (i volumi restano)
docker compose down -v            # elimina anche i volumi: CANCELLA I DATI
```

### Aggiornare un servizio

Dopo una modifica al codice dell'API:

```bash
docker compose up -d --build api
```

Compose ricostruisce l'immagine e ricrea **solo** i container la cui configurazione o immagine è cambiata. Gli altri
restano in funzione. Allo stesso modo, se modifichi il `compose.yaml` (per esempio una variabile), basta
`docker compose up -d`.

### Più file e ambienti

Compose unisce automaticamente `compose.yaml` e, se esiste, `compose.override.yaml`. È il modo classico di
aggiungere impostazioni di sviluppo (per esempio pubblicare la porta del database per ispezionarlo) senza toccare il
file principale. Con `-f` scegli esplicitamente i file: `docker compose -f compose.yaml -f compose.prod.yaml up -d`.

## Esercizi

!!! esercizio "Esercizio 1 · Lo stack completo"
    Nel repository `tasks-api` aggiungi `compose.yaml`, `nginx.conf`, `.env.example`; aggiungi `.env` al
    `.gitignore`. Avvia lo stack e verifica:

    1. con `docker compose ps` che tutti i servizi siano in esecuzione e `api` e `db` siano `healthy`;
    2. che `curl localhost:8080/health` risponda attraverso Nginx;
    3. che il database **non** sia raggiungibile dall'host sulla porta 5432;
    4. che i dati sopravvivano a `docker compose down` seguito da `docker compose up -d`.

??? success "Soluzione"
    ```bash
    cp .env.example .env && sed -i 's/cambiami/una-password-lunga/' .env
    echo ".env" >> .gitignore
    docker compose up -d --build
    docker compose ps
    curl -s localhost:8080/health
    curl -s -X POST localhost:8080/tasks -H 'Content-Type: application/json' -d '{"title":"Con Compose"}'
    nc -zv -w 2 localhost 5432          # Connection refused: la porta non è pubblicata
    docker compose down && docker compose up -d
    curl -s localhost:8080/tasks        # l'attività è ancora lì: il volume pgdata non è stato toccato
    ```
    **Spiegazione.** `down` elimina container e reti ma conserva i volumi con nome. Solo `down -v` li cancella.

!!! esercizio "Esercizio 2 · Cosa succede se il database cade?"
    Con lo stack avviato, ferma solo il database (`docker compose stop db`). Cosa risponde `/health`? E `/tasks`?
    Quale stato mostra `docker compose ps` per `api` dopo una trentina di secondi? Riavvia il database e osserva.

??? success "Soluzione"
    `/health` risponde `503` con `{"status": "error"}`, perché l'endpoint verifica il database. `/tasks` risponde
    `500`: l'applicazione non gestisce l'errore in modo specifico. Dopo alcuni controlli falliti `api` diventa
    `unhealthy`. Con `docker compose start db` tutto torna a funzionare, senza riavviare l'API: SQLAlchemy riapre le
    connessioni.

    **Spiegazione.** Ecco perché l'endpoint di salute deve controllare le dipendenze: un orchestratore (Kubernetes,
    un load balancer) può smettere di mandare traffico a un'istanza che non è in grado di servirlo. Docker da solo
    **non** riavvia un container `unhealthy`: si limita a segnalarlo.

!!! esercizio "Esercizio 3 · Un ambiente di sviluppo"
    Crea un `compose.override.yaml` che, solo in sviluppo:

    1. pubblichi la porta 5432 del database, ma solo su `127.0.0.1`;
    2. avvii l'API con il server di sviluppo di Flask e il ricaricamento automatico, montando il codice con un bind mount.

    Verifica con `docker compose config` il risultato dell'unione dei due file.

??? success "Soluzione"
    ```yaml title="compose.override.yaml"
    services:
      db:
        ports:
          - "127.0.0.1:5432:5432"
      api:
        command: ["flask", "--app", "app:create_app()", "run", "--host", "0.0.0.0", "--port", "8000", "--debug"]
        volumes:
          - ./app.py:/app/app.py:ro
    ```
    ```bash
    docker compose config | less
    docker compose up -d
    # modifica app.py: Flask si ricarica da solo
    docker compose logs -f api
    ```
    **Spiegazione.** `command` sostituisce il `CMD` dell'immagine. Il file di override non va usato in produzione:
    in quel caso si avvia con `docker compose -f compose.yaml up -d`, che ignora l'override.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Con quale nome il servizio api raggiunge il database in questo compose.yaml?","options":["localhost","db","postgres","tasks-api_db"],"answer":[1],"explanation":"Il nome del servizio è anche il suo nome DNS sulla rete."},
 {"q":"Cosa garantisce `depends_on` con `condition: service_healthy`?","options":["Solo che il container sia avviato","Che il servizio dipendente parta dopo che l'healthcheck dell'altro risulta sano","Che i servizi siano sulla stessa rete","Il riavvio automatico"],"answer":[1],"explanation":"Senza condizione, depends_on garantisce solo l'ordine di avvio dei container."},
 {"q":"Quale comando cancella anche i dati del database?","options":["docker compose stop","docker compose down","docker compose down -v","docker compose restart"],"answer":[2],"explanation":"-v elimina i volumi con nome."},
 {"q":"Dove metti la password del database per Compose?","options":["Direttamente nel compose.yaml in Git","In un file .env escluso da Git, con un .env.example versionato","Nel Dockerfile","Nel nome dell'immagine"],"answer":[1],"explanation":"I segreti non vanno mai nel repository."},
 {"q":"Hai modificato il codice dell'API. Quale comando aggiorna solo quel servizio?","options":["docker compose down && docker compose up","docker compose up -d --build api","docker compose restart","docker build ."],"answer":[1],"explanation":"Compose ricrea solo ciò che è cambiato."},
 {"q":"Quali file legge automaticamente Compose?","options":["compose.yaml","compose.override.yaml","docker.json","Containerfile.yaml","settings.env"],"answer":[0,1],"explanation":"Legge anche .env per sostituire le variabili."}
]}
</script>

## Fonti

- [Docker Docs: Docker Compose overview](https://docs.docker.com/compose/)
- [Compose file reference](https://docs.docker.com/reference/compose-file/)
- [Docker Docs: Control startup order](https://docs.docker.com/compose/how-tos/startup-order/)
- [Docker Docs: Variable interpolation](https://docs.docker.com/reference/compose-file/interpolation/)
- [Docker Docs: Merge Compose files](https://docs.docker.com/compose/how-tos/multiple-compose-files/merge/)
- [Docker Docs: Secrets in Compose](https://docs.docker.com/compose/how-tos/use-secrets/)
- [PostgreSQL: pg_isready](https://www.postgresql.org/docs/17/app-pg-isready.html)
- [Podman: podman-compose](https://docs.podman.io/en/latest/markdown/podman-compose.1.html)
