# 3.7 Podman in profondità

!!! obiettivi "Obiettivi"
    - Capire cosa significa davvero *rootless* e come Podman isola i container di ogni utente.
    - Raggruppare più container in un **pod** e farli parlare su `localhost`.
    - Trasformare un pod in un file YAML di Kubernetes e riavviarlo da quel file.
    - Far gestire i container a **systemd** con i **Quadlet**, con avvio al boot e aggiornamento automatico.
    - Scegliere consapevolmente tra Docker e Podman.

!!! laboratorio "Dove"
    Sul tuo computer con **Fedora** (Podman è già installato), su Ubuntu dopo averlo installato, oppure sulla VM
    `rhel-lab`. Il Dev Container non contiene Podman: questa lezione si fa fuori.

## Perché una lezione dedicata

Nelle lezioni precedenti hai usato Podman come "Docker senza demone": stessi comandi, stesse immagini OCI
([lezione 3.1](01-concetti.md#gli-standard-e-gli-strumenti)). Podman però ha tre caratteristiche che Docker non
ha, o ha in modo diverso, e che nel lavoro reale contano:

1. è nato **rootless**: ogni utente ha i suoi container, le sue immagini e i suoi volumi, senza privilegi di root;
2. ha i **pod**, lo stesso concetto di Kubernetes, e sa leggere e scrivere YAML di Kubernetes;
3. si integra con **systemd**: un container diventa un servizio come `nginx` o `sshd` ([lezione 1.4](../fase-1/04-processi-systemd.md)).

È lo strumento predefinito di Fedora e di tutta la famiglia Red Hat, ed è il motore che `ansible-navigator` usa
per gli *execution environment* dell'esame RHCE ([appendice RHCE](../appendici/rhce/argomenti.md)).

## Installazione

=== "Fedora"

    È già installato. Verifica:

    ```bash
    podman version
    ```

=== "Ubuntu"

    È nei repository ufficiali:

    ```bash
    sudo apt install -y podman
    podman version
    ```

    Ubuntu 24.04 distribuisce Podman 4.9: tutto quello che vedi in questa lezione funziona, tranne i file `.pod`
    dei Quadlet, arrivati con Podman 5 e segnalati più avanti.

=== "RHEL (Rocky/AlmaLinux)"

    Sulla VM `rhel-lab` lo hai installato nella [lezione 3.2](02-immagini-container.md): `sudo dnf install -y podman`.

## Rootless, sul serio

Avvia un container come utente normale e guarda chi lo esegue davvero:

```bash
podman run -d --name prova docker.io/library/nginx:1.30-alpine
podman top prova user huser       # (1)!
ps -o user,pid,cmd -C nginx | head -3   # (2)!
```

1.  `user` è l'utente **dentro** il container (`root` per il processo principale di Nginx), `huser` quello **sul
    sistema**: il tuo nome utente.
2.  Visto dal sistema, Nginx gira con il tuo utente, non come root.

Il trucco si chiama **user namespace** ([lezione 3.1](01-concetti.md)): dentro il container l'utente 0 (root) esiste,
ma corrisponde al tuo utente sul sistema. Gli altri UID del container vengono presi da un intervallo riservato a te,
scritto in `/etc/subuid` e `/etc/subgid`:

```console
$ grep "^$USER:" /etc/subuid
vincenzo:524288:65536
$ podman unshare cat /proc/self/uid_map
         0       1000          1
         1     524288      65536
```

Si legge così: l'UID 0 del container è il tuo UID 1000; gli UID da 1 in su del container corrispondono agli UID
da 524288 in su sul sistema, che non appartengono a nessun utente reale. Se un processo riesce a uscire dal
container, si ritrova con i permessi di un utente qualunque, non di root.

| | Rootless (`podman`) | Rootful (`sudo podman`) |
|---|---|---|
| Immagini e container | in `~/.local/share/containers/` | in `/var/lib/containers/` |
| Configurazione | `~/.config/containers/` | `/etc/containers/` |
| Porte sotto la 1024 | no ([lezione 3.4](04-volumi-reti.md#podman-e-le-reti)) | sì |
| Chi vede i container | solo il tuo utente | solo root |

!!! tip "Il proprietario dei file nei volumi"
    Un file creato da root dentro il container appare tuo sull'host; un file creato dall'utente `nginx` del
    container (UID 101) appare di proprietà di un UID alto, come 524388. Per entrare nel "punto di vista" del
    container usa `podman unshare`, per esempio `podman unshare ls -ln ~/.local/share/containers/storage/volumes`.

Fai pulizia: `podman rm -f prova`.

### I nomi completi delle immagini

Podman non presume che un'immagine venga da Docker Hub. Con un nome corto come `nginx` consulta l'elenco in
`/etc/containers/registries.conf` e, se ci sono più registry possibili, ti chiede quale usare. Negli script e nei
file di configurazione scrivi sempre il **nome completo**: `docker.io/library/nginx:1.30-alpine`. È più sicuro (non
scarichi per errore un'immagine omonima da un altro registry) ed è obbligatorio per l'aggiornamento automatico.

## I pod

Un **pod** è un gruppo di container che condividono la stessa rete: si parlano su `localhost` e le porte si
pubblicano una volta sola, sul pod. È esattamente il Pod di Kubernetes che studierai nella
[Fase 8](../fase-8/index.md). Ogni pod contiene un piccolo container *infra* che tiene in vita la rete comune.

```mermaid
flowchart LR
  H[Host :8080] --> P
  subgraph P[Pod tasks · localhost condiviso]
    N[proxy<br>Nginx :80] --> A[api<br>:8000]
    A --> D[db<br>PostgreSQL :5432]
  end
```

Ricreiamo lo stack della [lezione 3.5](05-compose.md#lo-stack-della-tasks-api) come pod. Lavora nel repository
`tasks-api`, che contiene il `Dockerfile` della [lezione 3.3](03-dockerfile.md).

```bash
podman build -t tasks-api:dev .                    # (1)!
podman volume create tasks-pgdata
podman pod create --name tasks -p 8080:80          # (2)!

podman run -d --pod tasks --name db \
  -e POSTGRES_DB=tasks -e POSTGRES_USER=tasks -e POSTGRES_PASSWORD=cambiami \
  -v tasks-pgdata:/var/lib/postgresql/data docker.io/library/postgres:17

podman run -d --pod tasks --name api \
  -e DATABASE_URL=postgresql+psycopg://tasks:cambiami@127.0.0.1:5432/tasks localhost/tasks-api:dev   # (3)!
```

1.  L'immagine costruita in locale si chiama `localhost/tasks-api:dev`: anche qui Podman usa nomi completi.
2.  Le porte si pubblicano quando crei il pod, non sui singoli container.
3.  L'API trova il database su `127.0.0.1`, perché condividono la rete. Nello stack Compose usava il nome `db`.

Per Nginx serve una configurazione che punti a `127.0.0.1:8000` invece che a `api:8000`:

```nginx title="nginx-pod.conf"
# Reverse proxy verso la Tasks API dentro un pod (lezione 3.7)
server {
    listen 80;
    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

```bash
podman run -d --pod tasks --name proxy \
  -v ./nginx-pod.conf:/etc/nginx/conf.d/default.conf:ro,Z docker.io/library/nginx:1.30-alpine
podman pod ps                    # il pod, con 4 container (3 + infra)
curl -s localhost:8080/health
```

L'opzione `:Z` serve su Fedora e RHEL, dove SELinux è attivo ([lezione 3.4](04-volumi-reti.md)); su Ubuntu è
innocua.

| Comando | Cosa fa |
|---|---|
| `podman pod ps` | elenca i pod |
| `podman ps --pod` | elenca i container con il pod a cui appartengono |
| `podman pod stop tasks` / `start tasks` | ferma o avvia tutti i container del pod |
| `podman pod rm -f tasks` | elimina il pod e i suoi container (non i volumi) |

## Dal pod al YAML di Kubernetes

Podman sa **descrivere** un pod in esecuzione come file YAML di Kubernetes, e sa **avviare** un pod a partire da
quel file:

```bash
podman kube generate tasks > tasks-pod.yaml     # dal pod al file
podman pod rm -f tasks                          # elimina il pod
podman kube play tasks-pod.yaml                 # lo ricrea dal file
podman kube down tasks-pod.yaml                 # lo ferma e lo elimina
```

Apri `tasks-pod.yaml`: troverai un oggetto `kind: Pod` con i tre container, le variabili d'ambiente e il volume come
`persistentVolumeClaim`. Non è ancora un file pronto per un cluster (la password è in chiaro e l'immagine
`localhost/...` esiste solo sul tuo computer), ma è un ottimo punto di partenza: nella [Fase 8](../fase-8/index.md)
riconoscerai ogni riga.

## I Quadlet: container gestiti da systemd

Con Docker è il demone a riavviare i container (le politiche `restart` della [lezione 3.2](02-immagini-container.md)).
Podman non ha un demone: delega a **systemd**, il gestore dei servizi che conosci dalla
[lezione 1.4](../fase-1/04-processi-systemd.md). Il modo moderno per farlo sono i **Quadlet**: file di testo simili alle unità
systemd che descrivono un container; al caricamento systemd li trasforma in servizi veri.

| Estensione | Descrive |
|---|---|
| `.container` | un container |
| `.volume` | un volume con nome |
| `.network` | una rete |
| `.kube` | un pod avviato da un file YAML di Kubernetes |
| `.pod` | un pod (Podman 5 o successivo) |

I file dell'utente stanno in `~/.config/containers/systemd/`; quelli di sistema, per i container rootful, in
`/etc/containers/systemd/`.

### Un sito statico come servizio

```bash
mkdir -p ~/sito ~/.config/containers/systemd
echo '<h1>Servito da un Quadlet</h1>' > ~/sito/index.html
```

```ini title="~/.config/containers/systemd/sito.container"
# Sito statico servito da Nginx come servizio utente (lezione 3.7)
[Unit]
Description=Sito statico del laboratorio

[Container]
Image=docker.io/library/nginx:1.30-alpine
ContainerName=sito
PublishPort=8081:80
Volume=%h/sito:/usr/share/nginx/html:ro,Z
AutoUpdate=registry

[Service]
Restart=on-failure

[Install]
WantedBy=default.target
```

| Riga | Significato |
|---|---|
| `Image=` | l'immagine, con nome completo |
| `PublishPort=` | come `-p` |
| `Volume=` | come `-v`; `%h` è la tua home |
| `AutoUpdate=registry` | consente l'aggiornamento automatico quando nel registry compare una nuova immagine con lo stesso tag |
| `Restart=on-failure` | systemd lo riavvia se si ferma con un errore |
| `WantedBy=default.target` | lo avvia all'avvio della tua sessione |

Carica e avvia il servizio:

```bash
systemctl --user daemon-reload        # (1)!
systemctl --user start sito
systemctl --user status sito
curl -s localhost:8081
journalctl --user -u sito -n 20       # i log del container, come per ogni servizio
```

1.  A ogni `daemon-reload` systemd rigenera i servizi a partire dai file Quadlet. Il servizio si chiama come il
    file, senza estensione: `sito.service`.

!!! warning "Non serve `systemctl enable`"
    I servizi generati dai Quadlet non si abilitano con `enable`: l'avvio automatico lo decide la sezione
    `[Install]` del file. Se il servizio non compare, controlla la sintassi con
    `/usr/libexec/podman/quadlet -dryrun -user`.

### Avvio senza che tu sia collegato

I servizi `--user` partono quando accedi e si fermano quando esci. Su un server, dove nessuno resta collegato,
abilita il **lingering**: i servizi dell'utente partono al boot e restano attivi.

```bash
sudo loginctl enable-linger $USER
```

### Aggiornamento automatico

Podman può controllare i registry e aggiornare i container con `AutoUpdate=registry`. Se il nuovo container non
parte, torna alla versione precedente.

```bash
podman auto-update --dry-run                              # cosa verrebbe aggiornato
systemctl --user enable --now podman-auto-update.timer    # controllo giornaliero
```

In produzione è comodo per i servizi semplici, ma va usato con giudizio: un tag mobile come `latest` può portarti
in casa una versione nuova senza che tu l'abbia provata ([lezione 3.6](06-registry-sicurezza.md#la-strategia-dei-tag)).

## Docker o Podman?

| | Docker Engine | Podman |
|---|---|---|
| Architettura | un demone (`dockerd`) che gestisce tutti i container | nessun demone: ogni container è un processo figlio |
| Senza root | modalità *rootless* opzionale | predefinito |
| Avvio automatico | politiche `restart` del demone | systemd, con i Quadlet |
| Più container insieme | Docker Compose | pod, `podman kube play`, oppure `podman compose` ([lezione 3.5](05-compose.md)) |
| Diffusione | il più usato in sviluppo e nelle pipeline | predefinito su Fedora e RHEL, comune in ambienti Red Hat e OpenShift |
| Immagini | OCI | OCI: le stesse |

Nel lavoro li incontrerai entrambi, e passare dall'uno all'altro richiede pochi minuti. Due ponti utili:

- il pacchetto `podman-docker` fornisce il comando `docker` che chiama Podman;
- `systemctl --user enable --now podman.socket` espone un'API compatibile con quella di Docker, che strumenti come
  `docker compose` possono usare.

## Esercizi

!!! esercizio "Esercizio 1 · Chi sei davvero?"
    Avvia `docker.io/library/alpine:3.22` con `podman run --rm -it ... sh`. Dentro il container esegui `id` e
    crea un file in un bind mount della tua home (`-v ~/prova:/dati:Z`). Fuori dal container controlla il
    proprietario del file con `ls -ln ~/prova`. Ripeti creando il file con l'utente `nobody` del container
    (`su nobody -s /bin/sh -c 'touch /dati/b'`, dopo aver dato i permessi con `chmod 777 /dati`). Spiega i due
    risultati.

??? success "Soluzione"
    ```bash
    mkdir -p ~/prova
    podman run --rm -it -v ~/prova:/dati:Z docker.io/library/alpine:3.22 sh
    # dentro il container
    id                     # uid=0(root) gid=0(root)
    touch /dati/a
    chmod 777 /dati
    su nobody -s /bin/sh -c 'touch /dati/b'
    exit
    ls -ln ~/prova
    ```
    **Spiegazione.** Il file `a` appartiene al tuo UID (per esempio 1000), perché il root del container corrisponde
    a te. Il file `b` appartiene a un UID alto (524288 + 65534 − 1 = 589821 con l'intervallo di esempio), perché
    l'utente `nobody` del container (UID 65534) viene mappato nell'intervallo di `/etc/subuid`. Nessuno dei due è
    root sul tuo sistema.

!!! esercizio "Esercizio 2 · La Tasks API in un pod"
    Ricrea lo stack in un pod come nella lezione, verifica che risponda su `localhost:8080`, poi esporta il pod in
    `tasks-pod.yaml`. Elimina il pod e ricrealo con `podman kube play`. I dati sono sopravvissuti?

??? success "Soluzione"
    ```bash
    curl -s -X POST localhost:8080/tasks -H 'Content-Type: application/json' -d '{"title":"Nel pod"}'
    podman kube generate tasks > tasks-pod.yaml
    podman pod rm -f tasks
    podman kube play tasks-pod.yaml
    curl -s localhost:8080/tasks          # l'attività c'è ancora
    ```
    **Spiegazione.** `pod rm` elimina i container ma non il volume `tasks-pgdata`. Il file YAML lo richiama come
    `persistentVolumeClaim` con lo stesso nome, e `kube play` riusa il volume esistente.

!!! esercizio "Esercizio 3 · Un servizio che sopravvive al riavvio"
    Crea il Quadlet `sito.container` della lezione, abilita il lingering e riavvia il computer (o la VM). Senza
    accedere alla sessione grafica, collegati via SSH e verifica che il sito risponda. Poi uccidi il processo
    principale del container con `podman kill sito` e osserva cosa fa systemd.

??? success "Soluzione"
    ```bash
    systemctl --user daemon-reload && systemctl --user start sito
    sudo loginctl enable-linger $USER
    sudo reboot
    # dopo il riavvio, via SSH
    curl -s localhost:8081
    podman kill sito
    sleep 5; systemctl --user status sito   # attivo di nuovo, con un riavvio registrato
    ```
    **Spiegazione.** Con il lingering i servizi dell'utente partono al boot. `podman kill` fa uscire il container
    con un errore, e `Restart=on-failure` fa ripartire il servizio: è systemd a riavviarlo, non Podman.

## Nei tre livelli

=== "Livello 1 · Un PC"

    Gli esercizi sulla workstation (o in `rhel-lab`) sono tutto ciò che serve.

=== "Livello 2 · Workstation + server"

    Metti in servizio su `srv1` il Quadlet del [sito statico](#un-sito-statico-come-servizio), rootless, con il tuo
    utente e `loginctl enable-linger`. Su `srv1` con Ubuntu i file `.pod` non sono disponibili (Podman 4.9), i
    `.container` sì.

    **Collaudo**

    | Prova | Risultato atteso |
    |---|---|
    | `ssh srv1 systemctl --user status sito` | `active (running)` |
    | riavvio di `srv1`, **senza** entrare con SSH, poi `curl -sI http://IP-SRV1:PORTA/` dalla workstation (apri la porta nel firewall) | `200`: il servizio è partito senza che tu fossi collegato |
    | `ssh srv1 podman auto-update --dry-run` | il container compare tra quelli controllati |

=== "Livello 3 · Mini data center"

    Metti in servizio il Quadlet del sito statico su **`nodo3`**, che nel progetto userà i Quadlet per tutto lo
    stack. AlmaLinux 10 ha Podman 5: puoi provare anche i file `.pod`. Poi lancia la stessa immagine su `nodo2` con
    Docker (`docker run -d -p 127.0.0.1:8081:80 ...`).

    **Collaudo**

    | Prova | Risultato atteso |
    |---|---|
    | riavvio di `nodo3` senza entrare, poi `ssh nodo2 curl -sI http://nodo3.lab:PORTA/` | `200` (porta aperta in firewalld) |
    | `ssh nodo3 podman images --digests` e `ssh nodo2 docker images --digests` | lo stesso digest per la stessa immagine |
    | `ssh nodo3 sudo ausearch -m AVC -ts boot` | nessun blocco di SELinux |

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"In un container rootless, a quale utente del sistema corrisponde l'utente root del container?","options":["A root","All'utente che ha avviato il container","All'utente nobody","A un utente casuale"],"answer":[1],"explanation":"Lo user namespace mappa l'UID 0 del container sul tuo UID; gli altri UID vengono da /etc/subuid."},
 {"q":"Dove trova Podman l'intervallo di UID riservato a un utente?","options":["/etc/passwd","/etc/subuid","/etc/containers/registries.conf","~/.bashrc"],"answer":[1],"explanation":"/etc/subuid e /etc/subgid assegnano a ogni utente un intervallo di UID e GID."},
 {"q":"In un pod, come raggiunge l'API il database?","options":["Con il nome db","Con 127.0.0.1 (localhost)","Con l'IP dell'host","Non può, servono due pod"],"answer":[1],"explanation":"I container di un pod condividono la stessa rete."},
 {"q":"Quale comando ricrea un pod a partire da un file YAML di Kubernetes?","options":["podman kube generate","podman kube play","podman pod create --yaml","podman compose up"],"answer":[1],"explanation":"generate crea il file, play lo esegue, down lo ferma."},
 {"q":"Dove metti un Quadlet per un container rootless del tuo utente?","options":["/etc/systemd/system/","~/.config/containers/systemd/","/var/lib/containers/","~/.local/bin/"],"answer":[1],"explanation":"I Quadlet di sistema vanno in /etc/containers/systemd/."},
 {"q":"Quali affermazioni sui Quadlet sono corrette?","options":["Dopo averli scritti serve systemctl --user daemon-reload","Si abilitano con systemctl enable","L'avvio automatico si decide nella sezione [Install]","Funzionano solo con Docker","Il servizio prende il nome del file"],"answer":[0,2,4],"explanation":"I servizi generati non si abilitano con enable: conta [Install]."},
 {"q":"Perché un servizio utente parta al boot anche senza che l'utente acceda, serve:","options":["sudo loginctl enable-linger UTENTE","systemctl enable --now podman","Restart=always","Eseguire il container con sudo"],"answer":[0],"explanation":"Il lingering mantiene attivo il gestore systemd dell'utente."}
]}
</script>

## Fonti

- [Podman: documentazione ufficiale](https://docs.podman.io/en/latest/)
- [Podman: tutorial rootless](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)
- [podman-pod-create(1)](https://docs.podman.io/en/latest/markdown/podman-pod-create.1.html)
- [podman-kube-generate(1)](https://docs.podman.io/en/latest/markdown/podman-kube-generate.1.html) e [podman-kube-play(1)](https://docs.podman.io/en/latest/markdown/podman-kube-play.1.html)
- [Quadlet: podman-systemd.unit(5)](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html)
- [podman-auto-update(1)](https://docs.podman.io/en/latest/markdown/podman-auto-update.1.html)
- [user_namespaces(7)](https://man7.org/linux/man-pages/man7/user_namespaces.7.html)
- [loginctl(1): enable-linger](https://www.freedesktop.org/software/systemd/man/latest/loginctl.html)
