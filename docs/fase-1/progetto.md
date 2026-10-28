# Progetto 1 · Server web messo in sicurezza

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai prendere un server Linux appena installato e portarlo a uno stato sicuro, documentato e riproducibile:
    utenti e SSH con chiavi, firewall, servizi systemd, reverse proxy HTTPS, protezione dagli attacchi a forza
    bruta e una procedura di intervento. È il lavoro quotidiano di un sysadmin, ed è un ottimo argomento di colloquio.

| | |
|---|---|
| **Durata** | 10-15 ore |
| **Consegna** | Repository pubblico GitHub `secure-web-server` |
| **Ambiente** | Una VM nuova, `web`, creata con Multipass, più la VM `client` |

## Il compito

!!! tip "Su quale distribuzione?"
    Svolgi il progetto sulla famiglia che hai scelto come percorso principale. Per un portfolio davvero completo,
    realizzalo su **entrambe**: lo stesso repository può contenere `setup-ubuntu.sh` e `setup-rhel.sh`. Nel percorso
    RHEL crea una terza VM `web` aggiungendola al `Vagrantfile` (per esempio con IP `192.168.56.20`) e sostituisci:
    `nftables` con **firewalld**, il gruppo `sudo` con **`wheel`**, il servizio `ssh` con **`sshd`**,
    `sites-available` con **`/etc/nginx/conf.d/`**, `unattended-upgrades` con **`dnf-automatic`**; installa
    fail2ban da **EPEL**; mantieni SELinux in **enforcing** e attiva il booleano `httpd_can_network_connect`.

Su una VM nuova chiamata `web` (`multipass launch 24.04 --name web --cpus 1 --memory 1G --disk 10G`):

### Requisiti funzionali

1. **Utenti**: un utente amministratore `admin` (con sudo) e un utente di servizio `webapp` senza shell di login.
2. **SSH**: accesso solo con chiave ed25519, root disattivato, accesso consentito solo al gruppo `ssh-users`.
3. **Aggiornamenti**: `unattended-upgrades` attivo.
4. **Applicazione**: un servizio systemd `webapp.service` che esegue come utente `webapp` un'applicazione in ascolto
   su `127.0.0.1:8000`. Per ora va bene `python3 -m http.server` su una directory dedicata `/srv/webapp/public`
   con una pagina `index.html`: nella [Fase 2](../fase-2/index.md) la sostituirai con la tua API. Il servizio deve ripartire in caso di errore.
5. **Nginx**: reverse proxy verso l'applicazione, HTTPS con certificato autofirmato per `web.lab`, reindirizzamento
   permanente da HTTP a HTTPS, intestazione HSTS. Nessun sito di default attivo.
6. **Firewall nftables**: politica `drop` in ingresso, aperte solo 22, 80 e 443; nuovi tentativi SSH limitati;
   pacchetti scartati registrati nel journal; configurazione permanente.
7. **fail2ban**: protezione di SSH con blocco dell'IP dopo 5 tentativi falliti in 10 minuti.
8. **Log**: rotazione dei log di Nginx verificata (`logrotate`), journal limitato a 200 MB.

### Requisiti di documentazione

Il repository deve contenere:

- `README.md` con: obiettivo, diagramma dell'architettura (Mermaid), elenco delle misure di sicurezza con il
  **perché** di ciascuna, istruzioni per riprodurre l'ambiente, come verificare che tutto funzioni;
- `config/` con **tutti** i file di configurazione usati (unità systemd, file Nginx, `nftables.conf`, file sshd,
  `jail.local` di fail2ban), con commenti;
- `setup.sh`: uno script che applica tutta la configurazione su una VM nuova (versione semplice: una sequenza di
  comandi; la renderai robusta nella [Fase 2](../fase-2/index.md));
- `RUNBOOK.md`: "Il sito non risponde", come nella [lezione 1.13](13-troubleshooting.md);
- `VERIFICA.md`: l'elenco dei test che hai eseguito, con i comandi e i risultati attesi.

!!! warning "Niente segreti nel repository"
    Non pubblicare **mai** chiavi private, né la chiave del certificato TLS. Nello script, genera il certificato
    sulla VM. Aggiungi un file `.gitignore` che escluda `*.key` e `*.pem`.

## Suggerimenti

- Configura fail2ban creando `/etc/fail2ban/jail.local` (non modificare `jail.conf`):

    ```ini
    [sshd]
    enabled  = true
    backend  = systemd
    maxretry = 5
    findtime = 10m
    bantime  = 1h
    ```

    Verifica con `sudo fail2ban-client status sshd`.
- Il limite del journal si imposta con `SystemMaxUse=200M` in un file in `/etc/systemd/journald.conf.d/`, seguito da
  `sudo systemctl restart systemd-journald`.
- Per lo script `setup.sh`, procedi un passo alla volta sulla VM, annota ogni comando che funziona e alla fine
  mettili in ordine nello script. Poi distruggi la VM, ricreala e lancia lo script: se tutto torna a funzionare,
  hai un'infrastruttura **riproducibile**. È il primo passo verso l'Infrastructure as Code.

## Come verificare

| Test (da `client`) | Risultato atteso |
|---|---|
| `curl -sI http://web.lab` | `301` con `Location: https://web.lab/` |
| `curl -skI https://web.lab` | `200` e intestazione `Strict-Transport-Security` |
| `nc -zvw3 web.lab 8000` | timeout |
| `ssh -o PubkeyAuthentication=no admin@web.lab` | `Permission denied (publickey)` |
| `ssh root@web.lab` | rifiutato |
| 6 tentativi SSH falliti di fila da `client` | `client` bloccato: `sudo fail2ban-client status sshd` mostra il suo IP |
| `sudo systemctl kill -s KILL webapp` su `web` | il servizio torna attivo da solo |
| riavvio di `web` | tutto funziona di nuovo senza interventi |

Dopo il test di fail2ban sblocca il client con `sudo fail2ban-client set sshd unbanip IP-CLIENT`.

## Criteri di valutazione

- [ ] Tutti i test della tabella passano, anche dopo un riavvio.
- [ ] `setup.sh` su una VM nuova porta allo stesso risultato.
- [ ] Il README spiega il perché di ogni scelta di sicurezza.
- [ ] Nessun segreto nel repository.
- [ ] I file di configurazione sono commentati.
- [ ] (Percorso RHEL) SELinux è in enforcing e `sudo ausearch -m AVC -ts today` non mostra blocchi.

??? success "Una possibile soluzione: i file principali"
    Prova a costruire il progetto da solo prima di aprire questa sezione: è il modo in cui impari davvero.

    ```ini title="config/webapp.service"
    [Unit]
    Description=Applicazione web (dietro Nginx)
    After=network.target

    [Service]
    User=webapp
    Group=webapp
    WorkingDirectory=/srv/webapp
    ExecStart=/usr/bin/python3 -m http.server 8000 --bind 127.0.0.1 --directory /srv/webapp/public
    Restart=on-failure
    RestartSec=3
    NoNewPrivileges=true
    ProtectSystem=strict
    ProtectHome=true
    PrivateTmp=true

    [Install]
    WantedBy=multi-user.target
    ```

    Le ultime quattro direttive di `[Service]` sono **protezioni di systemd**: il processo non può acquisire nuovi
    privilegi, vede il sistema in sola lettura, non vede le home degli utenti e ha un `/tmp` privato. Costano
    zero e limitano i danni se l'applicazione viene compromessa.

    ```bash title="setup.sh (estratto)"
    #!/usr/bin/env bash
    set -euo pipefail    # ferma lo script al primo errore (lo studierai nella Fase 2)

    apt-get update
    apt-get install -y nginx nftables fail2ban unattended-upgrades

    # utenti
    groupadd -f ssh-users
    id admin  &>/dev/null || useradd -m -s /bin/bash -G sudo,ssh-users admin
    id webapp &>/dev/null || useradd --system --home /srv/webapp --shell /usr/sbin/nologin webapp
    usermod -aG ssh-users ubuntu    # mantiene l'accesso di Multipass

    # applicazione
    install -d -o webapp -g webapp /srv/webapp/public
    echo '<h1>web.lab funziona</h1>' > /srv/webapp/public/index.html
    install -m 644 config/webapp.service /etc/systemd/system/webapp.service

    # ssh, nftables, nginx, fail2ban
    install -m 644 config/10-sicurezza.conf /etc/ssh/sshd_config.d/10-sicurezza.conf
    install -m 644 config/nftables.conf /etc/nftables.conf
    install -m 644 config/web.lab /etc/nginx/sites-available/web.lab
    install -m 644 config/jail.local /etc/fail2ban/jail.local
    ln -sf /etc/nginx/sites-available/web.lab /etc/nginx/sites-enabled/web.lab
    rm -f /etc/nginx/sites-enabled/default

    # certificato autofirmato, generato sulla macchina
    install -d -m 700 /etc/nginx/tls
    [ -f /etc/nginx/tls/web.lab.key ] || openssl req -x509 -newkey rsa:2048 -nodes -days 365 \
      -keyout /etc/nginx/tls/web.lab.key -out /etc/nginx/tls/web.lab.crt \
      -subj "/CN=web.lab" -addext "subjectAltName=DNS:web.lab"

    sshd -t && nginx -t && nft -c -f /etc/nftables.conf
    systemctl daemon-reload
    systemctl enable --now webapp nftables fail2ban
    systemctl reload ssh nginx
    nft -f /etc/nftables.conf
    ```

    Lo script è **idempotente** in molti punti: `groupadd -f`, i controlli `id ... ||` e `[ -f ... ] ||` evitano errori
    se lo lanci due volte. L'idempotenza è un concetto chiave dell'automazione: lo ritroverai con Ansible e Terraform.

## Fonti

- [Ubuntu Server docs: Security](https://documentation.ubuntu.com/server/how-to/security/)
- [systemd.exec(5): Sandboxing](https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html#Sandboxing)
- [fail2ban: documentazione e wiki](https://github.com/fail2ban/fail2ban/wiki)
- [CIS Benchmarks (linee guida di hardening)](https://www.cisecurity.org/cis-benchmarks)
- [logrotate(8)](https://man7.org/linux/man-pages/man8/logrotate.8.html), [journald.conf(5)](https://www.freedesktop.org/software/systemd/man/latest/journald.conf.html)
