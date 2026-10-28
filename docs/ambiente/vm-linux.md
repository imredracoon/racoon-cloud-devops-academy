# Le macchine virtuali Linux

!!! obiettivi "Obiettivi"
    - Capire perché alcune lezioni richiedono una macchina virtuale e non un container.
    - Creare la VM di laboratorio Ubuntu con Multipass e quella RHEL (AlmaLinux) con Vagrant.
    - Conoscere le alternative nel browser se il tuo computer non può ospitare una VM.

## Perché una macchina virtuale

Una **macchina virtuale** (VM) è un computer simulato da un programma: ha un suo processore virtuale, la sua
memoria, i suoi dischi e il suo sistema operativo completo. Nella [Fase 1](../fase-1/index.md) impari a gestire un server Linux vero:
servizi che partono all'avvio, dischi, firewall, utenti. Un container condivide parti del sistema con il computer
che lo ospita e non permette di esercitarsi su tutto questo; una VM sì.

## Due famiglie, due laboratori

Nel mondo dei server Linux dominano due famiglie ([lezione 0.2](../fase-0/02-sistemi-operativi.md)): **Debian/Ubuntu** e **Red Hat Enterprise Linux
(RHEL)** con le distribuzioni compatibili **Rocky Linux** e **AlmaLinux**. Le aziende usano l'una o l'altra, spesso
entrambe, e gli annunci di lavoro le chiedono tutte e due. Per questo nella [Fase 1](../fase-1/index.md) ogni argomento che cambia tra le
due famiglie è presentato in schede affiancate, e i laboratori si possono svolgere su entrambe.

| Percorso | Distribuzione | Strumento | Certificazioni collegate |
|---|---|---|---|
| **Debian/Ubuntu** | Ubuntu 24.04 LTS | Multipass | LFCS, Canonical Academy |
| **RHEL** | AlmaLinux 10 (compatibile con RHEL 10) | Vagrant | RHCSA, RHCE, LFCS |

!!! tip "Quale scegliere?"
    Se è la prima volta, segui **Ubuntu** come percorso principale e ripeti i laboratori su **AlmaLinux** quando
    la scheda RHEL mostra differenze. Se punti alle certificazioni Red Hat, fai il contrario. Conoscere entrambe è un
    vantaggio concreto nei colloqui.

## Percorso Debian/Ubuntu: Multipass

Usiamo **Multipass**, uno strumento gratuito e open source di Canonical (l'azienda che sviluppa Ubuntu) che crea
VM Ubuntu con un solo comando su Windows, macOS e Linux.

### 1. Installa Multipass

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

### 2. Crea la VM di laboratorio

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

### 3. Entra nella VM e verifica

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

### Comandi utili

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

## Percorso RHEL: AlmaLinux con Vagrant

Multipass crea solo VM Ubuntu. Per la famiglia Red Hat usiamo **Vagrant** di HashiCorp, uno strumento gratuito che
crea VM a partire da un file di testo, il `Vagrantfile`, e le **box** ufficiali di **AlmaLinux 10**, compatibile
binariamente con RHEL 10: quello che impari vale identico su Red Hat Enterprise Linux.

Il file [`lab-vm/rhel/Vagrantfile`](https://github.com/imredracoon/racoon-cloud-devops-academy/blob/main/lab-vm/rhel/Vagrantfile)
del corso crea due VM su una rete privata con indirizzi fissi:

| VM | Indirizzo | Ruolo |
|---|---|---|
| `rhel-lab` | 192.168.56.10 | il server, con tre dischi di laboratorio come la VM Ubuntu |
| `rhel-client` | 192.168.56.11 | il client per gli esercizi di rete |

### 1. Installa Vagrant e un hypervisor

=== "Windows e Linux"

    1. Installa [VirtualBox](https://www.virtualbox.org/wiki/Downloads) (gratuito).
    2. Installa [Vagrant](https://developer.hashicorp.com/vagrant/install).

    Su Windows VirtualBox convive con Hyper-V (usato da Docker Desktop e Multipass), ma può essere più lento.
    Su Linux puoi usare in alternativa **libvirt/KVM** con il plugin `vagrant-libvirt`.

=== "macOS (Intel)"

    Installa [VirtualBox](https://www.virtualbox.org/wiki/Downloads) e [Vagrant](https://developer.hashicorp.com/vagrant/install).

=== "macOS (Apple M1-M4)"

    Le box AlmaLinux per processori ARM sono disponibili per **VMware** e **Parallels**. Installa
    [VMware Fusion](https://www.vmware.com/products/desktop-hypervisor/workstation-and-fusion) (gratuito per uso
    personale), [Vagrant](https://developer.hashicorp.com/vagrant/install), la
    [Vagrant VMware Utility](https://developer.hashicorp.com/vagrant/install/vmware) e il plugin:

    ```bash
    vagrant plugin install vagrant-vmware-desktop
    ```

    Poi aggiungi `--provider=vmware_desktop` al comando `vagrant up`.

### 2. Crea le VM

Scarica la cartella `lab-vm/rhel` del corso (oppure clona il repository), aprila nel terminale e lancia:

```bash
vagrant up
```

La prima volta scarica la box (circa 1 GB) e prepara le macchine. Poi:

```bash
vagrant ssh rhel-lab          # entra nel server
cat /etc/os-release           # AlmaLinux 10
getenforce                    # Enforcing: SELinux è attivo (lezione 1.14)
losetup -a                    # i tre dischi di laboratorio
```

| Comando | Cosa fa |
|---|---|
| `vagrant status` | Stato delle VM |
| `vagrant halt` | Spegne le VM |
| `vagrant snapshot save rhel-lab pulita` | Salva uno snapshot |
| `vagrant snapshot restore rhel-lab pulita` | Lo ripristina |
| `vagrant destroy -f` | Cancella le VM |

!!! note "Red Hat Enterprise Linux vero, gratis"
    Red Hat offre la **Red Hat Developer Subscription for Individuals**: con un account gratuito su
    [developers.redhat.com](https://developers.redhat.com/) puoi installare RHEL originale per uso personale e
    di sviluppo. È utile soprattutto nelle ultime settimane prima dell'esame RHCSA.

## Alternativa nel browser: Killercoda

Se il tuo computer non può ospitare una VM, usa i **playground** gratuiti di [Killercoda](https://killercoda.com/playgrounds):
offrono una macchina Ubuntu completa nel browser per un'ora alla volta. Killercoda è la piattaforma degli stessi
autori di killer.sh, il simulatore ufficiale incluso negli esami della Linux Foundation.
Ricorda che ogni sessione si cancella alla fine: tieni i tuoi appunti e i tuoi script su GitHub.

Per gli esercizi sui dischi, nel playground crea i dischi di laboratorio a mano:

```bash
sudo mkdir -p /srv/lab-disks
for n in 1 2 3; do sudo truncate -s 1G /srv/lab-disks/disk$n.img; sudo losetup -fP /srv/lab-disks/disk$n.img; done
```

## Fonti

- [Multipass: documentazione ufficiale](https://canonical.com/multipass/docs)
- [Multipass: installazione](https://canonical.com/multipass/docs/latest/how-to-guides/install-multipass/)
- [cloud-init: documentazione ufficiale](https://cloudinit.readthedocs.io/)
- [Killercoda: playground](https://killercoda.com/playgrounds)
- [Vagrant: documentazione](https://developer.hashicorp.com/vagrant/docs)
- [AlmaLinux: box Vagrant ufficiali](https://wiki.almalinux.org/installation/vagrant-boxes.html)
- [Red Hat Developer Subscription for Individuals](https://developers.redhat.com/articles/faqs-no-cost-red-hat-enterprise-linux)
