# Argomenti specifici · Canonical Academy

Gli argomenti pubblicati da Canonical per ciascuno dei quattro esami, con le lezioni del corso.

## Using Linux Terminal

| Argomento | Lezioni |
|---|---|
| Navigare file e filesystem: comandi per file e directory, ricerca, confronto e modifica con espressioni regolari, pipe e redirezioni | [0.4](../../fase-0/04-terminale.md), [1.1](../../fase-1/01-filesystem.md), [1.2](../../fase-1/02-file-testo.md) |
| Risorse di sistema: trovare e ruotare i log, ricostruire partizioni e filesystem (`fdisk`, `fsck`, `parted`), formato e posizione di crontab, leggere i log per risolvere i problemi | [1.4](../../fase-1/04-processi-systemd.md), [1.6](../../fase-1/06-storage.md), [2.5](../../fase-2/05-pianificazione.md) |
| Accesso sicuro: chiavi SSH, impostazioni di sicurezza di sistema e utente, complessità e scadenza delle password, policy di sudo, utenti e gruppi, proprietà e permessi | [1.3](../../fase-1/03-permessi.md), [1.10](../../fase-1/10-ssh.md) |

## Using Ubuntu Desktop

| Argomento | Lezioni |
|---|---|
| Rete: configurare una connessione; livello 2 (indirizzi MAC, ARP, broadcast e multicast) | [0.3](../../fase-0/03-reti-concetti.md), [1.7](../../fase-1/07-tcpip.md) |
| Comunità open source: cicli di rilascio LTS e versioni, governance della comunità Ubuntu, interpretazione delle licenze | [0.2](../../fase-0/02-sistemi-operativi.md) e [approfondimento](#comunita-e-licenze) |
| Installazione: riparare GRUB2 e initramfs, aggiornare tra versioni LTS, chiavette avviabili, driver proprietari per GPU e Wi-Fi, volumi LVM, filesystem e snapshot | [1.6](../../fase-1/06-storage.md), [1.15](../../fase-1/15-avvio-tempo.md) |
| Applicazioni: elencare e aggiornare pacchetti, deb contro snap, sorgenti diverse, descrizioni dei pacchetti, `apt-get`, aggiornamenti automatici, Synaptic, Aptitude | [1.5](../../fase-1/05-pacchetti.md) |
| Sicurezza del desktop: componenti e usi della PKI, collegarsi a un dominio Active Directory, permessi con ACL | [1.12](../../fase-1/12-http-nginx.md), [1.3](../../fase-1/03-permessi.md) |

## Using Ubuntu Server

| Argomento | Lezioni |
|---|---|
| Servizi: systemd, backup remoti automatici con timer o cron, file server NFS o Samba, applicazione web con Apache o Nginx e un database | [1.4](../../fase-1/04-processi-systemd.md), [1.6](../../fase-1/06-storage.md), [1.12](../../fase-1/12-http-nginx.md), [2.5](../../fase-2/05-pianificazione.md), [3.5](../../fase-3/05-compose.md) |
| Distribuzione: avvio PXE, autoinstall, cloud-init, verificare i requisiti, avvio in modalità utente singolo | [1.15](../../fase-1/15-avvio-tempo.md), [6.10](../../fase-6/10-ansible.md) |
| Accesso sicuro: SSH, utenti normali e di sistema, firewall con iptables o ufw | [1.10](../../fase-1/10-ssh.md), [1.11](../../fase-1/11-firewall.md) |
| Processi: identificare e gestire processi e servizi, elencare unit attive, disabilitate e inattive | [1.4](../../fase-1/04-processi-systemd.md) |

## Using DevOps Principles

| Argomento | Lezioni |
|---|---|
| Automazione e aggiornamenti: script Bash o Python versionati con git, cloud-init e Ansible, aggiornamenti con Landscape, Ubuntu Pro, Livepatch ed ESM | [2.1](../../fase-2/01-git-base.md)-[2.7](../../fase-2/07-python-automazione.md), [6.10](../../fase-6/10-ansible.md) e [approfondimento](#ubuntu-pro-livepatch-esm-e-landscape) |
| Tecnologie di distribuzione: applicazioni web e database in container con Docker o LXD, virtualizzazione con virsh, virt-manager o LXD | [3.2](../../fase-3/02-immagini-container.md)-[3.5](../../fase-3/05-compose.md) e [approfondimento](#lxd-e-libvirt) |
| Monitoraggio e diagnosi: effetto della configurazione su calcolo e rete, `tcpdump`, `dig`, `ss`, iproute2, contesa delle risorse, problemi di connettività e di carico | [1.9](../../fase-1/09-diagnostica.md), [1.13](../../fase-1/13-troubleshooting.md), [1.15](../../fase-1/15-avvio-tempo.md) |

## Approfondimenti

### Comunità e licenze

Ubuntu pubblica una versione **LTS** ogni due anni ad aprile (24.04, 26.04...) con 5 anni di manutenzione standard,
estendibili con Ubuntu Pro; le versioni intermedie escono ogni sei mesi e sono supportate per 9 mesi. La comunità
è guidata dal Community Council e dal Technical Board. Le licenze **copyleft** (GPL) chiedono di distribuire le
opere derivate con la stessa licenza; quelle **permissive** (MIT, BSD, Apache 2.0) permettono il riuso anche in
software proprietario, con attribuzione. Fonti: [ciclo di rilascio](https://ubuntu.com/about/release-cycle),
[governance](https://ubuntu.com/community/governance), [licenze approvate OSI](https://opensource.org/licenses).

### Ubuntu Pro, Livepatch, ESM e Landscape

**Ubuntu Pro** (gratuito per uso personale su un numero limitato di macchine) si attiva con `sudo pro attach` e
include **ESM** (aggiornamenti di sicurezza estesi per i pacchetti di `main` e `universe`) e **Livepatch** (patch
critiche del kernel senza riavvio). `pro status` mostra i servizi attivi. **Landscape** gestisce, monitora e
aggiorna molte macchine Ubuntu da un'unica console. Fonti: [Ubuntu Pro](https://documentation.ubuntu.com/pro/),
[Landscape](https://documentation.ubuntu.com/landscape/).

### LXD e libvirt

**LXD** gestisce *container di sistema* (un Ubuntu completo con systemd, che condivide il kernel dell'host) e
macchine virtuali: `lxc launch ubuntu:24.04 web1`, `lxc exec web1 -- bash`, `lxc launch ubuntu:24.04 vm1 --vm`.
**libvirt** gestisce VM KVM con `virsh` (`list --all`, `start`, `shutdown`, `autostart`) o con l'interfaccia
grafica virt-manager. Fonti: [LXD](https://documentation.ubuntu.com/lxd/latest/),
[Ubuntu Server · libvirt](https://documentation.ubuntu.com/server/how-to/virtualisation/libvirt/).

!!! tip "La parte pratica"
    La simulazione di questa appendice copre la parte a domande. Per le prove pratiche usa i laboratori della Fase
    1 sul percorso Debian/Ubuntu e i task della [simulazione LFCS](../lfcs/simulazione.md), quasi tutti su Ubuntu.
