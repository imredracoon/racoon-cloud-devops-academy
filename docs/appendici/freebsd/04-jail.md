# 4. Le jail

!!! obiettivi "Obiettivi"
    - Confrontare le jail con i container Linux.
    - Creare una jail con il suo indirizzo IP e un server web dentro.
    - Gestirla con `jail.conf`, `service jail` e `jexec`, e proteggerla con gli snapshot ZFS.

## I container di FreeBSD

Le **jail** esistono dal 2000: sono il primo sistema di isolamento diffuso, più vecchio di Docker di oltre dieci
anni. Una jail è un pezzo di sistema con una propria radice del filesystem, i propri utenti (anche un root, che però
non può uscire dalla jail), i propri processi e un indirizzo IP. Il kernel è quello dell'host, condiviso, come per i
container ([lezione 3.1](../../fase-3/01-concetti.md)).

| | Container Linux (Docker, Podman) | Jail |
|---|---|---|
| Isolamento | namespace e cgroup | meccanismo unico del kernel, `jail(2)` |
| Contenuto tipico | un processo, da un'immagine stratificata | un intero sistema FreeBSD, gestito come una piccola macchina |
| Immagini e registry | standard OCI, Docker Hub, GHCR | nessuno standard: la radice è una cartella (spesso un dataset ZFS) |
| Rete | reti virtuali, porte pubblicate | un IP dell'host, oppure uno stack di rete proprio (*VNET*) |
| Orchestrazione | Kubernetes | strumenti come Bastille o AppJail, niente di paragonabile a Kubernetes |

Le jail somigliano più a VM leggerissime che a container "da un processo". Dalla 14.2 il progetto pubblica anche
immagini di container **OCI**, da usare con Podman e il runtime `ocijail`: una strada ancora poco diffusa.

## Una jail con un server web

Nella VM FreeBSD della [pagina 2](02-installazione.md). La jail avrà l'indirizzo `192.168.56.90`, fuori dall'intervallo DHCP
della rete `academy`, aggiunto alla scheda `vtnet0` dell'host.

**1. La radice, in un dataset ZFS:**

```sh
sudo zfs create -o mountpoint=/jail zroot/jail
sudo zfs create zroot/jail/web
fetch https://download.freebsd.org/releases/amd64/15.1-RELEASE/base.txz -o /tmp/base.txz
sudo tar -xf /tmp/base.txz -C /jail/web --unlink      # un sistema di base completo
sudo cp /etc/resolv.conf /etc/localtime /jail/web/etc/
```

**2. La configurazione:**

```text title="/etc/jail.conf"
# Impostazioni comuni a tutte le jail
exec.start = "/bin/sh /etc/rc";
exec.stop = "/bin/sh /etc/rc.shutdown";
exec.clean;
mount.devfs;
path = "/jail/${name}";
host.hostname = "${name}.lab";

web {
    ip4.addr = "vtnet0|192.168.56.90/24";
}
```

**3. Avvio e software:**

```sh
sudo sysrc jail_enable=YES jail_list="web"
sudo service jail start web
jls                                                    # le jail attive
sudo jexec web env ASSUME_ALWAYS_YES=yes pkg install nginx
sudo jexec web sysrc nginx_enable=YES
sudo jexec web service nginx start
```

`jexec web COMANDO` esegue un comando dentro la jail, come `docker exec`; `sudo jexec web sh` apre una shell.

**4. Uno snapshot della jail pulita**, per tornare indietro in un attimo:

```sh
sudo zfs snapshot zroot/jail/web@pulita
```

## Collaudo

| Prova | Risultato atteso |
|---|---|
| `curl -sI http://192.168.56.90/` dall'host della VM (workstation, `srv1` o `nodo1`) | `200` da Nginx: la jail ha il suo indirizzo |
| `sudo jexec web ps aux` | solo i processi della jail |
| `ps -ax -o jid,pid,command` sull'host FreeBSD | anche i processi della jail, con il suo numero nella colonna `JID` |
| `sudo jexec web ifconfig vtnet0` | solo l'indirizzo `192.168.56.90` |
| `sudo jexec web sh -c 'ls /jail'` | errore: dalla jail non si vede il resto del filesystem |
| riavvio della VM | la jail e Nginx ripartono da soli |
| `sudo service jail stop web`, `sudo zfs rollback zroot/jail/web@pulita`, `sudo service jail start web` | la jail torna come appena creata, senza Nginx |

!!! tip "Gli strumenti di gestione"
    A mano si capisce come funziona; con molte jail si usa un gestore. **Bastille** (`pkg install bastille`) crea,
    aggiorna e clona jail con un comando e usa modelli (*template*) simili a un Dockerfile.

## Fonti

- [FreeBSD Handbook: Jails and Containers](https://docs.freebsd.org/en/books/handbook/jails/)
- [jail.conf(5)](https://man.freebsd.org/cgi/man.cgi?query=jail.conf&sektion=5), [jexec(8)](https://man.freebsd.org/cgi/man.cgi?query=jexec&sektion=8)
- [BastilleBSD](https://bastillebsd.org/)
