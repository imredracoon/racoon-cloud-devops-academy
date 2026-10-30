# 5. pf e bhyve

!!! obiettivi "Obiettivi"
    - Leggere e scrivere un `pf.conf` e confrontarlo con nftables.
    - Attivare pf, controllare le regole e guardare i pacchetti bloccati.
    - Sapere cos'è bhyve e come si usa.

## pf, il firewall

**pf** (*packet filter*) è nato su OpenBSD ed è arrivato su FreeBSD: è il firewall di OPNsense e pfSense. Fa le
stesse cose di nftables ([lezione 1.11](../../fase-1/11-firewall.md)), con stato, NAT e tabelle di indirizzi, ma con una regola
importante: **vince l'ultima regola che corrisponde**, non la prima. Per fermarsi alla prima si aggiunge `quick`.

```text title="/etc/pf.conf"
# Firewall della VM FreeBSD: SSH e HTTP in ingresso, tutto in uscita.
ext_if = "vtnet0"
servizi_tcp = "{ ssh, http }"

set skip on lo0                          # nessun filtro sulla loopback

block in log all                         # politica predefinita: blocca e registra
pass out all keep state                  # tutto in uscita, con lo stato per le risposte
pass in on $ext_if proto tcp to port $servizi_tcp keep state
pass in on $ext_if inet proto icmp icmp-type echoreq
```

Il confronto con il firewall della lezione 1.11:

| nftables | pf |
|---|---|
| `policy drop` nella catena `input` | `block in all` come prima regola |
| `ct state established,related accept` | `keep state` sulle regole `pass` (predefinito) |
| vince la prima regola che corrisponde | vince l'**ultima**, salvo `quick` |
| `define` e insiemi `{ 22, 80 }` | macro (`ext_if = ...`) e liste `{ ssh, http }` |
| `log prefix` e il journal | `log` e l'interfaccia `pflog0`, da leggere con `tcpdump` |
| `nft -c -f` per controllare | `pfctl -nf` per controllare |

### Attivarlo

!!! warning "Fallo dalla console"
    All'avvio di pf le connessioni già aperte possono bloccarsi, compresa la tua sessione SSH. Apri la console della
    VM (virt-manager o `virsh console freebsd`) prima di cominciare.

```sh
sudo pfctl -nf /etc/pf.conf              # controlla la sintassi, senza applicare
sudo sysrc pf_enable=YES pflog_enable=YES
sudo service pf start && sudo service pflog start
sudo pfctl -sr                           # le regole attive
sudo pfctl -ss                           # gli stati: le connessioni in corso
sudo pfctl -si                           # statistiche
sudo tcpdump -n -e -ttt -i pflog0        # i pacchetti bloccati, in tempo reale
```

Dopo una modifica a `pf.conf` basta `sudo pfctl -f /etc/pf.conf`.

!!! esercizio "Esercizio · Il firewall alla prova"
    Con pf attivo e Nginx in ascolto, avvia sulla VM un secondo server di prova, `python3 -m http.server 9000`
    (installa prima `python3`). Dall'host verifica che le porte 22 e 80 rispondano e la 9000 no, e trova il
    pacchetto bloccato su `pflog0`. Che errore ottiene il client sulla 9000, e perché è lo stesso di nftables con
    `drop`? Poi cambia la regola `block` in `block return` e riprova.

??? success "Soluzione"
    ```sh
    # sulla VM
    sudo pkg install -y python3
    python3 -m http.server 9000 &
    sudo tcpdump -n -e -ttt -i pflog0
    # sull'host
    curl -sI http://IP-FREEBSD/ | head -1      # 200
    curl -m 5 http://IP-FREEBSD:9000/          # timeout
    ```
    `block` scarta in silenzio, come `drop` in nftables: il client aspetta fino al timeout. Con `block return` pf
    risponde con un RST (TCP) o un ICMP *unreachable*, e il client riceve subito *Connection refused*. La lezione
    [1.11](../../fase-1/11-firewall.md#laboratorio) spiega quando conviene l'uno o l'altro.

Il NAT si scrive nello stesso file, con regole come
`nat on $ext_if from 10.30.30.0/24 to any -> ($ext_if)`: è ciò che OPNsense genera per te nella
[pagina 6](06-opnsense.md).

## bhyve, l'hypervisor

**bhyve** è l'hypervisor di FreeBSD, l'equivalente di KVM ([macchine virtuali](../../ambiente/vm-linux.md)): usa le
estensioni di virtualizzazione del processore e fa girare FreeBSD, Linux e Windows. Il comando `bhyve` è di basso
livello, come `qemu` senza libvirt; nella pratica si usa un gestore:

- **vm-bhyve** (`pkg install vm-bhyve`): VM da riga di comando, con modelli per i sistemi più comuni;
- i NAS e gli apparecchi basati su FreeBSD che offrono VM le creano con bhyve dalla loro interfaccia web.

!!! info "bhyve dentro una VM"
    Per provare bhyve nella VM del laboratorio serve la **virtualizzazione annidata**: la VM deve vedere le
    estensioni del processore. Ricreala con l'opzione `--cpu host-passthrough` di `virt-install`. Funziona, ma è
    lenta: bhyve dà il meglio su una macchina fisica, per esempio un mini PC del
    [livello 3](../../ambiente/livelli.md) dedicato a FreeBSD.

## Fonti

- [FreeBSD Handbook: Firewalls (PF)](https://docs.freebsd.org/en/books/handbook/firewalls/)
- [pf.conf(5)](https://man.freebsd.org/cgi/man.cgi?query=pf.conf&sektion=5), [pfctl(8)](https://man.freebsd.org/cgi/man.cgi?query=pfctl&sektion=8)
- [FreeBSD Handbook: Virtualization (bhyve)](https://docs.freebsd.org/en/books/handbook/virtualization/)
- [vm-bhyve](https://github.com/churchers/vm-bhyve)
