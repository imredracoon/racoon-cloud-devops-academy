# Fase 1 · Linux e reti

!!! obiettivi "Al termine di questa fase saprai"
    - muoverti con sicurezza nel filesystem Linux e manipolare file e testo con le pipe;
    - gestire utenti, gruppi, permessi, processi e servizi systemd;
    - installare software e gestire dischi, filesystem e LVM;
    - lavorare con sicurezza sia su Debian/Ubuntu sia su RHEL (Rocky/AlmaLinux), conoscendone le differenze;
    - calcolare subnet, spiegare DNS, DHCP, routing e NAT e diagnosticare problemi di rete;
    - collegarti ai server con SSH e chiavi, proteggerli con un firewall nftables;
    - pubblicare un sito con Nginx e HTTPS e risolvere i guasti con metodo;
    - gestire SELinux e AppArmor, la sequenza di avvio, l'ora di sistema e il recupero dell'accesso.

| | |
|---|---|
| **Durata indicativa** | 2,5-3 mesi, circa 90 ore (più circa 20 ore per ripetere i laboratori sulla seconda famiglia) |
| **Prerequisiti** | [Fase 0](../fase-0/index.md) |
| **Ambiente** | Percorso Debian/Ubuntu: VM `lab` e `client` (Ubuntu 24.04). Percorso RHEL: VM `rhel-lab` e `rhel-client` (AlmaLinux 10). Tutte con KVM/libvirt ([macchine virtuali](../ambiente/vm-linux.md)) |
| **Progetto** | [Server web messo in sicurezza](progetto.md) |
| **Certificazioni collegate** | [LFCS](../appendici/lfcs/index.md), [RHCSA](../appendici/rhcsa/index.md), [Canonical Academy](../appendici/canonical/index.md) (tutte facoltative) |

## Perché Linux e reti prima di tutto

Ogni tecnologia che studierai dopo poggia su questi due pilastri. Un container è un insieme di processi Linux
isolati. Un'istanza nel cloud è una VM Linux. Un cluster Kubernetes è una rete di macchine Linux. Quando qualcosa
si rompe, la causa si trova quasi sempre in uno di questi strati: un permesso, un servizio fermo, un disco pieno,
un DNS sbagliato, una porta chiusa. Chi padroneggia Linux e reti risolve problemi che gli altri non capiscono.

## Due percorsi: Debian/Ubuntu e RHEL

In questa fase ogni argomento che cambia tra le due grandi famiglie di distribuzioni è presentato in **schede
affiancate**, come questa:

=== "Debian/Ubuntu"

    Ubuntu 24.04 LTS, gestore di pacchetti `apt`, firewall `ufw`, rete con Netplan, sicurezza con AppArmor.

=== "RHEL (Rocky/AlmaLinux)"

    AlmaLinux 10, compatibile con Red Hat Enterprise Linux 10: gestore di pacchetti `dnf`, firewall `firewalld`,
    rete con NetworkManager, sicurezza con SELinux.

La scelta della scheda resta memorizzata in tutte le pagine. I laboratori indicano come svolgerli su entrambe le
famiglie. Le istruzioni per creare le VM sono nella pagina [Macchine virtuali Linux](../ambiente/vm-linux.md).

## Preparare la seconda VM

Alcuni esercizi di rete (SSH, firewall, Nginx) sono più realistici con due macchine. Quando la lezione lo
richiede, crea una seconda VM leggera dalla cartella `~/academy/lab-vm`:

=== "Debian/Ubuntu"

    ```bash
    ./academy-vm crea client
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    ./academy-vm crea rhel-client
    ```

Le VM stanno sulla stessa rete privata e si raggiungono per nome; `./academy-vm elenco` mostra gli indirizzi
([tutte le VM del corso](../ambiente/vm-linux.md#tutte-le-vm-del-corso)).

## Le lezioni

1. [Il filesystem](01-filesystem.md)
2. [Lavorare con file e testo](02-file-testo.md)
3. [Utenti, gruppi e permessi](03-permessi.md)
4. [Processi e servizi con systemd](04-processi-systemd.md)
5. [Pacchetti e software](05-pacchetti.md)
6. [Dischi, filesystem e LVM](06-storage.md)
7. [TCP/IP e indirizzamento](07-tcpip.md)
8. [DNS, DHCP, routing e NAT](08-dns-routing.md)
9. [Diagnostica di rete](09-diagnostica.md)
10. [SSH e chiavi](10-ssh.md)
11. [Firewall con nftables](11-firewall.md)
12. [HTTP, TLS e Nginx](12-http-nginx.md)
13. [Metodo di troubleshooting](13-troubleshooting.md)
14. [SELinux e AppArmor](14-selinux-apparmor.md)
15. [Avvio, tempo e prestazioni](15-avvio-tempo.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

!!! tip "Fai uno snapshot"
    Prima di ogni lezione dalla [1.3](03-permessi.md) in poi, crea uno snapshot: `virsh shutdown lab` e poi
    `virsh snapshot-create-as lab prima-lezione-X` (oppure lo stesso con `rhel-lab`). Se rompi qualcosa, torni indietro in un minuto.
