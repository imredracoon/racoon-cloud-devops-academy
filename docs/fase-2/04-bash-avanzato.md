# 2.4 Bash: cicli, funzioni e robustezza

!!! obiettivi "Obiettivi"
    - Ripetere operazioni con `for`, `while` e `until` e leggere file riga per riga.
    - Organizzare il codice in funzioni.
    - Usare gli array.
    - Rendere uno script robusto: `set -euo pipefail`, `trap`, file temporanei, log.
    - Gestire le opzioni con `getopts`.

## I cicli

### for

```bash
for server in web01 web02 db01; do
    echo "Controllo $server"
done

for f in /var/log/*.log; do          # i caratteri jolly producono un elenco di file
    [[ -e "$f" ]] || continue        # se non ci sono file, il modello resta letterale: lo saltiamo
    echo "$f: $(wc -l < "$f") righe"
done

for i in {1..5}; do echo "$i"; done        # sequenza
for (( i = 0; i < 3; i++ )); do echo "$i"; done   # forma aritmetica, come in C
```

### while e until

```bash
tentativi=0
until curl -fs http://localhost:8000/health > /dev/null; do
    (( tentativi++ ))
    if (( tentativi >= 10 )); then
        echo "Il servizio non risponde" >&2
        exit 1
    fi
    sleep 2
done
echo "Servizio pronto"
```

`while` ripete **finché** la condizione è vera, `until` **finché è falsa**. `break` esce dal ciclo, `continue`
passa all'iterazione successiva.

### Leggere un file riga per riga

```bash
while IFS= read -r riga; do
    echo "Riga: $riga"
done < elenco.txt
```

È la forma corretta: `IFS=` conserva gli spazi iniziali e finali, `-r` conserva le barre rovesciate, e il file
viene passato con `<` al ciclo. Evita `for riga in $(cat file)`: divide sugli spazi, non sulle righe.

Lo stesso schema funziona con l'output di un comando:

```bash
while IFS=: read -r utente _ uid _; do
    (( uid >= 1000 )) && echo "$utente ($uid)"
done < /etc/passwd
```

Qui `IFS=:` divide ogni riga sui due punti e `_` raccoglie i campi che non ci interessano.

## Le funzioni

```bash
log() {
    local livello=$1; shift
    printf '%s [%s] %s\n' "$(date '+%F %T')" "$livello" "$*" >&2
}

ha_comando() {
    command -v "$1" > /dev/null
}

log INFO "Avvio"
ha_comando docker || log WARN "docker non è installato"
```

- Una funzione si definisce **prima** di usarla.
- Riceve argomenti come uno script: `$1`, `$2`, `"$@"`, `$#`.
- `local` rende una variabile visibile solo dentro la funzione: usalo **sempre**, altrimenti le variabili sono
  globali e si sovrascrivono a vicenda.
- `return N` imposta il codice di uscita della funzione (senza `return`, vale quello dell'ultimo comando).
- Per "restituire" un valore, la funzione lo **stampa** e chi la chiama lo cattura con `$( )`.

```bash
data_backup() { date +%Y%m%d-%H%M%S; }
archivio="backup-$(data_backup).tar.gz"
```

`printf` è preferibile a `echo` quando il testo è variabile: il suo comportamento è identico ovunque.

## Gli array

```bash
servizi=(nginx ssh cron)
echo "${servizi[0]}"            # nginx
echo "${#servizi[@]}"           # 3, il numero di elementi
servizi+=(chrony)               # aggiunge
for s in "${servizi[@]}"; do systemctl is-active --quiet "$s" && echo "$s attivo"; done

declare -A porte=([web]=80 [ssh]=22)    # array associativo (chiave → valore)
echo "${porte[web]}"
for k in "${!porte[@]}"; do echo "$k → ${porte[$k]}"; done
```

Gli array sono il modo giusto di costruire l'elenco di argomenti di un comando:

```bash
opzioni=(-a -v --delete)
[[ "$prova" == si ]] && opzioni+=(--dry-run)
rsync "${opzioni[@]}" sorgente/ destinazione/
```

## Script robusti

Per impostazione predefinita Bash è **permissivo**: se un comando fallisce, passa al successivo come se nulla
fosse. In uno script di automazione è pericoloso: un `cd` fallito seguito da un `rm -rf *` cancella la cartella
sbagliata. Per questo gli script professionali iniziano così:

```bash
#!/usr/bin/env bash
set -euo pipefail
```

| Opzione | Effetto |
|---|---|
| `-e` | esce appena un comando fallisce (con alcune eccezioni: dentro `if`, `while`, a sinistra di `&&` e `||`) |
| `-u` | considera un errore usare una variabile non definita (spesso un errore di battitura) |
| `-o pipefail` | una pipeline fallisce se fallisce **uno qualsiasi** dei comandi, non solo l'ultimo |

!!! warning "set -e non sostituisce il ragionamento"
    `-e` ha eccezioni poco intuitive (per esempio dentro le funzioni chiamate in un `if`). Usalo come rete di
    sicurezza, ma controlla esplicitamente gli errori dei comandi importanti. Con `-u`, per leggere un argomento
    che può mancare usa `${1:-}`.

### trap: fare pulizia in ogni caso

`trap` esegue un comando quando lo script riceve un segnale o termina:

```bash
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT        # eseguito all'uscita, sia normale sia per errore

trap 'echo "Errore alla riga $LINENO" >&2' ERR
```

### File temporanei sicuri

Non usare nomi fissi come `/tmp/dati.txt`: un altro utente potrebbe crearlo prima di te, oppure due esecuzioni
contemporanee si sovrascriverebbero. Usa sempre `mktemp`:

```bash
tmp=$(mktemp)                 # file con nome casuale e permessi 600
tmpdir=$(mktemp -d)           # directory
```

### Evitare due esecuzioni contemporanee

```bash
exec 9> /run/lock/backup.lock     # apre il file di lock sul descrittore 9
flock -n 9 || { echo "Già in esecuzione" >&2; exit 1; }
```

### Registrare nei log di sistema

```bash
logger -t backup "Backup completato"     # finisce nel journal: journalctl -t backup
```

## Le opzioni con getopts

```bash
uso() {
    cat >&2 <<EOF
Uso: $0 [-v] [-d DESTINAZIONE] SORGENTE
  -d  cartella di destinazione (predefinita: /backup)
  -v  output dettagliato
  -h  questo aiuto
EOF
    exit 2
}

dest=/backup
verbose=false
while getopts ":d:vh" opt; do
    case "$opt" in
        d) dest=$OPTARG ;;
        v) verbose=true ;;
        h) uso ;;
        :) echo "L'opzione -$OPTARG richiede un valore" >&2; uso ;;
        *) echo "Opzione sconosciuta: -$OPTARG" >&2; uso ;;
    esac
done
shift $(( OPTIND - 1 ))       # rimuove le opzioni: ora $1 è SORGENTE
[[ $# -eq 1 ]] || uso
```

Nella stringa `":d:vh"`, i due punti dopo una lettera indicano che l'opzione vuole un valore; quelli iniziali
attivano la gestione silenziosa degli errori (i casi `:` e `*`). Il blocco `<<EOF ... EOF` è un **here document**:
passa più righe di testo come input a un comando.

## Lo scheletro di uno script professionale

```bash title="modello.sh"
#!/usr/bin/env bash
# Descrizione: cosa fa lo script.
# Uso: modello.sh [-v] ARGOMENTO
set -euo pipefail

readonly SCRIPT=$(basename "$0")

log() { printf '%s %s: %s\n' "$(date '+%F %T')" "$SCRIPT" "$*" >&2; }
die() { log "ERRORE: $*"; exit 1; }

cleanup() {
    [[ -n "${tmpdir:-}" ]] && rm -rf "$tmpdir"
}
trap cleanup EXIT

main() {
    [[ $# -ge 1 ]] || die "manca l'argomento"
    tmpdir=$(mktemp -d)
    log "lavoro in $tmpdir su $1"
    # ...
}

main "$@"
```

Mettere la logica in una funzione `main` chiamata in fondo garantisce che lo script venga letto tutto prima di
essere eseguito, e lo rende più ordinato.

## Esercizi

!!! esercizio "Esercizio 1 · Backup con rotazione"
    Scrivi `backup.sh SORGENTE DESTINAZIONE` che:

    1. usi `set -euo pipefail` e verifichi che la sorgente sia una directory;
    2. crei in DESTINAZIONE un archivio `NOME-AAAAMMGG-HHMMSS.tar.gz` (dove NOME è il nome della cartella sorgente);
    3. tenga solo i **5** archivi più recenti di quella sorgente, cancellando i più vecchi;
    4. registri l'esito nel journal con `logger -t backup`.

    Provalo eseguendolo sette volte su `~/appunti-git` e controlla che restino cinque archivi.

??? success "Soluzione"
    ```bash
    #!/usr/bin/env bash
    # Backup di una directory con rotazione degli archivi.
    set -euo pipefail

    readonly CONSERVA=5

    die() { echo "ERRORE: $*" >&2; logger -t backup "ERRORE: $*"; exit 1; }

    [[ $# -eq 2 ]] || die "uso: $0 SORGENTE DESTINAZIONE"
    sorgente=${1%/}                   # toglie una eventuale / finale
    dest=$2
    [[ -d "$sorgente" ]] || die "$sorgente non è una directory"
    mkdir -p "$dest"

    nome=$(basename "$sorgente")
    archivio="$dest/$nome-$(date +%Y%m%d-%H%M%S).tar.gz"

    tar -czf "$archivio" -C "$(dirname "$sorgente")" "$nome"

    # elenco dal più recente al più vecchio, salta i primi 5, cancella il resto
    mapfile -t vecchi < <(ls -1t "$dest/$nome"-*.tar.gz | tail -n +$(( CONSERVA + 1 )))
    for f in "${vecchi[@]}"; do
        rm -f -- "$f"
    done

    logger -t backup "creato $archivio, rimossi ${#vecchi[@]} archivi vecchi"
    echo "$archivio"
    ```
    ```bash
    for i in {1..7}; do ./backup.sh ~/appunti-git ~/backup; sleep 1; done
    ls ~/backup               # 5 archivi
    journalctl -t backup -n 7
    ```
    **Spiegazione.** `${1%/}` è un'**espansione di parametro** che rimuove `/` dalla fine. `tar -C` entra nella
    cartella madre, così l'archivio contiene il percorso relativo `appunti-git/...` e non quello assoluto.
    `mapfile -t` legge le righe in un array; `< <(...)` è una **sostituzione di processo**, che passa l'output di un
    comando come se fosse un file. `--` dopo `rm` indica la fine delle opzioni: protegge da nomi che iniziano con `-`.
    Qui `ls` è accettabile perché i nomi degli archivi li generiamo noi e non contengono spazi.

!!! esercizio "Esercizio 2 · Verifica di più siti"
    Crea `siti.txt` con un URL per riga (righe vuote e righe che iniziano con `#` vanno ignorate). Scrivi
    `verifica-siti.sh` che, per ogni URL, stampi il codice HTTP e il tempo di risposta, e che alla fine esca con 1
    se almeno un sito non ha risposto con un codice 2xx o 3xx. Aggiungi un'opzione `-t SECONDI` per il timeout
    (predefinito 5).

??? success "Soluzione"
    ```bash
    #!/usr/bin/env bash
    set -euo pipefail

    timeout=5
    while getopts ":t:" opt; do
        case "$opt" in
            t) timeout=$OPTARG ;;
            *) echo "Uso: $0 [-t SECONDI] FILE" >&2; exit 2 ;;
        esac
    done
    shift $(( OPTIND - 1 ))
    file=${1:-siti.txt}

    errori=0
    while IFS= read -r url; do
        [[ -z "$url" || "$url" == \#* ]] && continue
        risultato=$(curl -s -o /dev/null -m "$timeout" -w '%{http_code} %{time_total}' "$url" || true)
        codice=${risultato%% *}
        tempo=${risultato#* }
        printf '%-40s %s %ss\n' "$url" "$codice" "$tempo"
        [[ "$codice" =~ ^[23] ]] || (( ++errori ))
    done < "$file"

    (( errori == 0 )) || { echo "$errori siti in errore" >&2; exit 1; }
    ```
    **Spiegazione.** `curl -w` stampa variabili come il codice HTTP e il tempo totale. `|| true` impedisce che
    `set -e` fermi lo script quando un sito non risponde (curl restituisce un errore). `${risultato%% *}` toglie tutto
    dal primo spazio in poi; `${risultato#* }` toglie tutto fino al primo spazio. Usiamo `(( ++errori ))` e non
    `(( errori++ ))`: quest'ultimo vale 0 la prima volta, e `(( 0 ))` ha codice di uscita 1, che con `set -e`
    farebbe terminare lo script.

!!! esercizio "Esercizio 3 · Caccia ai bug"
    Lo script seguente ha almeno quattro problemi. Trovali con l'aiuto di ShellCheck e correggili.

    ```bash
    #!/bin/sh
    cartella=$1
    cd $cartella
    for f in $(ls *.log); do
        rm $f
    done
    echo "Puliti i log in $cartella"
    ```

??? success "Soluzione"
    ```bash
    #!/usr/bin/env bash
    set -euo pipefail
    cartella=${1:?Uso: $0 CARTELLA}
    cd -- "$cartella"
    for f in ./*.log; do
        [[ -e "$f" ]] || continue
        rm -- "$f"
    done
    echo "Puliti i log in $cartella"
    ```
    **Spiegazione.** 1) Senza controllo dell'argomento, `cd` senza parametri va nella home e lo script cancella i
    log lì. 2) Se `cd` fallisce, lo script prosegue nella cartella sbagliata: `set -e` lo ferma. 3) `$(ls *.log)`
    si rompe con nomi che contengono spazi: si usa direttamente il modello `./*.log`. 4) `$cartella` e `$f` senza
    virgolette. 5) Lo shebang `/bin/sh` su Ubuntu è dash: con le funzioni di Bash sarebbe un problema.

## Nei tre livelli

!!! esercizio "Esercizio per i livelli 2 e 3 · fleet.sh"
    Scrivi `fleet.sh`, uno script che esegue lo stesso comando su più macchine e riassume l'esito. Le macchine sono
    elencate in un file, una per riga, con i nomi del tuo `~/.ssh/config`; le righe vuote e quelle che iniziano
    con `#` si saltano.

    === "Livello 2 · Workstation + server"

        ```text title="nodi.txt"
        srv1
        lab
        client
        ```

    === "Livello 3 · Mini data center"

        ```text title="nodi.txt"
        nodo1
        nodo2
        nodo3
        lab
        ```

    Requisiti: opzione `-f FILE` (predefinito `nodi.txt`) e `-h`; una riga per macchina con `OK` o `ERRORE` e
    l'output del comando; codice di uscita 1 se anche una sola macchina ha fallito; nessuna richiesta di password
    (se la chiave non basta, la macchina è in errore). Deve passare `shellcheck`.

    **Collaudo:** `./fleet.sh uptime` mostra tutte le macchine in `OK`. Spegni una VM (`virsh shutdown client` sul
    server) o un nodo: lo script la segnala in `ERRORE` dopo al massimo 5 secondi e `echo $?` stampa `1`.

??? success "Soluzione"
    ```bash title="fleet.sh"
    #!/usr/bin/env bash
    # Esegue un comando su tutte le macchine elencate in un file e riassume l'esito.
    # Uso: fleet.sh [-f FILE] COMANDO...
    set -euo pipefail

    file=nodi.txt
    while getopts "f:h" opt; do
      case $opt in
        f) file=$OPTARG ;;
        h) sed -n '2,3p' "$0"; exit 0 ;;
        *) exit 2 ;;
      esac
    done
    shift $((OPTIND - 1))
    [ $# -gt 0 ] || { echo "manca il comando da eseguire" >&2; exit 2; }
    [ -r "$file" ] || { echo "non riesco a leggere $file" >&2; exit 2; }

    falliti=0
    while read -r nodo; do
      [[ -z $nodo || $nodo == \#* ]] && continue
      if out=$(ssh -n -o BatchMode=yes -o ConnectTimeout=5 "$nodo" "$@" 2>&1); then
        printf '%-8s OK     %s\n' "$nodo" "$out"
      else
        printf '%-8s ERRORE %s\n' "$nodo" "$out"
        falliti=$((falliti + 1))
      fi
    done < "$file"

    exit $((falliti > 0))
    ```
    **Spiegazione.** `ssh -n` è indispensabile: senza, `ssh` leggerebbe il resto di `nodi.txt` dal suo input
    standard e il ciclo si fermerebbe alla prima macchina. `BatchMode=yes` vieta le richieste di password,
    `ConnectTimeout=5` evita di aspettare a lungo una macchina spenta. L'`if` intorno a `ssh` impedisce a
    `set -e` di fermare lo script al primo errore: è uno dei casi in cui si gestisce l'errore invece di subirlo.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è il modo corretto di leggere un file riga per riga?","options":["for r in $(cat file)","while IFS= read -r r; do ...; done < file","cat file | for r","read file"],"answer":[1],"explanation":"for con $(cat) divide sugli spazi e non sulle righe."},
 {"q":"Cosa fa `set -o pipefail`?","options":["Ferma lo script al primo errore","Fa fallire una pipeline se fallisce uno qualsiasi dei suoi comandi","Vieta le pipe","Stampa i comandi eseguiti"],"answer":[1],"explanation":"Senza pipefail conta solo il codice dell'ultimo comando della pipeline."},
 {"q":"Perché si usa `local` dentro le funzioni?","options":["Per rendere la funzione più veloce","Per limitare la variabile alla funzione ed evitare di sovrascrivere variabili globali","Per esportare la variabile","È obbligatorio"],"answer":[1],"explanation":"Senza local, le variabili di Bash sono globali."},
 {"q":"Qual è il modo sicuro di creare un file temporaneo?","options":["touch /tmp/dati","mktemp","echo > /tmp/$$","mkdir /tmp/x"],"answer":[1],"explanation":"mktemp crea un nome casuale con permessi 600, evitando collisioni e attacchi."},
 {"q":"A cosa serve `trap cleanup EXIT`?","options":["A eseguire cleanup quando lo script termina, anche per errore","A ignorare gli errori","A uscire subito","A catturare l'output"],"answer":[0],"explanation":"È il modo standard di garantire la pulizia delle risorse."},
 {"q":"Come si espandono tutti gli elementi di un array, conservando gli spazi?","options":["$array","${array}","\"${array[@]}\"","${#array[@]}"],"answer":[2],"explanation":"${#array[@]} restituisce il numero di elementi."}
]}
</script>

## Fonti

- [Bash manual: Looping Constructs](https://www.gnu.org/software/bash/manual/html_node/Looping-Constructs.html)
- [Bash manual: Shell Functions](https://www.gnu.org/software/bash/manual/html_node/Shell-Functions.html)
- [Bash manual: Arrays](https://www.gnu.org/software/bash/manual/html_node/Arrays.html)
- [Bash manual: The Set Builtin](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html)
- [Bash manual: Shell Parameter Expansion](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html)
- [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html)
- [mktemp(1)](https://man7.org/linux/man-pages/man1/mktemp.1.html), [flock(1)](https://man7.org/linux/man-pages/man1/flock.1.html), [logger(1)](https://man7.org/linux/man-pages/man1/logger.1.html)
