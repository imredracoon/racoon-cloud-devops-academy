# 0.7 Laboratorio: prepara l'ambiente

!!! obiettivi "Obiettivi"
    - Creare l'account GitHub che userai per tutto il corso.
    - Creare e usare la VM Linux di laboratorio.
    - Aprire il Dev Container del corso (sul tuo computer o in Codespaces).
    - Raccogliere dal terminale le informazioni di sistema di una macchina.

!!! laboratorio "Laboratorio · circa 3 ore"
    **Dove:** il tuo computer (oppure GitHub Codespaces e Killercoda).
    **Risultato:** account GitHub attivo, VM `lab` funzionante, Dev Container aperto, un file con le informazioni
    di sistema raccolte.

## Passo 1 · L'account GitHub

**GitHub** è il servizio più usato al mondo per ospitare progetti software gestiti con Git (che studierai nella
Fase 2). Per te avrà tre ruoli: ospita il corso, ospita i tuoi appunti e diventerà la **vetrina del tuo portfolio**.

1. Registrati su [github.com/signup](https://github.com/signup). Scegli un nome utente **professionale**:
   comparirà nel tuo CV (per esempio `nome-cognome`).
2. Attiva subito l'**autenticazione a due fattori**: Settings → Password and authentication → *Enable two-factor
   authentication*. Usa un'app di autenticazione sul telefono e conserva i **codici di recupero** in un posto sicuro.
3. Completa il profilo con nome e foto.

## Passo 2 · La VM Linux

Segui la pagina [Macchina virtuale Linux](../ambiente/vm-linux.md) fino alla verifica dei dischi. Se il tuo
computer non lo permette, usa il playground Killercoda indicato nella stessa pagina.

## Passo 3 · Esplora la VM

Entra nella VM (`multipass shell lab`) e lancia questi comandi uno alla volta. Prima di passare al successivo,
leggi l'output e la spiegazione.

```bash
hostname          # (1)!
whoami            # (2)!
cat /etc/os-release   # (3)!
uname -r          # (4)!
nproc             # (5)!
free -h           # (6)!
df -h /           # (7)!
ip -brief address # (8)!
uptime            # (9)!
```

1.  Il nome del computer: `lab`.
2.  L'utente con cui sei collegato: `ubuntu`.
3.  Nome e versione della distribuzione. `/etc/os-release` è un file di testo standard presente in quasi tutte le distribuzioni.
4.  La versione del **kernel** Linux.
5.  Il numero di CPU (core) disponibili: dovrebbe essere 2, come richiesto a Multipass.
6.  La memoria: totale, usata, libera. `-h` rende i numeri leggibili (Mi, Gi: unità binarie, ricordi la lezione 0.1?).
7.  Lo spazio sul disco della radice `/`.
8.  Le interfacce di rete e i loro indirizzi IP. `lo` è l'interfaccia *loopback*, con cui la macchina parla con sé stessa (`127.0.0.1`).
9.  Da quanto tempo la macchina è accesa e il carico medio della CPU.

!!! esercizio "Esercizio 1 · La scheda della macchina"
    Usando solo i comandi visti e la redirezione `>>` (lezione 0.4), crea nella home della VM il file
    `scheda-lab.txt` che contenga, in ordine: hostname, distribuzione e versione, versione del kernel, numero di CPU,
    memoria totale, indirizzo IP dell'interfaccia principale.

??? success "Soluzione"
    ```bash
    hostname > scheda-lab.txt
    grep PRETTY_NAME /etc/os-release >> scheda-lab.txt
    uname -r >> scheda-lab.txt
    nproc >> scheda-lab.txt
    free -h | grep Mem >> scheda-lab.txt
    ip -brief address >> scheda-lab.txt
    cat scheda-lab.txt
    ```
    **Spiegazione.** Il primo comando usa `>` per creare il file da zero, gli altri `>>` per aggiungere. Sono
    comparsi due strumenti nuovi che studierai nella Fase 1: `grep`, che mostra solo le righe che contengono una
    parola, e il simbolo `|` (*pipe*), che passa l'output di un comando come input al successivo.
    Va bene anche una soluzione senza `grep`: l'importante è che il file contenga le informazioni richieste.

## Passo 4 · Il Dev Container

Segui la pagina [Dev Container](../ambiente/devcontainer.md), con la strada A (Codespaces) o la strada B (sul tuo
computer), fino al comando `docker run --rm hello-world`.

!!! esercizio "Esercizio 2 · Confronta i due ambienti"
    Nel terminale del Dev Container lancia `cat /etc/os-release`, `nproc` e `free -h`. Confronta i risultati con la
    VM. Che distribuzione usa il Dev Container? Perché il numero di CPU e la memoria possono essere diversi?

??? success "Soluzione"
    Il Dev Container usa **Ubuntu 24.04**, la stessa versione della VM, perché la sua immagine di base è
    `mcr.microsoft.com/devcontainers/base:ubuntu-24.04` (vedi il file `.devcontainer/Dockerfile`).
    CPU e memoria sono diverse perché un container, a differenza di una VM, **condivide il kernel** e le risorse
    del sistema che lo ospita: vede la macchina di Codespaces o quella su cui gira Docker. Capirai bene questa
    differenza nella lezione 3.1.

## Passo 5 · Spegni e riaccendi

Abituati fin da subito a non lasciare accese le risorse che non usi:

```bash
exit                  # esci dalla VM
multipass stop lab    # spegni la VM
multipass list        # lo stato deve essere Stopped
```

Ferma anche il Codespace se lo stai usando (menu **Codespaces** in basso a sinistra → *Stop Current Codespace*).
Nel cloud questa abitudine ti farà risparmiare denaro.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale file contiene nome e versione della distribuzione Linux?","options":["/etc/hostname","/etc/os-release","/proc/cpuinfo","/home/ubuntu/.bashrc"],"answer":[1],"explanation":"/etc/os-release è il file standard con le informazioni sulla distribuzione."},
 {"q":"Quale comando mostra la versione del kernel?","options":["uname -r","nproc","free -h","hostname"],"answer":[0],"explanation":"uname mostra informazioni sul kernel; -r stampa la release."},
 {"q":"Cosa indica l'interfaccia `lo`?","options":["L'interfaccia Wi-Fi","L'interfaccia di loopback, con cui la macchina comunica con sé stessa","Un errore di configurazione","L'interfaccia verso Internet"],"answer":[1],"explanation":"La loopback ha indirizzo 127.0.0.1 ed è sempre presente."},
 {"q":"Perché il Dev Container può mostrare più CPU della VM?","options":["Perché i container sono più veloci","Perché condivide il kernel e le risorse del sistema che lo ospita","Perché Ubuntu nei container è diverso","Per un errore"],"answer":[1],"explanation":"Un container non ha un hardware virtuale suo: vede le risorse dell'host, salvo limiti impostati."},
 {"q":"Quale comando spegne la VM di laboratorio?","options":["multipass delete lab","multipass stop lab","exit","shutdown lab"],"answer":[1],"explanation":"exit esce solo dalla shell; delete la cancella."}
]}
</script>

## Fonti

- [GitHub Docs: Configuring two-factor authentication](https://docs.github.com/en/authentication/securing-your-account-with-two-factor-authentication-2fa/configuring-two-factor-authentication)
- [os-release(5), man page](https://man7.org/linux/man-pages/man5/os-release.5.html)
- [Multipass: documentazione](https://canonical.com/multipass/docs)
- [iproute2: ip(8)](https://man7.org/linux/man-pages/man8/ip.8.html)
