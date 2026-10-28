# 2.3 Bash: variabili e condizioni

!!! obiettivi "Obiettivi"
    - Scrivere ed eseguire uno script Bash.
    - Usare variabili, quoting, sostituzione di comando ed espansioni.
    - Leggere argomenti e input.
    - Capire il codice di uscita e prendere decisioni con `if`, `test` e `case`.

## Cos'è uno script

Uno **script** è un file di testo che contiene comandi della shell, eseguiti in sequenza. Tutto quello che scrivi
nel terminale può andare in uno script. È il primo passo dell'automazione: scrivi una volta, esegui mille volte.

```bash title="ciao.sh"
#!/usr/bin/env bash
# Il mio primo script
echo "Ciao, $USER! Oggi è $(date +%A)."
```

- La prima riga, lo **shebang** `#!`, dice al sistema quale interprete usare. `#!/usr/bin/env bash` cerca `bash`
  nel `PATH` ed è la forma più portabile.
- Le righe che iniziano con `#` sono **commenti**.

Per eseguirlo:

```bash
chmod +x ciao.sh      # rende il file eseguibile (lezione 1.3)
./ciao.sh             # ./ perché la cartella corrente non è nel PATH
bash ciao.sh          # alternativa: passa il file all'interprete, non serve il permesso x
```

!!! info "sh, bash e le altre shell"
    `sh` è la shell POSIX minimale; su Debian e Ubuntu `/bin/sh` è **dash**, su RHEL è un collegamento a **bash**.
    Uno script con `#!/bin/sh` che usa funzioni tipiche di Bash può quindi funzionare su RHEL e fallire su Ubuntu.
    In questo corso scriviamo script **Bash** e lo dichiariamo sempre nello shebang.

## Variabili

```bash
nome="Vincenzo"          # NESSUNO spazio attorno a =
eta=30
echo "Nome: $nome, età: ${eta}"
```

- I nomi distinguono maiuscole e minuscole. Per convenzione le variabili d'ambiente e le costanti sono in
  MAIUSCOLO, quelle dello script in minuscolo.
- `${nome}` e `$nome` sono equivalenti; le graffe servono quando subito dopo c'è altro testo: `"${nome}_backup"`.
- `nome = "x"` (con spazi) è un **errore**: Bash cerca un comando chiamato `nome`.

### Variabili d'ambiente

Una variabile normale esiste solo nella shell corrente. Con `export` diventa una **variabile d'ambiente** e la
ereditano i programmi avviati da quella shell:

```bash
export APP_ENV=sviluppo
python3 -c 'import os; print(os.environ["APP_ENV"])'     # sviluppo
env | grep APP_ENV
APP_ENV=prova ./script.sh        # imposta la variabile solo per quel comando
```

Variabili predefinite utili: `$HOME`, `$USER`, `$PATH`, `$PWD`, `$HOSTNAME`, `$RANDOM`.

### Sostituzione di comando

`$(comando)` viene sostituito con l'output del comando:

```bash
oggi=$(date +%F)                  # 2026-10-28
file_log=$(ls /var/log | wc -l)
echo "Backup del $oggi"
```

Esiste anche la vecchia forma con gli apici inversi `` `comando` ``: evitala, è meno leggibile e non si annida bene.

### Aritmetica

```bash
a=7; b=3
echo $(( a + b ))     # 10
echo $(( a / b ))     # 2: Bash lavora solo con numeri interi
echo $(( a % b ))     # 1, il resto
(( a++ ))             # incrementa a
```

Per i decimali serve un altro strumento (`bc`, `awk` o Python).

## Il quoting: la regola più importante

Bash, prima di eseguire un comando, **espande** le variabili e poi **divide** il risultato in parole sugli spazi.
Le virgolette controllano questo comportamento.

| Forma | Effetto |
|---|---|
| `"testo $var"` | **virgolette doppie**: le variabili e `$( )` vengono espansi, ma il risultato resta **una sola parola** |
| `'testo $var'` | **apici singoli**: tutto letterale, niente espansioni |
| `$var` senza virgolette | espansa e poi **divisa** sugli spazi, e i caratteri jolly (`*`) vengono espansi |

```bash
file="relazione finale.txt"
touch "$file"         # crea UN file: "relazione finale.txt"
rm $file              # ERRORE: cerca due file, "relazione" e "finale.txt"
rm "$file"            # corretto
```

!!! warning "Metti sempre le variabili tra virgolette doppie"
    `"$variabile"`, sempre, a meno che tu non voglia esplicitamente la divisione in parole. La maggior parte dei
    bug negli script Bash nasce da variabili senza virgolette e da nomi di file con spazi.

## Argomenti e input

Gli argomenti passati allo script sono disponibili in variabili speciali:

| Variabile | Contenuto |
|---|---|
| `$0` | il nome dello script |
| `$1`, `$2`, ... | il primo, il secondo argomento... (`${10}` dal decimo) |
| `$#` | il numero di argomenti |
| `"$@"` | tutti gli argomenti, ciascuno come parola separata (la forma corretta da usare) |
| `$?` | il codice di uscita dell'ultimo comando |
| `$$` | il PID dello script |

```bash title="saluta.sh"
#!/usr/bin/env bash
echo "Script: $0, argomenti ricevuti: $#"
echo "Primo: $1"
for arg in "$@"; do echo "- $arg"; done
```

```bash
./saluta.sh uno "due tre"     # 2 argomenti: "uno" e "due tre"
```

Per chiedere un dato all'utente:

```bash
read -rp "Come ti chiami? " nome
read -rsp "Password: " pw; echo     # -s non mostra ciò che digiti
```

L'opzione `-r` evita che le barre rovesciate vengano interpretate: usala sempre.

### Valori predefiniti

```bash
porta=${1:-8000}           # se $1 è vuoto o assente, usa 8000
: "${API_URL:?Variabile API_URL non impostata}"   # se manca, esce con errore
```

## Il codice di uscita

Ogni comando termina con un **codice di uscita** (*exit status*), un numero da 0 a 255:

- **0** significa **successo**;
- qualsiasi altro valore significa **errore** (il significato dipende dal comando).

```bash
ls /etc > /dev/null; echo $?          # 0
ls /inesistente 2> /dev/null; echo $? # 2
grep -q root /etc/passwd; echo $?     # 0 se trovato, 1 se no
```

Uno script restituisce il codice con `exit N`; senza `exit`, restituisce quello dell'ultimo comando.

### && e ||

```bash
mkdir -p backup && echo "cartella pronta"     # il secondo solo se il primo riesce
ping -c1 -W1 10.0.0.1 > /dev/null || echo "host irraggiungibile"   # solo se fallisce
```

## Le condizioni: if

```bash
if comando; then
    # eseguito se il comando restituisce 0
elif altro_comando; then
    # ...
else
    # ...
fi
```

`if` non valuta "vero o falso": esegue un comando e guarda il suo **codice di uscita**. Per confrontare valori si usa
il comando di test `[[ ... ]]` (specifico di Bash, più sicuro del vecchio `[ ... ]`):

```bash
if [[ -f /etc/os-release ]]; then
    . /etc/os-release            # carica le variabili del file nella shell
    echo "Distribuzione: $NAME"
fi
```

!!! warning "Gli spazi dentro le parentesi sono obbligatori"
    `[[ -f file ]]` è corretto, `[[-f file]]` è un errore. `[[` è una parola come le altre e va separata.

### Gli operatori di test

| File | Vero se |
|---|---|
| `-e f` | esiste |
| `-f f` | è un file regolare |
| `-d f` | è una directory |
| `-r f` / `-w f` / `-x f` | è leggibile / scrivibile / eseguibile |
| `-s f` | esiste e non è vuoto |

| Stringhe | Vero se |
|---|---|
| `-z "$s"` | è vuota |
| `-n "$s"` | non è vuota |
| `"$a" == "$b"` | sono uguali (dentro `[[ ]]` a destra puoi usare un modello: `== *.log`) |
| `"$a" != "$b"` | sono diverse |
| `"$s" =~ ^[0-9]+$` | corrisponde all'espressione regolare |

| Numeri | Vero se |
|---|---|
| `-eq`, `-ne` | uguali, diversi |
| `-lt`, `-le` | minore, minore o uguale |
| `-gt`, `-ge` | maggiore, maggiore o uguale |

Per i numeri puoi usare anche la forma aritmetica, più leggibile: `if (( uso > 80 )); then`.

Le condizioni si combinano con `&&`, `||` e `!` dentro `[[ ]]`:

```bash
if [[ -d "$dir" && -w "$dir" ]]; then echo "posso scrivere in $dir"; fi
if [[ ! -f "$conf" ]]; then echo "manca $conf"; fi
```

### case: scegliere tra molte opzioni

```bash
case "$1" in
    start)        echo "avvio" ;;
    stop|halt)    echo "arresto" ;;
    *.tar.gz)     echo "archivio" ;;
    *)            echo "uso: $0 {start|stop}"; exit 1 ;;
esac
```

`case` confronta il valore con dei **modelli** (con `*`, `?`, `|`); `*)` è il caso predefinito. È perfetto per
gli script con sotto-comandi.

## Uno script completo

```bash title="controlla-disco.sh"
#!/usr/bin/env bash
# Avvisa se l'uso del filesystem supera una soglia.
# Uso: controlla-disco.sh [PERCORSO] [SOGLIA]

percorso=${1:-/}
soglia=${2:-80}

if [[ ! -d "$percorso" ]]; then
    echo "Errore: $percorso non è una directory" >&2
    exit 2
fi

if [[ ! "$soglia" =~ ^[0-9]+$ ]] || (( soglia > 100 )); then
    echo "Errore: la soglia deve essere un numero tra 0 e 100" >&2
    exit 2
fi

uso=$(df --output=pcent "$percorso" | tail -n 1 | tr -dc '0-9')

if (( uso >= soglia )); then
    echo "ATTENZIONE: $percorso è al ${uso}% (soglia ${soglia}%)"
    exit 1
else
    echo "OK: $percorso è al ${uso}%"
fi
```

Nota i dettagli professionali: i messaggi di errore vanno su **standard error** (`>&2`), gli argomenti vengono
**validati** e lo script usa **codici di uscita** diversi per "problema trovato" (1) e "uso sbagliato" (2), così
un altro programma può distinguerli.

## Controllare gli script con ShellCheck

[ShellCheck](https://www.shellcheck.net/) è un analizzatore che trova errori comuni: variabili senza virgolette,
confronti sbagliati, comandi pericolosi. Usalo su **ogni** script.

```bash
sudo apt install -y shellcheck      # Debian/Ubuntu (nel Dev Container è nel repository Ubuntu)
sudo dnf install -y ShellCheck      # RHEL e derivate: dal repository EPEL
shellcheck controlla-disco.sh
```

## Esercizi

!!! esercizio "Esercizio 1 · Informazioni di sistema"
    Scrivi `info.sh` che stampi: nome host, distribuzione (dalla variabile `PRETTY_NAME` di `/etc/os-release`),
    versione del kernel, tempo di attività, numero di utenti con una shell di login reale (che non finisce con
    `nologin` o `false`). Deve passare ShellCheck senza avvisi.

??? success "Soluzione"
    ```bash
    #!/usr/bin/env bash
    # Mostra un riepilogo del sistema.
    # shellcheck source=/dev/null
    . /etc/os-release

    utenti=$(grep -Evc '(nologin|false)$' /etc/passwd)

    echo "Host:          $(hostname)"
    echo "Sistema:       $PRETTY_NAME"
    echo "Kernel:        $(uname -r)"
    echo "Attivo da:     $(uptime -p)"
    echo "Utenti login:  $utenti"
    ```
    **Spiegazione.** `.` (o `source`) esegue `/etc/os-release` nella shell corrente, che così acquisisce variabili come
    `PRETTY_NAME`. `grep -v` inverte la ricerca, `-c` conta le righe, `-E` abilita l'alternativa `|`. Il commento
    `# shellcheck source=/dev/null` dice a ShellCheck di non provare a leggere il file incluso.

!!! esercizio "Esercizio 2 · Validare gli argomenti"
    Scrivi `crea-utente.sh NOME` che:

    - esca con codice 2 e un messaggio d'uso su standard error se non riceve esattamente un argomento;
    - esca con codice 2 se il nome non è composto solo da lettere minuscole e cifre e non inizia con una lettera;
    - esca con codice 1 se l'utente esiste già;
    - altrimenti stampi `Creerei l'utente NOME` (non crearlo davvero).

??? success "Soluzione"
    ```bash
    #!/usr/bin/env bash
    if (( $# != 1 )); then
        echo "Uso: $0 NOME" >&2
        exit 2
    fi

    nome=$1

    if [[ ! "$nome" =~ ^[a-z][a-z0-9]*$ ]]; then
        echo "Nome non valido: $nome" >&2
        exit 2
    fi

    if id "$nome" &> /dev/null; then
        echo "L'utente $nome esiste già" >&2
        exit 1
    fi

    echo "Creerei l'utente $nome"
    ```
    **Spiegazione.** `id` restituisce 0 se l'utente esiste: lo usiamo direttamente come condizione, senza
    confrontare output. `&> /dev/null` scarta sia l'output sia gli errori.

!!! esercizio "Esercizio 3 · Un piccolo gestore di servizio"
    Scrivi `servizio.sh` che accetti un sotto-comando (`stato`, `avvia`, `ferma`) e il nome di un servizio, e che
    chiami il corrispondente `systemctl` (`status`, `start`, `stop`). Con un sotto-comando sconosciuto deve stampare
    l'uso ed uscire con 2. Usa `case`.

??? success "Soluzione"
    ```bash
    #!/usr/bin/env bash
    uso() { echo "Uso: $0 {stato|avvia|ferma} SERVIZIO" >&2; exit 2; }

    [[ $# -eq 2 ]] || uso
    servizio=$2

    case "$1" in
        stato) systemctl status "$servizio" --no-pager ;;
        avvia) sudo systemctl start "$servizio" ;;
        ferma) sudo systemctl stop "$servizio" ;;
        *)     uso ;;
    esac
    ```
    **Spiegazione.** `uso()` è una **funzione**: le studierai nella prossima lezione. `[[ ... ]] || uso` è una forma
    compatta di `if`: se il test fallisce, chiama `uso`. Lo script restituisce il codice di uscita di `systemctl`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale assegnazione è corretta in Bash?","options":["nome = \"Ada\"","nome=\"Ada\"","$nome=\"Ada\"","set nome \"Ada\""],"answer":[1],"explanation":"Nessuno spazio attorno a = e niente $ a sinistra."},
 {"q":"Cosa stampa `echo '$HOME'`?","options":["Il percorso della home","La stringa $HOME","Niente","Un errore"],"answer":[1],"explanation":"Tra apici singoli non avviene alcuna espansione."},
 {"q":"Quale codice di uscita indica successo?","options":["1","0","255","-1"],"answer":[1],"explanation":"0 è successo, qualsiasi altro valore è errore."},
 {"q":"Quale variabile contiene tutti gli argomenti come parole separate?","options":["$*","\"$@\"","$#","$0"],"answer":[1],"explanation":"\"$@\" conserva ogni argomento, anche se contiene spazi."},
 {"q":"Come confronti due numeri interi dentro [[ ]]?","options":["[[ $a > $b ]]","[[ $a -gt $b ]]","[[ $a => $b ]]","[[ $a gt $b ]]"],"answer":[1],"explanation":"Dentro [[ ]] il simbolo > confronta stringhe in ordine alfabetico, non numeri. In alternativa: (( a > b ))."},
 {"q":"Perché uno script dovrebbe scrivere gli errori con `>&2`?","options":["Per colorarli","Per inviarli su standard error, separati dall'output normale","Per scriverli in un file di log","Per terminare lo script"],"answer":[1],"explanation":"Chi usa lo script può così ridirigere output ed errori in modo indipendente."}
]}
</script>

## Fonti

- [GNU Bash Reference Manual](https://www.gnu.org/software/bash/manual/bash.html)
- [Bash manual: Quoting](https://www.gnu.org/software/bash/manual/html_node/Quoting.html)
- [Bash manual: Conditional Constructs](https://www.gnu.org/software/bash/manual/html_node/Conditional-Constructs.html)
- [Bash manual: Bash Conditional Expressions](https://www.gnu.org/software/bash/manual/html_node/Bash-Conditional-Expressions.html)
- [ShellCheck](https://www.shellcheck.net/) e il suo [wiki degli errori](https://www.shellcheck.net/wiki/)
- [Ubuntu Wiki: DashAsBinSh](https://wiki.ubuntu.com/DashAsBinSh)
