# Progetto 3 · Stack containerizzato

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai impacchettare un'applicazione reale in immagini sicure e ripetibili, descrivere un'architettura a più
    servizi come codice e gestirne dati, configurazione e aggiornamenti. È la competenza che quasi ogni annuncio
    DevOps dà per scontata.

| | |
|---|---|
| **Durata** | 10-14 ore |
| **Consegna** | Il repository `tasks-api`, esteso, e l'immagine pubblicata su GitHub Container Registry |
| **Ambiente** | Dev Container per lo sviluppo; VM `web` (o `rhel-lab` con Podman) per l'installazione |

## Il compito

### Nel repository `tasks-api`

1. **Dockerfile** con: immagine di base `slim` con tag preciso, utente non root, cache delle dipendenze ottimizzata,
   `HEALTHCHECK`, `CMD` in forma exec, etichetta `org.opencontainers.image.source`. `.dockerignore` completo.
2. **compose.yaml** con tre servizi: `proxy` (Nginx, unico con porta pubblicata), `api`, `db` (PostgreSQL con volume).
   Reti separate `frontend` e `backend`, avvio ordinato con gli healthcheck, servizio `api` indurito
   (`read_only`, `cap_drop`, `no-new-privileges`, limite di memoria).
3. **Configurazione** solo tramite `.env` (escluso da Git) con un `.env.example` documentato.
4. **Backup**: uno script `scripts/backup-db.sh` che esegue `pg_dump` nel container del database e salva un file
   `tasks-AAAAMMGG-HHMMSS.sql.gz` in una cartella indicata, conservando gli ultimi 7 (riusa ciò che hai scritto nella
   Fase 2). E uno script `scripts/restore-db.sh FILE` che ripristina.
5. **Immagine pubblicata** su `ghcr.io/TUO-UTENTE/tasks-api` con i tag `1.1.0` e `sha-HASH`.
6. **Scansione**: nel README, l'esito della scansione Trivy (HIGH e CRITICAL) dell'immagine pubblicata, con una breve
   valutazione.
7. **README** aggiornato: architettura (diagramma Mermaid), avvio in un comando, variabili, backup e ripristino,
   aggiornamento a una nuova versione.

### Sul server

Sulla VM `web` (percorso Debian/Ubuntu, con Docker) oppure su `rhel-lab` (percorso RHEL, con Podman e
`podman-compose`):

- lo stack gira usando l'immagine **scaricata dal registry**, non costruita sul server: crea un
  `compose.prod.yaml` che usa `image: ghcr.io/TUO-UTENTE/tasks-api:1.1.0` al posto di `build`;
- lo stack riparte da solo dopo un riavvio della VM;
- il backup gira ogni notte con un timer di systemd (lezione 2.5).

## Suggerimenti

- Per lo script di backup:

    ```bash
    docker compose exec -T db pg_dump -U tasks -d tasks | gzip > "$dest/tasks-$(date +%Y%m%d-%H%M%S).sql.gz"
    ```

    `-T` disattiva l'allocazione del terminale: è necessaria quando l'output va in una pipe o lo script gira da un
    timer, senza terminale.

- Per il ripristino, il database deve esistere e conviene fermare l'API durante l'operazione:

    ```bash
    docker compose stop api
    gunzip -c "$1" | docker compose exec -T db psql -U tasks -d tasks
    docker compose start api
    ```

    Per un ripristino pulito, crea il dump con `pg_dump --clean --if-exists`: il file conterrà i comandi per
    cancellare le tabelle esistenti prima di ricrearle.

- Sul percorso RHEL, se vuoi andare oltre, trasforma i tre servizi in **Quadlet** (`.container`, `.network`,
  `.volume` in `/etc/containers/systemd/`): systemd li gestirà come servizi, con dipendenze e riavvii.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `docker compose up -d` da un clone pulito con un `.env` | stack attivo, `api` e `db` healthy |
| `curl localhost:8080/health` | `{"status": "ok", ...}` |
| `nc -zv localhost 5432` e `nc -zv localhost 8000` dall'host | entrambe rifiutate |
| `docker compose exec api id` | utente non root |
| `docker compose exec api touch /app/x` | `Read-only file system` |
| backup, `docker compose down -v`, `up -d`, ripristino | le attività tornano |
| riavvio della VM | lo stack riparte da solo |
| `docker history ghcr.io/TUO-UTENTE/tasks-api:1.1.0` | nessun segreto negli strati |

## Criteri di valutazione

- [ ] Tutti i test della tabella passano.
- [ ] Nessun segreto nel repository né nell'immagine.
- [ ] Il README permette a un estraneo di avviare lo stack in meno di cinque minuti.
- [ ] La scansione di sicurezza è documentata e commentata.
- [ ] La storia Git racconta il lavoro: una pull request per ogni passo importante.

??? success "Una possibile soluzione: compose.prod.yaml e il timer"
    ```yaml title="compose.prod.yaml"
    # Uso: docker compose -f compose.yaml -f compose.prod.yaml up -d
    services:
      api:
        image: ghcr.io/TUO-UTENTE/tasks-api:1.1.0
        build: !reset null
        pull_policy: always
    ```

    `!reset` è un'etichetta YAML di Compose che **annulla** un valore ereditato dal file precedente: qui toglie
    `build`, così sul server l'immagine viene solo scaricata.

    ```ini title="/etc/systemd/system/tasks-backup.service"
    [Unit]
    Description=Backup del database della Tasks API

    [Service]
    Type=oneshot
    WorkingDirectory=/srv/tasks-api
    ExecStart=/srv/tasks-api/scripts/backup-db.sh /srv/backup
    ```

    ```ini title="/etc/systemd/system/tasks-backup.timer"
    [Unit]
    Description=Backup notturno della Tasks API

    [Timer]
    OnCalendar=*-*-* 02:00:00
    Persistent=true

    [Install]
    WantedBy=timers.target
    ```

    Lo stack riparte dopo il riavvio grazie a `restart: unless-stopped` e al servizio `docker` abilitato all'avvio
    (`systemctl is-enabled docker`).

## Fonti

- [Docker Docs: Use Compose in production](https://docs.docker.com/compose/how-tos/production/)
- [Compose file reference: merge e `!reset`](https://docs.docker.com/reference/compose-file/merge/)
- [PostgreSQL: pg_dump](https://www.postgresql.org/docs/17/app-pgdump.html)
- [Podman: Quadlet, podman-systemd.unit(5)](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html)
