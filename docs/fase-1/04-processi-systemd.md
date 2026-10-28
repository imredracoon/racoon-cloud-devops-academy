# 1.4 Processi e servizi con systemd

!!! obiettivi "Obiettivi"
    - Osservare i processi con `ps`, `top` e `/proc` e capire PID, PPID e stati.
    - Inviare segnali ai processi e gestire priorità e job in background.
    - Gestire servizi con `systemctl` e leggere i log con `journalctl`.
    - Scrivere un'unità systemd per un servizio personalizzato.

## Processi

Ogni processo ha un **PID** (*Process ID*) unico e un **PPID**, il PID del processo che lo ha creato (il
**genitore**). Tutti i processi discendono dal primo, con PID 1, che nelle distribuzioni moderne è **systemd**.

```bash
ps aux                     # tutti i processi, con utente, CPU, memoria e comando
ps -ef --forest            # mostra l'albero genitore-figlio
ps -o pid,ppid,user,stat,cmd -p $$   # dettagli della shell corrente ($$ = PID della shell)
pgrep -a ssh               # cerca processi per nome
pstree -p | head           # albero dei processi
```

Nella colonna `STAT` di `ps` trovi lo stato: `R` in esecuzione, `S` in attesa (*sleeping*), `D` in attesa non
interrompibile (di solito I/O su disco), `T` fermato, `Z` **zombie** (terminato ma non ancora "raccolto" dal genitore).

### top e htop

`top` mostra i processi in tempo reale. Le righe in alto contano molto:

- **load average**: il numero medio di processi in esecuzione o in attesa negli ultimi 1, 5 e 15 minuti.
  Va confrontato con il numero di CPU: un carico di 4 su una macchina con 4 CPU significa "pieno".
- **%Cpu(s)**: `us` (programmi utente), `sy` (kernel), `id` (inattività), `wa` (attesa di I/O: se è alto il
  collo di bottiglia è il disco).
- **Mem / Swap**: memoria usata, libera e disponibile.

Dentro `top`: ++shift+p++ ordina per CPU, ++shift+m++ per memoria, `k` invia un segnale, `q` esce.
`htop` è una versione più comoda (installala con `sudo apt install htop`).

## Segnali

Un **segnale** è una notifica che il kernel consegna a un processo.

| Segnale | Numero | Significato |
|---|---|---|
| `SIGTERM` | 15 | "Per favore termina": il processo può chiudere in modo ordinato. È il default di `kill`. |
| `SIGKILL` | 9 | Terminazione immediata, non può essere ignorata. Da usare solo come ultima risorsa. |
| `SIGHUP` | 1 | Storicamente "terminale chiuso"; molti servizi lo usano per **ricaricare la configurazione**. |
| `SIGINT` | 2 | Interruzione, quello inviato da ++ctrl+c++. |
| `SIGSTOP` / `SIGCONT` | 19 / 18 | Sospende / riprende. |

```bash
kill 1234            # invia SIGTERM al PID 1234
kill -9 1234         # SIGKILL
kill -HUP 1234       # SIGHUP
pkill -f "python app.py"   # per nome o riga di comando
```

## Job in primo piano e in background

```bash
sleep 300 &          # avvia in background: la shell resta libera
jobs                 # elenca i job della shell
fg %1                # riporta il job 1 in primo piano
# premi Ctrl+Z per sospendere un processo in primo piano, poi:
bg %1                # lo fa continuare in background
nohup ./script.sh &  # continua anche se chiudi il terminale
```

Per i servizi veri non si usano `&` e `nohup`: si usa systemd.

## Priorità

`nice` va da -20 (massima priorità) a 19 (minima). Un utente normale può solo **abbassare** la priorità.

```bash
nice -n 10 tar -czf backup.tar.gz /home   # avvia con priorità bassa
sudo renice -n 5 -p 1234                  # cambia la priorità di un processo esistente
```

## systemd

**systemd** è il sistema di avvio e di gestione dei servizi di quasi tutte le distribuzioni moderne. Gestisce
**unità** (*unit*) di vari tipi: `.service` (servizi), `.timer` (pianificazioni, lezione 2.5), `.socket`,
`.mount`, `.target` (gruppi di unità, l'equivalente dei vecchi *runlevel*).

### systemctl

```bash
systemctl status ssh             # stato, PID, ultime righe di log
sudo systemctl start nginx       # avvia ora
sudo systemctl stop nginx        # ferma ora
sudo systemctl restart nginx     # ferma e riavvia
sudo systemctl reload nginx      # ricarica la configurazione senza interrompere il servizio
sudo systemctl enable nginx      # avvio automatico al boot (non lo avvia ora)
sudo systemctl enable --now nginx   # abilita e avvia
sudo systemctl disable nginx     # niente avvio automatico
sudo systemctl mask nginx        # impedisce qualsiasi avvio, anche manuale
systemctl is-active nginx        # active / inactive / failed
systemctl is-enabled nginx       # enabled / disabled
systemctl list-units --type=service --state=failed   # servizi in errore
systemctl get-default            # target di avvio: graphical.target o multi-user.target
```

!!! tip "start ≠ enable"
    `start` agisce **adesso**; `enable` agisce **al prossimo avvio**. Dimenticare `enable` è la causa classica del
    "dopo il riavvio il servizio non c'è più".

### journalctl: i log di systemd

systemd raccoglie i log di tutti i servizi nel **journal**.

```bash
journalctl -u ssh                 # log del servizio ssh
journalctl -u ssh -f              # in tempo reale
journalctl -u nginx --since "1 hour ago"
journalctl -p err -b              # solo errori (priority err e superiori) dall'ultimo avvio
journalctl -b -1                  # log dell'avvio precedente (se il journal è persistente)
journalctl --disk-usage           # spazio occupato
sudo journalctl --vacuum-time=7d  # elimina i log più vecchi di 7 giorni
```

Molti servizi scrivono anche file di testo classici in `/var/log` (per esempio `/var/log/nginx/`), e il servizio
`rsyslog` copia i messaggi del journal in file di testo. I nomi dei file cambiano tra le famiglie:

| Contenuto | Debian/Ubuntu | RHEL (Rocky/AlmaLinux) |
|---|---|---|
| Messaggi generali di sistema | `/var/log/syslog` | `/var/log/messages` |
| Accessi, sudo, SSH | `/var/log/auth.log` | `/var/log/secure` |
| Nome del servizio SSH | `ssh` | `sshd` |
| Server web Apache | `apache2` | `httpd` |
| Log di audit (SELinux) | — | `/var/log/audit/audit.log` |

## Scrivere un servizio

Le unità installate dai pacchetti stanno in `/usr/lib/systemd/system/` (o `/lib/systemd/system/`); quelle
dell'amministratore in **`/etc/systemd/system/`**, che ha la precedenza.

Creiamo un piccolo servizio che scrive l'ora in un file ogni 10 secondi. Prima lo script:

```bash
sudo tee /usr/local/bin/orologio.sh > /dev/null <<'EOF'
#!/usr/bin/env bash
while true; do
  date >> /var/lib/orologio/ora.log
  sleep 10
done
EOF
sudo chmod 755 /usr/local/bin/orologio.sh
```

Poi l'unità:

```ini title="/etc/systemd/system/orologio.service"
[Unit]
Description=Scrive l'ora ogni 10 secondi
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/orologio.sh
User=orologio
StateDirectory=orologio
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

| Direttiva | Significato |
|---|---|
| `After=` | Ordine di avvio: parte dopo la rete. Non crea dipendenza (per quella c'è `Requires=`/`Wants=`). |
| `Type=simple` | Il processo di `ExecStart` è il servizio stesso. |
| `User=` | Gira con un utente dedicato senza privilegi: buona pratica di sicurezza. |
| `StateDirectory=` | systemd crea `/var/lib/orologio` con il proprietario giusto. |
| `Restart=on-failure` | Riavvia il servizio se termina con un errore. |
| `WantedBy=multi-user.target` | Con `enable`, il servizio parte all'avvio normale del sistema. |

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin orologio
sudo systemctl daemon-reload            # rilegge le unità dopo ogni modifica
sudo systemctl enable --now orologio
systemctl status orologio
tail -f /var/lib/orologio/ora.log
```

Per modificare un'unità fornita da un pacchetto **senza toccare l'originale** si usa un *drop-in*:
`sudo systemctl edit nginx` crea `/etc/systemd/system/nginx.service.d/override.conf`.

## Esercizi

!!! esercizio "Esercizio 1 · Il processo più affamato"
    Avvia in background un processo che consuma CPU: `yes > /dev/null &`. Trova il suo PID con due metodi diversi,
    osservalo in `top`, abbassane la priorità a 19 e terminalo con un segnale ordinato.

??? success "Soluzione"
    ```bash
    yes > /dev/null &
    echo $!                  # $! è il PID dell'ultimo processo avviato in background
    pgrep -a yes
    top                      # in cima, con circa 100% di CPU; esci con q
    renice -n 19 -p $(pgrep yes)
    kill $(pgrep yes)        # SIGTERM
    ```
    **Spiegazione.** `$( ... )` (*command substitution*) inserisce l'output di un comando dentro un altro. Un utente
    normale può alzare il valore di nice del proprio processo senza sudo.

!!! esercizio "Esercizio 2 · Il servizio che si rialza"
    Con il servizio `orologio` attivo, uccidi il suo processo con `SIGKILL`. Cosa fa systemd? Come lo verifichi?
    Poi fermalo con `systemctl stop`: si riavvia anche in questo caso? Perché?

??? success "Soluzione"
    ```bash
    sudo kill -9 $(systemctl show -p MainPID --value orologio)
    sleep 6; systemctl status orologio          # di nuovo active, con un nuovo PID
    journalctl -u orologio -n 10                # "Main process exited, code=killed, status=9/KILL" e poi il riavvio
    sudo systemctl stop orologio                # resta fermo
    ```
    **Spiegazione.** `Restart=on-failure` riavvia quando il processo termina in modo anomalo (codice d'errore o
    segnale non pulito). Uno `stop` richiesto dall'amministratore non è un fallimento, quindi systemd non riavvia.

!!! esercizio "Esercizio 3 · Diagnosi"
    Modifica l'unità `orologio.service` scrivendo per errore `ExecStart=/usr/local/bin/orologi.sh`, ricarica e
    riavvia. Usa solo `systemctl` e `journalctl` per scoprire il problema, poi correggilo.

??? success "Soluzione"
    ```bash
    sudo systemctl daemon-reload && sudo systemctl restart orologio
    systemctl status orologio        # failed, status=203/EXEC
    journalctl -u orologio -n 20     # "Failed to locate executable /usr/local/bin/orologi.sh"
    ```
    Correggi il percorso, poi `daemon-reload` e `restart`.

    **Spiegazione.** Il codice `203/EXEC` significa che systemd non è riuscito a eseguire il programma: percorso
    errato, permesso di esecuzione mancante o interprete inesistente. Leggere `status` e il journal è sempre il
    primo passo con un servizio in errore.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale processo ha PID 1 nelle distribuzioni moderne?","options":["bash","init.d","systemd","kernel"],"answer":[2],"explanation":"systemd è il primo processo in spazio utente e il genitore di tutti gli altri."},
 {"q":"Quale segnale chiede a un processo di terminare in modo ordinato?","options":["SIGKILL","SIGTERM","SIGSTOP","SIGCONT"],"answer":[1],"explanation":"SIGTERM (15) è il default di kill e può essere gestito dal processo."},
 {"q":"Un servizio funziona ma dopo il riavvio del server non parte. Cosa manca probabilmente?","options":["systemctl start","systemctl enable","systemctl reload","daemon-reload"],"answer":[1],"explanation":"enable crea il collegamento per l'avvio automatico."},
 {"q":"Dopo aver modificato un file .service, quale comando serve prima del restart?","options":["systemctl daemon-reload","systemctl reset","journalctl --rotate","systemctl enable"],"answer":[0],"explanation":"systemd deve rileggere le definizioni delle unità."},
 {"q":"Come vedi in tempo reale i log del servizio nginx?","options":["journalctl -u nginx -f","systemctl log nginx","tail nginx","journalctl -p nginx"],"answer":[0],"explanation":"-u filtra per unità, -f segue i nuovi messaggi."},
 {"q":"Load average 8,0 su un server con 2 CPU indica:","options":["Server scarico","Server molto sovraccarico","Server con 8 utenti collegati","Memoria piena"],"answer":[1],"explanation":"In media 8 processi competono per 2 CPU."},
 {"q":"Dove vanno messe le unità systemd create dall'amministratore?","options":["/usr/lib/systemd/system","/etc/systemd/system","/var/lib/systemd","/home/utente"],"answer":[1],"explanation":"/etc/systemd/system ha precedenza sulle unità dei pacchetti."}
]}
</script>

## Fonti

- [systemd: documentazione ufficiale](https://systemd.io/)
- [systemctl(1)](https://www.freedesktop.org/software/systemd/man/latest/systemctl.html), [journalctl(1)](https://www.freedesktop.org/software/systemd/man/latest/journalctl.html)
- [systemd.service(5)](https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html), [systemd.exec(5)](https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html)
- [signal(7), man page](https://man7.org/linux/man-pages/man7/signal.7.html)
- [ps(1)](https://man7.org/linux/man-pages/man1/ps.1.html) e [top(1)](https://man7.org/linux/man-pages/man1/top.1.html)
