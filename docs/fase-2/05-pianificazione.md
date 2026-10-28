# 2.5 Pianificare i job

!!! obiettivi "Obiettivi"
    - Pianificare attività ricorrenti con cron e leggere la sintassi delle cinque colonne.
    - Pianificare con i timer di systemd e capire quando preferirli.
    - Eseguire un'attività una sola volta con `at` e con `systemd-run`.
    - Verificare che i job siano davvero eseguiti e trovare i loro log.

Questa lezione si svolge sulla VM `lab` (o `rhel-lab`): il Dev Container non esegue systemd.

## Due strumenti per lo stesso problema

Backup notturni, pulizia dei file temporanei, rinnovo dei certificati, report settimanali: in ogni server ci sono
attività da eseguire a orari fissi. Linux offre due strumenti:

- **cron**, lo strumento storico, presente ovunque, semplice;
- i **timer di systemd**, più recenti, integrati con i servizi e con il journal.

Devi conoscerli entrambi: troverai cron in moltissimi sistemi esistenti, e i timer sono sempre più usati (anche
dalla distribuzione stessa: `systemctl list-timers` lo mostra).

=== "Debian/Ubuntu"

    Il demone si chiama `cron` (pacchetto `cron`) ed è installato e attivo per impostazione predefinita.

    ```bash
    systemctl status cron
    ```

=== "RHEL (Rocky/AlmaLinux)"

    Il demone si chiama `crond` (pacchetto `cronie`). Su un'installazione minimale potrebbe mancare:

    ```bash
    sudo dnf install -y cronie
    sudo systemctl enable --now crond
    ```

## cron

### Il crontab dell'utente

Ogni utente ha un proprio **crontab**, un elenco di job. Non si modifica il file direttamente, ma con il comando
`crontab`:

```bash
crontab -e        # modifica (al primo uso chiede l'editor)
crontab -l        # mostra
crontab -r        # cancella TUTTO il crontab, senza conferma: attenzione
sudo crontab -u alice -l    # il crontab di un altro utente
```

### La sintassi

Ogni riga è un job, con cinque campi per l'orario seguiti dal comando:

```text
┌───────────── minuto (0-59)
│ ┌─────────── ora (0-23)
│ │ ┌───────── giorno del mese (1-31)
│ │ │ ┌─────── mese (1-12)
│ │ │ │ ┌───── giorno della settimana (0-7: 0 e 7 = domenica)
│ │ │ │ │
* * * * *  comando
```

| Simbolo | Significato | Esempio |
|---|---|---|
| `*` | ogni valore | `* * * * *` ogni minuto |
| `,` | elenco | `0 8,20 * * *` alle 8:00 e alle 20:00 |
| `-` | intervallo | `0 9-18 * * 1-5` ogni ora dalle 9 alle 18, da lunedì a venerdì |
| `/` | passo | `*/15 * * * *` ogni 15 minuti |

Esempi:

```text
30 2 * * *     /home/vincenzo/bin/backup.sh       # ogni notte alle 2:30
0 */6 * * *    /usr/local/bin/pulizia.sh          # alle 0, 6, 12, 18
0 9 * * 1      /usr/local/bin/report.sh           # il lunedì alle 9:00
@reboot        /usr/local/bin/all-avvio.sh        # a ogni avvio
```

!!! warning "La trappola del giorno"
    Se specifichi **sia** il giorno del mese **sia** il giorno della settimana, cron esegue il job quando
    **uno qualsiasi** dei due corrisponde. `0 0 13 * 5` gira ogni 13 del mese **e** ogni venerdì, non solo il
    venerdì 13.

### L'ambiente di cron è minimale

L'errore più comune: uno script funziona nel terminale e non funziona in cron. Il motivo è che cron esegue i job
con un **ambiente ridotto**: `PATH` contiene solo `/usr/bin:/bin`, la shell è `/bin/sh`, non vengono letti
`.bashrc` né `.profile`, la cartella di lavoro è la home.

Regole pratiche:

- usa **percorsi assoluti** per script, file e comandi non standard;
- imposta le variabili in cima al crontab;
- **ridireziona l'output** in un file o nel journal, altrimenti cron prova a inviarlo per email e, senza un server
  di posta, va perso.

```text
SHELL=/bin/bash
PATH=/usr/local/bin:/usr/bin:/bin
MAILTO=""

30 2 * * * /home/vincenzo/bin/backup.sh /home/vincenzo/progetti /srv/backup 2>&1 | logger -t backup
```

`%` ha un significato speciale nel crontab (diventa un "a capo"): in un comando come `date +%F` va scritto `\%`.
È un'altra ragione per mettere la logica in uno script e lasciare nel crontab solo la sua chiamata.

### I crontab di sistema

Oltre a quelli degli utenti esistono file di sistema, che hanno **una colonna in più**: l'utente con cui eseguire.

```text title="/etc/cron.d/pulizia-tmp"
# m h dom mon dow utente comando
15 3 * * * root /usr/local/bin/pulizia-tmp.sh
```

| Posizione | Uso |
|---|---|
| `/etc/crontab` | il crontab principale di sistema (meglio non modificarlo) |
| `/etc/cron.d/` | file dei pacchetti o tuoi, con la colonna utente |
| `/etc/cron.daily/`, `weekly/`, `monthly/`, `hourly/` | script eseguiti con quella frequenza, senza orario preciso |

Su Debian/Ubuntu gli script in `cron.daily/` non devono avere estensione: `run-parts` ignora i nomi che contengono
un punto, quindi `backup.sh` non verrebbe eseguito, `backup` sì.

### Chi può usare cron

I file `/etc/cron.allow` e `/etc/cron.deny` limitano gli utenti. Se `cron.allow` esiste, solo gli utenti elencati
possono usare `crontab`.

## I timer di systemd

Un timer di systemd è una coppia di unità:

- un **servizio** (`.service`) che descrive **cosa** eseguire;
- un **timer** (`.timer`) con lo stesso nome che descrive **quando**.

```ini title="/etc/systemd/system/backup.service"
[Unit]
Description=Backup dei progetti

[Service]
Type=oneshot
User=vincenzo
ExecStart=/home/vincenzo/bin/backup.sh /home/vincenzo/progetti /srv/backup
```

```ini title="/etc/systemd/system/backup.timer"
[Unit]
Description=Backup notturno dei progetti

[Timer]
OnCalendar=*-*-* 02:30:00
Persistent=true
RandomizedDelaySec=10min

[Install]
WantedBy=timers.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now backup.timer     # si abilita il TIMER, non il servizio
systemctl list-timers backup.timer           # prossima e ultima esecuzione
sudo systemctl start backup.service          # esecuzione manuale immediata, per provarlo
journalctl -u backup.service                 # l'output, già nel journal
```

- `Type=oneshot`: il servizio esegue un compito e termina.
- `Persistent=true`: se la macchina era spenta all'orario previsto, il job parte appena possibile dopo l'avvio.
  cron non ha questa funzione (esiste `anacron`, uno strumento separato).
- `RandomizedDelaySec` sposta l'avvio di un tempo casuale: utile quando cento server non devono colpire lo stesso
  sistema di backup nello stesso secondo.

### La sintassi di OnCalendar

Il formato è `GiornoSettimana Anno-Mese-Giorno Ora:Minuto:Secondo`, con `*` come jolly e `,` e `..` per elenchi e
intervalli:

| Espressione | Significato |
|---|---|
| `daily` | ogni giorno a mezzanotte (anche `hourly`, `weekly`, `monthly`) |
| `*-*-* 02:30:00` | ogni giorno alle 2:30 |
| `Mon..Fri 09:00` | dal lunedì al venerdì alle 9 |
| `*-*-01 00:00:00` | il primo di ogni mese |
| `*:0/15` | ogni 15 minuti |

Puoi verificare un'espressione prima di usarla:

```bash
systemd-analyze calendar "Mon..Fri 09:00" --iterations=3
```

Esistono anche timer **relativi**: `OnBootSec=5min` (5 minuti dopo l'avvio), `OnUnitActiveSec=1h` (un'ora dopo
l'ultima esecuzione).

### cron o timer?

| | cron | Timer systemd |
|---|---|---|
| Configurazione | una riga | due file |
| Log | da gestire a mano | automatici nel journal |
| Job persi a macchina spenta | persi (salvo anacron) | recuperati con `Persistent=true` |
| Limiti di risorse, sandbox, dipendenze | no | sì, tutte le opzioni dei servizi |
| Esecuzione manuale identica a quella pianificata | difficile (ambiente diverso) | `systemctl start` |
| Diffusione | ovunque, anche fuori da Linux | sistemi con systemd |

Per un nuovo job su un server moderno, il timer è in genere la scelta migliore. cron resta perfetto per cose
semplici e lo incontrerai ovunque.

## Esecuzioni singole

```bash
echo "/usr/local/bin/riavvia-app.sh" | at 23:00     # una volta, alle 23 (lezione 1.15)
sudo systemd-run --on-active=30min /usr/bin/systemctl restart nginx   # timer temporaneo
systemctl list-timers                                # compare come run-....timer
```

## Esercizi

!!! esercizio "Esercizio 1 · Leggere le espressioni"
    Che cosa significano?

    1. `0 22 * * 1-5`
    2. `*/10 8-17 * * *`
    3. `0 0 1 1,7 *`
    4. `OnCalendar=Sat,Sun 10:00`
    5. `OnCalendar=*-*-* *:00/30:00`

??? success "Soluzione"
    1. Alle 22:00, dal lunedì al venerdì.
    2. Ogni 10 minuti tra le 8:00 e le 17:50.
    3. A mezzanotte il primo gennaio e il primo luglio.
    4. Sabato e domenica alle 10:00.
    5. Ogni ora, ai minuti 0 e 30.

    **Spiegazione.** Puoi sempre verificare i timer con `systemd-analyze calendar`. Per cron, siti come
    crontab.guru sono comodi, ma il riferimento resta `man 5 crontab`.

!!! esercizio "Esercizio 2 · Il backup con cron"
    Copia sulla VM lo script `backup.sh` della lezione 2.4 in `~/bin/`. Pianificalo con cron ogni 2 minuti
    (solo per la prova) per salvare `~/progetti` in `~/backup`, con l'output nel journal. Verifica che funzioni,
    poi rimuovi il job.

??? success "Soluzione"
    ```bash
    mkdir -p ~/bin ~/progetti ~/backup && echo prova > ~/progetti/file.txt
    chmod +x ~/bin/backup.sh
    crontab -e
    ```
    ```text
    */2 * * * * $HOME/bin/backup.sh $HOME/progetti $HOME/backup 2>&1 | logger -t backup-cron
    ```
    ```bash
    # dopo qualche minuto
    ls ~/backup
    journalctl -t backup-cron -t backup --since "10 min ago"
    crontab -e      # cancella la riga
    ```
    **Spiegazione.** In cron `$HOME` è impostata, quindi si può usare. Se il job non parte: controlla il log del
    demone (`journalctl -u cron` su Ubuntu, `journalctl -u crond` o `/var/log/cron` su RHEL) e prova lo script con
    un ambiente vuoto: `env -i HOME="$HOME" /bin/sh -c '...'`.

!!! esercizio "Esercizio 3 · Lo stesso job come timer"
    Trasforma il job dell'esercizio 2 in un servizio e un timer di systemd che girino ogni giorno alle 3:00 con
    `Persistent=true`, eseguiti come il tuo utente. Installa lo script in `/usr/local/bin/backup.sh`, la posizione corretta per gli script di sistema.
    Prima di attivarlo, esegui il servizio manualmente e controlla il
    journal. Mostra quando avverrà la prossima esecuzione.

??? success "Soluzione"
    ```bash
    sudo install -m 755 ~/bin/backup.sh /usr/local/bin/backup.sh
    sudo tee /etc/systemd/system/backup-progetti.service > /dev/null <<EOF
    [Unit]
    Description=Backup dei progetti

    [Service]
    Type=oneshot
    User=$USER
    ExecStart=/usr/local/bin/backup.sh $HOME/progetti $HOME/backup
    EOF

    sudo tee /etc/systemd/system/backup-progetti.timer > /dev/null <<'EOF'
    [Unit]
    Description=Backup giornaliero dei progetti

    [Timer]
    OnCalendar=*-*-* 03:00:00
    Persistent=true

    [Install]
    WantedBy=timers.target
    EOF

    sudo systemctl daemon-reload
    sudo systemctl start backup-progetti.service
    journalctl -u backup-progetti.service -n 20
    sudo systemctl enable --now backup-progetti.timer
    systemctl list-timers backup-progetti.timer
    ```
    **Spiegazione.** Nel primo here document `EOF` senza apici fa espandere `$USER` e `$HOME` subito, scrivendo nel
    file i valori reali (systemd non conosce le variabili della tua shell). Nel secondo `'EOF'` tra apici evita
    qualsiasi espansione.

    **Perché non lasciare lo script nella home?** Sul percorso RHEL, SELinux impedisce a systemd di eseguire file
    con l'etichetta `user_home_t`: il servizio fallirebbe con `status=203/EXEC` e `Permission denied`, e
    `sudo ausearch -m AVC -ts recent` mostrerebbe il blocco. `/usr/local/bin/` ha l'etichetta corretta (`bin_t`).
    Su Ubuntu funzionerebbe anche dalla home, ma la buona pratica è la stessa ovunque.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa significa `15 4 * * 0` in un crontab?","options":["Alle 4:15 di ogni domenica","Alle 15:04 di ogni giorno","Ogni 15 minuti alle 4","Il 15 di aprile"],"answer":[0],"explanation":"minuto 15, ora 4, ogni giorno del mese, ogni mese, giorno della settimana 0 = domenica."},
 {"q":"Uno script funziona nel terminale ma non in cron. Qual è la causa più probabile?","options":["cron non supporta Bash","L'ambiente di cron è minimale: PATH diverso, nessun profilo caricato","cron funziona solo per root","Lo script è troppo lungo"],"answer":[1],"explanation":"Usa percorsi assoluti e imposta le variabili necessarie."},
 {"q":"Nei timer di systemd, quale unità si abilita con `systemctl enable`?","options":["Il .service","Il .timer","Entrambe","Nessuna, basta daemon-reload"],"answer":[1],"explanation":"Il timer attiva il servizio con lo stesso nome all'orario previsto."},
 {"q":"Cosa fa `Persistent=true` in un timer?","options":["Rende il timer permanente dopo il riavvio","Esegue il job appena possibile se l'esecuzione prevista è stata persa a macchina spenta","Ripete il job all'infinito","Conserva i log"],"answer":[1],"explanation":"Sostituisce la funzione di anacron."},
 {"q":"Quali file di cron richiedono la colonna con il nome utente?","options":["Il crontab personale modificato con crontab -e","/etc/crontab","I file in /etc/cron.d/","Gli script in /etc/cron.daily/","/etc/cron.allow"],"answer":[1,2],"explanation":"I crontab di sistema indicano con quale utente eseguire il comando."},
 {"q":"Come verifichi quando scatterà un'espressione OnCalendar?","options":["systemctl calendar","systemd-analyze calendar \"...\"","crontab -t","timedatectl"],"answer":[1],"explanation":"Mostra la forma normalizzata e le prossime esecuzioni."}
]}
</script>

## Fonti

- [crontab(5)](https://man7.org/linux/man-pages/man5/crontab.5.html) e [cron(8)](https://manpages.debian.org/stable/cron/cron.8.en.html)
- [systemd.timer(5)](https://www.freedesktop.org/software/systemd/man/latest/systemd.timer.html)
- [systemd.time(7): formato delle date](https://www.freedesktop.org/software/systemd/man/latest/systemd.time.html)
- [systemd-run(1)](https://www.freedesktop.org/software/systemd/man/latest/systemd-run.html)
- [Red Hat Enterprise Linux 10: documentazione](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10)
- [Debian: run-parts(8)](https://manpages.debian.org/stable/debianutils/run-parts.8.en.html)
