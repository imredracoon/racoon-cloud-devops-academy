# Red Hat Certified System Administrator (RHCSA, EX200)

La certificazione Linux più riconosciuta nelle aziende che usano Red Hat Enterprise Linux. L'esame è **pratico**:
ricevi una o più macchine e devi portarle allo stato richiesto. Conta solo il risultato, e deve **sopravvivere al
riavvio**.

| | |
|---|---|
| **Ente** | Red Hat |
| **Codice** | EX200 |
| **Sistema** | Red Hat Enterprise Linux 10 |
| **Formato** | pratico, in un centro d'esame o da remoto (Remote Exam) |
| **Durata** | circa 3 ore (la pagina dell'esame non la indica: verificala al momento della prenotazione) |
| **Aiuti** | nessun accesso a Internet né appunti; solo la documentazione installata nel sistema (man, info, `/usr/share/doc`) |
| **Costo** | varia per paese: vedi la pagina ufficiale |
| **Validità** | 3 anni; si estende ottenendo altre certificazioni Red Hat |
| **Pagina ufficiale** | [EX200 su redhat.com](https://www.redhat.com/en/services/training/ex200-red-hat-certified-system-administrator-rhcsa-exam) |

Dati verificati il 28 ottobre 2026.

## Quando sostenerla

Dopo la [Fase 1](../../fase-1/index.md) seguita sul percorso **RHEL (Rocky/AlmaLinux)**, più le lezioni [2.3](../../fase-2/03-bash-base.md)-[2.5](../../fase-2/05-pianificazione.md) su
script e pianificazione. Allenati su **Rocky Linux 10 o AlmaLinux 10**: sono compatibili con RHEL 10 e gratuite.
In alternativa Red Hat offre l'abbonamento gratuito *Red Hat Developer* per usare RHEL in laboratorio.

## Piano di studio (8 settimane)

| Settimana | Attività |
|---|---|
| 1 | Strumenti essenziali, documentazione, permessi ([1.1](../../fase-1/01-filesystem.md)-[1.3](../../fase-1/03-permessi.md)) |
| 2 | Software: repository, `dnf`, Flatpak ([1.5](../../fase-1/05-pacchetti.md)); utenti, gruppi, sudo ([1.3](../../fase-1/03-permessi.md)) |
| 3 | Processi, servizi, log, journal persistente, tuned ([1.4](../../fase-1/04-processi-systemd.md), [1.15](../../fase-1/15-avvio-tempo.md)) |
| 4 | Avvio: target, GRUB, recupero della password di root ([1.15](../../fase-1/15-avvio-tempo.md)) |
| 5 | Storage: partizioni GPT, LVM, swap, VFAT/ext4/XFS, NFS, autofs ([1.6](../../fase-1/06-storage.md)) |
| 6 | Rete con `nmcli`, risoluzione dei nomi, firewalld, ora ([1.7](../../fase-1/07-tcpip.md), [1.8](../../fase-1/08-dns-routing.md), [1.11](../../fase-1/11-firewall.md), [1.15](../../fase-1/15-avvio-tempo.md)) |
| 7 | SELinux, chiavi SSH, script Bash, at, cron e timer ([1.10](../../fase-1/10-ssh.md), [1.14](../../fase-1/14-selinux-apparmor.md), [2.3](../../fase-2/03-bash-base.md)-[2.5](../../fase-2/05-pianificazione.md)) |
| 8 | Esami pratici a tempo con riavvio finale, ripasso dei punti deboli |

!!! tip "Impara a cercare nella documentazione locale"
    All'esame non hai Internet. Esercitati a trovare le risposte con `man -k parola`, `man 5 fstab`,
    `ls /usr/share/doc/` e gli esempi delle pagine man (per esempio `man semanage-fcontext` e `man nmcli-examples`).

## Il percorso Red Hat

- **RHCSA** (EX200) è la base e il prerequisito delle certificazioni successive.
- **RHCE in Ansible** si ottiene con RHCSA + EX294: vedi l'[appendice EX294](../rhce/index.md).
- **RHCE in Enterprise Linux** si ottiene con RHCSA + EX342, un esame avanzato su RHEL non trattato in questo corso.

## In questa appendice

- [Argomenti specifici](argomenti.md)
- [Simulazione pratica](simulazione.md): 43 task originali con verifica e soluzione
