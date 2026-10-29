# 1.9 Diagnostica di rete

!!! obiettivi "Obiettivi"
    - Verificare la raggiungibilità con `ping` e il percorso con `traceroute`/`tracepath`.
    - Scoprire quali servizi sono in ascolto e quali connessioni sono aperte con `ss`.
    - Provare porte e servizi con `nc` e `curl`.
    - Catturare e leggere il traffico con `tcpdump`.

## La cassetta degli attrezzi

| Domanda | Strumento |
|---|---|
| Ho un indirizzo e una rotta? | `ip a`, `ip route` |
| L'host risponde? | `ping` |
| Che strada fanno i pacchetti? | `traceroute`, `tracepath`, `mtr` |
| Il nome si risolve? | `dig`, `getent hosts` ([lezione 1.8](08-dns-routing.md)) |
| Chi ascolta su quale porta, sul mio host? | `ss -tulpn` |
| La porta remota è raggiungibile? | `nc -zv`, `curl` |
| Cosa passa davvero sul cavo? | `tcpdump` |

## ping

```bash
ping -c 4 1.1.1.1          # 4 pacchetti ICMP echo
ping -c 4 www.example.com  # verifica anche la risoluzione DNS
```

Leggi il tempo (`time=`) e le perdite (`packet loss`). Attenzione: **molti host e quasi tutti i servizi cloud
bloccano ICMP** per default. Un ping che non risponde non prova che l'host sia spento; un ping che risponde non
prova che il servizio funzioni.

## traceroute e tracepath

Mostrano i router attraversati fino alla destinazione, sfruttando il campo TTL dei pacchetti IP.

```bash
tracepath -n 1.1.1.1           # non richiede privilegi
sudo traceroute -n -T -p 443 www.example.com   # usa pacchetti TCP verso la porta 443
mtr -rwc 10 1.1.1.1            # ping e traceroute insieme, con statistiche (pacchetto mtr-tiny)
```

Gli asterischi `* * *` indicano un router che non risponde: è normale per molti router e non è necessariamente un
problema, se i passaggi successivi rispondono.

## ss: chi ascolta e chi è collegato

```bash
ss -tulpn          # porte in ascolto: t=TCP u=UDP l=listening p=processo n=numeri (senza nomi)
ss -tn             # connessioni TCP stabilite
ss -tn state established '( dport = :443 )'   # filtro: connessioni verso la porta 443
```

```console
$ sudo ss -tlpn
State  Recv-Q Send-Q Local Address:Port  Peer Address:Port Process
LISTEN 0      4096   127.0.0.53%lo:53         0.0.0.0:*     users:(("systemd-resolve",pid=512,...))
LISTEN 0      4096         0.0.0.0:22         0.0.0.0:*     users:(("sshd",pid=801,...))
LISTEN 0      511        127.0.0.1:8080       0.0.0.0:*     users:(("python3",pid=1203,...))
```

!!! tip "L'indirizzo di ascolto conta"
    - `0.0.0.0:22` (o `[::]:22`) ascolta su **tutte** le interfacce: raggiungibile dall'esterno.
    - `127.0.0.1:8080` ascolta solo sulla **loopback**: raggiungibile solo dalla macchina stessa.

    Un'applicazione "che funziona in locale ma non da fuori" ascolta spesso solo su `127.0.0.1`. È uno degli errori
    più frequenti in assoluto, anche con i container.

`netstat` è il predecessore di `ss`: lo troverai in vecchie guide, con opzioni simili (`netstat -tulpn`).

## nc e curl: provare un servizio

```bash
nc -zv 10.84.37.51 22          # la porta 22 è aperta? (z = solo prova, v = dettagli)
nc -zvw 3 10.84.37.51 80       # con timeout di 3 secondi
nc -l 9000                     # mette in ascolto sulla porta 9000 (su un host)
nc 10.84.37.50 9000            # si collega (dall'altro host) e ciò che scrivi arriva di là
```

`curl` è lo strumento per HTTP:

```bash
curl http://example.com                 # scarica la pagina
curl -I https://example.com             # solo le intestazioni della risposta
curl -v https://example.com             # dettagli: DNS, connessione, TLS, richiesta e risposta
curl -sS -o /dev/null -w "%{http_code} %{time_total}s\n" https://example.com   # codice e tempo
curl -H "Host: sito.lab" http://127.0.0.1/     # invia un'intestazione Host specifica (lezione 1.12)
curl --resolve sito.lab:443:10.84.37.50 https://sito.lab/   # forza l'IP senza toccare il DNS
```

## tcpdump: guardare i pacchetti

`tcpdump` cattura il traffico che passa su un'interfaccia. È lo strumento definitivo quando "non si capisce cosa
succede".

```bash
sudo tcpdump -i any -n port 53             # traffico DNS su tutte le interfacce
sudo tcpdump -i ens3 -n host 1.1.1.1       # traffico da/verso un host
sudo tcpdump -i ens3 -n 'tcp port 80 and (tcp[tcpflags] & tcp-syn != 0)'   # solo i SYN verso/da la porta 80
sudo tcpdump -i ens3 -n -c 20 -w cattura.pcap port 443   # salva 20 pacchetti su file
```

Il file `.pcap` si può aprire con **Wireshark** sul tuo computer per un'analisi grafica.

!!! warning "Riservatezza"
    Il traffico catturato può contenere dati personali e credenziali (se non cifrato). Cattura solo sulle tue
    macchine di laboratorio o con autorizzazione, e cancella le catture quando non servono più.

## Laboratorio

!!! laboratorio "Laboratorio · Due VM che si parlano · 1 ora"
    **Dove:** VM `lab` e VM `client` (vedi l'introduzione della fase per crearla).

    1. Trova gli IP delle due VM: sul tuo computer `./academy-vm elenco`, dentro ogni VM `ip -brief address`.
    2. Su `lab` avvia un server web di prova sulla porta 8000 **solo sulla loopback**:
       `python3 -m http.server 8000 --bind 127.0.0.1`.
    3. Da `client` prova a raggiungerlo con `curl` e con `nc`. Che errore ottieni?
    4. Su `lab`, in un secondo terminale, verifica con `ss` su quale indirizzo ascolta.
    5. Correggi il problema e riprova.
    6. Mentre fai una richiesta da `client`, cattura il traffico su `lab` con `tcpdump` e individua l'handshake TCP
       e la richiesta HTTP.

??? success "Soluzione"
    ```bash
    # su client (sostituisci con l'IP di lab)
    curl http://10.84.37.50:8000/        # curl: (7) Failed to connect ... Connection refused
    nc -zv 10.84.37.50 8000              # Connection refused

    # su lab
    ss -tlpn | grep 8000                 # 127.0.0.1:8000

    # su lab: ferma il server (Ctrl+C) e riavvialo su tutte le interfacce
    python3 -m http.server 8000 --bind 0.0.0.0

    # su lab, altro terminale
    sudo tcpdump -i any -n -A port 8000

    # su client
    curl http://10.84.37.50:8000/
    ```
    Nella cattura vedrai i pacchetti con flag `[S]` (SYN), `[S.]` (SYN-ACK), `[.]` (ACK), poi `[P.]` con il testo
    `GET / HTTP/1.1` (l'opzione `-A` mostra il contenuto in ASCII) e la risposta `HTTP/1.0 200 OK`.

    **Spiegazione.** *Connection refused* arriva perché sull'interfaccia esterna nessuno ascolta sulla porta 8000:
    il kernel di `lab` risponde con un RST. Ascoltando su `0.0.0.0` il servizio diventa raggiungibile. HTTP non è
    cifrato, per questo puoi leggere la richiesta: con HTTPS vedresti solo dati illeggibili.

## Nei tre livelli

Il laboratorio si fa nelle VM in tutti i livelli. In più, una prova che rende visibile il **NAT** della
[lezione 1.8](08-dns-routing.md#nat): catturi lo stesso `ping` sulle due schede del router, prima e dopo la
traduzione degli indirizzi. Apri due terminali sul router e lancia il `ping` dalla macchina indicata.

=== "Livello 1 · Un PC"

    | Router | Terminale 1 | Terminale 2 | Il `ping` parte da |
    |---|---|---|---|
    | la workstation | `sudo tcpdump -ni virbr56 icmp` | `sudo tcpdump -ni wlp2s0 icmp` (la tua scheda verso Internet: `ip route get 1.1.1.1`) | `ssh lab ping -c 3 1.1.1.1` |

=== "Livello 2 · Workstation + server"

    | Router | Terminale 1 | Terminale 2 | Il `ping` parte da |
    |---|---|---|---|
    | `srv1` | `sudo tcpdump -ni virbr56 icmp` | `sudo tcpdump -ni eno1 icmp` (la scheda verso il router di casa) | `ssh lab ping -c 3 1.1.1.1` |

=== "Livello 3 · Mini data center"

    | Router | Terminale 1 | Terminale 2 | Il `ping` parte da |
    |---|---|---|---|
    | `nodo1` | `sudo tcpdump -ni lab0 icmp` | `sudo tcpdump -ni eno1 icmp` | `ssh nodo2 ping -c 3 1.1.1.1` |

    Ripeti con `ssh nodo3 ping -c 3 1.1.1.1`: sulla scheda `eno1` i pacchetti dei due nodi sono indistinguibili.

**Cosa devi vedere.** Sulla scheda interna il pacchetto ha come sorgente l'indirizzo privato (`192.168.56.21`
oppure `10.10.10.12`); su quella esterna lo stesso pacchetto ha come sorgente l'indirizzo del router. Le risposte
fanno il percorso inverso: è il *connection tracking* che sa a chi restituirle. Annota nel README di `academy-lab`
le righe catturate e spiega la differenza.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale comando mostra le porte TCP e UDP in ascolto con il processo associato?","options":["ss -tulpn","ip route","ping -a","dig -p"],"answer":[0],"explanation":"t, u, l, p, n: TCP, UDP, in ascolto, processo, numerico."},
 {"q":"Un servizio ascolta su 127.0.0.1:5000. Da un'altra macchina:","options":["È raggiungibile","Non è raggiungibile: ascolta solo sulla loopback","È raggiungibile solo con ping","Dipende dal DNS"],"answer":[1],"explanation":"127.0.0.1 accetta solo connessioni locali."},
 {"q":"Un ping a un server cloud non risponde. Cosa puoi concludere?","options":["Il server è spento","Nulla di certo: ICMP potrebbe essere bloccato","Il DNS non funziona","La porta 22 è chiusa"],"answer":[1],"explanation":"Molti ambienti bloccano ICMP: verifica le porte dei servizi."},
 {"q":"Quale opzione di curl mostra i dettagli di connessione, TLS e intestazioni?","options":["-I","-v","-o","-s"],"answer":[1],"explanation":"-v (verbose) mostra tutto lo scambio."},
 {"q":"Quale strumento cattura i pacchetti che passano su un'interfaccia?","options":["tcpdump","traceroute","ss","nc"],"answer":[0],"explanation":"tcpdump cattura e mostra il traffico; i file si analizzano anche con Wireshark."}
]}
</script>

## Fonti

- [ss(8)](https://man7.org/linux/man-pages/man8/ss.8.html), [ping(8)](https://man7.org/linux/man-pages/man8/ping.8.html), [tracepath(8)](https://man7.org/linux/man-pages/man8/tracepath.8.html)
- [tcpdump: documentazione e pcap-filter](https://www.tcpdump.org/manpages/)
- [curl: documentazione ufficiale](https://curl.se/docs/)
- [nc(1) OpenBSD netcat](https://man.openbsd.org/nc.1)
- [Wireshark User's Guide](https://www.wireshark.org/docs/wsug_html_chunked/)
