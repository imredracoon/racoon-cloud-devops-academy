# 1.11 Firewall con nftables

!!! obiettivi "Obiettivi"
    - Capire come funziona un firewall con stato e il principio del "nega tutto per default".
    - Scrivere un set di regole nftables per un server web e renderlo permanente.
    - Configurare NAT e port forwarding con nftables.
    - Usare `ufw` per i casi semplici.

## Cos'è un firewall

Un **firewall** decide quali pacchetti possono entrare, uscire o attraversare una macchina. Nel kernel Linux questo
compito è svolto da **netfilter**; lo strumento moderno per configurarlo è **nftables** (comando `nft`), che ha
sostituito `iptables`. Su Ubuntu, Debian e RHEL nftables è il default; anche `iptables`, se presente, è spesso solo
un'interfaccia sopra nftables.

### Firewall con stato

netfilter tiene traccia delle connessioni (*connection tracking*). Ogni pacchetto ha uno **stato**:

| Stato | Significato |
|---|---|
| `new` | Primo pacchetto di una nuova connessione (per TCP, il SYN). |
| `established` | Appartiene a una connessione già accettata. |
| `related` | Collegato a una connessione esistente (per esempio un errore ICMP relativo). |
| `invalid` | Non corrisponde a nulla di sensato: si scarta. |

Grazie allo stato basta una regola per accettare **tutte le risposte** al traffico che hai già permesso.

### La politica giusta

La regola d'oro: **nega tutto ciò che non è esplicitamente permesso** (*default deny*). Per un server web in
ingresso si permettono solo SSH, HTTP e HTTPS.

## La struttura di nftables

- Una **tabella** (*table*) raggruppa le catene; ha una famiglia: `inet` vale per IPv4 e IPv6 insieme.
- Una **catena** (*chain*) di tipo base si aggancia a un punto del percorso dei pacchetti (*hook*):
  `input` (pacchetti per questa macchina), `output` (generati da questa macchina), `forward` (in transito),
  `prerouting` e `postrouting` (usati per il NAT). Ha una **politica** predefinita: `accept` o `drop`.
- Le **regole** si valutano in ordine; la prima che decide (`accept`, `drop`, `reject`) chiude la valutazione.

```bash
sudo nft list ruleset            # mostra tutte le regole attive
```

## Un firewall per un server web

```nft title="/etc/nftables.conf"
#!/usr/sbin/nft -f

flush ruleset

table inet filtro {
    chain input {
        type filter hook input priority filter; policy drop;

        ct state established,related accept   # risposte a connessioni già permesse
        ct state invalid drop
        iif "lo" accept                        # traffico locale sulla loopback

        # ICMP: utile per la diagnostica e indispensabile per IPv6
        ip protocol icmp icmp type { echo-request, destination-unreachable, time-exceeded } accept
        meta l4proto ipv6-icmp accept

        tcp dport 22 ct state new limit rate 10/minute accept   # SSH, con limite ai nuovi tentativi
        tcp dport { 80, 443 } accept                            # web

        counter log prefix "nft-drop: " level info drop        # registra e scarta il resto
    }

    chain forward {
        type filter hook forward priority filter; policy drop;
    }

    chain output {
        type filter hook output priority filter; policy accept;
    }
}
```

Le righe che cominciano con `#` sono commenti. L'ultima regola conta i pacchetti scartati e li registra nel
journal (`journalctl -k | grep nft-drop`).

```bash
sudo nft -c -f /etc/nftables.conf      # controlla la sintassi senza applicare (-c = check)
sudo nft -f /etc/nftables.conf         # applica
sudo systemctl enable nftables         # carica /etc/nftables.conf a ogni avvio
sudo nft list ruleset
```

!!! danger "Non chiuderti fuori, di nuovo"
    Applicare un firewall da remoto è il secondo modo classico per perdere l'accesso a un server. Prima di
    applicare: controlla che la regola per SSH ci sia, tieni aperta la sessione corrente e prova un nuovo accesso.
    Un trucco dei professionisti: pianifica un "ripristino automatico" prima di applicare, per esempio
    `sudo systemd-run --on-active=5m nft flush ruleset`, e cancellalo (`systemctl stop` sul timer creato) solo
    quando hai verificato che tutto funziona.

### Modificare le regole al volo

```bash
sudo nft add rule inet filtro input tcp dport 8080 accept            # aggiunge in fondo alla catena
sudo nft -a list chain inet filtro input                              # mostra le regole con il loro handle
sudo nft delete rule inet filtro input handle 15                     # cancella una regola per handle
```

Le modifiche al volo **non sono permanenti**: riportale in `/etc/nftables.conf`.

## NAT con nftables

Per far uscire su Internet una rete privata attraverso questa macchina (che deve avere `ip_forward` attivo,
[lezione 1.8](08-dns-routing.md)):

```nft
table ip nat {
    chain postrouting {
        type nat hook postrouting priority srcnat; policy accept;
        oifname "ens3" ip saddr 10.10.0.0/24 masquerade      # SNAT verso l'IP dell'interfaccia di uscita
    }
    chain prerouting {
        type nat hook prerouting priority dstnat; policy accept;
        iifname "ens3" tcp dport 8080 dnat to 10.10.0.20:80   # port forwarding (DNAT)
    }
}
```

Ricorda che la catena `forward` della tabella filtro deve **permettere** il traffico in transito, altrimenti il NAT
è inutile.

## I gestori delle distribuzioni: ufw e firewalld

Scrivere regole nftables a mano ti fa capire cosa succede davvero. Nel lavoro quotidiano però ogni famiglia di
distribuzioni ha il suo **gestore** di firewall, che genera le regole nftables per te.

=== "Debian/Ubuntu"

    **ufw** (*Uncomplicated Firewall*) è pensato per i casi semplici:

    ```bash
    sudo ufw default deny incoming
    sudo ufw default allow outgoing
    sudo ufw allow OpenSSH          # o: sudo ufw allow 22/tcp
    sudo ufw allow 80,443/tcp
    sudo ufw limit 22/tcp           # limita i tentativi di connessione ripetuti
    sudo ufw enable
    sudo ufw status numbered
    sudo ufw delete 3               # cancella la regola numero 3
    ```

=== "RHEL (Rocky/AlmaLinux)"

    **firewalld** è attivo per default su RHEL ed è negli obiettivi della RHCSA. Ragiona per **zone**: ogni
    interfaccia appartiene a una zona (di solito `public`) e ogni zona ha i suoi servizi e porte permessi. Tutto ciò
    che non è permesso viene rifiutato.

    ```bash
    sudo firewall-cmd --state
    sudo firewall-cmd --get-active-zones
    sudo firewall-cmd --list-all                    # cosa permette la zona predefinita
    sudo firewall-cmd --add-service=http            # solo runtime: sparisce al reload/riavvio
    sudo firewall-cmd --permanent --add-service=http --add-service=https   # permanente: non attivo finché...
    sudo firewall-cmd --reload                      # ...non ricarichi
    sudo firewall-cmd --permanent --add-port=8080/tcp
    sudo firewall-cmd --permanent --remove-service=cockpit
    sudo firewall-cmd --get-services                # servizi predefiniti (definiti in /usr/lib/firewalld/services)
    sudo firewall-cmd --runtime-to-permanent        # salva come permanente la configurazione runtime attuale
    ```

    Regole più specifiche si scrivono con le **rich rule**:

    ```bash
    sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.56.11" service name="ssh" accept'
    ```

    NAT e port forwarding:

    ```bash
    sudo firewall-cmd --permanent --add-masquerade
    sudo firewall-cmd --permanent --add-forward-port=port=8080:proto=tcp:toport=80:toaddr=10.10.0.20
    sudo firewall-cmd --reload
    ```

!!! warning "Runtime e permanente"
    Con firewalld un comando senza `--permanent` ha effetto subito ma si perde al reload; con `--permanent` viene
    salvato ma non è attivo fino al `--reload`. Agli esami Red Hat le configurazioni devono sopravvivere al riavvio:
    usa sempre `--permanent` seguito da `--reload`, e verifica con `--list-all`.

**Non usare insieme** un gestore (ufw o firewalld) e un `/etc/nftables.conf` scritto a mano: si sovrascriverebbero
a vicenda. Su un server scegli un solo strumento.

!!! note "Nel cloud"
    I fornitori cloud offrono firewall di rete gestiti: i **security group** di AWS e le **regole firewall VPC** di
    Google Cloud ([Fase 4](../fase-4/index.md) e 5). Si usano **insieme** al firewall dell'host, secondo il principio della *difesa in
    profondità*.

## Laboratorio

!!! laboratorio "Laboratorio · Proteggi lab · 1 ora"
    **Dove:** VM `lab` e VM `client`. Fai prima uno snapshot di `lab`.

    1. Su `lab` avvia due server di prova: `python3 -m http.server 80` (con sudo) e `python3 -m http.server 9000`.
    2. Da `client` verifica che entrambe le porte rispondano.
    3. Applica il firewall dell'esempio.
    4. Da `client` verifica: SSH e porta 80 funzionano, la porta 9000 no. Che tipo di errore ottieni sulla 9000, e
       perché è diverso da quello della [lezione 1.9](09-diagnostica.md)?
    5. Trova nel journal di `lab` la riga che registra il pacchetto scartato.
    6. Rendi il firewall permanente e verifica dopo un riavvio.

    **Percorso RHEL:** su `rhel-lab` ottieni lo stesso risultato con **firewalld** (lascia attivo firewalld e non
    usare `/etc/nftables.conf`): solo SSH e HTTP nella zona `public`, porta 9000 chiusa. Nota che firewalld
    **rifiuta** (reject) invece di scartare: che errore ottiene `rhel-client` sulla porta 9000?

??? success "Soluzione"
    ```bash
    # su lab
    sudo python3 -m http.server 80 &
    python3 -m http.server 9000 &

    # su client
    curl -sI http://IP-LAB/ | head -1          # HTTP/1.0 200 OK
    curl -sI http://IP-LAB:9000/ | head -1     # HTTP/1.0 200 OK

    # su lab
    sudo nano /etc/nftables.conf               # incolla le regole
    sudo nft -c -f /etc/nftables.conf && sudo nft -f /etc/nftables.conf

    # su client
    ssh ubuntu@IP-LAB true && echo "ssh ok"    # se hai autorizzato una chiave nella lezione 1.10
    curl -sI http://IP-LAB/ | head -1          # funziona
    curl -m 5 http://IP-LAB:9000/              # Connection timed out dopo 5 secondi

    # su lab
    journalctl -k --since "5 min ago" | grep nft-drop    # ... DPT=9000 ...
    sudo systemctl enable nftables
    sudo reboot
    # dopo il riavvio
    sudo nft list ruleset | head
    ```
    **Spiegazione.** La regola finale fa **`drop`**: il pacchetto viene scartato senza risposta, quindi il client
    aspetta fino al **timeout**. Nella [lezione 1.9](09-diagnostica.md) ricevevi *Connection refused* perché il pacchetto arrivava al
    kernel e nessuno ascoltava. Se usassi `reject` invece di `drop`, il client riceverebbe subito un rifiuto.
    `drop` rivela meno informazioni all'esterno; `reject` rende più rapida la diagnosi nelle reti interne.

    (Se `ssh lab` smette di funzionare, entra con `virsh console lab` e controlla la regola per la porta 22.)

??? success "Soluzione: percorso RHEL"
    ```bash
    # su rhel-lab
    sudo firewall-cmd --permanent --add-service=http
    sudo firewall-cmd --reload
    sudo firewall-cmd --list-all        # services: cockpit dhcpv6-client http ssh
    sudo python3 -m http.server 80 &
    python3 -m http.server 9000 &

    # su rhel-client
    curl -sI http://192.168.56.10/ | head -1     # 200 OK
    curl -m 5 http://192.168.56.10:9000/         # errore immediato: "No route to host"
    ```
    **Spiegazione.** firewalld risponde ai pacchetti non permessi con un messaggio ICMP *host-prohibited*, che il
    client mostra come *No route to host*: è un rifiuto (`reject`), quindi arriva subito, a differenza del timeout
    del `drop`. Il comportamento si può cambiare impostando il *target* della zona a `DROP`.

## Nei tre livelli

!!! warning "Un firewall su un host con le VM"
    Su `srv1` e `nodo1` libvirt crea regole proprie per la rete delle VM. Per questo, sulle macchine fisiche, il
    tuo `/etc/nftables.conf` **non** inizia con `flush ruleset` (cancellerebbe anche le regole di libvirt) ma
    ricrea solo la propria tabella, e lascia passare il traffico dell'interfaccia `virbr56`.

=== "Livello 1 · Un PC"

    Il laboratorio con `lab` e `client` è tutto ciò che serve. Il firewall della workstation lascialo com'è.

=== "Livello 2 · Workstation + server"

    Fai il laboratorio nelle VM, poi proteggi `srv1`.

    === "srv1 con Ubuntu Server"

        ```nft title="/etc/nftables.conf"
        #!/usr/sbin/nft -f
        # Firewall di srv1. Niente "flush ruleset": libvirt ha le sue regole per le VM.
        table inet homelab
        delete table inet homelab

        table inet homelab {
            chain input {
                type filter hook input priority filter; policy drop;
                ct state established,related accept
                ct state invalid drop
                iif "lo" accept
                meta l4proto { icmp, ipv6-icmp } accept
                iifname "virbr56" accept                                # le VM chiedono DHCP e DNS a srv1
                tcp dport 22 accept
                ip saddr 192.168.1.0/24 tcp dport 2049 accept           # NFS verso la workstation (lezione 1.6)
                limit rate 5/minute log prefix "srv1-drop: "
            }
        }
        ```

        ```bash
        sudo nft -c -f /etc/nftables.conf && sudo nft -f /etc/nftables.conf
        sudo systemctl enable nftables
        ```

    === "srv1 con AlmaLinux"

        firewalld è già attivo. libvirt mette da solo `virbr56` nella zona `libvirt`, quindi le VM non vanno toccate:

        ```bash
        sudo firewall-cmd --get-active-zones                      # public (la scheda di rete) e libvirt (virbr56)
        sudo firewall-cmd --permanent --remove-service=cockpit
        sudo firewall-cmd --permanent --add-service=nfs           # se hai fatto NFS nella lezione 1.6
        sudo firewall-cmd --reload && sudo firewall-cmd --list-all
        ```

    **Facoltativo, con la seconda rete:** consenti SSH **solo dalla rete di gestione**, come nei data center veri.
    Con nftables cambia la riga di SSH in `iifname "gest0" tcp dport 22 accept`; con firewalld metti la connessione
    nella zona `internal` (`sudo nmcli connection modify gest0 connection.zone internal`) e togli `ssh` dalla zona
    `public`. Poi, in `~/.ssh/config` della workstation, cambia l'`HostName` di `srv1` in `10.20.20.1`.

    **Collaudo** (dalla workstation; su `srv1` avvia prima `python3 -m http.server 9000`)

    | Prova | Risultato atteso |
    |---|---|
    | `ssh srv1 hostname` | `srv1` |
    | `curl -m 5 http://IP-SRV1:9000/` | timeout con nftables, `No route to host` con firewalld |
    | `ssh lab ping -c 2 1.1.1.1` | risponde: le VM escono ancora su Internet |
    | `ssh lab hostname` | `lab` |
    | `ssh srv1 journalctl -k --since -10min \| grep srv1-drop` (Ubuntu) | la riga del pacchetto per la porta 9000 |
    | riavvio di `srv1`, poi le prove 1, 3 e 4 | tutto come prima |
    | (con la rete di gestione) `ssh -o HostName=IP-SRV1-DI-CASA srv1` | timeout: SSH dalla rete di casa è chiuso |

=== "Livello 3 · Mini data center"

    Fai il laboratorio **su `nodo2`** (al posto di `lab`, firewall nftables) e **su `nodo3`** (al posto di
    `rhel-lab`, firewalld), usando l'altro nodo come client. Poi la parte più importante: il firewall del gateway.
    `nodo1` deve lasciar uscire il laboratorio, non far entrare nulla che non sia stato pubblicato, e
    **pubblicare** con il port forwarding (DNAT) il servizio web di `nodo2`.

    ```nft title="/etc/nftables.conf su nodo1 (sostituisce la prima versione della lezione 1.8)"
    #!/usr/sbin/nft -f
    # Firewall di nodo1, gateway del laboratorio.
    # Niente "flush ruleset": cancellerebbe anche le regole che libvirt crea per le VM.
    table inet homelab
    delete table inet homelab

    define WAN = "eno1"        # la scheda integrata, verso il router di casa
    define LAB = "lab0"        # l'adattatore USB, verso lo switch
    define VM  = "virbr56"     # la rete delle VM di libvirt

    table inet homelab {
        chain input {
            type filter hook input priority filter; policy drop;
            ct state established,related accept
            ct state invalid drop
            iif "lo" accept
            meta l4proto { icmp, ipv6-icmp } accept
            iifname $VM accept                                     # DHCP e DNS per le VM
            iifname $LAB udp dport { 53, 67 } accept               # DNS e DHCP per i nodi
            iifname $LAB tcp dport { 53, 2049 } accept             # DNS e NFS (lezione 1.6)
            tcp dport 22 accept                                    # SSH da casa e dal laboratorio
            limit rate 5/minute log prefix "nodo1-in-drop: "
        }

        chain forward {
            type filter hook forward priority filter; policy drop;
            ct state established,related accept
            ct state invalid drop
            iifname $LAB oifname $WAN accept                       # il laboratorio esce
            iifname $WAN oifname $LAB ct status dnat accept        # da casa entra solo ciò che è pubblicato
            iifname $VM accept                                     # le VM: decidono le regole di libvirt
            oifname $VM accept
            limit rate 5/minute log prefix "nodo1-fwd-drop: "
        }

        chain prerouting {
            type nat hook prerouting priority dstnat; policy accept;
            iifname $WAN tcp dport 8080 dnat ip to 10.10.10.12:80   # pubblica il web di prova di nodo2
        }

        chain postrouting {
            type nat hook postrouting priority srcnat; policy accept;
            oifname $WAN ip saddr 10.10.10.0/24 masquerade
        }
    }
    ```

    ```bash
    sudo nft -c -f /etc/nftables.conf && sudo nft -f /etc/nftables.conf
    sudo nft list ruleset | less          # la tua tabella homelab e, separate, le regole di libvirt
    ```

    Tieni aperta una seconda sessione su `nodo1` mentre applichi le regole. Se perdi l'accesso, schermo e
    tastiera su `nodo1` e `sudo nft delete table inet homelab` ti riportano indietro (il NAT si ferma finché non
    correggi).

    **Collaudo** (su `nodo2` avvia prima `sudo python3 -m http.server 80`)

    | Prova | Risultato atteso |
    |---|---|
    | `curl -sI http://IP-NODO1:8080/` dalla workstation | `200`: il DNAT pubblica `nodo2` |
    | `curl -m 5 http://10.10.10.12/` dalla workstation | errore: la rete del laboratorio resta privata |
    | `nc -zvw3 IP-NODO1 53` dalla workstation | timeout: il DNS del laboratorio non è esposto a casa |
    | `ssh nodo2 getent hosts ubuntu.com` e `ssh nodo3 curl -sI https://almalinux.org` | rispondono: il laboratorio esce |
    | `ssh lab ping -c 2 1.1.1.1` | risponde: le VM di `nodo1` escono ancora |
    | `ssh nodo2 'curl -m 5 http://nodo3.lab:9000/'` con un server sulla 9000 di `nodo3` | `No route to host`: lo blocca firewalld di `nodo3` |
    | `ssh nodo1 journalctl -k --since -10min \| grep nodo1-` | i pacchetti scartati, con interfaccia di ingresso e porta |
    | riavvio di `nodo1`, poi le prove 1 e 4 | tutto come prima |

    Nel log del server di prova su `nodo2` le richieste arrivano dall'indirizzo della **workstation**, non da quello
    di `nodo1`: il DNAT cambia la destinazione, non la sorgente.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa significa `policy drop` su una catena input?","options":["Tutto ciò che non è esplicitamente accettato viene scartato","Tutto viene accettato","Le regole vengono cancellate","Il firewall è disattivato"],"answer":[0],"explanation":"È il principio default deny."},
 {"q":"A cosa serve la regola `ct state established,related accept`?","options":["A permettere le risposte alle connessioni già accettate","A bloccare SSH","A permettere tutto il traffico in uscita","A registrare i log"],"answer":[0],"explanation":"Il connection tracking permette di accettare il traffico di ritorno con una sola regola."},
 {"q":"Il client va in timeout invece di ricevere *Connection refused*. Cosa è più probabile?","options":["Nessun processo ascolta sulla porta","Un firewall scarta (drop) i pacchetti","Il DNS è sbagliato","Il server è sovraccarico di CPU"],"answer":[1],"explanation":"drop non invia risposta; un servizio assente provoca un RST immediato."},
 {"q":"Quale comando controlla la sintassi di un file nftables senza applicarlo?","options":["nft -c -f file","nft list ruleset","nft flush ruleset","nft -a file"],"answer":[0],"explanation":"-c esegue solo il controllo."},
 {"q":"Quale tipo di NAT permette a una rete privata di uscire su Internet con l'IP del router?","options":["DNAT","masquerade (SNAT)","Port forwarding","Hairpin"],"answer":[1],"explanation":"masquerade è una forma di SNAT che usa l'indirizzo dell'interfaccia di uscita."},
 {"q":"In quale catena si filtra il traffico che attraversa la macchina senza esserle destinato?","options":["input","output","forward","prerouting"],"answer":[2],"explanation":"forward riguarda i pacchetti in transito."}
]}
</script>

## Fonti

- [nftables wiki (progetto netfilter)](https://wiki.nftables.org/)
- [nft(8), man page](https://www.netfilter.org/projects/nftables/manpage.html)
- [Debian wiki: nftables](https://wiki.debian.org/nftables)
- [Ubuntu Server docs: Firewall (ufw)](https://documentation.ubuntu.com/server/how-to/security/firewalls/)
- [firewalld: documentazione ufficiale](https://firewalld.org/documentation/)
- [Red Hat: Using and configuring firewalld (RHEL 10)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/configuring_firewalls_and_packet_filters/index)
- [Red Hat: Getting started with nftables](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/configuring_firewalls_and_packet_filters/getting-started-with-nftables_firewall-packet-filters)
