# Argomenti specifici · RHCSA (EX200)

I dieci gruppi di obiettivi pubblicati da Red Hat per RHEL 10. Per ogni obiettivo: lezioni del corso, sempre nella
scheda **RHEL (Rocky/AlmaLinux)**.

| Obiettivi | Cosa devi saper fare | Lezioni |
|---|---|---|
| **1 · Strumenti essenziali** | shell e redirezioni, `grep` ed espressioni regolari, SSH, cambio utente, `tar` con gzip e bzip2, creare e modificare file, collegamenti fisici e simbolici, permessi ugo/rwx, documentazione (man, info, `/usr/share/doc`) | [0.4](../../fase-0/04-terminale.md), [1.1](../../fase-1/01-filesystem.md)-[1.3](../../fase-1/03-permessi.md), [1.10](../../fase-1/10-ssh.md) |
| **2 · Gestire il software** | repository RPM, installare e rimuovere pacchetti con `dnf`, repository e pacchetti **Flatpak** | [1.5](../../fase-1/05-pacchetti.md) |
| **3 · Script di shell** | condizioni (`if`, `test`, `[ ]`), cicli, argomenti `$1 $2`, output dei comandi negli script | [2.3](../../fase-2/03-bash-base.md), [2.4](../../fase-2/04-bash-avanzato.md) |
| **4 · Sistemi in esecuzione** | avvio, riavvio e spegnimento, target manuali, interrompere l'avvio per accedere al sistema, processi che consumano CPU e memoria, priorità, **profili di tuning**, log e journal, **journal persistente**, servizi di rete, trasferimento sicuro di file | [1.4](../../fase-1/04-processi-systemd.md), [1.10](../../fase-1/10-ssh.md), [1.15](../../fase-1/15-avvio-tempo.md) |
| **5 · Storage locale** | partizioni su dischi GPT, physical volume, volume group, logical volume, mount all'avvio per UUID o etichetta, aggiungere partizioni, LV e swap senza perdere dati | [1.6](../../fase-1/06-storage.md) |
| **6 · File system** | creare e montare VFAT, ext4 e XFS, NFS, **autofs**, estendere LV, diagnosticare problemi di permessi | [1.6](../../fase-1/06-storage.md), [1.3](../../fase-1/03-permessi.md) |
| **7 · Mantenere i sistemi** | `at`, cron e timer systemd, servizi all'avvio, target predefinito, client NTP, installare e aggiornare pacchetti da CDN, repository remoto o file locale, modificare il bootloader | [1.4](../../fase-1/04-processi-systemd.md), [1.5](../../fase-1/05-pacchetti.md), [1.15](../../fase-1/15-avvio-tempo.md), [2.5](../../fase-2/05-pianificazione.md) |
| **8 · Rete di base** | indirizzi IPv4 e IPv6, risoluzione dei nomi, servizi di rete all'avvio, restrizioni con firewalld e `firewall-cmd` | [1.7](../../fase-1/07-tcpip.md), [1.8](../../fase-1/08-dns-routing.md), [1.11](../../fase-1/11-firewall.md) |
| **9 · Utenti e gruppi** | creare, eliminare e modificare utenti; password e scadenze; gruppi e appartenenze; accesso da superutente | [1.3](../../fase-1/03-permessi.md) |
| **10 · Sicurezza** | firewalld, permessi predefiniti (umask), autenticazione con chiavi SSH, modalità enforcing e permissive di SELinux, contesti di file e processi, ripristino dei contesti, etichette delle porte, booleani, diagnosi delle violazioni | [1.10](../../fase-1/10-ssh.md), [1.11](../../fase-1/11-firewall.md), [1.14](../../fase-1/14-selinux-apparmor.md), [1.3](../../fase-1/03-permessi.md) |

## Novità della versione RHEL 10

Rispetto alla versione basata su RHEL 9 sono entrati **Flatpak** e i **profili di tuning**, e la gestione dei
container non è più tra gli obiettivi. Controlla sempre l'elenco aggiornato sulla
[pagina ufficiale](https://www.redhat.com/en/services/training/ex200-red-hat-certified-system-administrator-rhcsa-exam).

## Dove sbagliano più spesso i candidati

- **Persistenza**: `setenforce`, `ip addr add`, `mount` e `firewall-cmd` senza `--permanent` spariscono al riavvio.
- **Recupero di root senza `/.autorelabel`**: il sistema si avvia ma nessuno riesce più ad accedere.
- **`chcon` invece di `semanage fcontext`**: il contesto si perde al primo relabel.
- **`/etc/fstab`**: un errore blocca l'avvio e fa perdere tutti i punti. Prova con `mount -a` e
  `findmnt --verify`, e riavvia prima della fine.
- **Nomi delle connessioni**: `nmcli` lavora sulle connessioni, che possono chiamarsi diversamente dalle interfacce.
