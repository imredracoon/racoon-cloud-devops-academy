# 1.3 Utenti, gruppi e permessi

!!! obiettivi "Obiettivi"
    - Capire utenti, gruppi, UID e GID e i file che li descrivono.
    - Creare e gestire utenti e gruppi.
    - Leggere e modificare i permessi in forma simbolica e ottale.
    - Usare `sudo` in modo sicuro e conoscere `umask`, permessi speciali e ACL.

## Utenti e gruppi

Linux è un sistema **multiutente**: ogni processo e ogni file appartengono a un **utente**. Gli utenti sono
raggruppati in **gruppi** per assegnare permessi a più persone insieme.

- Ogni utente ha un numero, l'**UID**; ogni gruppo un **GID**. Il sistema ragiona con i numeri; i nomi servono a noi.
- L'utente **root** (UID 0) è l'amministratore e può fare qualsiasi cosa.
- Gli **utenti di sistema** (di solito UID sotto 1000) servono a far girare i servizi con privilegi ridotti:
  per esempio `www-data` per il server web.
- Gli utenti "umani" partono di norma da UID 1000.

| File | Contenuto |
|---|---|
| `/etc/passwd` | Elenco degli utenti (lo hai visto nella lezione 1.2). Leggibile da tutti. |
| `/etc/shadow` | Le password, conservate come **hash** (impronta crittografica). Leggibile solo da root. |
| `/etc/group` | Elenco dei gruppi e dei loro membri. |

```console
$ id
uid=1000(ubuntu) gid=1000(ubuntu) groups=1000(ubuntu),4(adm),27(sudo),...
$ id root
uid=0(root) gid=0(root) groups=0(root)
```

## Diventare amministratore: sudo

Lavorare sempre come root è pericoloso: un errore di battitura può distruggere il sistema. Per questo si lavora
come utente normale e si usa **`sudo`** (*superuser do*) solo per i singoli comandi che richiedono privilegi:

```bash
sudo apt update              # esegue il comando come root
sudo -i                      # apre una shell di root (esci con exit): da usare con parsimonia
sudo -u www-data whoami      # esegue un comando come un altro utente
```

Chi può usare `sudo` è definito in `/etc/sudoers` e nei file in `/etc/sudoers.d/`. Su Ubuntu basta far parte del
gruppo `sudo`; su Red Hat del gruppo `wheel`. Il file **va modificato solo con `sudo visudo`**, che controlla la
sintassi prima di salvare: un errore in quel file può impedire a tutti di diventare amministratori.
Ogni uso di sudo viene registrato nei log (lezione 1.4), così si sa chi ha fatto cosa.

## Gestire utenti e gruppi

```bash
sudo useradd -m -s /bin/bash anna      # crea l'utente con home (-m) e shell bash
sudo passwd anna                       # imposta la password
sudo groupadd sviluppo                 # crea un gruppo
sudo usermod -aG sviluppo anna         # aggiunge anna al gruppo (-a = append: senza, la toglie dagli altri!)
sudo usermod -aG sudo anna             # le dà i privilegi sudo
groups anna                            # mostra i gruppi di anna
sudo chage -l anna                     # informazioni sulla scadenza della password
sudo usermod -L anna                   # blocca l'account (-U per sbloccarlo)
sudo userdel -r anna                   # cancella l'utente e la sua home
su - anna                              # diventa anna (chiede la password di anna)
```

!!! warning "L'errore classico"
    `usermod -G sviluppo anna` (senza `-a`) **sostituisce** tutti i gruppi secondari di anna con `sviluppo`. Se anna
    era nel gruppo `sudo`, perde i privilegi. Usa sempre `-aG`. Inoltre l'appartenenza a un nuovo gruppo vale dal
    **prossimo login** dell'utente.

Su Debian e Ubuntu esiste anche `adduser`, un comando interattivo più amichevole; `useradd` è disponibile su tutte
le distribuzioni ed è quello da conoscere per l'automazione.

## I permessi

Ogni file ha un **proprietario**, un **gruppo** e tre serie di permessi:

```text
-rwxr-x--- 1 anna sviluppo 1024 Oct 28 10:00 script.sh
 │  │  │       │     │
 │  │  │       │     └ gruppo del file
 │  │  │       └ proprietario
 │  │  └ altri (others): ---
 │  └ gruppo (group):    r-x
 └ proprietario (user):  rwx
```

| Permesso | Su un file | Su una directory |
|---|---|---|
| `r` read | leggere il contenuto | elencare i nomi dei file (`ls`) |
| `w` write | modificare il contenuto | creare, cancellare, rinominare file al suo interno |
| `x` execute | eseguire come programma | **entrare** nella directory e accedere ai file (`cd`) |

!!! note "Il permesso w sulla directory"
    Per cancellare un file conta il permesso `w` sulla **directory** che lo contiene, non sul file. È una domanda
    classica nei colloqui.

### Notazione ottale

Ogni permesso ha un valore: `r` = 4, `w` = 2, `x` = 1. Si sommano per ogni serie:

| Ottale | Simbolico | Uso tipico |
|---|---|---|
| `755` | `rwxr-xr-x` | programmi, directory pubbliche |
| `644` | `rw-r--r--` | file normali leggibili da tutti |
| `700` | `rwx------` | directory private |
| `600` | `rw-------` | file privati, chiavi SSH |
| `640` | `rw-r-----` | configurazioni con segreti, leggibili da un gruppo di servizio |

### chmod, chown, chgrp

```bash
chmod 640 config.ini            # forma ottale
chmod u+x script.sh             # aggiunge l'esecuzione al proprietario (u=user, g=group, o=others, a=all)
chmod go-w file                 # toglie la scrittura a gruppo e altri
chmod -R g+rX progetto/         # ricorsivo; X maiuscola = esecuzione solo per directory (e file già eseguibili)
sudo chown anna file            # cambia proprietario
sudo chown anna:sviluppo file   # cambia proprietario e gruppo
sudo chgrp sviluppo file        # cambia solo il gruppo
```

### umask

La **umask** stabilisce quali permessi **togliere** ai file appena creati. I file nascono al massimo con `666`
e le directory con `777`; con la umask tipica `022` si ottiene `644` per i file e `755` per le directory.

```console
$ umask
0022
```

## Permessi speciali

| Permesso | Ottale | Effetto | Esempio |
|---|---|---|---|
| **setuid** | `4000` | Un programma eseguibile gira con i privilegi del **proprietario** del file. | `/usr/bin/passwd` (`-rwsr-xr-x`): un utente normale può cambiare la propria password scrivendo in `/etc/shadow`. |
| **setgid** | `2000` | Su una directory, i nuovi file **ereditano il gruppo** della directory. | Cartelle condivise di un team. |
| **sticky bit** | `1000` | Su una directory, ognuno può cancellare **solo i propri file**. | `/tmp` (`drwxrwxrwt`). |

```bash
sudo chmod 2770 /srv/progetto     # setgid + rwx per proprietario e gruppo
```

## ACL: permessi più granulari

Le **ACL** (*Access Control List*) permettono di dare permessi a utenti o gruppi specifici oltre alle tre serie
standard. Un `+` alla fine dei permessi in `ls -l` indica che il file ha una ACL.

```bash
setfacl -m u:anna:rw report.txt     # dà lettura e scrittura ad anna
getfacl report.txt                  # mostra le ACL
setfacl -x u:anna report.txt        # rimuove la voce di anna
```

## Laboratorio

!!! laboratorio "Laboratorio · Una cartella condivisa per il team · 45 minuti"
    **Dove:** VM `lab`. Fai prima uno snapshot.

    Un team di due sviluppatori, `marco` e `giulia`, deve condividere la cartella `/srv/team`:

    1. entrambi devono poter creare e modificare file nella cartella;
    2. i file creati da uno devono essere modificabili dall'altro;
    3. gli altri utenti del sistema non devono vedere nulla;
    4. nessuno dei due deve poter cancellare i file dell'altro.

??? success "Soluzione"
    ```bash
    sudo groupadd team
    sudo useradd -m -s /bin/bash -G team marco
    sudo useradd -m -s /bin/bash -G team giulia
    sudo mkdir /srv/team
    sudo chown root:team /srv/team
    sudo chmod 3770 /srv/team        # 2000 setgid + 1000 sticky + 770
    ls -ld /srv/team                 # drwxrws--T root team
    ```
    Per il punto 2 i file devono nascere scrivibili dal gruppo: la umask di default `022` li crea `644`. Imposta la
    umask `002` per i due utenti aggiungendo `umask 002` al loro `~/.profile`, che viene letto a ogni login:

    ```bash
    echo "umask 002" | sudo tee -a /home/marco/.profile /home/giulia/.profile
    ```

    (Non usare `~/.bashrc`: su Ubuntu termina subito se la shell non è interattiva, quindi la umask non verrebbe
    applicata ai comandi lanciati da script.)

    **Verifica.**
    ```bash
    sudo -iu marco bash -c 'echo ciao > /srv/team/nota.txt; ls -l /srv/team'   # gruppo team, rw-rw-r--
    sudo -iu giulia bash -c 'echo modifica >> /srv/team/nota.txt'             # funziona
    sudo -iu giulia bash -c 'rm /srv/team/nota.txt'                           # Operation not permitted
    ls /srv/team                                                              # come ubuntu: Permission denied
    ```

    **Spiegazione.** Il **setgid** fa ereditare il gruppo `team` ai nuovi file; la **umask 002** li rende scrivibili
    dal gruppo; lo **sticky bit** impedisce di cancellare i file altrui; `770` esclude gli altri. La `T` maiuscola
    in `ls` indica sticky bit senza permesso `x` per gli altri. (`ubuntu` non vede il contenuto perché non è nel
    gruppo; con `sudo` naturalmente può tutto.)

## Esercizi

!!! esercizio "Esercizio 1 · Da simbolico a ottale"
    Converti: `rw-r-----`, `rwxr-xr-x`, `rwx--x--x`, `r--------`.

??? success "Soluzione"
    `640`, `755`, `711`, `400`. Per ogni serie somma 4 (r), 2 (w), 1 (x).

!!! esercizio "Esercizio 2 · Chi può fare cosa?"
    La directory `/dati` ha permessi `drwxr-x--x root staff`. Il file `/dati/report.txt` ha `-rw-rw-rw- root root`.
    L'utente `luca`, che **non** è nel gruppo `staff`, può: (a) elencare il contenuto di `/dati`? (b) leggere
    `/dati/report.txt` se conosce il nome? (c) cancellarlo?

??? success "Soluzione"
    (a) **No**: per gli altri la directory ha solo `x`, non `r`.
    (b) **Sì**: con `x` può attraversare la directory, e il file è leggibile da tutti.
    (c) **No**: per cancellare serve `w` sulla **directory**, che per gli altri non c'è.

!!! esercizio "Esercizio 3 · Un file di segreti"
    Crea `/etc/app/secret.env`, di proprietà di root e del gruppo `www-data`, leggibile da root e dal gruppo e
    da nessun altro, non modificabile dal gruppo.

??? success "Soluzione"
    ```bash
    sudo mkdir -p /etc/app
    echo "DB_PASSWORD=cambiami" | sudo tee /etc/app/secret.env > /dev/null
    sudo chown root:www-data /etc/app/secret.env
    sudo chmod 640 /etc/app/secret.env
    ls -l /etc/app/secret.env      # -rw-r----- root www-data
    ```
    **Spiegazione.** `sudo echo ... > file` non funziona, perché la redirezione è eseguita dalla tua shell, che non è
    root. Per questo si usa `| sudo tee file`: è `tee`, eseguito come root, a scrivere il file.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Dove sono conservati gli hash delle password?","options":["/etc/passwd","/etc/shadow","/etc/group","/etc/sudoers"],"answer":[1],"explanation":"/etc/shadow è leggibile solo da root."},
 {"q":"Quale comando aggiunge l'utente anna al gruppo docker senza toglierla dagli altri gruppi?","options":["usermod -G docker anna","usermod -aG docker anna","groupadd anna docker","chgrp docker anna"],"answer":[1],"explanation":"-a (append) è indispensabile."},
 {"q":"Che permessi corrispondono a 750?","options":["rwxr-x---","rwxrw----","rw-r-x---","rwxr-xr-x"],"answer":[0],"explanation":"7 = rwx, 5 = r-x, 0 = ---."},
 {"q":"Per cancellare un file serve il permesso:","options":["w sul file","w sulla directory che lo contiene","x sul file","r sulla directory"],"answer":[1],"explanation":"Cancellare significa modificare l'elenco della directory."},
 {"q":"Perché /tmp ha lo sticky bit?","options":["Per renderla più veloce","Perché ognuno possa cancellare solo i propri file","Per cifrare i file","Per far ereditare il gruppo"],"answer":[1],"explanation":"/tmp è scrivibile da tutti, lo sticky bit evita che si cancellino i file altrui."},
 {"q":"Come si modifica in modo sicuro /etc/sudoers?","options":["sudo nano /etc/sudoers","sudo visudo","chmod 777 /etc/sudoers","echo >> /etc/sudoers"],"answer":[1],"explanation":"visudo controlla la sintassi prima di salvare."},
 {"q":"Con umask 027, con quali permessi nasce un nuovo file?","options":["640","750","644","600"],"answer":[0],"explanation":"666 meno 027 (bit per bit) = 640."}
]}
</script>

## Fonti

- [chmod(1), man page](https://man7.org/linux/man-pages/man1/chmod.1.html)
- [useradd(8)](https://man7.org/linux/man-pages/man8/useradd.8.html) e [usermod(8)](https://man7.org/linux/man-pages/man8/usermod.8.html)
- [sudoers(5), documentazione di Sudo](https://www.sudo.ws/docs/man/sudoers.man/)
- [umask(2), man page](https://man7.org/linux/man-pages/man2/umask.2.html)
- [acl(5), man page](https://man7.org/linux/man-pages/man5/acl.5.html)
- [Ubuntu Server docs: User management](https://documentation.ubuntu.com/server/how-to/security/user-management/)
