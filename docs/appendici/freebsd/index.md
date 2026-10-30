# FreeBSD

FreeBSD è un sistema operativo Unix libero, parente stretto di Linux ma con una storia e un'organizzazione
diverse. Questa appendice è **breve e facoltativa**: non fa parte del programma e non prepara a un esame. Serve a
conoscere un altro Unix e a riconoscerlo quando lo incontri, soprattutto in firewall e apparecchi di rete come
**OPNsense** e **pfSense**, che puoi mettere nel tuo [homelab](../../ambiente/livelli.md).

!!! obiettivi "Obiettivi"
    - Spiegare in cosa FreeBSD è diverso da Linux e dove si usa nel lavoro.
    - Installare FreeBSD in una VM e amministrarlo: pacchetti, servizi, aggiornamenti.
    - Usare le sue tecnologie più note: ZFS, le jail, il firewall pf e l'hypervisor bhyve.
    - Mettere OPNsense come firewall davanti a una rete del laboratorio, nei tre livelli del corso.

## Perché conoscerlo

Nelle offerte di lavoro Cloud e DevOps FreeBSD compare poco: Kubernetes e i container che usi nel corso sono
pensati per Linux. Però FreeBSD è dietro a prodotti che un sistemista incontra davvero:

- **firewall e router**: OPNsense e pfSense sono basati su FreeBSD, e molte piccole e medie aziende li usano come
  firewall perimetrale;
- **rete**: il sistema operativo Junos dei router e switch Juniper è nato su FreeBSD;
- **distribuzione di contenuti**: i server della CDN di Netflix (*Open Connect*) usano FreeBSD;
- **storage**: ZFS è nato su Solaris ed è integrato in FreeBSD da molti anni.

Capire un Unix diverso ti fa anche capire meglio Linux: vedi quali scelte sono "Unix" e quali sono solo di Linux.

## Come è organizzata

| Pagina | Di cosa parla |
|---|---|
| [1. FreeBSD e Linux](01-freebsd-e-linux.md) | la storia, le differenze, la tabella dei comandi equivalenti |
| [2. Installare e amministrare](02-installazione.md) | la VM di laboratorio, pkg e ports, `rc.conf` e i servizi, gli aggiornamenti |
| [3. ZFS](03-zfs.md) | pool, dataset, snapshot, copie con send e receive, boot environment |
| [4. Le jail](04-jail.md) | i "container" di FreeBSD, una jail con un server web |
| [5. pf e bhyve](05-pf-bhyve.md) | il firewall pf e l'hypervisor di FreeBSD |
| [6. Laboratorio: OPNsense nell'homelab](06-opnsense.md) | un firewall vero davanti a una rete del laboratorio, nei tre livelli |
| [Verifica](verifica.md) | quiz finale |

!!! abstract "Prima di cominciare, ripassa"
    - [Sistemi operativi e Linux](../../fase-0/02-sistemi-operativi.md) (lezione 0.2).
    - [Processi e servizi con systemd](../../fase-1/04-processi-systemd.md) (lezione 1.4) e [pacchetti](../../fase-1/05-pacchetti.md) (lezione 1.5).
    - [Firewall con nftables](../../fase-1/11-firewall.md) (lezione 1.11) e [i container](../../fase-3/01-concetti.md) (lezione 3.1).
    - [Macchine virtuali Linux](../../ambiente/vm-linux.md): KVM e libvirt, che userai per la VM FreeBSD.

!!! warning "Le versioni"
    Versioni verificate il 30 ottobre 2026: **FreeBSD 15.1** (giugno 2026) è la versione di produzione, la serie
    14 è ancora supportata; **OPNsense 26.7** (luglio 2026) è basata su FreeBSD 15.1. Quando leggi, controlla la
    [pagina dei rilasci](https://www.freebsd.org/releases/).

## Fonti

- [FreeBSD Handbook](https://docs.freebsd.org/en/books/handbook/)
- [FreeBSD: informazioni sui rilasci](https://www.freebsd.org/releases/)
- [OPNsense: note della serie 26.7](https://wiki.opnsense.org/releases/CE_26.7.html)
- [Netflix Open Connect](https://openconnect.netflix.com/)
