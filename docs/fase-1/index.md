# Fase 1 · Linux e reti

!!! obiettivi "Al termine di questa fase saprai"
    - muoverti con sicurezza nel filesystem Linux e manipolare file e testo con le pipe;
    - gestire utenti, gruppi, permessi, processi e servizi systemd;
    - installare software e gestire dischi, filesystem e LVM;
    - calcolare subnet, spiegare DNS, DHCP, routing e NAT e diagnosticare problemi di rete;
    - collegarti ai server con SSH e chiavi, proteggerli con un firewall nftables;
    - pubblicare un sito con Nginx e HTTPS e risolvere i guasti con metodo.

| | |
|---|---|
| **Durata indicativa** | 2,5 mesi, circa 90 ore |
| **Prerequisiti** | Fase 0 |
| **Ambiente** | VM Linux di laboratorio (`lab`); per alcuni esercizi di rete una seconda VM |
| **Progetto** | [Server web messo in sicurezza](progetto.md) |
| **Certificazione collegata** | [LFCS](../appendici/lfcs/index.md) (facoltativa) |

## Perché Linux e reti prima di tutto

Ogni tecnologia che studierai dopo poggia su questi due pilastri. Un container è un insieme di processi Linux
isolati. Un'istanza nel cloud è una VM Linux. Un cluster Kubernetes è una rete di macchine Linux. Quando qualcosa
si rompe, la causa si trova quasi sempre in uno di questi strati: un permesso, un servizio fermo, un disco pieno,
un DNS sbagliato, una porta chiusa. Chi padroneggia Linux e reti risolve problemi che gli altri non capiscono.

## Preparare la seconda VM

Alcuni esercizi di rete (SSH, firewall, Nginx) sono più realistici con due macchine. Quando la lezione lo
richiede, crea una seconda VM leggera:

```bash
multipass launch 24.04 --name client --cpus 1 --memory 1G --disk 8G
```

Le due VM si vedono sulla stessa rete virtuale di Multipass: `multipass list` mostra i loro indirizzi IP.

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

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

!!! tip "Fai uno snapshot"
    Prima di ogni lezione dalla 1.3 in poi, spegni la VM e crea uno snapshot (`multipass stop lab`,
    `multipass snapshot lab`). Se rompi qualcosa, puoi tornare indietro in un minuto.
