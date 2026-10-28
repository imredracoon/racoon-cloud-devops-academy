# 1.5 Pacchetti e software

!!! obiettivi "Obiettivi"
    - Capire cosa sono pacchetti, repository e dipendenze.
    - Installare, aggiornare, cercare e rimuovere software con `apt` (e conoscere l'equivalente `dnf`).
    - Aggiungere un repository esterno in modo sicuro, verificandone la firma.
    - Tenere aggiornato un server, anche in automatico.

## Pacchetti e repository

Un **pacchetto** è un archivio che contiene un programma, i suoi file di configurazione e le informazioni per
installarlo: versione, descrizione e **dipendenze**, cioè gli altri pacchetti di cui ha bisogno.

I pacchetti si scaricano da **repository**: server della distribuzione (o di terzi) che li pubblicano **firmati
digitalmente**. Il gestore di pacchetti verifica la firma, così sei sicuro che il pacchetto non sia stato alterato.

| Famiglia | Formato | Strumento di basso livello | Gestore |
|---|---|---|---|
| Debian, Ubuntu | `.deb` | `dpkg` | `apt` |
| RHEL, Rocky, Fedora, Amazon Linux | `.rpm` | `rpm` | `dnf` |

Installare software scaricando programmi a caso da Internet è una cattiva pratica sui server: si perdono
aggiornamenti automatici, verifica delle firme e gestione delle dipendenze.

## apt in pratica

```bash
sudo apt update                  # aggiorna l'ELENCO dei pacchetti disponibili (non installa nulla)
apt list --upgradable            # cosa si può aggiornare
sudo apt upgrade                 # aggiorna i pacchetti installati
sudo apt full-upgrade            # aggiorna anche se servono rimozioni o nuove dipendenze
apt search nginx                 # cerca
apt show nginx                   # dettagli: versione, dipendenze, descrizione
sudo apt install nginx           # installa (con le dipendenze)
sudo apt install -y htop tree    # -y risponde sì alle domande: utile negli script
sudo apt remove nginx            # rimuove il programma, lascia la configurazione
sudo apt purge nginx             # rimuove anche la configurazione
sudo apt autoremove              # rimuove le dipendenze non più necessarie
apt list --installed | grep nginx
dpkg -L nginx                    # quali file ha installato il pacchetto
dpkg -S /usr/sbin/nginx          # a quale pacchetto appartiene un file
```

!!! tip "update e upgrade"
    `update` scarica il catalogo aggiornato; `upgrade` installa le nuove versioni. Si eseguono quasi sempre in
    coppia: `sudo apt update && sudo apt upgrade`. L'operatore `&&` esegue il secondo comando **solo se** il primo
    è riuscito.

I repository configurati si trovano in `/etc/apt/sources.list.d/` (su Ubuntu 24.04 nel formato *deb822*, file
`ubuntu.sources`).

## L'equivalente su Red Hat

| Operazione | Debian/Ubuntu | RHEL/Rocky/Fedora |
|---|---|---|
| Aggiornare tutto | `apt update && apt upgrade` | `dnf upgrade` |
| Cercare | `apt search x` | `dnf search x` |
| Installare | `apt install x` | `dnf install x` |
| Rimuovere | `apt remove x` | `dnf remove x` |
| Informazioni | `apt show x` | `dnf info x` |
| Chi fornisce un file | `dpkg -S /percorso` | `dnf provides /percorso` |
| Elenco repository | `/etc/apt/sources.list.d/` | `dnf repolist`, `/etc/yum.repos.d/` |

## Aggiungere un repository esterno

A volte un software (per esempio Docker o la Google Cloud CLI) va installato dal repository ufficiale del
produttore. Lo schema sicuro, usato anche nel Dockerfile del Dev Container del corso, è sempre lo stesso:

1. scaricare la **chiave pubblica** del produttore e salvarla in `/usr/share/keyrings/` o `/etc/apt/keyrings/`;
2. aggiungere il repository indicando con `signed-by` che va verificato **solo** con quella chiave;
3. `apt update` e installazione.

```bash
# Esempio dalla documentazione ufficiale della Google Cloud CLI
curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
  | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" \
  | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list
sudo apt update && sudo apt install google-cloud-cli
```

!!! warning "Solo fonti ufficiali"
    Aggiungi repository solo dalla documentazione ufficiale del produttore. Un repository malevolo può installare
    qualsiasi cosa con i privilegi di root.

## Altri formati: snap

Ubuntu offre anche gli **snap**, pacchetti autocontenuti con le proprie dipendenze, aggiornati automaticamente
(`snap install`, `snap list`). Li hai usati per Multipass su Linux. Sui server si preferiscono di norma i pacchetti
della distribuzione.

## Aggiornamenti di sicurezza automatici

Un server non aggiornato è un server vulnerabile. Su Ubuntu il pacchetto **unattended-upgrades** installa in
automatico gli aggiornamenti di sicurezza:

```bash
sudo apt install unattended-upgrades
sudo dpkg-reconfigure --priority=low unattended-upgrades   # attiva gli aggiornamenti automatici
cat /etc/apt/apt.conf.d/20auto-upgrades
```

Alcuni aggiornamenti (per esempio del kernel) richiedono un riavvio: su Ubuntu ne trovi l'avviso nel file
`/var/run/reboot-required`.

## Esercizi

!!! esercizio "Esercizio 1 · Investigazione"
    1. Installa il pacchetto `tree`.
    2. Trova a quale pacchetto appartiene il comando `/usr/bin/curl`.
    3. Elenca i file installati dal pacchetto `tree`.
    4. Rimuovi `tree` cancellando anche la configurazione e le dipendenze non più usate.

??? success "Soluzione"
    ```bash
    sudo apt update && sudo apt install -y tree
    dpkg -S /usr/bin/curl          # curl: /usr/bin/curl
    dpkg -L tree
    sudo apt purge -y tree && sudo apt autoremove -y
    ```

!!! esercizio "Esercizio 2 · Versioni"
    Usa `apt policy nginx` (o `apt-cache policy nginx`) e spiega cosa significano le righe *Installed* e *Candidate*.

??? success "Soluzione"
    *Installed* è la versione installata (`(none)` se non c'è); *Candidate* è quella che apt installerebbe ora,
    scelta tra i repository configurati secondo la loro priorità. Sotto trovi l'elenco dei repository che offrono il
    pacchetto. È il comando da usare quando "apt installa una versione diversa da quella che mi aspettavo".

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa fa `sudo apt update`?","options":["Aggiorna tutti i pacchetti installati","Aggiorna l'elenco dei pacchetti disponibili dai repository","Aggiorna il kernel","Aggiorna apt stesso"],"answer":[1],"explanation":"Per installare le nuove versioni serve apt upgrade."},
 {"q":"Qual è la differenza tra `apt remove` e `apt purge`?","options":["Nessuna","purge rimuove anche i file di configurazione","remove rimuove anche le dipendenze","purge funziona solo con gli snap"],"answer":[1],"explanation":"remove lascia la configurazione in /etc."},
 {"q":"Quale comando usa il gestore di pacchetti di Rocky Linux per installare nginx?","options":["apt install nginx","dnf install nginx","zypper install nginx","snap install nginx"],"answer":[1],"explanation":"La famiglia Red Hat usa dnf."},
 {"q":"Perché i repository usano firme digitali?","options":["Per comprimere i pacchetti","Per garantire che i pacchetti provengano dal produttore e non siano stati alterati","Per velocizzare i download","Per limitare il numero di installazioni"],"answer":[1],"explanation":"apt rifiuta pacchetti con firme non valide."},
 {"q":"Quale pacchetto installa automaticamente gli aggiornamenti di sicurezza su Ubuntu?","options":["auto-apt","unattended-upgrades","cron-apt","snapd"],"answer":[1],"explanation":"unattended-upgrades è il meccanismo standard."}
]}
</script>

## Fonti

- [Ubuntu Server docs: Package management](https://documentation.ubuntu.com/server/how-to/software/package-management/)
- [apt(8), man page](https://manpages.ubuntu.com/manpages/noble/man8/apt.8.html)
- [Debian Reference: Debian package management](https://www.debian.org/doc/manuals/debian-reference/ch02.en.html)
- [DNF documentation](https://dnf.readthedocs.io/)
- [Ubuntu Server docs: Automatic updates](https://documentation.ubuntu.com/server/how-to/software/automatic-updates/)
- [Google Cloud: installare la gcloud CLI](https://cloud.google.com/sdk/docs/install)
