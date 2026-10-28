# 1.2 Lavorare con file e testo

!!! obiettivi "Obiettivi"
    - Leggere file di qualsiasi dimensione con `less`, `head`, `tail`.
    - Capire standard input, output ed errore e usare redirezioni e pipe.
    - Filtrare e trasformare testo con `grep`, `cut`, `sort`, `uniq`, `wc`, `tr`, `sed`, `awk`.
    - Archiviare e comprimere con `tar` e `gzip`.

## Leggere i file

| Comando | Uso |
|---|---|
| `cat file` | Stampa tutto il file. Adatto ai file brevi. |
| `less file` | Sfoglia file lunghi: frecce per scorrere, `/parola` per cercare, `G` per andare in fondo, `q` per uscire. |
| `head -n 20 file` | Le prime 20 righe. |
| `tail -n 20 file` | Le ultime 20 righe. |
| `tail -f file` | Segue il file in tempo reale mentre cresce: fondamentale per i log. ++ctrl+c++ per uscire. |
| `wc -l file` | Conta le righe (`-w` parole, `-c` byte). |

## I tre canali standard

Ogni processo nasce con tre canali di comunicazione, identificati da un numero (*file descriptor*):

| Numero | Nome | Default |
|---|---|---|
| 0 | **stdin**, standard input | la tastiera |
| 1 | **stdout**, standard output | lo schermo |
| 2 | **stderr**, standard error | lo schermo |

Output normale ed errori vanno entrambi sullo schermo, ma su canali **separati**: questo permette di trattarli in
modo diverso.

### Redirezioni

```bash
ls /etc > elenco.txt            # stdout su file (sovrascrive)
ls /etc >> elenco.txt           # stdout su file (aggiunge)
ls /nonesiste 2> errori.txt     # solo stderr su file
find / -name passwd 2>/dev/null # scarta gli errori
comando > tutto.log 2>&1        # stdout su file e stderr dove va stdout
comando &> tutto.log            # forma breve di Bash, stesso effetto
wc -l < /etc/passwd             # il file diventa lo stdin del comando
```

!!! warning "L'ordine conta"
    `comando > file 2>&1` manda entrambi i canali nel file. `comando 2>&1 > file` no: lo stderr viene collegato
    a dove punta stdout **in quel momento** (lo schermo), e solo dopo stdout viene spostato sul file.

### Pipe

La **pipe** `|` collega lo stdout di un comando allo stdin del successivo. È il cuore della filosofia Unix: tanti
piccoli strumenti, ognuno fa una cosa bene, combinati tra loro.

```bash
ls /etc | wc -l                 # quanti elementi ci sono in /etc
cat /etc/passwd | sort | head -3
```

`tee` scrive contemporaneamente su schermo e su file, utile per salvare un output mentre lo guardi:

```bash
ls -l /etc | tee elenco.txt | wc -l
```

## Gli strumenti per il testo

Useremo come esempio il file `/etc/passwd`, che elenca gli utenti del sistema. Ogni riga ha sette campi separati
da `:`:

```text
ubuntu:x:1000:1000:Ubuntu:/home/ubuntu:/bin/bash
│      │ │    │    │      │            └ shell
│      │ │    │    │      └ directory home
│      │ │    │    └ descrizione
│      │ │    └ ID del gruppo principale (GID)
│      │ └ ID dell'utente (UID)
│      └ la password non è qui (x): è in /etc/shadow
└ nome utente
```

### grep: cercare righe

```bash
grep bash /etc/passwd           # righe che contengono "bash"
grep -i ERROR app.log           # senza distinguere maiuscole/minuscole
grep -v nologin /etc/passwd     # righe che NON contengono "nologin"
grep -n root /etc/passwd        # con il numero di riga
grep -c bash /etc/passwd        # conta le righe trovate
grep -r "Listen" /etc/ 2>/dev/null   # cerca in tutti i file di una directory
grep -E "^(root|ubuntu):" /etc/passwd  # espressione regolare estesa
```

### Espressioni regolari essenziali

Un'**espressione regolare** (*regex*) descrive uno schema di testo:

| Simbolo | Significato | Esempio |
|---|---|---|
| `^` | inizio riga | `^root` |
| `$` | fine riga | `bash$` |
| `.` | un carattere qualsiasi | `r..t` |
| `*` | zero o più volte il precedente | `ab*c` |
| `[abc]` | uno dei caratteri | `[0-9]` |
| `+`, `?`, `|`, `( )` | uno o più, zero o uno, alternativa, gruppo (con `grep -E`) | `(GET|POST)` |

### cut, sort, uniq, tr

```bash
cut -d: -f1 /etc/passwd             # primo campo, separatore ":"
cut -d: -f1,7 /etc/passwd           # campi 1 e 7
sort file                           # ordina alfabeticamente
sort -n file                        # ordina numericamente
sort -t: -k3 -n /etc/passwd         # ordina per terzo campo numerico
sort file | uniq                    # elimina i duplicati consecutivi (per questo si ordina prima)
sort file | uniq -c | sort -rn      # conta le occorrenze e mostra le più frequenti
echo "ciao" | tr 'a-z' 'A-Z'        # trasforma caratteri: CIAO
```

### sed: modificare il testo

`sed` (*stream editor*) trasforma il testo riga per riga. L'uso più comune è la sostituzione:

```bash
sed 's/vecchio/nuovo/' file         # sostituisce la prima occorrenza di ogni riga (stampa a schermo)
sed 's/vecchio/nuovo/g' file        # tutte le occorrenze
sed -i 's/vecchio/nuovo/g' file     # modifica il file sul posto
sed -n '5,10p' file                 # stampa solo le righe 5-10
sed '/^#/d' file                    # elimina le righe di commento
```

!!! tip "Prima guarda, poi modifica"
    Prova sempre il comando **senza** `-i` e controlla l'output. Solo quando è corretto aggiungi `-i`.
    Puoi anche usare `-i.bak` per conservare una copia di sicurezza del file originale.

### awk: lavorare per colonne

`awk` è un piccolo linguaggio per elaborare testo a colonne. Le colonne si chiamano `$1`, `$2`…; `$0` è la riga intera.

```bash
awk -F: '{print $1, $3}' /etc/passwd            # nome e UID
awk -F: '$3 >= 1000 {print $1}' /etc/passwd     # utenti "umani" (UID >= 1000)
df -h | awk 'NR>1 {print $5, $6}'               # uso % e punto di montaggio, saltando l'intestazione
```

## Un esempio reale: analizzare un log

Crea un piccolo log di un server web per esercitarti:

```bash
cat > access.log <<'EOF'
203.0.113.5 - - [28/Oct/2026:10:00:01 +0000] "GET / HTTP/1.1" 200 512
198.51.100.7 - - [28/Oct/2026:10:00:02 +0000] "GET /login HTTP/1.1" 200 1024
203.0.113.5 - - [28/Oct/2026:10:00:03 +0000] "POST /login HTTP/1.1" 401 64
192.0.2.44 - - [28/Oct/2026:10:00:04 +0000] "GET /admin HTTP/1.1" 403 128
203.0.113.5 - - [28/Oct/2026:10:00:05 +0000] "POST /login HTTP/1.1" 401 64
198.51.100.7 - - [28/Oct/2026:10:00:06 +0000] "GET /missing HTTP/1.1" 404 256
203.0.113.5 - - [28/Oct/2026:10:00:07 +0000] "POST /login HTTP/1.1" 200 2048
EOF
```

La sintassi `<<'EOF' ... EOF` si chiama **here document**: passa al comando tutte le righe fino a `EOF`.
Gli indirizzi `203.0.113.x`, `198.51.100.x` e `192.0.2.x` sono riservati alla documentazione (RFC 5737).

```bash
awk '{print $1}' access.log | sort | uniq -c | sort -rn   # richieste per IP
awk '{print $9}' access.log | sort | uniq -c               # richieste per codice di stato
grep '" 401 ' access.log | awk '{print $1}' | sort | uniq -c   # login falliti per IP
```

## Archiviare e comprimere

```bash
tar -czf backup-etc.tar.gz /etc/hostname /etc/hosts   # crea (c) un archivio compresso con gzip (z) nel file (f)
tar -tzf backup-etc.tar.gz                            # elenca (t) il contenuto
tar -xzf backup-etc.tar.gz -C /tmp                    # estrae (x) nella directory /tmp
gzip file.log        # comprime in file.log.gz (e rimuove l'originale)
gunzip file.log.gz   # decomprime
zcat file.log.gz     # legge un file compresso senza decomprimerlo
```

## Esercizi

!!! esercizio "Esercizio 1 · Utenti e shell"
    Usando `/etc/passwd`, mostra un elenco ordinato alfabeticamente dei soli **nomi** degli utenti la cui shell è
    `/bin/bash`.

??? success "Soluzione"
    ```bash
    grep '/bin/bash$' /etc/passwd | cut -d: -f1 | sort
    ```
    oppure, solo con awk:
    ```bash
    awk -F: '$7 == "/bin/bash" {print $1}' /etc/passwd | sort
    ```
    **Spiegazione.** `$` ancora la ricerca alla fine della riga, così non trovi per errore shell come
    `/bin/bashx`. `cut` estrae il primo campo.

!!! esercizio "Esercizio 2 · Analisi del log"
    Con il file `access.log` dell'esempio:

    1. conta quante richieste hanno ricevuto un codice di errore (4xx);
    2. mostra il percorso più richiesto;
    3. salva in `sospetti.txt` gli IP con almeno due login falliti (401).

??? success "Soluzione"
    ```bash
    awk '$9 ~ /^4/' access.log | wc -l
    awk '{print $7}' access.log | sort | uniq -c | sort -rn | head -1
    grep '" 401 ' access.log | awk '{print $1}' | sort | uniq -c | awk '$1 >= 2 {print $2}' > sospetti.txt
    ```
    Risultati: 4 errori; `/login` con 4 richieste; `sospetti.txt` contiene `203.0.113.5`.

    **Spiegazione.** In awk `~` confronta con un'espressione regolare. Nella terza riga, dopo `uniq -c` la colonna 1
    è il conteggio e la 2 è l'IP. Questo tipo di analisi è la base di strumenti come fail2ban, che userai nel progetto.

!!! esercizio "Esercizio 3 · Redirezioni"
    Esegui `ls /etc /nonesiste` in modo che l'elenco di `/etc` finisca in `ok.txt` e il messaggio di errore in
    `ko.txt`. Poi esegui lo stesso comando mandando entrambi in `tutto.txt`.

??? success "Soluzione"
    ```bash
    ls /etc /nonesiste > ok.txt 2> ko.txt
    ls /etc /nonesiste > tutto.txt 2>&1
    ```

!!! esercizio "Esercizio 4 · sed"
    Copia `/etc/ssh/sshd_config` nella tua home. Con un solo comando `sed`, mostra il file senza righe di commento
    e senza righe vuote.

??? success "Soluzione"
    ```bash
    cp /etc/ssh/sshd_config ~/
    sed -e '/^#/d' -e '/^$/d' ~/sshd_config
    ```
    Alternativa con grep: `grep -Ev '^(#|$)' ~/sshd_config`.

    **Spiegazione.** `^$` corrisponde a una riga in cui inizio e fine coincidono, cioè vuota. È il modo più rapido
    per vedere solo le impostazioni attive di un file di configurazione.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale comando segue un file di log in tempo reale?","options":["head -f","tail -f","less -r","cat -f"],"answer":[1],"explanation":"tail -f mostra le nuove righe man mano che vengono aggiunte."},
 {"q":"Quale numero identifica lo standard error?","options":["0","1","2","3"],"answer":[2],"explanation":"0 stdin, 1 stdout, 2 stderr."},
 {"q":"Cosa fa `cmd > out.txt 2>&1`?","options":["Solo stdout nel file","Solo stderr nel file","Sia stdout sia stderr nel file","Legge out.txt come input"],"answer":[2],"explanation":"Prima stdout va nel file, poi stderr viene duplicato su stdout."},
 {"q":"Perché si usa `sort` prima di `uniq`?","options":["uniq elimina solo i duplicati consecutivi","uniq funziona solo su file ordinati per legge","Per velocizzare la pipe","Non serve"],"answer":[0],"explanation":"uniq confronta ogni riga con la precedente."},
 {"q":"Quale comando sostituisce tutte le occorrenze di http con https in un file, modificandolo?","options":["sed 's/http/https/' file","sed -i 's/http/https/g' file","grep -r http https file","tr http https file"],"answer":[1],"explanation":"-i modifica sul posto, g sostituisce tutte le occorrenze nella riga."},
 {"q":"Quale opzione di tar crea un archivio compresso con gzip?","options":["-xzf","-tzf","-czf","-cvf"],"answer":[2],"explanation":"c = create, z = gzip, f = file."}
]}
</script>

## Fonti

- [GNU Bash manual: Redirections](https://www.gnu.org/software/bash/manual/html_node/Redirections.html)
- [GNU grep manual](https://www.gnu.org/software/grep/manual/grep.html)
- [GNU sed manual](https://www.gnu.org/software/sed/manual/sed.html)
- [The GNU Awk User's Guide](https://www.gnu.org/software/gawk/manual/gawk.html)
- [GNU tar manual](https://www.gnu.org/software/tar/manual/tar.html)
- [RFC 5737: indirizzi IPv4 riservati alla documentazione](https://www.rfc-editor.org/rfc/rfc5737)
- [passwd(5), man page](https://man7.org/linux/man-pages/man5/passwd.5.html)
