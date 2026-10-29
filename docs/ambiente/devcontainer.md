# Il Dev Container

!!! obiettivi "Obiettivi"
    - Capire cos'è un Dev Container e perché lo usiamo.
    - Aprire l'ambiente di laboratorio del corso sul tuo computer **oppure** nel browser con GitHub Codespaces.
    - Verificare che tutti gli strumenti funzionino.

## Cos'è e perché lo usiamo

Un **Dev Container** (*development container*, contenitore di sviluppo) è un ambiente di lavoro completo descritto
in un file di testo. Il file elenca il sistema operativo di base e tutti gli strumenti da installare, con le loro
versioni. Chiunque apra il progetto ottiene **esattamente lo stesso ambiente**, che usi Ubuntu, Fedora o il browser.

Questo risolve il problema più comune di chi studia informatica: "a me non funziona". Con il Dev Container,
se un comando funziona nel corso, funziona anche da te.

La descrizione dell'ambiente del corso si trova nella cartella
[`.devcontainer`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/.devcontainer) del repository.
Contiene, tra gli altri:

| Strumento | A cosa serve | Fase |
|---|---|---|
| Git e GitHub CLI | Versionare il codice | 2 |
| Python 3.12, Flask, pytest | Automazione e applicazione d'esempio | 2 |
| Docker (*docker-in-docker*) | Container | 3 |
| AWS CLI e Google Cloud CLI | Gestire il cloud da riga di comando | 4-5 |
| Terraform 1.12 e Ansible | Infrastructure as Code | 6 |
| kubectl, Helm e kind | Kubernetes | 8 |

## Strada A: GitHub Codespaces (nessuna installazione)

**GitHub Codespaces** avvia il Dev Container su un server di GitHub e te lo mostra nel browser, con l'editor
VS Code. È la strada più semplice.

1. Crea un account gratuito su [github.com](https://github.com/signup), se non lo hai già.
2. Apri il [repository del corso](https://github.com/imredracoon/racoon-cloud-devops-academy) e premi **Fork**
   in alto a destra: così hai una copia tua, dove puoi salvare il tuo lavoro.
3. Nella tua copia premi il pulsante verde **Code**, scegli la scheda **Codespaces** e poi **Create codespace on main**.
4. Attendi qualche minuto: la prima volta GitHub costruisce l'ambiente. Alla fine, nel pannello in basso
   (il **terminale**) compare l'elenco delle versioni degli strumenti.

!!! warning "La quota gratuita"
    Gli account personali GitHub Free includono ogni mese **120 ore di calcolo** e **15 GB di spazio** per
    Codespaces. Le ore si consumano moltiplicate per i core della macchina: una macchina da 4 core consuma 4 ore
    di quota per ogni ora di utilizzo, quindi hai circa **30 ore di lavoro al mese** su 4 core.
    Per non sprecare quota:

    - **ferma** il codespace quando hai finito (menu **Codespaces** in basso a sinistra → *Stop Current Codespace*);
    - cancella i codespace che non usi più da [github.com/codespaces](https://github.com/codespaces).

    GitHub ferma comunque un codespace inattivo dopo 30 minuti. Con quota esaurita il servizio si blocca, non
    addebita nulla, a meno che tu non abbia impostato un limite di spesa.

## Strada B: Dev Container sul tuo computer

Questa strada non ha limiti di ore, ma richiede di installare due programmi.

### 1. Installa Docker Engine

Docker è il programma che esegue i container. Lo studierai a fondo nella [Fase 3](../fase-3/index.md). Installiamo
**Docker Engine** dal repository ufficiale di Docker, che offre sempre la versione aggiornata.

=== "Ubuntu"

    ```bash
    # 1. strumenti e chiave di firma del repository Docker
    sudo apt update
    sudo apt install -y ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # 2. il repository
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
      https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
      | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    # 3. Docker Engine con i plugin Buildx e Compose
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    ```

    Su Ubuntu il servizio parte da solo dopo l'installazione.

=== "Fedora"

    ```bash
    # 1. il repository Docker
    sudo dnf config-manager addrepo --from-repofile=https://download.docker.com/linux/fedora/docker-ce.repo

    # 2. Docker Engine con i plugin Buildx e Compose
    sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # 3. avvia il servizio ora e a ogni accensione
    sudo systemctl enable --now docker
    ```

    Durante l'installazione `dnf` chiede di accettare la chiave di firma: controlla che l'impronta coincida con
    quella della [guida ufficiale](https://docs.docker.com/engine/install/fedora/) e rispondi `y`.

    !!! note "E Podman?"
        Fedora include già **Podman**, che studierai nella [lezione 3.7](../fase-3/07-podman.md). Per il Dev Container
        del corso usiamo comunque Docker: l'ambiente contiene a sua volta Docker (*docker-in-docker*), una
        configurazione che con Docker funziona senza ritocchi.

Poi, su entrambe, aggiungi il tuo utente al gruppo `docker` per usare Docker senza `sudo`, ed esci e rientra
nella sessione perché la modifica abbia effetto:

```bash
sudo usermod -aG docker $USER
```

!!! warning "Il gruppo `docker` equivale a root"
    Chi appartiene al gruppo `docker` può ottenere i privilegi di amministratore sul computer. Va bene sul tuo
    computer personale di studio; su un server condiviso no. Ne riparlerai nella [lezione 3.6](../fase-3/06-registry-sicurezza.md).

### 2. Installa Visual Studio Code e l'estensione Dev Containers

=== "Ubuntu"

    ```bash
    sudo snap install code --classic
    ```

=== "Fedora"

    ```bash
    sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc
    printf '%s\n' '[code]' 'name=Visual Studio Code' 'baseurl=https://packages.microsoft.com/yumrepos/vscode' \
      'enabled=1' 'gpgcheck=1' 'gpgkey=https://packages.microsoft.com/keys/microsoft.asc' \
      | sudo tee /etc/yum.repos.d/vscode.repo > /dev/null
    sudo dnf install -y code
    ```

In alternativa puoi scaricare il pacchetto `.deb` o `.rpm` dalla [pagina ufficiale](https://code.visualstudio.com/docs/setup/linux).
Poi apri VS Code, vai nel pannello **Extensions** (icona dei quadratini a sinistra) e installa
**Dev Containers** di Microsoft.

### 3. Apri il corso nel container

1. Fai il **Fork** del repository del corso come nella strada A.
2. In VS Code premi ++f1++, scrivi `Dev Containers: Clone Repository in Container Volume` e incolla l'indirizzo
   della tua copia, per esempio `https://github.com/TUO-UTENTE/racoon-cloud-devops-academy`.
3. Attendi la costruzione dell'ambiente: la prima volta richiede 5-10 minuti.

## Verifica che tutto funzioni

Apri il terminale (menu **Terminal → New Terminal**) e scrivi:

```bash
docker run --rm hello-world
```

Se compare il messaggio `Hello from Docker!`, l'ambiente è pronto. Il comando scarica una piccola immagine di
prova e la esegue: capirai ogni parola nella [Fase 3](../fase-3/index.md).

## Problemi frequenti

??? question "La costruzione si interrompe con un errore di rete"
    Riprova: premi ++f1++ e scegli `Dev Containers: Rebuild Container`. Se il problema persiste, controlla di non
    essere dietro una rete aziendale che blocca i download.

??? question "`permission denied while trying to connect to the Docker daemon socket`"
    Il tuo utente non è ancora nel gruppo `docker`, oppure non sei uscito e rientrato dopo `usermod`. Controlla con
    `groups`: deve comparire `docker`. Se non vuoi chiudere la sessione, apri una shell con il gruppo attivo:
    `newgrp docker`.

??? question "Il computer rallenta durante la costruzione"
    Con Docker Engine i container usano direttamente la memoria del tuo computer, senza un limite fisso. Chiudi le
    applicazioni pesanti mentre l'ambiente si costruisce e verifica di avere almeno 8 GB di RAM
    ([requisiti](index.md#requisiti-del-computer)).

## Fonti

- [Development Containers, specifica ufficiale](https://containers.dev/)
- [Visual Studio Code: Developing inside a Container](https://code.visualstudio.com/docs/devcontainers/containers)
- [GitHub Docs: GitHub Codespaces](https://docs.github.com/en/codespaces)
- [GitHub Docs: come viene misurato l'uso di Codespaces](https://docs.github.com/en/billing/concepts/product-billing/github-codespaces)
- [Docker Docs: installare Docker Engine su Ubuntu](https://docs.docker.com/engine/install/ubuntu/)
- [Docker Docs: installare Docker Engine su Fedora](https://docs.docker.com/engine/install/fedora/)
- [Docker Docs: passi dopo l'installazione su Linux](https://docs.docker.com/engine/install/linux-postinstall/)
- [Visual Studio Code on Linux](https://code.visualstudio.com/docs/setup/linux)
