# Le macchine virtuali Linux

!!! obiettivi "Obiettivi"
    - Capire perché alcune lezioni richiedono una macchina virtuale e non un container.
    - Installare KVM/QEMU e libvirt, la virtualizzazione integrata in Linux.
    - Creare le VM di laboratorio Ubuntu e AlmaLinux con un solo script e gestirle con `virsh`.
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

| Percorso | Distribuzione | VM | Certificazioni collegate |
|---|---|---|---|
| **Debian/Ubuntu** | Ubuntu 24.04 LTS | `lab`, `client`, `web` | LFCS, Canonical Academy |
| **RHEL** | AlmaLinux 10 (compatibile binariamente con RHEL 10) | `rhel-lab`, `rhel-client`, `rhel-web` | RHCSA, RHCE, LFCS |

Tutte le VM, delle due famiglie, si creano con **lo stesso strumento** e vivono sulla stessa rete privata.

!!! tip "Quale scegliere?"
    Se è la prima volta, segui **Ubuntu** come percorso principale e ripeti i laboratori su **AlmaLinux** quando
    la scheda RHEL mostra differenze. Se punti alle certificazioni Red Hat, fai il contrario. Conoscere entrambe è un
    vantaggio concreto nei colloqui.

## Lo strumento: KVM, QEMU e libvirt

Linux ha la virtualizzazione **dentro il kernel**: si chiama **KVM** (*Kernel-based Virtual Machine*) ed è la stessa
tecnologia su cui girano molti cloud pubblici. Le VM sono processi normali del tuo computer, quindi sono leggere e
veloci. Tre componenti lavorano insieme:

| Componente | Ruolo |
|---|---|
| **KVM** | il modulo del kernel che usa le istruzioni di virtualizzazione del processore |
| **QEMU** | il programma che simula il resto della macchina: dischi, schede di rete, console |
| **libvirt** | il gestore che crea, avvia e ferma le VM e le loro reti; si comanda con `virsh` o con la grafica di **virt-manager** |

Per non farti creare le VM a mano, il corso include lo script
[`lab-vm/academy-vm`](https://github.com/imredracoon/racoon-cloud-devops-academy/blob/main/lab-vm/academy-vm):
scarica le **immagini cloud ufficiali** di Ubuntu e AlmaLinux (le stesse che si usano su AWS e Google Cloud) e le
configura al primo avvio con **cloud-init**, lo standard con cui si preparano automaticamente i server nel cloud.

### 1. Controlla che la virtualizzazione sia attiva

```bash
grep -cE 'vmx|svm' /proc/cpuinfo
```

Deve stampare un numero maggiore di zero. Se stampa `0`, attiva la virtualizzazione (*Intel VT-x* o *AMD-V*, a volte
chiamata *SVM*) nel BIOS/UEFI del computer: il manuale o il sito del produttore indicano dove si trova l'opzione.

### 2. Installa KVM, libvirt e gli strumenti

=== "Ubuntu"

    ```bash
    sudo apt install -y qemu-system-x86 libvirt-daemon-system virtinst virt-manager git
    ```

=== "Fedora"

    ```bash
    sudo dnf install -y @virtualization git
    sudo systemctl enable --now libvirtd
    ```

Poi, su entrambe:

```bash
sudo usermod -aG libvirt $USER                                   # (1)!
echo 'export LIBVIRT_DEFAULT_URI=qemu:///system' >> ~/.bashrc     # (2)!
```

1.  Il gruppo `libvirt` ti permette di gestire le VM senza `sudo`.
2.  libvirt ha due "modalità": quella **di sistema**, con reti vere tra le VM, e quella **di sessione**, ridotta.
    Usiamo quella di sistema, e questa riga dice a `virsh` di usarla sempre.

**Esci e rientra nella sessione** (o riavvia), poi verifica:

```bash
virsh list --all          # nessun errore, e per ora nessuna VM
virt-host-validate qemu   # le righe sulla virtualizzazione devono dire PASS
```

`virt-host-validate` può segnalare qualche `WARN` su funzioni facoltative (per esempio IOMMU): per il corso non
servono.

### 3. Scarica i file del corso

```bash
git clone https://github.com/imredracoon/racoon-cloud-devops-academy.git ~/academy
cd ~/academy/lab-vm
```

`git clone` scarica una copia del repository: lo studierai nella [Fase 2](../fase-2/index.md). Per aggiornare i file in futuro basta
`git pull` dentro `~/academy`.

### 4. Crea la rete e la prima VM

```bash
./academy-vm rete         # una volta sola
./academy-vm crea lab     # la VM Ubuntu
```

| Cosa fa lo script | Dettagli |
|---|---|
| crea la rete privata `academy` | `192.168.56.0/24`, con accesso a Internet e un indirizzo fisso per ogni VM |
| scarica l'immagine di base | solo la prima volta, in `/var/lib/libvirt/images/academy/`: per questo ti chiede la password di `sudo` |
| crea il disco della VM | registra solo le differenze rispetto all'immagine di base, quindi occupa poco |
| prepara la configurazione | nome, la tua chiave SSH (se non l'hai, la crea), pacchetti e tre dischi di laboratorio da 1 GB |
| aggiunge la VM a `~/.ssh/config` | così ti basta scrivere `ssh lab` |

La prima volta richiede qualche minuto. Quando compare `Pronta: ssh lab`, entra:

```bash
ssh lab
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

### 5. Le VM della famiglia Red Hat

Con lo stesso comando crei le VM AlmaLinux 10, compatibile binariamente con RHEL 10: quello che impari vale identico
su Red Hat Enterprise Linux.

```bash
./academy-vm crea rhel-lab rhel-client
ssh rhel-lab
cat /etc/os-release           # AlmaLinux 10
getenforce                    # Enforcing: SELinux è attivo (lezione 1.14)
losetup -a                    # i tre dischi di laboratorio
```

### Tutte le VM del corso

Non serve crearle tutte subito: il corso ti dice quando ne serve una nuova.

| VM | Sistema | Indirizzo | Utente | Risorse | Dove si usa |
|---|---|---|---|---|---|
| `lab` | Ubuntu 24.04 | 192.168.56.21 | `ubuntu` | 2 vCPU, 2 GB, 20 GB | Fase 0, Fase 1, LFCS |
| `client` | Ubuntu 24.04 | 192.168.56.22 | `ubuntu` | 1 vCPU, 1 GB, 10 GB | Fase 1, esercizi di rete |
| `web` | Ubuntu 24.04 | 192.168.56.23 | `ubuntu` | 1 vCPU, 1 GB, 10 GB | [Progetto 1](../fase-1/progetto.md) |
| `rhel-lab` | AlmaLinux 10 | 192.168.56.10 | `almalinux` | 2 vCPU, 2 GB, 20 GB | Fase 1, RHCSA, RHCE |
| `rhel-client` | AlmaLinux 10 | 192.168.56.11 | `almalinux` | 1 vCPU, 1,5 GB, 10 GB | esercizi di rete, RHCE |
| `rhel-web` | AlmaLinux 10 | 192.168.56.12 | `almalinux` | 1 vCPU, 1,5 GB, 10 GB | [Progetto 1](../fase-1/progetto.md), percorso RHEL |
| `k8s-cp` | Ubuntu 24.04 | 192.168.56.31 | `ubuntu` | 2 vCPU, 2 GB, 20 GB | [lezione 8.12](../fase-8/12-kubeadm.md): control plane |
| `k8s-w1`, `k8s-w2` | Ubuntu 24.04 | 192.168.56.32-33 | `ubuntu` | 2 vCPU, 2 GB, 20 GB | [lezione 8.12](../fase-8/12-kubeadm.md): worker |

!!! note "Hai creato la rete prima dell'arrivo delle VM `k8s-*`?"
    Se `./academy-vm rete` l'hai eseguito prima che lo script conoscesse le VM `k8s-cp`, `k8s-w1` e `k8s-w2`,
    ricrea la rete una volta, con le VM spente: `virsh net-destroy academy && virsh net-undefine academy &&
    ./academy-vm rete`. Dopo `git pull` lo script e la rete sono aggiornati.

Le VM si vedono tra loro **per nome**: da `lab` puoi scrivere `ping client`, da `rhel-client` `ssh rhel-lab`. Il
gateway della rete, cioè il tuo computer, è `192.168.56.1`.

### Livelli 2 e 3: le VM sul server

Nei [livelli 2 e 3](livelli.md) le VM non girano sulla workstation ma su `srv1` (livello 2) o su `nodo1` (livello 3).
Tutti i comandi di questa pagina si eseguono **sul server**, dopo esserci entrati con `ssh srv1` o `ssh nodo1`; lì
servono solo i pacchetti senza grafica (`virt-manager` non serve). La [lezione 0.8](../fase-0/08-homelab.md) ti guida passo
passo. Due comodità per lavorare dalla workstation:

- **SSH diretto alle VM**: sul server esegui `./academy-vm ssh-config srv1` (o `nodo1`) e copia l'output in fondo al file
  `~/.ssh/config` della **workstation**. Da quel momento `ssh lab` funziona anche dalla workstation: la connessione
  passa dal server (*ProxyJump*). Lo script autorizza nelle VM anche le chiavi che hai già autorizzato sul server.
- **virt-manager remoto**: sulla workstation apri virt-manager, *File → Aggiungi connessione*, scegli *QEMU/KVM*,
  spunta *Connetti a un host remoto via SSH* e scrivi il tuo utente e `srv1`. Vedi e gestisci le VM del server come
  se fossero sulla workstation.

### Comandi utili

| Comando | Cosa fa |
|---|---|
| `./academy-vm elenco` | le VM del corso, con indirizzo e stato |
| `virsh list --all` | tutte le VM e il loro stato |
| `virsh start lab` | accende la VM |
| `virsh shutdown lab` | la spegne in modo ordinato (libera memoria sul tuo computer) |
| `virsh destroy lab` | la spegne di colpo, come staccare la spina: non cancella nulla |
| `virsh console lab` | la **console seriale**: entri anche se la rete della VM non funziona; esci con ++ctrl+bracket-right++ |
| `virsh snapshot-create-as lab pulita` | salva una "fotografia" della VM con il nome `pulita` |
| `virsh snapshot-list lab` | elenca le fotografie |
| `virsh snapshot-revert lab pulita` | torna alla fotografia |
| `./academy-vm elimina lab` | cancella definitivamente la VM e i suoi dischi |

!!! tip "Rompere senza paura"
    Prima di un laboratorio rischioso (dischi, firewall, rete), spegni la VM e fai uno **snapshot**. Se qualcosa va
    storto, ripristinalo. Se proprio vuoi ricominciare da zero, cancella la VM e ricreala con `./academy-vm crea`.

!!! note "La console di emergenza"
    `virsh console` chiede utente e password: l'utente è quello della tabella (`ubuntu` o `almalinux`), la password
    è `academy`. Vale **solo** per la console: via SSH si entra unicamente con la chiave. È l'equivalente della
    console che i provider cloud offrono quando un server non risponde più in rete.

### Se preferisci la grafica

Apri **Gestore macchine virtuali** (*virt-manager*) dal menu delle applicazioni: vedrai le VM del corso, potrai
avviarle, spegnerle, guardarne la console e i grafici di CPU e memoria. È la stessa libvirt di `virsh`, solo con le
finestre.

!!! info "E GNOME Boxes?"
    GNOME Boxes usa lo stesso motore KVM, ma lavora nella modalità *di sessione* di libvirt: non vede le VM del corso
    e non gestisce reti private tra VM. Va benissimo per provare una distribuzione al volo da un file ISO, non per
    i laboratori.

!!! note "Red Hat Enterprise Linux vero, gratis"
    Red Hat offre la **Red Hat Developer Subscription for Individuals**: con un account gratuito su
    [developers.redhat.com](https://developers.redhat.com/) puoi installare RHEL originale per uso personale e
    di sviluppo. È utile soprattutto nelle ultime settimane prima dell'esame RHCSA.

### Problemi frequenti

??? question "`error: failed to connect to the hypervisor` o `Permission denied`"
    Il tuo utente non è ancora nel gruppo `libvirt`, oppure non sei uscito e rientrato dopo `usermod`. Controlla con
    `groups`. Verifica anche che `echo $LIBVIRT_DEFAULT_URI` stampi `qemu:///system`.

??? question "Lo script resta fermo su «Attendo che … completi la configurazione»"
    La prima volta la VM scarica gli aggiornamenti e i pacchetti: può servire qualche minuto. Se dopo 10 minuti è
    ancora fermo, apri un altro terminale ed entra con `virsh console NOME` per vedere cosa succede: cloud-init
    scrive i suoi messaggi in `/var/log/cloud-init-output.log`.

??? question "AlmaLinux 10 e i processori meno recenti"
    AlmaLinux 10, come RHEL 10, richiede un processore con le istruzioni **x86-64-v3** (in pratica: AVX2, presenti
    sui processori Intel e AMD dal 2013-2015 in poi). Lo script se ne accorge e, se il tuo processore non le ha,
    scarica la variante **x86-64-v2** di AlmaLinux, pensata proprio per questi casi.

??? question "Ho ricreato una VM e `ssh` segnala un cambio di chiave"
    Lo script rimuove la vecchia chiave dal file `~/.ssh/known_hosts_academy` quando ricrei una VM. Se il messaggio
    compare comunque, eseguilo a mano: `ssh-keygen -R 192.168.56.21 -f ~/.ssh/known_hosts_academy` (con l'indirizzo
    della VM).

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

- [KVM: documentazione del kernel Linux](https://docs.kernel.org/virt/kvm/index.html)
- [libvirt: documentazione](https://libvirt.org/docs.html) e [virsh(1)](https://libvirt.org/manpages/virsh.html)
- [virt-install(1)](https://github.com/virt-manager/virt-manager/blob/main/man/virt-install.rst)
- [Ubuntu Server: libvirt](https://documentation.ubuntu.com/server/how-to/virtualisation/libvirt/)
- [Fedora Docs: Getting started with virtualization](https://docs.fedoraproject.org/en-US/quick-docs/virtualization-getting-started/)
- [Ubuntu: immagini cloud ufficiali](https://cloud-images.ubuntu.com/)
- [AlmaLinux: immagini cloud generiche](https://wiki.almalinux.org/cloud/Generic-cloud.html)
- [cloud-init: documentazione ufficiale](https://cloudinit.readthedocs.io/)
- [Killercoda: playground](https://killercoda.com/playgrounds)
- [Red Hat Developer Subscription for Individuals](https://developers.redhat.com/articles/faqs-no-cost-red-hat-enterprise-linux)
