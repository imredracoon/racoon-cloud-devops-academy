# 1.10 SSH e chiavi

!!! obiettivi "Obiettivi"
    - Capire come SSH protegge la comunicazione e come funziona l'autenticazione a chiave pubblica.
    - Generare una coppia di chiavi e usarla per collegarti a un server.
    - Configurare il client con `~/.ssh/config` e trasferire file con `scp`/`rsync`.
    - Mettere in sicurezza il server SSH.

## Cos'è SSH

**SSH** (*Secure Shell*) permette di aprire un terminale su un computer remoto attraverso una connessione
**cifrata**. È il modo standard per amministrare server Linux, anche nel cloud. L'implementazione più diffusa è
**OpenSSH**: il client è `ssh`, il server è `sshd`, in ascolto sulla porta **22**.

## Chiavi asimmetriche in due minuti

La crittografia **asimmetrica** usa una **coppia di chiavi** legate matematicamente:

- la **chiave privata** resta solo sul tuo computer, protetta, e non si condivide mai;
- la **chiave pubblica** si può dare a chiunque: la copi sui server a cui vuoi accedere.

Quando ti colleghi, il server ti invia una "sfida" che solo chi possiede la chiave privata può firmare
correttamente. Così dimostri la tua identità **senza mai inviare un segreto in rete**. È molto più sicuro di una
password, che può essere indovinata o rubata.

C'è anche una seconda verifica, nel verso opposto: il **server** ha una propria coppia di chiavi (*host key*).
Alla prima connessione il client ti mostra l'impronta (*fingerprint*) della chiave del server e la salva in
`~/.ssh/known_hosts`. Se in futuro la chiave cambia, il client ti avvisa: potrebbe essere un attacco
*man-in-the-middle* (o un server reinstallato).

## Generare le chiavi

Sul tuo computer (o nella VM `client` per esercitarti):

```bash
ssh-keygen -t ed25519 -C "nome.cognome@laptop"
```

- `-t ed25519` sceglie l'algoritmo moderno consigliato.
- `-C` aggiunge un commento per riconoscere la chiave.
- Ti viene chiesta una **passphrase**: proteggila, così una chiave rubata è inutilizzabile.

Vengono creati `~/.ssh/id_ed25519` (privata, permessi `600`) e `~/.ssh/id_ed25519.pub` (pubblica).
Per non digitare la passphrase ogni volta si usa l'**ssh-agent** (`eval "$(ssh-agent -s)"` e `ssh-add`).

## Autorizzare la chiave sul server

Sul server, le chiavi pubbliche autorizzate per un utente stanno in `~/.ssh/authorized_keys`, una per riga.

```bash
ssh-copy-id utente@server          # copia la chiave pubblica e imposta i permessi corretti
ssh utente@server                  # ora entri senza password (al massimo la passphrase della chiave)
```

Se non puoi usare `ssh-copy-id`, aggiungi a mano il contenuto di `id_ed25519.pub` al file `authorized_keys` e
controlla i permessi: `~/.ssh` deve essere `700` e `authorized_keys` `600`, altrimenti sshd ignora le chiavi.

!!! note "Nel cloud"
    Quando crei una VM su AWS o Google Cloud scegli (o carichi) una chiave pubblica: il fornitore la inserisce in
    `authorized_keys` al primo avvio tramite cloud-init. È lo stesso meccanismo, automatizzato. In alternativa
    esistono accessi senza porta 22 aperta, come AWS Systems Manager Session Manager e IAP su Google Cloud (Fase 5).

## Usare SSH

```bash
ssh ubuntu@10.84.37.50              # collegamento
ssh -p 2222 ubuntu@server           # porta diversa
ssh -i ~/.ssh/chiave-lab ubuntu@server   # chiave specifica
ssh ubuntu@server 'uptime; df -h /'      # esegue comandi remoti e torna indietro
```

### Il file ~/.ssh/config

Evita di ricordare indirizzi, utenti e opzioni:

```text title="~/.ssh/config"
Host lab
    HostName 10.84.37.50
    User ubuntu
    IdentityFile ~/.ssh/id_ed25519

Host bastion
    HostName 203.0.113.20
    User admin

Host interno
    HostName 10.0.10.15
    User ubuntu
    ProxyJump bastion
```

Ora basta `ssh lab`. `ProxyJump` passa attraverso un **bastion host**: un server esposto che fa da unico punto di
ingresso verso server privati, schema molto usato nel cloud.

### Copiare file

```bash
scp file.txt lab:/tmp/              # copia verso il server
scp lab:/var/log/syslog ./          # copia dal server
rsync -avz --progress cartella/ lab:/srv/backup/   # sincronizza, trasferendo solo le differenze
```

### Tunnel

```bash
ssh -L 8080:127.0.0.1:80 lab        # la porta 8080 del tuo PC porta alla porta 80 locale di lab
```

È il modo sicuro per raggiungere un servizio che ascolta solo su `127.0.0.1` del server, per esempio un
pannello di amministrazione o un database.

## Mettere in sicurezza il server SSH

La configurazione del server è in `/etc/ssh/sshd_config` e nei file `/etc/ssh/sshd_config.d/*.conf`, che conviene
usare per le proprie modifiche. Le impostazioni consigliate:

```text title="/etc/ssh/sshd_config.d/10-sicurezza.conf"
PasswordAuthentication no
KbdInteractiveAuthentication no
PermitRootLogin no
PubkeyAuthentication yes
MaxAuthTries 3
AllowGroups ssh-users
```

| Impostazione | Effetto |
|---|---|
| `PasswordAuthentication no` | Solo chiavi: gli attacchi che provano milioni di password diventano inutili. |
| `PermitRootLogin no` | root non può entrare direttamente; si entra come utente e si usa sudo. |
| `MaxAuthTries 3` | Pochi tentativi per connessione. |
| `AllowGroups` | Solo i membri del gruppo indicato possono collegarsi. |

```bash
sudo sshd -t                     # controlla la sintassi: nessun output = tutto bene
sudo systemctl reload ssh        # su Ubuntu il servizio si chiama ssh; su RHEL sshd
```

!!! danger "Non chiuderti fuori"
    Prima di disattivare le password, **verifica di riuscire a entrare con la chiave**. Prima di ricaricare la
    configurazione, lascia **aperta** la sessione corrente e prova il nuovo accesso da un **secondo terminale**.
    Se qualcosa non va, la prima sessione ti permette di correggere. Nella VM di laboratorio hai comunque
    `multipass shell` come accesso di emergenza; su un server remoto potresti non averlo.

Il pacchetto **fail2ban** legge i log e blocca temporaneamente gli IP che sbagliano troppe volte l'accesso: lo
configurerai nel progetto della fase.

### Le differenze tra le due famiglie

=== "Debian/Ubuntu"

    - Pacchetto `openssh-server`, servizio **`ssh`** (`systemctl reload ssh`).
    - Su Ubuntu 24.04 il servizio può essere attivato tramite *socket activation* (`ssh.socket`): se cambi la
      porta in `sshd_config`, esegui anche `sudo systemctl daemon-reload` e `sudo systemctl restart ssh.socket`.
    - Il firewall (ufw o nftables) deve permettere la porta.

=== "RHEL (Rocky/AlmaLinux)"

    - Pacchetto `openssh-server`, servizio **`sshd`** (`systemctl reload sshd`).
    - Se cambi la porta, **SELinux** deve sapere che la nuova porta è di SSH, altrimenti sshd non riesce ad
      ascoltare (lezione 1.14):

        ```bash
        sudo semanage port -a -t ssh_port_t -p tcp 2222
        sudo firewall-cmd --permanent --add-port=2222/tcp && sudo firewall-cmd --reload
        ```

    - La configurazione di sicurezza va anche qui in `/etc/ssh/sshd_config.d/`.

## Laboratorio

!!! laboratorio "Laboratorio · Accesso solo con chiave · 1 ora"
    **Dove:** VM `client` (il "tuo computer") e VM `lab` (il "server").

    1. Su `client` genera una chiave ed25519 con passphrase.
    2. Su `lab` crea l'utente `deploy` con shell bash e aggiungilo al nuovo gruppo `ssh-users`; aggiungi anche
       `ubuntu` al gruppo, per non perdere l'accesso.
    3. Autorizza la chiave di `client` per l'utente `deploy`.
    4. Configura `~/.ssh/config` su `client` per entrare con `ssh lab-deploy`.
    5. Applica la configurazione di sicurezza vista sopra e verifica che: con la chiave entri; un utente senza
       chiave venga rifiutato; root non possa entrare.

    **Percorso RHEL:** ripeti il laboratorio con `rhel-client` (192.168.56.11) e `rhel-lab` (192.168.56.10), usando
    il servizio `sshd` e il gruppo `wheel` al posto di `sudo`. Ricorda che l'utente `vagrant` deve restare nel
    gruppo `ssh-users`, altrimenti `vagrant ssh` smette di funzionare.

??? success "Soluzione"
    ```bash
    # --- su client ---
    ssh-keygen -t ed25519 -C "deploy@client"
    cat ~/.ssh/id_ed25519.pub            # copia questa riga

    # --- su lab (multipass shell lab) ---
    sudo groupadd ssh-users
    sudo useradd -m -s /bin/bash -G ssh-users deploy
    sudo usermod -aG ssh-users ubuntu
    sudo install -d -m 700 -o deploy -g deploy /home/deploy/.ssh
    echo "INCOLLA-QUI-LA-CHIAVE-PUBBLICA" | sudo tee /home/deploy/.ssh/authorized_keys
    sudo chown deploy:deploy /home/deploy/.ssh/authorized_keys
    sudo chmod 600 /home/deploy/.ssh/authorized_keys

    # --- su client ---
    cat >> ~/.ssh/config <<'EOF'
    Host lab-deploy
        HostName 10.84.37.50
        User deploy
    EOF
    ssh lab-deploy                       # entra (chiede la passphrase della chiave)

    # --- su lab: sicurezza ---
    sudo tee /etc/ssh/sshd_config.d/10-sicurezza.conf > /dev/null <<'EOF'
    PasswordAuthentication no
    KbdInteractiveAuthentication no
    PermitRootLogin no
    PubkeyAuthentication yes
    MaxAuthTries 3
    AllowGroups ssh-users
    EOF
    sudo sshd -t && sudo systemctl reload ssh

    # --- verifiche da client ---
    ssh lab-deploy 'whoami'                                   # deploy
    ssh -o PubkeyAuthentication=no deploy@10.84.37.50         # Permission denied (publickey)
    ssh root@10.84.37.50                                      # Permission denied
    ```
    **Spiegazione.** `install -d` crea la directory con proprietario e permessi in un solo comando. Il secondo test
    forza il client a non usare la chiave: senza password ammesse, il server rifiuta. Su Ubuntu i file in
    `sshd_config.d` sono inclusi all'inizio di `sshd_config`, e per sshd vale **il primo valore letto** di ogni
    opzione: per questo le tue impostazioni vincono su quelle successive.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale chiave va copiata sul server?","options":["La privata","La pubblica","Entrambe","Nessuna, basta la passphrase"],"answer":[1],"explanation":"La chiave privata non lascia mai il tuo computer."},
 {"q":"In quale file il server cerca le chiavi pubbliche autorizzate per un utente?","options":["~/.ssh/known_hosts","~/.ssh/authorized_keys","/etc/ssh/ssh_host_keys","~/.ssh/config"],"answer":[1],"explanation":"known_hosts sta sul client e contiene le chiavi dei server."},
 {"q":"Il client avvisa che la chiave del server è cambiata. Cosa fai?","options":["Ignori l'avviso","Verifichi con l'amministratore se il server è stato reinstallato prima di accettare la nuova chiave","Cancelli la tua chiave privata","Disattivi SSH"],"answer":[1],"explanation":"Potrebbe essere un attacco man-in-the-middle."},
 {"q":"Quali impostazioni rafforzano sshd?","options":["PasswordAuthentication no","PermitRootLogin no","PermitEmptyPasswords yes","AllowGroups ssh-users","Port 22 obbligatoria"],"answer":[0,1,3],"explanation":"Password vuote ammesse sarebbero un grave rischio; cambiare porta non è un requisito di sicurezza."},
 {"q":"Come raggiungi in sicurezza un database che ascolta solo su 127.0.0.1 del server?","options":["Apri la porta del database a Internet","Con un tunnel SSH: ssh -L","Con scp","Con ping"],"answer":[1],"explanation":"Il port forwarding locale porta il traffico dentro la connessione cifrata."},
 {"q":"Quale opzione di ~/.ssh/config permette di passare attraverso un bastion host?","options":["ProxyJump","ForwardAgent","LocalForward","HostKeyAlias"],"answer":[0],"explanation":"ProxyJump salta attraverso l'host indicato."}
]}
</script>

## Fonti

- [OpenSSH: manuali ufficiali](https://www.openssh.com/manual.html)
- [sshd_config(5)](https://man.openbsd.org/sshd_config) e [ssh_config(5)](https://man.openbsd.org/ssh_config)
- [ssh-keygen(1)](https://man.openbsd.org/ssh-keygen)
- [Ubuntu Server docs: OpenSSH server](https://documentation.ubuntu.com/server/how-to/security/openssh-server/)
- [RFC 4251: The Secure Shell (SSH) Protocol Architecture](https://www.rfc-editor.org/rfc/rfc4251)
