# Hardware, server e data center

Il cloud è fatto di computer veri: server montati in armadi, collegati da switch e alimentati da UPS, dentro
edifici progettati per non spegnersi mai. Questa appendice apre quelle porte. È facoltativa e non prepara a un
esame: serve a capire **cosa c'è sotto** le macchine virtuali che hai creato nella [Fase 4](../../fase-4/index.md)
e a muoverti con sicurezza se un giorno lavorerai in una sala server, in un'azienda con infrastruttura propria
(*on-premises*) o in un provider.

!!! obiettivi "Obiettivi"
    - Riconoscere i dispositivi di una rete e capire il ruolo di ciascuno.
    - Sapere cosa si trova in un armadio rack e a cosa serve.
    - Leggere la scheda tecnica di un server e scegliere l'hardware in base al lavoro che deve fare.
    - Conoscere le tipologie di server per servizio e il software che le realizza su Ubuntu e RHEL.
    - Installare e mantenere un server fisico, dal firmware al sistema operativo.
    - Capire come si progetta un data center e quali tipi esistono in base allo scopo.

## Come è organizzata

| Pagina | Di cosa parla |
|---|---|
| [1. Dispositivi di rete](01-dispositivi-rete.md) | modem, router, switch, firewall, access point, load balancer, cavi e transceiver |
| [2. Il rack](02-rack.md) | armadio, unità rack, PDU, UPS, KVM, patch panel, console server, raffreddamento |
| [3. L'hardware di un server](03-hardware-server.md) | CPU, RAM, dischi, RAID, schede di rete, BMC, formati, con modelli reali |
| [4. Server di infrastruttura](04-server-infrastruttura.md) | DNS, DHCP, NTP, directory, CA, provisioning, proxy, VPN, firewall, bastion, posta, stampa, VoIP |
| [5. Server per i dati](05-server-dati.md) | file server, NAS, SAN e storage a blocchi, object storage, backup, database, cache, code di messaggi |
| [6. Server per calcolo e applicazioni](06-server-applicazioni.md) | web e applicazioni, virtualizzazione, container, GPU e AI, HPC, CI/CD, Git, mirror, monitoraggio, log, media e giochi |
| [7. Installare un server fisico](07-installazione.md) | dal disimballaggio al sistema operativo, a mano e in automatico |
| [8. Manutenzione](08-manutenzione.md) | aggiornamenti, firmware, salute dell'hardware, backup, ciclo di vita |
| [9. Architettura di un data center](09-data-center.md) | topologie, progettazione, ridondanza, tipologie di data center e loro scopo |
| [Verifica](verifica.md) | quiz finale |

## Come leggerla

Molti concetti li hai già studiati nel corso. Qui sono **solo richiamati**, con il collegamento alla lezione
che li spiega: se un termine non ti è chiaro, apri il collegamento e torna qui dopo.

!!! abstract "Prima di cominciare, ripassa"
    - [Come funziona un computer](../../fase-0/01-computer.md) e in particolare [il server](../../fase-0/01-computer.md#il-server) (lezione 0.1).
    - [Reti e Internet: i concetti](../../fase-0/03-reti-concetti.md) (lezione 0.3).
    - [TCP/IP e indirizzamento](../../fase-1/07-tcpip.md) (lezione 1.7) e [DNS, DHCP, routing e NAT](../../fase-1/08-dns-routing.md) (lezione 1.8).
    - [Dischi, filesystem e LVM](../../fase-1/06-storage.md) (lezione 1.6).

I comandi sono mostrati per le due famiglie del corso, in schede affiancate come nelle lezioni:
**Ubuntu** (famiglia Debian) e **RHEL** (con Rocky Linux e AlmaLinux come alternative gratuite compatibili).

!!! warning "I modelli cambiano"
    I prodotti citati sono **esempi reali** verificati sulle pagine ufficiali dei produttori il 28 ottobre 2026.
    I server si rinnovano ogni due o tre anni: quando leggi, controlla la generazione attuale. I principi,
    invece, restano validi a lungo.

## Fonti

- [Ubuntu Server documentation](https://documentation.ubuntu.com/server/)
- [Red Hat Enterprise Linux 10, documentazione](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10)
