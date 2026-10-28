# Argomenti specifici · LFCS

Le competenze ufficiali, dominio per dominio, con le lezioni del corso.

## Operations Deployment (25%)

| Competenza | Lezioni |
|---|---|
| Parametri del kernel, persistenti e non | [1.15](../../fase-1/15-avvio-tempo.md) |
| Diagnosticare e gestire processi e servizi | [1.4](../../fase-1/04-processi-systemd.md) |
| Pianificare job | [2.5](../../fase-2/05-pianificazione.md) |
| Trovare, installare, verificare e mantenere pacchetti e repository | [1.5](../../fase-1/05-pacchetti.md) |
| Ripristino da guasti hardware, del sistema o dei filesystem | [1.6](../../fase-1/06-storage.md), [1.13](../../fase-1/13-troubleshooting.md), [1.15](../../fase-1/15-avvio-tempo.md) |
| Macchine virtuali con libvirt | [Approfondimento](#macchine-virtuali-con-libvirt), [0.7](../../fase-0/07-laboratorio.md) |
| Motori di container e gestione dei container | [3.2](../../fase-3/02-immagini-container.md)-[3.5](../../fase-3/05-compose.md) |
| Controllo d'accesso obbligatorio con SELinux | [1.14](../../fase-1/14-selinux-apparmor.md) |

## Networking (25%)

| Competenza | Lezioni |
|---|---|
| IPv4, IPv6 e risoluzione dei nomi | [1.7](../../fase-1/07-tcpip.md), [1.8](../../fase-1/08-dns-routing.md) |
| Sincronizzazione dell'ora | [1.15](../../fase-1/15-avvio-tempo.md) |
| Monitorare e diagnosticare la rete | [1.9](../../fase-1/09-diagnostica.md) |
| Server e client OpenSSH | [1.10](../../fase-1/10-ssh.md) |
| Filtraggio dei pacchetti, redirezione di porta e NAT | [1.8](../../fase-1/08-dns-routing.md), [1.11](../../fase-1/11-firewall.md) |
| Routing statico | [1.8](../../fase-1/08-dns-routing.md) |
| Bridge e bonding | [Approfondimento](#bridge-e-bonding) |
| Reverse proxy e bilanciatori | [1.12](../../fase-1/12-http-nginx.md) |

## Storage (20%)

| Competenza | Lezioni |
|---|---|
| LVM | [1.6](../../fase-1/06-storage.md) |
| Filesystem virtuale e mount | [1.1](../../fase-1/01-filesystem.md), [1.6](../../fase-1/06-storage.md) |
| Creare, gestire e riparare filesystem | [1.6](../../fase-1/06-storage.md) |
| Filesystem remoti e dispositivi a blocchi di rete | [1.6](../../fase-1/06-storage.md), [approfondimento](#dispositivi-a-blocchi-di-rete) |
| Swap | [1.6](../../fase-1/06-storage.md) |
| Montaggio automatico | [1.6](../../fase-1/06-storage.md) |
| Prestazioni dello storage | [1.15](../../fase-1/15-avvio-tempo.md), [approfondimento](#prestazioni-dello-storage) |

## Essential Commands (20%)

| Competenza | Lezioni |
|---|---|
| Git di base | [2.1](../../fase-2/01-git-base.md) |
| Creare, configurare e riparare servizi | [1.4](../../fase-1/04-processi-systemd.md) |
| Prestazioni del sistema e dei servizi | [1.15](../../fase-1/15-avvio-tempo.md) |
| Vincoli di applicazioni e servizi (limiti, porte, dipendenze) | [1.4](../../fase-1/04-processi-systemd.md), [1.13](../../fase-1/13-troubleshooting.md) |
| Problemi di spazio su disco | [1.6](../../fase-1/06-storage.md), [1.13](../../fase-1/13-troubleshooting.md) |
| Certificati SSL/TLS | [1.12](../../fase-1/12-http-nginx.md) |

## Users and Groups (10%)

| Competenza | Lezioni |
|---|---|
| Utenti e gruppi locali | [1.3](../../fase-1/03-permessi.md) |
| Profili d'ambiente personali e di sistema | [Approfondimento](#profili-e-limiti) |
| Limiti delle risorse per utente | [Approfondimento](#profili-e-limiti) |
| ACL | [1.3](../../fase-1/03-permessi.md) |
| Utenti e gruppi da LDAP | [Approfondimento](#utenti-da-ldap-con-sssd) |

## Approfondimenti

Questi argomenti dell'esame non hanno una lezione dedicata nel corso. Qui trovi l'essenziale; la
[simulazione](simulazione.md) contiene task completi con soluzione spiegata.

### Bridge e bonding

Un **bridge** è uno switch virtuale: collega più interfacce (fisiche, di VM o di container) nello stesso segmento di
livello 2. Un **bond** unisce più interfacce fisiche in una sola, per tolleranza ai guasti (`active-backup`) o per
banda (`802.3ad`, che richiede la configurazione dello switch). Su Ubuntu entrambi si definiscono con netplan:

```yaml
network:
  version: 2
  ethernets:
    enp0s8: {dhcp4: false}
    enp0s9: {dhcp4: false}
  bonds:
    bond0:
      interfaces: [enp0s8, enp0s9]
      parameters:
        mode: active-backup
        primary: enp0s8
      addresses: [192.168.70.10/24]
```

Verifica con `cat /proc/net/bonding/bond0` e, per i bridge, `bridge link`. Su RHEL si usa
`nmcli con add type bond ...` e `nmcli con add type bridge ...`.
Fonte: [documentazione di netplan](https://netplan.readthedocs.io/en/stable/netplan-yaml/).

### Macchine virtuali con libvirt

libvirt gestisce le VM KVM; `virsh` è il suo client da riga di comando.

```bash
sudo apt install -y qemu-kvm libvirt-daemon-system virtinst
virsh list --all                      # tutte le VM
virsh start prova; virsh shutdown prova; virsh destroy prova   # avvio, arresto ordinato, arresto forzato
virsh autostart prova                 # avvio insieme all'host
virsh dominfo prova; virsh edit prova # informazioni e modifica della definizione XML
virsh setvcpus prova 2 --config --maximum && virsh setvcpus prova 2 --config
virsh setmaxmem prova 2G --config && virsh setmem prova 2G --config
```

`--config` rende la modifica permanente (effettiva al prossimo avvio). Fonte:
[Ubuntu Server · libvirt](https://documentation.ubuntu.com/server/how-to/virtualisation/libvirt/).

### Profili e limiti

- `/etc/profile` e i file `/etc/profile.d/*.sh` valgono per le shell di login di tutti gli utenti;
  `~/.profile` (o `~/.bash_profile`) e `~/.bashrc` per il singolo utente. `/etc/skel` contiene i file copiati nella
  home dei nuovi utenti.
- I limiti per utente o gruppo stanno in `/etc/security/limits.conf` e `/etc/security/limits.d/*.conf`
  (`@gruppo hard nproc 100`), applicati da pam_limits alle nuove sessioni; `ulimit -a` mostra quelli correnti. Per
  i servizi systemd si usano invece `LimitNOFILE=`, `MemoryMax=` e simili nella unit.

Fonte: [limits.conf(5)](https://manpages.ubuntu.com/manpages/noble/en/man5/limits.conf.5.html).

### Utenti da LDAP con SSSD

SSSD collega il sistema a una directory LDAP (o Active Directory) per utenti, gruppi e autenticazione, con cache
locale. Passi essenziali: installare `sssd-ldap`, scrivere `/etc/sssd/sssd.conf` (permessi 600) con
`id_provider = ldap`, `ldap_uri` e `ldap_search_base`, riavviare `sssd`, abilitare la creazione della home
(`pam-auth-update --enable mkhomedir`), verificare con `getent passwd UTENTE`. Fonte:
[Ubuntu Server · SSSD con LDAP](https://documentation.ubuntu.com/server/how-to/sssd/with-ldap/).

### Dispositivi a blocchi di rete

Oltre ai filesystem di rete (NFS, SMB), un disco può arrivare dalla rete **come dispositivo a blocchi**: con iSCSI
(`open-iscsi`, `iscsiadm -m discovery` e `-m node --login`) o con NBD (`nbd-client`). Il sistema lo vede come un
normale `/dev/sdX` o `/dev/nbdX`, da partizionare e formattare; in `/etc/fstab` va l'opzione `_netdev`.
Fonte: [Ubuntu Server · iSCSI](https://documentation.ubuntu.com/server/how-to/storage/iscsi-initiator-or-client/).

### Prestazioni dello storage

`iostat -dx 1` (pacchetto `sysstat`) mostra per ogni disco operazioni al secondo, latenza (`r_await`, `w_await`) e
occupazione (`%util`); `iotop` mostra i processi che fanno più I/O; `vmstat 1` indica nella colonna `wa` il tempo
in cui la CPU attende il disco. Fonte: [iostat(1)](https://manpages.ubuntu.com/manpages/noble/en/man1/iostat.1.html).

## Dove sbagliano più spesso i candidati

- Modifiche **non persistenti**: `sysctl -w`, `ip addr add`, `mount` senza `/etc/fstab` spariscono al riavvio.
- **`/etc/fstab` sbagliato**: prova sempre con `mount -a` (e `findmnt --verify`).
- **SSH**: dopo una modifica, controlla con `sshd -t` e tieni aperta la sessione corrente.
- **Firewall**: abilita la regola SSH prima di attivare il firewall.
