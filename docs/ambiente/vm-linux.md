# La macchina virtuale Linux

!!! obiettivi "Obiettivi"
    - Capire perché alcune lezioni richiedono una macchina virtuale e non un container.
    - Creare la VM di laboratorio con Multipass, con un solo comando.
    - Conoscere le alternative nel browser se il tuo computer non può ospitare una VM.

## Perché una macchina virtuale

Una **macchina virtuale** (VM) è un computer simulato da un programma: ha un suo processore virtuale, la sua
memoria, i suoi dischi e il suo sistema operativo completo. Nella Fase 1 impari a gestire un server Linux vero:
servizi che partono all'avvio, dischi, firewall, utenti. Un container condivide parti del sistema con il computer
che lo ospita e non permette di esercitarsi su tutto questo; una VM sì.

Usiamo **Multipass**, uno strumento gratuito e open source di Canonical (l'azienda che sviluppa Ubuntu) che crea
VM Ubuntu con un solo comando su Windows, macOS e Linux.

## 1. Installa Multipass

=== "Windows"

    Requisiti: Windows 10 o 11 (Home, Pro o Enterprise), versione 1809 o successiva, con la funzionalità
    **Piattaforma macchina virtuale** attiva (Pannello di controllo → Programmi → Attiva o disattiva funzionalità
    di Windows).

    Scarica il file `.msi` dalla [pagina ufficiale](https://canonical.com/multipass/install) ed eseguilo come
    amministratore.

=== "macOS"

    Requisiti: macOS 14 Sonoma o successivo, Mac con processore Apple (M) o Intel, account amministratore.

    Scarica il file `.pkg` dalla [pagina ufficiale](https://canonical.com/multipass/install) ed eseguilo.

=== "Linux"

    Multipass si installa come pacchetto *snap*:

    ```bash
    sudo snap install multipass
    ```

    Se la tua distribuzione non ha `snapd`, installalo seguendo la [guida di Snapcraft](https://snapcraft.io/docs/installing-snapd).

Apri un terminale (su Windows: **PowerShell**; su macOS: **Terminale**) e verifica:

```bash
multipass version
```

## 2. Crea la VM di laboratorio

Il repository del corso contiene il file
[`lab-vm/cloud-init.yaml`](https://github.com/imredracoon/racoon-cloud-devops-academy/blob/main/lab-vm/cloud-init.yaml):
descrive come preparare la VM (pacchetti da installare e tre dischi virtuali da 1 GB per gli esercizi sui dischi).
**cloud-init** è lo standard con cui si configurano automaticamente i server al primo avvio, anche nel cloud.

1. Scarica il file: apri il link qui sopra, premi l'icona **Download raw file** e salvalo in una cartella
   che ricordi, per esempio `Documenti/academy`.
2. Nel terminale, spostati in quella cartella (`cd Documenti/academy`) e lancia:

```bash
multipass launch 24.04 --name lab --cpus 2 --memory 2G --disk 20G --cloud-init cloud-init.yaml
```

| Parte del comando | Significato |
|---|---|
| `multipass launch 24.04` | Crea una VM con Ubuntu 24.04 LTS. |
| `--name lab` | La chiama `lab`. |
| `--cpus 2 --memory 2G --disk 20G` | Le assegna 2 processori, 2 GB di memoria e 20 GB di disco. |
| `--cloud-init cloud-init.yaml` | Applica la configurazione del corso al primo avvio. |

La prima volta richiede qualche minuto, perché scarica l'immagine di Ubuntu.

## 3. Entra nella VM e verifica

```bash
multipass shell lab
```

Ora il terminale è **dentro** la VM: il prompt diventa `ubuntu@lab:~$`. Verifica che i dischi di laboratorio
esistano:

```console
$ losetup -a
/dev/loop3: []: (/srv/lab-disks/disk1.img)
/dev/loop4: []: (/srv/lab-disks/disk2.img)
/dev/loop5: []: (/srv/lab-disks/disk3.img)
```

I numeri dei `loop` possono essere diversi da te: conta che ci siano tre righe. Per uscire dalla VM scrivi `exit`.

## Comandi utili

| Comando | Cosa fa |
|---|---|
| `multipass list` | Elenca le VM e il loro stato. |
| `multipass stop lab` | Spegne la VM (libera memoria sul tuo computer). |
| `multipass start lab` | La riaccende. |
| `multipass snapshot lab` | Salva una "fotografia" della VM da spenta, per tornare indietro se rompi qualcosa. |
| `multipass restore lab.snapshot1` | Ripristina la fotografia. |
| `multipass delete --purge lab` | Cancella definitivamente la VM. |

!!! tip "Rompere senza paura"
    Prima di un laboratorio rischioso (dischi, firewall), spegni la VM e fai uno **snapshot**. Se qualcosa va storto,
    ripristinalo. Se proprio vuoi ricominciare da zero, cancella la VM e ricreala con il comando del passo 2.

## Alternativa nel browser: Killercoda

Se il tuo computer non può ospitare una VM, usa i **playground** gratuiti di [Killercoda](https://killercoda.com/playgrounds):
offrono una macchina Ubuntu completa nel browser per un'ora alla volta. Killercoda è la piattaforma degli stessi
autori di killer.sh, il simulatore ufficiale incluso negli esami della Linux Foundation.
Ricorda che ogni sessione si cancella alla fine: tieni i tuoi appunti e i tuoi script su GitHub.

Per gli esercizi sui dischi, nel playground crea i dischi di laboratorio a mano:

```bash
sudo mkdir -p /srv/lab-disks
for n in 1 2 3; do sudo truncate -s 1G /srv/lab-disks/disk$n.img; sudo losetup -f /srv/lab-disks/disk$n.img; done
```

## Fonti

- [Multipass: documentazione ufficiale](https://canonical.com/multipass/docs)
- [Multipass: installazione](https://canonical.com/multipass/docs/latest/how-to-guides/install-multipass/)
- [cloud-init: documentazione ufficiale](https://cloudinit.readthedocs.io/)
- [Killercoda: playground](https://killercoda.com/playgrounds)
