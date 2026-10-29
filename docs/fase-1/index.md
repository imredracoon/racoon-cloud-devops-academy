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

## Nei tre livelli

Gli argomenti sono gli stessi per tutti; cambia **dove** li provi. Scegli la scheda del tuo
[livello](../ambiente/livelli.md): il sito la ricorda in tutte le pagine.

=== "Livello 1 · Un PC"

    Tutti i laboratori si svolgono nelle VM sulla workstation, come descritto sopra.

=== "Livello 2 · Workstation + server"

    Le VM `lab`, `client` e le altre girano su `srv1` ([lezione 0.8](../fase-0/08-homelab.md)): i comandi sono identici,
    entri con `ssh lab` passando dal server. Fai gli **esercizi nelle VM**, dove hai gli snapshot, e ripeti sul
    server vero le parti indicate nelle sezioni *Nei tre livelli* delle lezioni:

    | Lezione | Sul server `srv1` |
    |---|---|
    | [1.6](06-storage.md#nei-tre-livelli) | una chiavetta USB come disco vero, NFS verso la workstation |
    | [1.7](07-tcpip.md#nei-tre-livelli) | la seconda rete (facoltativa): indirizzi statici su un cavo diretto |
    | [1.8](08-dns-routing.md#nei-tre-livelli) | `srv1` è già un router: la rete NAT delle VM |
    | [1.9](09-diagnostica.md#nei-tre-livelli) | vedere il NAT con `tcpdump` su due interfacce |
    | [1.10](10-ssh.md#nei-tre-livelli) | SSH solo con chiave sul server |
    | [1.11](11-firewall.md#nei-tre-livelli) | il firewall del server, senza rompere le VM |
    | [1.12](12-http-nginx.md#nei-tre-livelli) | Nginx sul server, visto dal browser della workstation |
    | [Progetto](progetto.md#nei-tre-livelli) | il server web sicuro è `srv1` |

=== "Livello 3 · Mini data center"

    I nodi fisici prendono il posto delle VM principali: **`nodo2`** è il server del percorso Debian/Ubuntu (al
    posto di `lab`), **`nodo3`** quello del percorso RHEL (al posto di `rhel-lab`). Le VM su `nodo1` restano a
    disposizione per gli esercizi distruttivi, con gli snapshot. In questa fase costruisci la rete del data center:

    | Lezione | Nel mini data center |
    |---|---|
    | [1.6](06-storage.md#nei-tre-livelli) | chiavetta USB come disco vero su `nodo2` e `nodo3`; `nodo1` server NFS |
    | [1.7](07-tcpip.md#nei-tre-livelli) | il piano degli indirizzi `10.10.10.0/24` e la scheda `lab0` di `nodo1` |
    | [1.8](08-dns-routing.md#nei-tre-livelli) | **`nodo1` diventa il gateway**: inoltro, NAT, DHCP e DNS per la rete del laboratorio |
    | [1.9](09-diagnostica.md#nei-tre-livelli) | il NAT visto con `tcpdump` sulle due schede di `nodo1` |
    | [1.10](10-ssh.md#nei-tre-livelli) | SSH solo con chiave su tutti i nodi |
    | [1.11](11-firewall.md#nei-tre-livelli) | il firewall del gateway, con port forwarding verso `nodo2` |
    | [1.12](12-http-nginx.md#nei-tre-livelli) | Nginx su `nodo2`, pubblicato da `nodo1` |
    | [Progetto](progetto.md#nei-tre-livelli) | il server web sicuro su `nodo2` (e `nodo3`), pubblicato dal gateway |

!!! warning "Sulle macchine fisiche non ci sono snapshot"
    Prima di modificare un file di configurazione su `srv1` o su un nodo, fanne una copia
    (`sudo cp -a /etc/ssh/sshd_config.d /root/sshd_config.d.bak`) e tieni **aperta una seconda sessione SSH** finché
    non hai verificato che la modifica funziona. Se perdi l'accesso, schermo e tastiera collegati al server sono la
    tua console di emergenza.

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
