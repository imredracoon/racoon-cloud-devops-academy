# G.6 Go nei container

!!! obiettivi "Obiettivi"
    - Capire perché un binario Go statico è ideale per i container.
    - Scegliere l'immagine di base: `scratch`, `distroless/static` o `alpine`.
    - Scrivere un Dockerfile multi-stage per Go con cache dei moduli e versione incorporata nel binario.
    - Costruire immagini per amd64 e arm64 con la cross-compilazione, senza emulazione.
    - Confrontare dimensioni e vulnerabilità con l'immagine Python della Tasks API.

!!! info "Quando studiarla"
    Nella [Fase 3](../fase-3/index.md), dopo la [lezione 3.3](../fase-3/03-dockerfile.md). Qui applichi le build
    multi-stage e gli utenti non root a un programma compilato.

## Il binario statico

Un programma Python nel container ha bisogno dell'interprete, delle librerie del sistema e dei pacchetti `pip`:
l'immagine della Tasks API parte da `python:3.12-slim`, che pesa più di 100 MB prima ancora del tuo codice. Un
programma Go compilato con `CGO_ENABLED=0` è invece un **binario statico**: contiene tutto ciò che gli serve e
non chiede nulla al sistema, nemmeno la libreria C.

```console
$ CGO_ENABLED=0 go build -o probe ./cmd/probe
$ ldd probe
	not a dynamic executable
```

!!! info "Cos'è cgo"
    **cgo** permette a Go di chiamare codice C. Alcuni pacchetti della libreria standard lo usano quando è
    disponibile (ad esempio per la risoluzione DNS tramite la libc). Con `CGO_ENABLED=0` Go usa le proprie
    implementazioni in Go puro e il binario diventa statico. Per gli strumenti di questo percorso non perdi nulla.

Un binario statico può girare in un'immagine **senza sistema operativo**. Restano solo tre cose da fornire:

| Cosa | Perché serve | Chi la fornisce |
|---|---|---|
| Certificati delle CA | per verificare i certificati TLS dei siti HTTPS ([lezione 1.12](../fase-1/12-http-nginx.md)) | `/etc/ssl/certs/ca-certificates.crt` |
| Fusi orari | per `time.LoadLocation("Europe/Rome")` | `/usr/share/zoneinfo`, oppure l'import di `time/tzdata` |
| Un utente non root | per non girare come root ([lezione 3.6](../fase-3/06-registry-sicurezza.md)) | `/etc/passwd` con un utente |

## Quale immagine di base

| Base | Dimensione | Shell | Contenuto | Quando usarla |
|---|---|---|---|---|
| `scratch` | 0 MB | no | niente: un'immagine vuota | quando aggiungi tu certificati e utente |
| `gcr.io/distroless/static-debian13:nonroot` | circa 2 MB | no | certificati CA, fusi orari, utente `nonroot` (UID 65532) | **la scelta consigliata** |
| `alpine` | circa 8 MB | sì | un piccolo Linux con il gestore `apk` | se ti serve una shell per il debug |

Le immagini **distroless** del progetto Google Container Tools contengono solo ciò che serve a eseguire un
programma: niente shell, niente gestore di pacchetti. Meno software significa meno vulnerabilità da segnalare e
meno strumenti a disposizione di chi riuscisse a entrare nel container.

## Il Dockerfile

```dockerfile title="Dockerfile"
FROM --platform=$BUILDPLATFORM golang:1.27 AS build   # (1)!
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download                                  # (2)!

COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/probe ./cmd/probe   # (3)!

FROM gcr.io/distroless/static-debian13:nonroot       # (4)!
COPY --from=build /out/probe /probe
COPY probe.example.yaml /etc/probe/probe.yaml
EXPOSE 9101
USER nonroot:nonroot
ENTRYPOINT ["/probe"]
CMD ["serve", "-config", "/etc/probe/probe.yaml"]
```

1.  Lo stadio di build gira sempre sull'architettura della macchina che costruisce (`BUILDPLATFORM`), anche quando
    l'immagine finale è per un'altra: vedi la sezione sul multi-architettura.
2.  `go.mod` e `go.sum` si copiano **da soli** prima del codice: finché non cambiano, Docker riusa lo strato con i
    moduli già scaricati e ogni build successiva è molto più rapida. È lo stesso trucco che hai usato con
    `requirements.txt` nella [lezione 3.3](../fase-3/03-dockerfile.md).
3.  `-trimpath` toglie dal binario i percorsi della macchina di build (build riproducibili). `-s -w` toglie i
    simboli di debug. `-X main.version=...` scrive il valore nella variabile `version` del pacchetto `main`: così
    `probe version` dice esattamente quale rilascio sta girando.
4.  Lo stadio finale contiene solo il binario e il file di configurazione d'esempio.

```console
$ docker build --build-arg VERSION=v0.3.0 -t probe:v0.3.0 .
$ docker images probe
REPOSITORY   TAG      IMAGE ID       SIZE
probe        v0.3.0   4f1c2a9d8e7b   13.4MB
$ docker run --rm probe:v0.3.0 version
v0.3.0
```

Il file `.dockerignore` del progetto esclude `.git`, i binari locali e la cartella della pipeline, così il contesto
inviato a Docker resta piccolo e un binario compilato sul tuo PC non finisce per errore nell'immagine.

### Configurazione dall'esterno

L'immagine contiene una configurazione d'esempio; quella vera si **monta** al momento dell'avvio, come hai fatto
con Nginx nella [lezione 3.4](../fase-3/04-volumi-reti.md):

```bash
docker run -d --name probe -p 9101:9101 \
  -v "$PWD/probe.yaml:/etc/probe/probe.yaml:ro" probe:v0.3.0
curl -s localhost:9101/results
docker stop probe        # SIGTERM: nei log compare "arresto in corso"
docker logs probe
```

!!! tip "GOMAXPROCS nei container"
    Il runtime di Go decide quanti thread usare in base ai core disponibili. Da Go 1.25 tiene conto anche dei
    **limiti di CPU del cgroup** impostati con `docker run --cpus` o con i `resources.limits` di Kubernetes, quindi
    non serve più impostare `GOMAXPROCS` a mano come si faceva un tempo.

## Immagini multi-architettura

Le istanze cloud con processori ARM (AWS Graviton, Google Axion) costano meno a parità di prestazioni, e i mini PC
o i Raspberry Pi dell'homelab possono essere ARM. Un'immagine **multi-architettura** è un indice che punta a una
variante per ogni processore: chi la scarica riceve quella giusta in automatico.

Con un linguaggio interpretato, costruire la variante arm64 su un PC x86-64 richiede di **emulare** l'intero
processore con QEMU, ed è lento. Con Go non serve: lo stadio di build gira nativo e il compilatore produce
direttamente il binario per l'altra architettura (`GOARCH=$TARGETARCH`). Docker Buildx imposta da solo
`BUILDPLATFORM`, `TARGETOS` e `TARGETARCH`.

```bash
docker buildx create --use --name multi           # una volta: un builder che supporta più piattaforme
docker buildx build --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=v0.3.0 -t ghcr.io/TUONOME/racoon-probe:v0.3.0 --push .
docker buildx imagetools inspect ghcr.io/TUONOME/racoon-probe:v0.3.0
```

Con Podman ([lezione 3.7](../fase-3/07-podman.md)) il risultato si ottiene con un *manifest*:

```bash
podman build --platform linux/amd64,linux/arm64 --manifest racoon-probe:v0.3.0 .
podman manifest push racoon-probe:v0.3.0 docker://ghcr.io/TUONOME/racoon-probe:v0.3.0
```

## Laboratorio

!!! laboratorio "Laboratorio · L'immagine di racoon-probe"
    1. Scrivi il Dockerfile e il `.dockerignore` del tuo `racoon-probe`.
    2. Costruisci l'immagine con `VERSION` uguale all'ultimo tag Git (`git describe --tags --always`).
    3. Avvia la Tasks API con il suo `compose.yaml` della [Fase 3](../fase-3/05-compose.md) e aggiungi un servizio
       `probe` che controlla `http://api:8000/health` e `db:5432` **sulla rete interna** di Compose.
    4. Confronta con `docker images` la dimensione dell'immagine di probe con quella della Tasks API, e con
       `trivy image` il numero di vulnerabilità ([lezione 3.6](../fase-3/06-registry-sicurezza.md)).
    5. Pubblica l'immagine multi-architettura su GitHub Container Registry.

??? success "Soluzione guidata"
    ```yaml title="compose.yaml (servizio da aggiungere)"
      probe:
        build: ../racoon-probe          # percorso del tuo repository racoon-probe
        volumes:
          - ./probe.yaml:/etc/probe/probe.yaml:ro
        ports:
          - "9101:9101"
        networks: [backend]
        restart: unless-stopped
    ```
    ```yaml title="probe.yaml"
    interval: 15s
    timeout: 3s
    targets:
      - name: tasks-api
        type: http
        url: http://api:8000/health
      - name: postgres
        type: tcp
        address: db:5432
    ```
    ```console
    $ docker compose up -d --build && curl -s localhost:9101/results | jq '.[] | {name: .target.name, ok}'
    {"name": "tasks-api", "ok": true}
    {"name": "postgres", "ok": true}
    $ docker images --format '{{.Repository}}\t{{.Size}}' | grep -E 'probe|api'
    tasks-api-api   190MB
    tasks-api-probe 13.4MB
    ```
    **Spiegazione.** Nella rete di Compose i servizi si trovano per **nome** grazie al DNS interno, quindi probe
    usa `api` e `db`, non `localhost`. L'immagine di probe è più di dieci volte più piccola e `trivy` trova di solito
    zero vulnerabilità del sistema operativo, perché un sistema operativo, di fatto, non c'è. Le dimensioni
    esatte dipendono dalle versioni delle immagini di base.

## Esercizi

!!! esercizio "Esercizio 1 · Da distroless a scratch"
    Cambia lo stadio finale in `FROM scratch`, ricostruisci e avvia probe con un target `https://go.dev`. Che errore
    ottieni? Correggilo senza tornare a distroless.

??? success "Soluzione"
    L'errore è `tls: failed to verify certificate: x509: certificate signed by unknown authority`: in `scratch` non
    ci sono i certificati delle CA. Si copiano dallo stadio di build, che li ha, insieme a un utente non root:

    ```dockerfile
    FROM scratch
    COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
    COPY --from=build /out/probe /probe
    USER 65532:65532
    ENTRYPOINT ["/probe"]
    ```
    **Spiegazione.** Con `scratch` non esiste `/etc/passwd`, quindi l'utente si indica con il numero (UID:GID).
    Questo è esattamente il lavoro che distroless fa per te, con i fusi orari in più: per questo è la scelta
    consigliata.

!!! esercizio "Esercizio 2 · Quale versione sta girando?"
    Senza entrare nel container (non c'è una shell!), come scopri quale versione di probe sta girando in un
    container avviato da qualcun altro? Trova almeno due modi.

??? success "Soluzione"
    ```bash
    docker exec probe /probe version                         # esegue il binario stesso, senza shell
    docker logs probe | head -1                              # il log di avvio contiene "versione"
    docker inspect probe --format '{{.Config.Image}}'        # il tag dell'immagine, se è stato usato bene
    ```
    **Spiegazione.** `docker exec` non richiede una shell: esegue direttamente il programma indicato. Con
    `go version -m probe` sul binario puoi leggere anche la versione di Go e dei moduli usati per compilarlo: Go le
    incorpora sempre, ed è ciò che usano gli scanner di vulnerabilità.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa ottieni compilando con CGO_ENABLED=0?","options":["Un binario più lento","Un binario statico che non dipende dalla libc del sistema","Un binario solo per Windows","Un'immagine Docker"],"answer":[1],"explanation":"Può girare in scratch o distroless/static."},
 {"q":"Perché nel Dockerfile si copiano go.mod e go.sum prima del resto del codice?","options":["Per obbligo di Docker","Per riusare la cache dello strato con i moduli scaricati finché le dipendenze non cambiano","Per ridurre l'immagine finale","Per firmare l'immagine"],"answer":[1],"explanation":"Lo stesso schema di requirements.txt in Python."},
 {"q":"In un'immagine FROM scratch un client HTTPS fallisce con x509: certificate signed by unknown authority. Perché?","options":["Il sito è insicuro","Mancano i certificati delle CA","Manca la shell","Il binario non è statico"],"answer":[1],"explanation":"Vanno copiati, oppure si usa distroless che li include."},
 {"q":"Come si costruisce in modo efficiente l'immagine arm64 di un programma Go su un PC x86-64?","options":["Emulando arm64 con QEMU per tutta la build","Facendo girare lo stadio di build su BUILDPLATFORM e compilando con GOARCH=$TARGETARCH","Non si può","Usando un'immagine alpine"],"answer":[1],"explanation":"La cross-compilazione di Go evita l'emulazione."},
 {"q":"A cosa serve -ldflags \"-X main.version=v1.2.0\"?","options":["A ottimizzare il binario","A scrivere un valore nella variabile version del pacchetto main al momento della build","A firmare il binario","A scegliere la versione di Go"],"answer":[1],"explanation":"Così il binario sa quale rilascio è."},
 {"q":"Quale vantaggio di sicurezza danno le immagini distroless?","options":["Cifrano i dati","Non contengono shell né gestore di pacchetti: meno vulnerabilità e meno strumenti per un attaccante","Bloccano la rete","Eseguono come root in modo sicuro"],"answer":[1],"explanation":"Contengono solo ciò che serve a eseguire il programma."}
]}
</script>

## Fonti

- [Docker Docs: Multi-stage builds](https://docs.docker.com/build/building/multi-stage/) e [Multi-platform builds](https://docs.docker.com/build/building/multi-platform/)
- [Docker Docs: Building for Go](https://docs.docker.com/guides/golang/build-images/)
- [GoogleContainerTools/distroless](https://github.com/GoogleContainerTools/distroless)
- [Command cgo](https://pkg.go.dev/cmd/cgo) e [Command go: build flags](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Go 1.25 Release Notes](https://go.dev/doc/go1.25), GOMAXPROCS e cgroup
- [Podman: podman-build](https://docs.podman.io/en/latest/markdown/podman-build.1.html) e [podman-manifest](https://docs.podman.io/en/latest/markdown/podman-manifest.1.html)
