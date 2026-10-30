# 1. FreeBSD e Linux

!!! obiettivi "Obiettivi"
    - Sapere da dove viene FreeBSD e come è organizzato il progetto.
    - Riconoscere le differenze principali rispetto a una distribuzione Linux.
    - Tradurre i comandi Linux che conosci nei loro equivalenti FreeBSD.

## Un sistema, non una distribuzione

Linux è un **kernel**: una distribuzione come Ubuntu o AlmaLinux mette insieme quel kernel con strumenti presi
da tanti progetti diversi (GNU, systemd, e così via) ([lezione 0.2](../../fase-0/02-sistemi-operativi.md)). FreeBSD invece è un
**sistema operativo completo**: kernel, comandi di base, librerie e documentazione sono sviluppati insieme, nello
stesso repository, dallo stesso progetto. Discende da **BSD**, lo Unix dell'Università di Berkeley degli anni '70
e '80.

Le conseguenze pratiche:

- c'è un solo FreeBSD, non tante distribuzioni: una guida vale per tutti;
- il **sistema di base** (`/bin`, `/sbin`, `/usr/bin`, `/etc`) e il **software aggiunto** sono separati: tutto
  quello che installi con i pacchetti finisce in **`/usr/local`**, anche i file di configurazione
  (`/usr/local/etc/nginx/nginx.conf`, non `/etc/nginx/`);
- la documentazione ufficiale, l'**Handbook**, è molto completa e aggiornata: è il primo posto dove cercare;
- la licenza è la **BSD**, più permissiva della GPL di Linux: chi modifica FreeBSD non è obbligato a pubblicare le
  modifiche. Per questo molte aziende lo usano come base per prodotti chiusi.

## Le differenze che noterai subito

| Argomento | Linux (Ubuntu, RHEL) | FreeBSD |
|---|---|---|
| Avvio e servizi | systemd, unità `.service` | script `rc.d`, servizi accesi in `/etc/rc.conf` |
| Shell predefinita | bash | `sh` per tutti, anche per root (bash si installa come pacchetto) |
| Pacchetti | `apt`, `dnf` | `pkg` per i pacchetti già pronti, *ports* per compilare dai sorgenti |
| Nomi delle schede di rete | `enp1s0`, `eth0` (dal nome del dispositivo) | dal nome del driver: `em0`, `igb0`, `vtnet0` (virtio) |
| Nomi dei dischi | `sda`, `nvme0n1`, `vda` | `ada0` (SATA), `nda0` (NVMe), `vtbd0` (virtio) |
| Firewall | nftables, firewalld, ufw | **pf** (e anche ipfw) |
| Container | namespace e cgroup (Docker, Podman) | **jail** |
| Virtualizzazione | KVM | **bhyve** |
| Filesystem | ext4, XFS | **ZFS** o UFS |
| Log | journald e `/var/log` | `syslogd` e `/var/log/messages` |

FreeBSD può anche **eseguire programmi compilati per Linux** con un livello di compatibilità (il *Linuxulator*):
utile per qualche programma chiuso, non per far girare Docker.

## I comandi equivalenti

| Cosa vuoi fare | Linux | FreeBSD |
|---|---|---|
| Avviare un servizio | `systemctl start nginx` | `service nginx start` |
| Attivarlo all'avvio | `systemctl enable nginx` | `sysrc nginx_enable=YES` |
| Stato dei servizi | `systemctl status nginx` | `service nginx status` |
| Installare un pacchetto | `apt install nginx` / `dnf install nginx` | `pkg install nginx` |
| Cercare un pacchetto | `apt search` / `dnf search` | `pkg search` |
| A quale pacchetto appartiene un file | `dpkg -S` / `rpm -qf` | `pkg which` |
| Indirizzi di rete | `ip address` | `ifconfig` |
| Tabella di routing | `ip route` | `netstat -rn` |
| Porte in ascolto | `ss -tlpn` | `sockstat -l` |
| Elenco dei dischi | `lsblk` | `geom disk list`, `gpart show` |
| Moduli del kernel | `lsmod`, `modprobe` | `kldstat`, `kldload` |
| Aggiungere un utente | `useradd` | `adduser` (interattivo) o `pw useradd` |
| Versione del sistema | `cat /etc/os-release` | `freebsd-version -ku` |
| Aggiornare il sistema | `apt upgrade` / `dnf upgrade` | `freebsd-update` o `pkg upgrade` ([pagina 2](02-installazione.md#gli-aggiornamenti)) |
| Log del sistema | `journalctl` | `tail /var/log/messages` |

Molti comandi sono uguali (`ls`, `cp`, `grep`, `ps`, `top`, `sysctl`, `ssh`), ma le **opzioni** possono cambiare:
gli strumenti di base di FreeBSD non sono quelli GNU. Per esempio `sed -i` vuole un argomento
(`sed -i '' 's/a/b/' file`), e `ps aux` funziona ma `ps -ef` dà un risultato diverso. Nel dubbio: `man comando`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Dove installa i file di configurazione un pacchetto su FreeBSD?","options":["/etc","/usr/local/etc","/opt/etc","/var/etc"],"answer":[1],"explanation":"Il software aggiunto vive in /usr/local, separato dal sistema di base."},
 {"q":"Come si attiva all'avvio il servizio nginx su FreeBSD?","options":["systemctl enable nginx","sysrc nginx_enable=YES","chkconfig nginx on","pkg enable nginx"],"answer":[1],"explanation":"sysrc scrive la variabile in /etc/rc.conf."},
 {"q":"Una scheda di rete si chiama vtnet0. Cosa ti dice il nome?","options":["È la prima scheda Wi-Fi","Usa il driver virtio: sei in una VM","È una VLAN","È una scheda Intel"],"answer":[1],"explanation":"In FreeBSD le interfacce prendono il nome dal driver."},
 {"q":"Perché molte aziende scelgono FreeBSD come base per prodotti chiusi?","options":["È più veloce di Linux","La licenza BSD non obbliga a pubblicare le modifiche","Ha più pacchetti","È l'unico con ZFS"],"answer":[1],"explanation":"La GPL di Linux chiede di distribuire i sorgenti modificati; la BSD no."}
]}
</script>

## Fonti

- [FreeBSD Handbook: Introduction](https://docs.freebsd.org/en/books/handbook/introduction/)
- [FreeBSD Handbook: FreeBSD Basics](https://docs.freebsd.org/en/books/handbook/basics/)
- [FreeBSD Handbook: Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/)
- [sysrc(8)](https://man.freebsd.org/cgi/man.cgi?query=sysrc&sektion=8)
