# Linux Foundation Certified System Administrator (LFCS)

La certificazione Linux pratica della Linux Foundation. Lavori da riga di comando su sistemi reali e risolvi task
di amministrazione: servizi, rete, storage, utenti, sicurezza, container e macchine virtuali.

| | |
|---|---|
| **Ente** | Linux Foundation |
| **Formato** | task pratici da riga di comando, online con sorvegliante |
| **Durata** | 2 ore |
| **Distribuzione** | dichiarato indipendente dalla distribuzione |
| **Costo** | 445 USD (esame) |
| **Incluso** | 12 mesi per sostenere l'esame, un nuovo tentativo, due sessioni del simulatore Killer.sh |
| **Validità** | 2 anni |
| **Pagina ufficiale** | [training.linuxfoundation.org/certification/linux-foundation-certified-sysadmin-lfcs](https://training.linuxfoundation.org/certification/linux-foundation-certified-sysadmin-lfcs/) |

Dati verificati il 28 ottobre 2026. La simulazione usa una soglia indicativa del 67%.

## LFCS, RHCSA o Canonical?

| | LFCS | RHCSA | Canonical Academy |
|---|---|---|---|
| Sistema | indipendente dalla distribuzione | RHEL 10 | Ubuntu |
| Tipo | pratico | pratico | ibrido |
| Riconoscimento | buono, internazionale | molto alto nelle aziende che usano Red Hat | recente, ancora poco noto |
| Aiuto durante l'esame | documentazione del sistema | solo documentazione del sistema, niente Internet | documentazione Ubuntu e man |

Ne basta **una**. Scegli in base al sistema che userai al lavoro.

## Quando sostenerla

Dopo la [Fase 1](../../fase-1/index.md), seguendo il percorso **Debian/Ubuntu** e leggendo le schede RHEL per i
temi come SELinux. Alcuni argomenti (container, Git, script) li approfondisci nelle Fasi 2 e 3: se vuoi la massima
sicurezza, sostienila dopo la [Fase 3](../../fase-3/index.md).

## Piano di studio (6 settimane)

| Settimana | Attività |
|---|---|
| 1 | Comandi essenziali, Git di base, servizi ([1.1](../../fase-1/01-filesystem.md)-[1.4](../../fase-1/04-processi-systemd.md), [2.1](../../fase-2/01-git-base.md)) |
| 2 | Utenti, gruppi, ACL, limiti, LDAP ([1.3](../../fase-1/03-permessi.md)) |
| 3 | Storage: partizioni, LVM, swap, NFS, autofs ([1.6](../../fase-1/06-storage.md)) |
| 4 | Rete: indirizzi, routing, DNS, ora, SSH, firewall, NAT, proxy ([1.7](../../fase-1/07-tcpip.md)-[1.12](../../fase-1/12-http-nginx.md)) |
| 5 | Operazioni: parametri del kernel, processi, pianificazione, pacchetti, SELinux, container, VM ([1.4](../../fase-1/04-processi-systemd.md), [1.5](../../fase-1/05-pacchetti.md), [1.14](../../fase-1/14-selinux-apparmor.md), [1.15](../../fase-1/15-avvio-tempo.md), [2.5](../../fase-2/05-pianificazione.md), [3.2](../../fase-3/02-immagini-container.md)) |
| 6 | Esami pratici a tempo, poi le sessioni Killer.sh incluse |

## In questa appendice

- [Argomenti specifici](argomenti.md)
- [Simulazione pratica](simulazione.md): 41 task originali con verifica e soluzione
