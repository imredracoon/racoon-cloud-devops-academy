# 1.1 Il filesystem

!!! obiettivi "Obiettivi"
    - Conoscere la struttura standard delle directory Linux (FHS) e lo scopo delle principali.
    - Capire il principio "tutto è un file" e i tipi di file.
    - Usare link simbolici e link fisici.
    - Trovare file con `find` e capire cos'è un inode.

## Un unico albero

In Windows ogni disco ha una lettera (`C:`, `D:`). In Linux esiste **un solo albero** di directory che parte dalla
**radice** `/`. Dischi, chiavette e cartelle di rete vengono "agganciati" (**montati**) in un punto dell'albero,
come vedrai nella [lezione 1.6](06-storage.md).

La struttura è definita dal **Filesystem Hierarchy Standard** (FHS), mantenuto dalla Linux Foundation:

| Directory | Contenuto |
|---|---|
| `/` | La radice di tutto. |
| `/bin`, `/usr/bin` | Comandi per tutti gli utenti (`ls`, `cp`…). Nelle distribuzioni moderne `/bin` è un link a `/usr/bin`. |
| `/sbin`, `/usr/sbin` | Comandi di amministrazione del sistema. |
| `/etc` | **Configurazione** del sistema e dei servizi. Solo file di testo. |
| `/home` | Le directory personali degli utenti (`/home/ubuntu`). |
| `/root` | La home dell'amministratore `root`. |
| `/var` | Dati **variabili**: log (`/var/log`), code di posta, cache, database, siti web (`/var/www`). |
| `/tmp` | File temporanei, spesso cancellati al riavvio. |
| `/usr` | Programmi, librerie e documentazione installati dal sistema. |
| `/opt` | Software aggiuntivo installato fuori dal gestore di pacchetti. |
| `/srv` | Dati serviti dal sistema (per esempio i dischi di laboratorio in `/srv/lab-disks`). |
| `/dev` | File che rappresentano i **dispositivi** (`/dev/sda`, `/dev/null`). |
| `/proc`, `/sys` | Filesystem **virtuali** generati dal kernel: informazioni su processi e hardware. |
| `/boot` | Kernel e file di avvio. |
| `/mnt`, `/media` | Punti di montaggio temporanei e per supporti rimovibili. |

!!! tip "Da ricordare per il lavoro"
    Configurazione in `/etc`, log in `/var/log`, dati utente in `/home`. Queste tre risposte risolvono metà delle
    domande "dove trovo…?".

## Tutto è un file

In Linux quasi tutto è rappresentato come un file: i documenti, ma anche i dispositivi, le informazioni sui processi
e i canali di comunicazione. Questo permette di usare gli stessi strumenti per tutto.

```console
$ cat /proc/cpuinfo | head -5     # informazioni sulla CPU, generate dal kernel
$ cat /proc/meminfo | head -3     # informazioni sulla memoria
$ cat /sys/class/net/lo/mtu       # un parametro dell'interfaccia di loopback
65536
```

Il primo carattere dell'output di `ls -l` indica il **tipo di file**:

| Carattere | Tipo |
|---|---|
| `-` | file regolare |
| `d` | directory |
| `l` | link simbolico |
| `c` / `b` | dispositivo a caratteri / a blocchi (in `/dev`) |
| `s` / `p` | socket / pipe con nome (canali di comunicazione tra processi) |

```console
$ ls -l /dev/null /dev/loop0 /etc/hostname
crw-rw-rw- 1 root root 1, 3 Oct 28 10:00 /dev/null
brw-rw---- 1 root disk 7, 0 Oct 28 10:00 /dev/loop0
-rw-r--r-- 1 root root    4 Oct 28 10:00 /etc/hostname
```

`/dev/null` è un "buco nero": tutto ciò che ci scrivi sparisce. Lo userai spesso per scartare output indesiderato.

Il comando `file` indica il contenuto di un file, indipendentemente dall'estensione (che in Linux non ha valore
speciale):

```console
$ file /etc/hostname /usr/bin/ls
/etc/hostname: ASCII text
/usr/bin/ls:   ELF 64-bit LSB pie executable, x86-64, ...
```

## Inode e link

Ogni file ha un **inode**: una struttura che contiene i suoi **metadati** (proprietario, permessi, dimensione,
date, posizione dei dati sul disco). Il **nome** del file non è nell'inode: è una voce nella directory che punta
all'inode. `ls -i` mostra il numero di inode.

### Link fisici (hard link)

Un **hard link** è un secondo nome per lo stesso inode. I due nomi sono equivalenti: modificando uno, cambia anche
l'altro, perché è lo stesso file. Il file viene cancellato davvero solo quando sparisce l'ultimo nome.

```bash
echo "dati" > originale.txt
ln originale.txt secondo-nome.txt
ls -li originale.txt secondo-nome.txt   # stesso numero di inode, contatore dei link = 2
```

Limiti: non si possono creare hard link a directory né tra filesystem diversi.

### Link simbolici (symlink)

Un **link simbolico** è un piccolo file che contiene il **percorso** di un altro file, come un collegamento sul
desktop. Funziona anche con directory e tra filesystem diversi; se cancelli l'originale, il link resta "rotto".

```bash
ln -s /etc/os-release info-sistema
ls -l info-sistema          # info-sistema -> /etc/os-release
cat info-sistema            # legge il file puntato
readlink -f info-sistema    # mostra il percorso finale
```

I symlink sono ovunque nei server: per esempio Nginx attiva i siti con un link da `sites-enabled` a
`sites-available` ([lezione 1.12](12-http-nginx.md)).

## Trovare i file: find

`find` cerca file e directory in base a nome, tipo, dimensione, data, proprietario e altro.

```bash
find /etc -name "*.conf"               # file che finiscono per .conf sotto /etc
find /etc -iname "HOSTS"               # nome senza distinguere maiuscole/minuscole
find /var/log -type f -size +1M        # file regolari più grandi di 1 MiB
find ~ -type d -name ".*"              # directory nascoste nella home
find /tmp -type f -mmin -60            # file modificati negli ultimi 60 minuti
find /home -user ubuntu -name "*.txt"  # file .txt di proprietà di ubuntu
find . -name "*.bak" -delete           # cancella i .bak trovati (attenzione!)
find /etc -name "*.conf" -exec ls -l {} \;   # esegue un comando su ogni risultato
```

Cercando in directory di sistema compariranno errori *Permission denied*: li puoi scartare aggiungendo
`2>/dev/null` (capirai il perché nella prossima lezione).

Per trovare **dove si trova un comando** usa invece `which` o `type`:

```console
$ which ls
/usr/bin/ls
$ type cd
cd is a shell builtin
```

`cd` non è un programma su disco ma un comando interno della shell (*builtin*).

## Esercizi

!!! esercizio "Esercizio 1 · Dove si trova?"
    Senza cercare su Internet, indica in quale directory ti aspetti di trovare: (a) la configurazione del server SSH,
    (b) i log di sistema, (c) il programma `cp`, (d) i file di un sito web, (e) le informazioni sulla memoria
    generate dal kernel. Poi verifica nella VM.

??? success "Soluzione"
    (a) `/etc/ssh/sshd_config`, (b) `/var/log` (per esempio `/var/log/syslog`; con systemd anche il journal, [lezione 1.4](04-processi-systemd.md)),
    (c) `/usr/bin/cp` (`which cp`), (d) `/var/www`, (e) `/proc/meminfo`.

!!! esercizio "Esercizio 2 · Link"
    1. Crea `~/link-lab/dati.txt` con il testo `versione 1`.
    2. Crea un hard link `copia-fisica.txt` e un link simbolico `scorciatoia.txt`, entrambi verso `dati.txt`.
    3. Cancella `dati.txt`. Cosa succede leggendo i due link? Spiega perché.

??? success "Soluzione"
    ```bash
    mkdir -p ~/link-lab && cd ~/link-lab
    echo "versione 1" > dati.txt
    ln dati.txt copia-fisica.txt
    ln -s dati.txt scorciatoia.txt
    rm dati.txt
    cat copia-fisica.txt     # versione 1
    cat scorciatoia.txt      # No such file or directory
    ```
    **Spiegazione.** L'hard link è un altro nome dello stesso inode: i dati restano finché esiste almeno un nome.
    Il link simbolico contiene solo il percorso `dati.txt`, che non esiste più: il link è rotto (`ls -l` lo mostra
    spesso in rosso).

!!! esercizio "Esercizio 3 · find"
    Scrivi i comandi per trovare:

    1. tutti i file `.conf` direttamente in `/etc`, senza scendere nelle sottodirectory;
    2. i file in `/var/log` più grandi di 100 KiB;
    3. i link simbolici presenti in `/etc`, scartando gli errori.

??? success "Soluzione"
    ```bash
    find /etc -maxdepth 1 -name "*.conf"
    sudo find /var/log -type f -size +100k
    find /etc -type l 2>/dev/null
    ```
    **Spiegazione.** `-maxdepth 1` limita la profondità della ricerca. Per `find`, `k` indica unità da 1024 byte.
    `sudo` ([lezione 1.3](03-permessi.md)) serve perché alcuni log sono leggibili solo dall'amministratore.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Dove si trovano normalmente i file di configurazione di sistema?","options":["/var","/etc","/usr","/opt"],"answer":[1],"explanation":"/etc contiene la configurazione."},
 {"q":"Cosa contiene /proc?","options":["I programmi installati","I log","Un filesystem virtuale generato dal kernel con informazioni su processi e sistema","I file temporanei"],"answer":[2],"explanation":"/proc non occupa spazio su disco: è generato dal kernel."},
 {"q":"Cancelli il file a cui punta un link simbolico. Cosa succede al link?","options":["Viene cancellato automaticamente","Diventa un link rotto","Diventa un hard link","Continua a mostrare i dati"],"answer":[1],"explanation":"Il symlink contiene solo un percorso, ora inesistente."},
 {"q":"Cosa è memorizzato nell'inode?","options":["Il nome del file","I metadati del file (proprietario, permessi, dimensione, posizione dei dati)","Solo il contenuto","Il percorso completo"],"answer":[1],"explanation":"Il nome sta nella directory, non nell'inode."},
 {"q":"Quale comando trova i file più grandi di 50 MiB nella home?","options":["find ~ -size +50M","ls -S ~ 50M","find ~ -name 50M","grep -r 50M ~"],"answer":[0],"explanation":"-size +50M significa più grande di 50 MiB."},
 {"q":"Nell'output di `ls -l` una riga inizia con `l`. Di cosa si tratta?","options":["Un file di log","Un link simbolico","Una libreria","Un file bloccato"],"answer":[1],"explanation":"Il primo carattere indica il tipo: l = symlink."}
]}
</script>

## Fonti

- [Filesystem Hierarchy Standard 3.0 (Linux Foundation)](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html)
- [hier(7), man page](https://man7.org/linux/man-pages/man7/hier.7.html)
- [find(1), man page](https://man7.org/linux/man-pages/man1/find.1.html)
- [ln(1), man page](https://man7.org/linux/man-pages/man1/ln.1.html)
- [inode(7), man page](https://man7.org/linux/man-pages/man7/inode.7.html)
