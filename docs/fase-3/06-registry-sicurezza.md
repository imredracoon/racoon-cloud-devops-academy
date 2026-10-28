# 3.6 Registry e sicurezza

!!! obiettivi "Obiettivi"
    - Scegliere una strategia di tag per le immagini.
    - Pubblicare un'immagine su GitHub Container Registry e conoscere i registry dei cloud.
    - Analizzare un'immagine alla ricerca di vulnerabilità con Trivy.
    - Eseguire i container con il minimo dei privilegi.
    - Conoscere firma delle immagini e SBOM, i pilastri della sicurezza della supply chain.

## I registry

Un **registry** conserva e distribuisce le immagini. Quelli che incontrerai:

| Registry | Indirizzo | Note |
|---|---|---|
| Docker Hub | `docker.io` | il predefinito, con le immagini ufficiali |
| GitHub Container Registry | `ghcr.io` | integrato con GitHub e GitHub Actions |
| Quay.io | `quay.io` | di Red Hat, usato da molti progetti open source |
| Red Hat Registry | `registry.access.redhat.com` | le immagini UBI |
| Amazon ECR | `ACCOUNT.dkr.ecr.REGIONE.amazonaws.com` | il registry privato di AWS (Fase 5) |
| Google Artifact Registry | `REGIONE-docker.pkg.dev` | il registry di Google Cloud (Fase 5) |

In azienda le immagini delle proprie applicazioni stanno quasi sempre in un registry **privato**, vicino a dove
girano: ECR su AWS, Artifact Registry su Google Cloud.

## La strategia dei tag

Un tag è un'etichetta **mobile**: lo stesso tag può essere spostato su un'altra immagine. Una buona strategia
combina più tag sulla stessa immagine:

| Tag | Esempio | Uso |
|---|---|---|
| versione semantica | `1.4.2` | i rilasci: è ciò che si installa in produzione |
| versione parziale | `1.4`, `1` | chi vuole ricevere le correzioni in automatico |
| hash del commit Git | `sha-3f9a2c1` | sapere esattamente da quale codice viene un'immagine |
| `latest` | | comodità per chi prova; mai in produzione |

```bash
docker tag tasks-api:dev ghcr.io/TUO-UTENTE/tasks-api:1.0.0
docker tag tasks-api:dev ghcr.io/TUO-UTENTE/tasks-api:sha-$(git rev-parse --short HEAD)
```

Per la massima certezza, in produzione si fa riferimento al **digest**: `ghcr.io/utente/tasks-api@sha256:...`. Il
digest identifica il contenuto e non può cambiare.

## Pubblicare su GitHub Container Registry

1. Serve un token con il permesso di scrivere i pacchetti. Con la GitHub CLI:

    ```bash
    gh auth refresh -h github.com -s write:packages
    gh auth token | docker login ghcr.io -u TUO-UTENTE --password-stdin
    ```

    In alternativa crea un *personal access token (classic)* con lo scope `write:packages` in Settings → Developer
    settings. `--password-stdin` evita che il token finisca nella cronologia della shell.

2. Aggiungi al Dockerfile un'etichetta che collega l'immagine al repository del codice:

    ```dockerfile
    LABEL org.opencontainers.image.source="https://github.com/TUO-UTENTE/tasks-api"
    ```

3. Costruisci, etichetta e pubblica:

    ```bash
    docker build -t ghcr.io/TUO-UTENTE/tasks-api:1.0.0 .
    docker push ghcr.io/TUO-UTENTE/tasks-api:1.0.0
    ```

4. Per impostazione predefinita il pacchetto è **privato**. Per renderlo pubblico: sul tuo profilo GitHub →
   Packages → `tasks-api` → Package settings → Change visibility.

Nella Fase 7 sarà la pipeline di GitHub Actions a costruire e pubblicare l'immagine a ogni rilascio, senza token
personali.

## Analisi delle vulnerabilità con Trivy

Ogni immagine contiene centinaia di pacchetti: quelli della distribuzione di base e le librerie della tua
applicazione. Per alcuni di essi esistono **vulnerabilità note**, catalogate con un identificativo **CVE**
(*Common Vulnerabilities and Exposures*) e una gravità. Uno **scanner** confronta il contenuto dell'immagine con i
database delle vulnerabilità.

**Trivy**, progetto open source di Aqua Security, è tra i più usati. Il modo più semplice per usarlo è... un
container:

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy:0.75.0 image --severity HIGH,CRITICAL tasks-api:dev
```

Montare il socket di Docker permette a Trivy di leggere le immagini locali. (Controlla l'ultima versione nella
[pagina dei rilasci](https://github.com/aquasecurity/trivy/releases).)

Come leggere il risultato:

- la tabella elenca, per ogni pacchetto, la CVE, la **gravità**, la versione installata e quella che **corregge** il
  problema (*Fixed Version*);
- se c'è una versione corretta, **ricostruisci** l'immagine (con `docker build --pull` per scaricare l'immagine di
  base aggiornata) o aggiorna la libreria;
- se la correzione non esiste ancora, valuta se la vulnerabilità ti riguarda davvero (il componente è usato? è
  esposto?). `--ignore-unfixed` mostra solo ciò che puoi correggere subito.

!!! tip "Meno software, meno vulnerabilità"
    Prova a confrontare il numero di vulnerabilità di `python:3.12` (immagine completa) e `python:3.12-slim`. La
    differenza è il motivo per cui si scelgono immagini minimali.

Trivy analizza anche file di configurazione (Dockerfile, Kubernetes, Terraform) e cerca segreti dimenticati:
`trivy config .` e `trivy fs --scanners secret .`. Lo userai di nuovo nella pipeline della Fase 7.

## Eseguire con il minimo dei privilegi

Anche un'immagine pulita va eseguita con prudenza. Le opzioni principali di `docker run`, valide anche in Compose:

```bash
docker run -d --name api \
  --user 10001:10001 \
  --read-only --tmpfs /tmp \
  --cap-drop ALL \
  --security-opt no-new-privileges:true \
  --memory 256m --cpus 0.5 --pids-limit 100 \
  -p 8000:8000 tasks-api:dev
```

| Opzione | Protezione |
|---|---|
| `--user` (o `USER` nel Dockerfile) | il processo non è root |
| `--read-only` | il filesystem del container è in sola lettura: un attaccante non può modificare i programmi |
| `--tmpfs /tmp` | una zona scrivibile in memoria, se l'applicazione ne ha bisogno |
| `--cap-drop ALL` | toglie tutte le capability di Linux; con `--cap-add` si aggiunge solo ciò che serve |
| `no-new-privileges` | impedisce di acquisire privilegi con programmi setuid (come `sudo`) |
| `--memory`, `--cpus`, `--pids-limit` | un container impazzito non blocca l'intero host |

Da **evitare**:

- `--privileged`: dà al container quasi tutti i poteri dell'host;
- montare `/var/run/docker.sock` in container di cui non ti fidi completamente: chi controlla il socket controlla
  l'host (Trivy, sopra, è un'eccezione consapevole, per uno strumento noto e un'esecuzione temporanea);
- `--network host` senza un motivo preciso.

In Compose le stesse opzioni si scrivono così:

```yaml
services:
  api:
    read_only: true
    tmpfs: [/tmp]
    cap_drop: [ALL]
    security_opt: ["no-new-privileges:true"]
```

## La supply chain del software

Un attacco alla **supply chain** colpisce non la tua applicazione, ma ciò da cui dipende: una libreria, un'immagine
di base, uno strumento di build. Due pratiche si stanno diffondendo come standard:

- **SBOM** (*Software Bill of Materials*): l'elenco completo dei componenti di un'immagine, con le versioni. Quando
  esce una nuova vulnerabilità, sai subito quali immagini sono coinvolte. Trivy lo genera con
  `trivy image --format cyclonedx tasks-api:dev`.
- **Firma delle immagini**: con strumenti come **Cosign** (progetto Sigstore) chi pubblica firma l'immagine, e chi la
  usa verifica che provenga davvero da lui e non sia stata alterata.

Non serve padroneggiarli ora: sappi che esistono e a cosa servono, perché sono domande frequenti nei colloqui.

## Esercizi

!!! esercizio "Esercizio 1 · Pubblica la tua immagine"
    Pubblica la Tasks API su GitHub Container Registry con due tag: `1.0.0` e `sha-` seguito dall'hash breve del
    commit. Rendi pubblico il pacchetto. Poi elimina l'immagine locale e avviala scaricandola dal registry.

??? success "Soluzione"
    ```bash
    gh auth refresh -h github.com -s write:packages
    gh auth token | docker login ghcr.io -u TUO-UTENTE --password-stdin
    IMG=ghcr.io/TUO-UTENTE/tasks-api
    docker build -t $IMG:1.0.0 -t $IMG:sha-$(git rev-parse --short HEAD) .
    docker push $IMG:1.0.0
    docker push $IMG:sha-$(git rev-parse --short HEAD)
    # su GitHub: Packages → tasks-api → Package settings → Change visibility → Public
    docker rmi $IMG:1.0.0
    docker run -d --rm -p 8000:8000 $IMG:1.0.0      # la scarica dal registry
    ```
    **Spiegazione.** `docker build` accetta più `-t`. Il secondo push è velocissimo: gli strati sono già nel registry,
    cambia solo il tag. Il nome dell'immagine deve essere tutto in **minuscolo**.

!!! esercizio "Esercizio 2 · Confronta le immagini di base"
    Costruisci due versioni della Tasks API, una da `python:3.12` e una da `python:3.12-slim`. Confronta dimensioni e
    numero di vulnerabilità HIGH e CRITICAL con Trivy. Quale scegli e perché?

??? success "Soluzione"
    ```bash
    sed 's/python:3.12-slim/python:3.12/' Dockerfile > Dockerfile.full
    docker build -t tasks-api:full -f Dockerfile.full .
    docker build -t tasks-api:slim .
    docker images | grep tasks-api
    for t in full slim; do
      docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:0.75.0 \
        image --quiet --severity HIGH,CRITICAL --format json tasks-api:$t \
        | python3 -c 'import json,sys; d=json.load(sys.stdin); print(sum(len(r.get("Vulnerabilities") or []) for r in d.get("Results",[])))'
    done
    ```
    **Spiegazione.** L'immagine completa contiene compilatori, librerie di sviluppo e molti altri pacchetti: è diverse
    volte più grande e ha molte più vulnerabilità. Si sceglie `slim`, salvo che servano strumenti di compilazione, e
    in quel caso si usa una build multi-stage. I numeri esatti cambiano ogni giorno, con la pubblicazione di nuove CVE
    e di nuove immagini.

!!! esercizio "Esercizio 3 · Indurisci lo stack"
    Applica al servizio `api` del tuo `compose.yaml` le opzioni `read_only`, `cap_drop`, `no-new-privileges` e un
    limite di memoria di 256 MB. Verifica che lo stack funzioni ancora e che dentro il container non si possa
    scrivere in `/app`.

??? success "Soluzione"
    ```yaml
      api:
        # ... resto invariato
        read_only: true
        tmpfs: [/tmp]
        cap_drop: [ALL]
        security_opt: ["no-new-privileges:true"]
        mem_limit: 256m
    ```
    ```bash
    docker compose up -d
    curl -s localhost:8080/health
    docker compose exec api touch /app/prova      # Read-only file system
    ```
    **Spiegazione.** L'API usa PostgreSQL, quindi non scrive su disco: il filesystem in sola lettura non le crea
    problemi. Gunicorn ha bisogno di una cartella temporanea scrivibile per i suoi file di controllo dei worker, e la
    trova in `/tmp` grazie a `tmpfs`. Se un'applicazione fallisce con `read_only`, i log dicono quale percorso le serve
    scrivibile: lo aggiungi come `tmpfs` o come volume.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale riferimento garantisce di usare sempre esattamente la stessa immagine?","options":["Il tag latest","Il tag di versione","Il digest sha256","Il nome del repository"],"answer":[2],"explanation":"I tag possono essere spostati, il digest identifica il contenuto."},
 {"q":"Trivy segnala una CVE HIGH con una Fixed Version disponibile nella libreria di sistema dell'immagine di base. Cosa fai?","options":["Ignori la segnalazione","Ricostruisci l'immagine scaricando l'immagine di base aggiornata (docker build --pull)","Disattivi Trivy","Aggiungi --privileged"],"answer":[1],"explanation":"Le immagini vanno ricostruite regolarmente per incorporare le correzioni."},
 {"q":"Quali opzioni riducono i privilegi di un container? (scegli 2)","options":["--cap-drop ALL","--read-only","--privileged","-v /var/run/docker.sock:/var/run/docker.sock","--network host"],"answer":[0,1],"explanation":"--privileged e il socket di Docker danno al container il controllo dell'host."},
 {"q":"Cos'è uno SBOM?","options":["Un tipo di firewall","L'elenco completo dei componenti software di un'immagine, con le versioni","Un registry","Un formato di log"],"answer":[1],"explanation":"Permette di sapere subito quali immagini sono colpite da una nuova vulnerabilità."},
 {"q":"Perché si usa `docker login --password-stdin`?","options":["È più veloce","Evita che il token finisca nella cronologia della shell e nell'elenco dei processi","È obbligatorio per ghcr.io","Cripta l'immagine"],"answer":[1],"explanation":"Un segreto passato come argomento è visibile in giro."},
 {"q":"Un'immagine pubblicata su ghcr.io da un account personale è, per impostazione predefinita:","options":["Pubblica","Privata","Cancellata dopo 30 giorni","Visibile solo nell'organizzazione Docker"],"answer":[1],"explanation":"La visibilità si cambia nelle impostazioni del pacchetto."}
]}
</script>

## Fonti

- [GitHub Docs: Working with the Container registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Docker Docs: docker image tag](https://docs.docker.com/reference/cli/docker/image/tag/)
- [OCI: annotazioni standard delle immagini](https://github.com/opencontainers/image-spec/blob/main/annotations.md)
- [Trivy: documentazione](https://trivy.dev/latest/docs/)
- [Docker Docs: Docker security](https://docs.docker.com/engine/security/)
- [OWASP: Docker Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Docker_Security_Cheat_Sheet.html)
- [CISA: Software Bill of Materials](https://www.cisa.gov/sbom)
- [Sigstore: Cosign](https://docs.sigstore.dev/cosign/signing/overview/)
