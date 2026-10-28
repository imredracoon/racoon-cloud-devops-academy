# Il Dev Container

!!! obiettivi "Obiettivi"
    - Capire cos'è un Dev Container e perché lo usiamo.
    - Aprire l'ambiente di laboratorio del corso sul tuo computer **oppure** nel browser con GitHub Codespaces.
    - Verificare che tutti gli strumenti funzionino.

## Cos'è e perché lo usiamo

Un **Dev Container** (*development container*, contenitore di sviluppo) è un ambiente di lavoro completo descritto
in un file di testo. Il file elenca il sistema operativo di base e tutti gli strumenti da installare, con le loro
versioni. Chiunque apra il progetto ottiene **esattamente lo stesso ambiente**, che usi Windows, macOS o Linux.

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

### 1. Installa Docker

Docker è il programma che esegue i container. Lo studierai a fondo nella [Fase 3](../fase-3/index.md).

=== "Windows"

    1. Scarica e installa **Docker Desktop** dalla [pagina ufficiale](https://docs.docker.com/desktop/setup/install/windows-install/).
    2. Durante l'installazione lascia selezionata l'opzione **WSL 2**: Docker userà il sottosistema Linux di Windows.
    3. Riavvia il computer se richiesto, poi apri Docker Desktop e attendi che l'icona della balena diventi verde.

=== "macOS"

    1. Scarica **Docker Desktop** dalla [pagina ufficiale](https://docs.docker.com/desktop/setup/install/mac-install/),
       scegliendo la versione per il tuo processore (*Apple silicon* per i Mac M1/M2/M3/M4, *Intel* per gli altri).
    2. Trascina Docker nella cartella Applicazioni, aprilo e accetta le richieste di autorizzazione.

=== "Linux"

    Installa **Docker Engine** seguendo la [guida ufficiale](https://docs.docker.com/engine/install/) per la tua
    distribuzione, poi aggiungi il tuo utente al gruppo `docker` come indicato nei
    [passi post-installazione](https://docs.docker.com/engine/install/linux-postinstall/).

!!! note "Licenza di Docker Desktop"
    Docker Desktop è gratuito per uso personale, studio e piccole aziende. Le condizioni sono nella
    [pagina ufficiale](https://docs.docker.com/subscription/desktop-license/).

### 2. Installa Visual Studio Code e l'estensione Dev Containers

1. Scarica e installa [Visual Studio Code](https://code.visualstudio.com/).
2. Apri VS Code, vai nel pannello **Extensions** (icona dei quadratini a sinistra) e installa
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

??? question "Docker dice che non ha abbastanza memoria"
    In Docker Desktop apri **Settings → Resources** e assegna almeno 6 GB di memoria. Su Windows con WSL 2 il
    limite si imposta nel file `.wslconfig`, come spiegato nella
    [documentazione Microsoft](https://learn.microsoft.com/windows/wsl/wsl-config).

??? question "Su Windows compare un errore sulla virtualizzazione"
    La virtualizzazione deve essere attiva nel BIOS/UEFI del computer. La procedura varia per ogni produttore:
    cerca "abilitare virtualizzazione" insieme al modello del tuo computer sul sito del produttore.

## Fonti

- [Development Containers, specifica ufficiale](https://containers.dev/)
- [Visual Studio Code: Developing inside a Container](https://code.visualstudio.com/docs/devcontainers/containers)
- [GitHub Docs: GitHub Codespaces](https://docs.github.com/en/codespaces)
- [GitHub Docs: come viene misurato l'uso di Codespaces](https://docs.github.com/en/billing/concepts/product-billing/github-codespaces)
- [Docker Docs: installazione](https://docs.docker.com/get-started/get-docker/)
