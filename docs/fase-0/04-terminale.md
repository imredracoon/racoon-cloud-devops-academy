# 0.4 Il terminale: primi passi

!!! obiettivi "Obiettivi"
    - Capire cosa sono terminale, shell e prompt.
    - Leggere la struttura di un comando: comando, opzioni, argomenti.
    - Muoversi tra le cartelle e creare, copiare, spostare e cancellare file.
    - Trovare aiuto da soli con `man` e `--help`.

## Perché il terminale

I server non hanno quasi mai un'interfaccia grafica: si gestiscono scrivendo comandi. Il terminale è più veloce
del mouse, permette di **ripetere** e **automatizzare** le operazioni e funziona anche attraverso una connessione
remota lenta. È lo strumento principale di questo mestiere.

## Terminale, shell e prompt

- Il **terminale** è la finestra in cui scrivi i comandi e leggi i risultati.
- La **shell** è il programma che legge i comandi, li interpreta e li esegue. Su Linux la più diffusa è **Bash**
  (*Bourne Again SHell*); su macOS la shell predefinita è **zsh**, molto simile per quello che vedremo.
- Il **prompt** è il testo che la shell mostra quando aspetta un comando, per esempio:

```text
ubuntu@lab:~$
```

| Parte | Significato |
|---|---|
| `ubuntu` | il nome dell'utente |
| `lab` | il nome del computer (*hostname*) |
| `~` | la cartella in cui ti trovi; `~` indica la tua cartella personale (*home*) |
| `$` | sei un utente normale; per l'amministratore (`root`) compare `#` |

## Dove esercitarti adesso

Per questa lezione ti serve un terminale Linux. Hai tre possibilità:

1. **Killercoda** (consigliato ora, zero installazioni): apri il
   [playground Ubuntu](https://killercoda.com/playgrounds/scenario/ubuntu) e accedi con un account gratuito.
2. **macOS**: apri l'app **Terminale**. Quasi tutti i comandi di questa lezione funzionano uguali.
3. **Windows**: aspetta la lezione 0.7, dove installerai la VM Linux; oppure usa Killercoda.

## La struttura di un comando

```text
ls -l /etc
│   │  └── argomento: su cosa lavora il comando
│   └───── opzione: modifica il comportamento (inizia con - o --)
└───────── comando: il programma da eseguire
```

- Le **opzioni brevi** sono una lettera con un trattino (`-l`) e si possono unire: `ls -la` equivale a `ls -l -a`.
- Le **opzioni lunghe** sono parole con due trattini: `ls --all`.
- Linux distingue **maiuscole e minuscole**: `File.txt` e `file.txt` sono due file diversi, `LS` non è un comando.

## Orientarsi: pwd, ls, cd

```console
$ pwd
/home/ubuntu
```

`pwd` (*print working directory*) mostra la **cartella corrente**: quella in cui ti trovi.

```console
$ ls
Documenti  note.txt
$ ls -l
total 4
drwxr-xr-x 2 ubuntu ubuntu 4096 Oct 28 10:00 Documenti
-rw-r--r-- 1 ubuntu ubuntu   12 Oct 28 10:01 note.txt
$ ls -a
.  ..  .bashrc  Documenti  note.txt
```

`ls` elenca il contenuto della cartella. Con `-l` mostra i dettagli (la `d` iniziale indica una directory, il
resto lo capirai nella lezione 1.3); con `-a` mostra anche i file **nascosti**, il cui nome inizia con un punto.

`cd` (*change directory*) cambia cartella:

```bash
cd Documenti     # entra nella cartella Documenti
cd ..            # sale di un livello: .. è la cartella "genitore"
cd               # senza argomenti torna nella home
cd /etc          # va nella cartella /etc, indicata con il percorso completo
cd -             # torna alla cartella precedente
```

Il testo dopo `#` è un **commento**: la shell lo ignora. Nel corso lo usiamo per spiegare i comandi.

### Percorsi assoluti e relativi

- Un percorso **assoluto** parte dalla radice `/` del filesystem: `/home/ubuntu/Documenti`. Funziona ovunque ti trovi.
- Un percorso **relativo** parte dalla cartella corrente: `Documenti/lettera.txt`. Il punto `.` indica la cartella
  corrente, i due punti `..` quella superiore.

## Creare, copiare, spostare, cancellare

```bash
mkdir progetti                 # crea la cartella progetti
mkdir -p progetti/a/b          # crea anche le cartelle intermedie mancanti
touch note.txt                 # crea un file vuoto (o aggiorna la data se esiste)
echo "prima riga" > note.txt   # scrive testo nel file, sostituendo il contenuto
echo "seconda" >> note.txt     # aggiunge in fondo al file
cat note.txt                   # mostra il contenuto
cp note.txt copia.txt          # copia un file
cp -r progetti backup          # copia una cartella con tutto il contenuto (-r = ricorsivo)
mv copia.txt archivio.txt      # rinomina (o sposta) un file
mv archivio.txt progetti/      # sposta il file nella cartella progetti
rm archivio.txt                # cancella un file (errore: ora è in progetti/)
rm progetti/archivio.txt       # cancella il file giusto
rmdir progetti/a/b             # cancella una cartella vuota
rm -r backup                   # cancella una cartella e tutto il suo contenuto
```

!!! danger "Nel terminale non c'è il cestino"
    `rm` cancella **definitivamente**. In particolare `rm -r` su una cartella sbagliata può distruggere molto
    lavoro, e con i privilegi di amministratore anche il sistema. Prima di premere Invio rileggi sempre il percorso.
    Un'abitudine prudente è usare `rm -ri`, che chiede conferma per ogni file.

## Scorciatoie che fanno risparmiare ore

| Tasto | Effetto |
|---|---|
| ++tab++ | **Completa** nomi di comandi, file e cartelle. Premilo due volte per vedere le alternative. |
| ++arrow-up++ / ++arrow-down++ | Scorre i comandi già eseguiti. |
| ++ctrl+r++ | Cerca nella cronologia dei comandi. |
| ++ctrl+c++ | Interrompe il comando in esecuzione. |
| ++ctrl+l++ | Pulisce lo schermo (come il comando `clear`). |
| ++ctrl+a++ / ++ctrl+e++ | Va all'inizio / alla fine della riga. |

Usa ++tab++ sempre: scrivi meno, e se il completamento non funziona vuol dire che il percorso è sbagliato.

## Trovare aiuto da soli

Nessuno ricorda tutte le opzioni di tutti i comandi. I professionisti sanno **dove cercarle**:

```bash
ls --help        # breve riepilogo delle opzioni
man ls           # il manuale completo; si esce con il tasto q
man -k copy      # cerca nei manuali le pagine che parlano di "copy"
```

Dentro `man` ti muovi con le frecce, cerchi una parola con `/parola` e Invio, passi al risultato successivo con
`n` ed esci con `q`.

!!! note "Il manuale nel Dev Container e nei container minimi"
    Alcune immagini minime non includono le pagine di manuale per risparmiare spazio. La VM di laboratorio del
    corso le ha. In ogni caso le trovi anche online su [man7.org](https://man7.org/linux/man-pages/).

## Esercizi

Svolgili in ordine, nello stesso terminale.

!!! esercizio "Esercizio 1 · Una struttura di cartelle"
    Partendo dalla tua home, crea questa struttura con il minor numero di comandi possibile:

    ```text
    corso/
    ├── fase0/
    │   └── appunti.txt   (contiene la riga: Lezione 0.4)
    └── fase1/
    ```

??? success "Soluzione"
    ```bash
    cd
    mkdir -p corso/fase0 corso/fase1
    echo "Lezione 0.4" > corso/fase0/appunti.txt
    ```
    **Spiegazione.** `mkdir -p` accetta più argomenti e crea anche la cartella `corso` che non esiste ancora.
    `echo ... >` crea il file e ci scrive dentro in un solo passaggio. Per controllare il risultato puoi usare
    `ls -R corso` oppure, se installato, `tree corso`.

!!! esercizio "Esercizio 2 · Percorsi"
    Ti trovi in `/home/ubuntu/corso/fase1`. Scrivi il comando per mostrare il contenuto di `appunti.txt` usando:

    1. un percorso relativo;
    2. un percorso assoluto;
    3. il simbolo `~`.

??? success "Soluzione"
    ```bash
    cat ../fase0/appunti.txt
    cat /home/ubuntu/corso/fase0/appunti.txt
    cat ~/corso/fase0/appunti.txt
    ```
    **Spiegazione.** `..` sale da `fase1` a `corso`, poi si scende in `fase0`. Il percorso assoluto parte dalla
    radice. `~` è una scorciatoia per la home dell'utente corrente.

!!! esercizio "Esercizio 3 · Copia, rinomina, cancella"
    1. Copia `appunti.txt` in `fase1` con il nome `appunti-1.txt`.
    2. Aggiungi in fondo ad `appunti-1.txt` la riga `Lezione 1.1`.
    3. Cancella la cartella `fase0` con tutto il contenuto.
    4. Verifica che resti solo `fase1/appunti-1.txt` con due righe.

??? success "Soluzione"
    ```bash
    cd ~/corso
    cp fase0/appunti.txt fase1/appunti-1.txt
    echo "Lezione 1.1" >> fase1/appunti-1.txt
    rm -r fase0
    ls -R
    cat fase1/appunti-1.txt
    ```
    **Spiegazione.** `>>` aggiunge senza cancellare. Se avessi usato `>`, il file avrebbe contenuto solo la nuova
    riga: è uno degli errori più comuni, fai attenzione.

!!! esercizio "Esercizio 4 · Cerca nel manuale"
    Usando solo `man ls`, trova l'opzione che ordina i file per dimensione e quella che mostra le dimensioni in
    formato leggibile (K, M, G). Poi usale insieme su `/etc`.

??? success "Soluzione"
    ```bash
    ls -lSh /etc
    ```
    **Spiegazione.** Nel manuale, cercando `/size`, trovi `-S` (*sort by file size, largest first*) e `-h`
    (*human-readable*), che ha effetto insieme a `-l`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale comando mostra la cartella in cui ti trovi?","options":["ls","cd","pwd","whoami"],"answer":[2],"explanation":"pwd significa print working directory."},
 {"q":"Cosa fa `cd ..`?","options":["Torna nella home","Sale alla cartella superiore","Cancella la cartella corrente","Mostra i file nascosti"],"answer":[1],"explanation":".. indica la cartella genitore."},
 {"q":"Qual è la differenza tra `>` e `>>`?","options":["Nessuna","`>` sostituisce il contenuto del file, `>>` aggiunge in fondo","`>` aggiunge, `>>` sostituisce","`>>` funziona solo con le cartelle"],"answer":[1],"explanation":"Con > il file viene svuotato prima di scrivere."},
 {"q":"Quali sono percorsi assoluti?","options":["/etc/hosts","../note.txt","~/note.txt","Documenti/lettera.txt","/home/ubuntu"],"answer":[0,4],"explanation":"I percorsi assoluti iniziano con /. ~ viene espanso dalla shell in un percorso assoluto, ma come scritto non lo è."},
 {"q":"Hai lanciato un comando che non termina. Come lo interrompi?","options":["Ctrl+C","Ctrl+L","Tab","q"],"answer":[0],"explanation":"Ctrl+C invia al processo il segnale di interruzione."},
 {"q":"Come si esce dalla pagina di manuale aperta con `man`?","options":["exit","Ctrl+L","q","Esc"],"answer":[2],"explanation":"Il tasto q chiude il visualizzatore del manuale."}
]}
</script>

## Fonti

- [GNU Bash Reference Manual](https://www.gnu.org/software/bash/manual/bash.html)
- [GNU Coreutils manual](https://www.gnu.org/software/coreutils/manual/coreutils.html)
- [The Linux man-pages project](https://man7.org/linux/man-pages/)
- [Ubuntu: The Linux command line for beginners](https://ubuntu.com/tutorials/command-line-for-beginners)
