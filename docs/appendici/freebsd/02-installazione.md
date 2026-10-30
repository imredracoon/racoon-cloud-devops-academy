# 2. Installare e amministrare

!!! obiettivi "Obiettivi"
    - Creare una VM FreeBSD con KVM e libvirt, nel livello del laboratorio che usi.
    - Installare software con `pkg` e con i ports.
    - Gestire i servizi con `rc.conf`, `sysrc` e `service`.
    - Aggiornare il sistema di base e i pacchetti.

## La VM di laboratorio

Lo script `academy-vm` crea solo VM Linux. FreeBSD si installa dall'immagine ufficiale, con `virt-install` e lo
stesso KVM/libvirt delle [macchine virtuali Linux](../../ambiente/vm-linux.md), sulla rete `academy`.

=== "Livello 1 · Un PC"

    Lavori sulla workstation, con la grafica di virt-manager a portata di mano.

=== "Livello 2 · Workstation + server"

    La VM nasce su `srv1`: lancia `virt-install` lì, e apri la console dalla workstation con
    [virt-manager collegato a `srv1`](../../ambiente/vm-linux.md#livelli-2-e-3-le-vm-sul-server).

=== "Livello 3 · Mini data center"

    Come nel livello 2, su `nodo1`, che ospita le VM del laboratorio.

Scarica `FreeBSD-15.1-RELEASE-amd64-disc1.iso` dalla [pagina ufficiale](https://www.freebsd.org/where/) (sezione
*Installer Images*, architettura `amd64`), poi:

```bash
virt-install --name freebsd --memory 2048 --vcpus 2 \
  --disk size=20 --disk size=2 --disk size=2 \
  --cdrom ~/Scaricati/FreeBSD-15.1-RELEASE-amd64-disc1.iso \
  --osinfo freebsd14.0 \
  --network network=academy,model=virtio
```

I due dischi da 2 GB in più ti servono per ZFS nella [pagina 3](03-zfs.md). Se `virt-install --osinfo list | grep freebsd`
mostra anche `freebsd15.0`, usa quello: è solo un'indicazione per scegliere le impostazioni predefinite della VM.

### L'installazione

Si apre l'installatore testuale, **bsdinstall**. Le scelte per il laboratorio:

| Schermata | Scelta |
|---|---|
| Keymap | *Italian* |
| Hostname | `freebsd.lab` |
| Metodo di installazione | **Distribution Sets**, il metodo classico usato in questa appendice |
| Componenti | quelli proposti |
| Rete | `vtnet0`, IPv4 con DHCP, niente IPv6 |
| Partizionamento | **Auto (ZFS)**, tipo *stripe* sul solo disco `vtbd0` (i due dischi piccoli restano liberi) |
| Password di root | una password robusta |
| Fuso orario | *Europe*, *Italy* |
| Servizi all'avvio | `sshd` e `ntpd` |
| Utente | crea il tuo utente e, alla domanda *Invite into other groups?*, scrivi `wheel` |

Al riavvio, trova l'indirizzo della VM sull'host (`virsh net-dhcp-leases academy`) ed entra con SSH dalla
workstation. Nei livelli 2 e 3 aggiungi a `~/.ssh/config` un blocco `Host freebsd` con quell'indirizzo e
`ProxyJump srv1` (o `nodo1`), come quelli generati da `academy-vm ssh-config`.

Il tuo utente è nel gruppo `wheel`, che su FreeBSD è l'unico a poter usare `su`. Per `sudo`:

```sh
su -                                   # diventi root con la sua password
pkg install -y sudo
visudo                                 # togli il commento alla riga: %wheel ALL=(ALL:ALL) ALL
```

## I pacchetti: pkg

`pkg` installa pacchetti già compilati dal progetto, come `apt` e `dnf`. La prima volta si installa da solo.

```sh
pkg search nginx                       # cerca
sudo pkg install -y nginx              # installa
pkg info nginx                         # dettagli di un pacchetto installato
pkg info -l nginx                      # i suoi file: tutti sotto /usr/local
pkg which /usr/local/sbin/nginx        # a quale pacchetto appartiene un file
sudo pkg update && sudo pkg upgrade    # aggiorna l'elenco e i pacchetti
sudo pkg audit -F                      # pacchetti installati con vulnerabilità note
sudo pkg autoremove                    # toglie le dipendenze non più usate
```

Dopo l'installazione, leggi il messaggio che `pkg` stampa (lo ritrovi con `pkg info -D nginx`): spesso dice come
attivare il servizio o dove sta la configurazione.

## I ports: compilare dai sorgenti

I **ports** sono la raccolta delle "ricette" con cui si costruiscono i pacchetti: oltre 30.000 programmi. Servono
quando vuoi un'opzione di compilazione diversa da quella predefinita.

```sh
sudo pkg install -y git
sudo git clone --depth 1 https://git.FreeBSD.org/ports.git /usr/ports
cd /usr/ports/sysutils/tmux
sudo make config                       # le opzioni di compilazione, in un menu
sudo make install clean                # scarica i sorgenti, compila, installa, ripulisce
```

!!! warning "pkg e ports insieme"
    Mescolare pacchetti e ports funziona, ma un `pkg upgrade` può sostituire un programma che avevi compilato con
    opzioni tue. Sui server di produzione si sceglie uno dei due; chi ha bisogno di opzioni personalizzate costruisce
    i propri pacchetti con **poudriere** e li installa con `pkg`.

## I servizi: rc.conf

All'avvio FreeBSD esegue gli script in `/etc/rc.d` (sistema di base) e `/usr/local/etc/rc.d` (pacchetti). Quali
servizi partono lo decide un solo file, **`/etc/rc.conf`**, fatto di variabili:

```sh title="/etc/rc.conf (esempio)"
hostname="freebsd.lab"
ifconfig_vtnet0="DHCP"
sshd_enable="YES"
ntpd_enable="YES"
nginx_enable="YES"
zfs_enable="YES"
```

Si può modificare con un editor, ma è meglio usare **`sysrc`**, che controlla la sintassi:

```sh
sudo sysrc nginx_enable=YES            # attiva all'avvio
sudo service nginx start               # avvia adesso
service nginx status
service -e                             # i servizi attivati all'avvio
sysrc -a                               # tutte le variabili impostate
sudo sysrc -x nginx_enable             # toglie la variabile
```

Rispetto a systemd ([lezione 1.4](../../fase-1/04-processi-systemd.md)) è tutto più semplice: niente dipendenze dichiarate né
riavvio automatico dei servizi che cadono. Gli script si leggono come normali script `sh`.

Anche la rete si configura in `rc.conf`. Un indirizzo statico:

```sh
sudo sysrc ifconfig_vtnet0="inet 192.168.56.80/24"
sudo sysrc defaultrouter="192.168.56.1"
sudo service netif restart && sudo service routing restart
```

## Gli aggiornamenti

Il sistema di base e i pacchetti si aggiornano separatamente.

```sh
sudo freebsd-update fetch              # scarica le correzioni di sicurezza del sistema di base
sudo freebsd-update install            # le installa
sudo pkg upgrade                       # aggiorna i pacchetti
```

Per passare a una nuova versione (per esempio dalla 15.0 alla 15.1) si usa `freebsd-update upgrade -r 15.1-RELEASE`,
seguendo le istruzioni dell'annuncio di rilascio, che indicano anche i riavvii necessari.

!!! info "pkgbase: il sistema di base in pacchetti"
    Dalla 15.0 l'installatore propone anche un secondo metodo, **pkgbase**: il sistema di base è diviso in
    pacchetti e si aggiorna con `pkg`, come il resto del software. Per sapere quale metodo usa un sistema:
    `pkg which /usr/bin/uname`. Se risponde che il file non viene da un pacchetto, il sistema usa i
    *distribution sets* e si aggiorna con `freebsd-update`.

!!! esercizio "Esercizio · Il primo servizio"
    Nella VM installa `nginx`, attivalo all'avvio, avvialo e verifica dalla workstation (o da `srv1`, `nodo1`) che
    `curl -sI http://IP-FREEBSD/` risponda `200`. Poi trova il file di configurazione di Nginx e il suo log degli
    accessi. Riavvia la VM e controlla che Nginx riparta.

??? success "Soluzione"
    ```sh
    sudo pkg install -y nginx
    sudo sysrc nginx_enable=YES
    sudo service nginx start
    sockstat -4l | grep nginx                    # ascolta sulla porta 80
    ls /usr/local/etc/nginx/nginx.conf
    tail /var/log/nginx/access.log
    sudo shutdown -r now
    ```
    La configurazione è in `/usr/local/etc/nginx/`, come tutto il software aggiunto. Il log degli accessi, invece,
    sta in `/var/log/nginx/`: i dati che cambiano restano in `/var` anche su FreeBSD.

## Fonti

- [FreeBSD Handbook: Installing FreeBSD](https://docs.freebsd.org/en/books/handbook/bsdinstall/)
- [FreeBSD Handbook: Installing Applications, Packages and Ports](https://docs.freebsd.org/en/books/handbook/ports/)
- [FreeBSD Handbook: The FreeBSD Booting Process e rc](https://docs.freebsd.org/en/books/handbook/boot/)
- [FreeBSD Handbook: Updating and Upgrading FreeBSD](https://docs.freebsd.org/en/books/handbook/cutting-edge/)
- [FreeBSD: dove scaricare](https://www.freebsd.org/where/)
- [rc.conf(5)](https://man.freebsd.org/cgi/man.cgi?query=rc.conf&sektion=5), [pkg(8)](https://man.freebsd.org/cgi/man.cgi?query=pkg&sektion=8)
