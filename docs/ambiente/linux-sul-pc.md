# Non hai ancora Linux?

!!! obiettivi "Obiettivi"
    - Capire perché il corso usa Linux anche sul tuo computer.
    - Scegliere tra Ubuntu e Fedora.
    - Installare Linux in sicurezza, oppure seguire il corso dal browser.

## Perché Linux anche sul tuo computer

I server del cloud, i container e Kubernetes girano quasi sempre su **Linux**. Se lo usi ogni giorno anche sul tuo
computer, ogni comando che scrivi è allenamento: aprire un terminale, installare un pacchetto, leggere un log
diventano gesti naturali, e nei colloqui si vede. Per questo il corso dà per scontato un computer con
**Ubuntu** o **Fedora** e non riporta istruzioni per altri sistemi operativi.

## Ubuntu o Fedora?

Vanno bene entrambe: il corso le tratta alla pari, con schede affiancate dove i comandi cambiano.

| | Ubuntu Desktop | Fedora Workstation |
|---|---|---|
| Chi la sviluppa | Canonical | la comunità Fedora, con il sostegno di Red Hat |
| Famiglia | Debian (pacchetti `.deb`, comando `apt`) | Red Hat (pacchetti `.rpm`, comando `dnf`) |
| Ritmo delle versioni | una ogni sei mesi; ogni due anni una **LTS** con 5 anni di aggiornamenti | una ogni sei mesi, ognuna supportata circa 13 mesi |
| Ambiente grafico | GNOME | GNOME |
| Container già pronti | Docker si installa a parte | **Podman** è già installato |
| Si sposa con | il percorso Debian/Ubuntu e le certificazioni LFCS e Canonical | il percorso RHEL e le certificazioni RHCSA e RHCE |

!!! tip "Se non sai quale scegliere"
    Prendi **Ubuntu 24.04 LTS** se vuoi la massima stabilità e la documentazione più ampia. Prendi **Fedora** se
    punti alle certificazioni Red Hat o se vuoi strumenti sempre aggiornati. Le differenze tra le due famiglie sono
    spiegate nella [lezione 0.2](../fase-0/02-sistemi-operativi.md) e nella [lezione 1.5](../fase-1/05-pacchetti.md).

## Come installarlo

Prima di tutto, tre regole:

1. **Fai un backup** di tutti i tuoi file su un disco esterno o in un servizio di archiviazione online.
   L'installazione modifica il disco: un errore può cancellare i dati.
2. Ti serve una **chiavetta USB da almeno 8 GB**, che verrà cancellata.
3. Controlla che il computer rispetti i [requisiti dell'ambiente di laboratorio](index.md#requisiti-del-computer).

### 1. Prepara la chiavetta

=== "Ubuntu"

    1. Scarica l'immagine ISO di **Ubuntu Desktop** dalla [pagina ufficiale](https://ubuntu.com/download/desktop).
    2. Scrivila sulla chiavetta con [balenEtcher](https://etcher.balena.io/) o con un altro programma di scrittura
       di immagini, come indicato nella [guida ufficiale](https://ubuntu.com/tutorials/install-ubuntu-desktop).

=== "Fedora"

    1. Scarica **Fedora Media Writer** dalla [pagina ufficiale](https://fedoraproject.org/workstation/download):
       il programma scarica Fedora Workstation e scrive la chiavetta in un solo passaggio.
    2. Avvialo, scegli *Download automatically*, poi **Fedora Workstation** e la chiavetta.

### 2. Prova senza installare

Riavvia il computer dalla chiavetta: di solito si sceglie il dispositivo di avvio premendo un tasto all'accensione
(spesso ++f12++, ++f11++ o ++esc++; il manuale del computer indica quello giusto). Sia Ubuntu sia Fedora partono in
modalità **live**: il sistema gira dalla chiavetta e il disco non viene toccato. Controlla che funzionino rete
Wi-Fi, audio e schermo.

### 3. Installa

Dalla sessione live avvia l'installazione e scegli una delle due strade:

| Strada | Cosa succede | Per chi |
|---|---|---|
| **Usa tutto il disco** | Linux sostituisce il sistema attuale. È la scelta più semplice e pulita. | chi ha un computer dedicato allo studio, o ha fatto il backup e vuole passare del tutto a Linux |
| **Accanto al sistema attuale** (*dual boot*) | Il disco viene diviso in due e a ogni accensione scegli quale sistema avviare. | chi deve continuare a usare il sistema attuale |

!!! warning "Dual boot e cifratura del disco"
    Se il sistema attuale ha il disco **cifrato**, l'installatore non può ridurne lo spazio. Prima di procedere
    procurati la chiave di ripristino della cifratura e segui le indicazioni della guida ufficiale della tua
    distribuzione. Nel dubbio, usa un computer o un disco dedicato.

Segui la guida ufficiale passo per passo:

- [Ubuntu: Install Ubuntu Desktop](https://ubuntu.com/tutorials/install-ubuntu-desktop)
- [Fedora: Installation Guide](https://docs.fedoraproject.org/en-US/fedora/latest/getting-started/)

### 4. Dopo l'installazione

Aggiorna subito il sistema dal terminale:

=== "Ubuntu"

    ```bash
    sudo apt update && sudo apt full-upgrade -y
    ```

=== "Fedora"

    ```bash
    sudo dnf upgrade --refresh -y
    ```

Non preoccuparti se oggi questi comandi ti sembrano oscuri: li capirai nella [lezione 0.4](../fase-0/04-terminale.md) e nella
[lezione 1.5](../fase-1/05-pacchetti.md).

## Non puoi installare Linux? Il corso dal browser

Se il computer non è tuo o non vuoi modificarlo, puoi seguire quasi tutto il corso nel **browser**, con due servizi
gratuiti:

| Servizio | Cosa ti dà | Fasi |
|---|---|---|
| [**GitHub Codespaces**](devcontainer.md#strada-a-github-codespaces-nessuna-installazione) | il Dev Container del corso, con editor e terminale, su un server di GitHub | 2-8 |
| [**Killercoda**](vm-linux.md#alternativa-nel-browser-killercoda) | una macchina Ubuntu completa per un'ora alla volta | 0-1 |

Il limite è il tempo: Codespaces ha una quota mensile gratuita e Killercoda cancella la macchina a fine sessione.
Appena puoi, passa a un computer con Ubuntu o Fedora.

## Fonti

- [Ubuntu: Install Ubuntu Desktop](https://ubuntu.com/tutorials/install-ubuntu-desktop)
- [Ubuntu: rilasci e cicli di supporto](https://ubuntu.com/about/release-cycle)
- [Fedora Workstation: download e Fedora Media Writer](https://fedoraproject.org/workstation/download)
- [Fedora Docs: Getting Started](https://docs.fedoraproject.org/en-US/fedora/latest/getting-started/)
- [Fedora: ciclo di vita delle versioni](https://docs.fedoraproject.org/en-US/releases/lifecycle/)
