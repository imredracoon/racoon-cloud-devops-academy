# 1.14 SELinux e AppArmor

!!! obiettivi "Obiettivi"
    - Capire la differenza tra controllo d'accesso discrezionale (DAC) e obbligatorio (MAC).
    - Gestire SELinux su RHEL: modalità, contesti, booleani, etichette delle porte, diagnosi dei blocchi.
    - Conoscere AppArmor su Ubuntu: profili, modalità, diagnosi.
    - Risolvere un blocco senza disattivare la protezione.

## Perché non bastano i permessi

I permessi della [lezione 1.3](03-permessi.md) sono un controllo **discrezionale** (DAC): il proprietario di un file decide chi può
accedervi, e root può tutto. Se un attaccante prende il controllo di un servizio che gira come root, ottiene
accesso a tutto il sistema.

Il controllo d'accesso **obbligatorio** (MAC, *Mandatory Access Control*) aggiunge un secondo livello di regole,
definite dall'amministratore in una **policy**, che valgono **anche per root**. Un server web compromesso, anche
se gira come root, può leggere solo i file del sito e collegarsi solo alle porte previste.

Linux implementa il MAC con i **Linux Security Modules**. Le due implementazioni più diffuse:

| | SELinux | AppArmor |
|---|---|---|
| Distribuzioni | **RHEL**, Rocky, AlmaLinux, Fedora | **Ubuntu**, Debian, SUSE |
| Modello | Etichette (**contesti**) su ogni file, processo e porta | **Profili** per programma, basati sui percorsi |
| Default | Attivo in modalità *enforcing* su tutto il sistema | Profili attivi per alcuni programmi |

!!! danger "Non disattivarli"
    Molte guide in rete risolvono i problemi con "disattiva SELinux". È l'equivalente di togliere la porta blindata
    perché la chiave si inceppa. In un'azienda seria, e all'esame RHCSA, la soluzione corretta è **capire il blocco e
    correggere la configurazione** lasciando la protezione attiva.

## SELinux (percorso RHEL)

Svolgi questa parte su `rhel-lab`.

### Le modalità

| Modalità | Comportamento |
|---|---|
| `enforcing` | La policy è applicata: le azioni non permesse sono bloccate e registrate. |
| `permissive` | Nulla è bloccato, ma le violazioni vengono registrate. Utile per diagnosticare. |
| `disabled` | SELinux spento. Riattivarlo richiede di rietichettare tutto il filesystem. |

```bash
getenforce                      # Enforcing
sestatus                        # dettagli
sudo setenforce 0               # permissive fino al riavvio (1 = enforcing)
grep ^SELINUX= /etc/selinux/config   # modalità all'avvio: la modifichi qui per renderla permanente
```

### I contesti

Ogni file, processo e porta ha un **contesto** nella forma `utente:ruolo:tipo:livello`. Per l'amministrazione
quotidiana conta il **tipo** (che finisce in `_t`): la policy dice quali tipi di processo possono accedere a quali
tipi di file e porte.

```bash
ls -Z /var/www/html             # system_u:object_r:httpd_sys_content_t:s0
ps -eZ | grep nginx             # system_u:system_r:httpd_t:s0  ... nginx
id -Z                           # il contesto della tua shell
sudo semanage port -l | grep http_port_t   # porte etichettate per i server web
```

Un processo con tipo `httpd_t` (Nginx o Apache) può leggere file `httpd_sys_content_t` e ascoltare su porte
`http_port_t`. Se provi a servire file con un altro tipo, o ad ascoltare su una porta diversa, SELinux blocca.

### Ripristinare e cambiare i contesti

Un errore classico: **sposti** (`mv`) un file da un'altra directory. Il file conserva il contesto di origine
(per esempio `user_home_t`) e il server web non può leggerlo. Con `cp`, invece, il nuovo file prende il contesto
della directory di destinazione.

```bash
restorecon -Rv /var/www/html                   # riporta i file al contesto previsto dalla policy
sudo semanage fcontext -a -t httpd_sys_content_t "/srv/sito(/.*)?"   # regola permanente per una nuova directory
sudo restorecon -Rv /srv/sito                  # applica la regola
```

`chcon` cambia il contesto direttamente, ma la modifica si perde al primo `restorecon` o alla rietichettatura:
per una soluzione permanente usa sempre `semanage fcontext` seguito da `restorecon`.

### Booleani

I **booleani** sono interruttori che attivano parti della policy senza scriverne di nuove:

```bash
getsebool -a | grep httpd                      # elenco dei booleani di httpd
sudo setsebool -P httpd_can_network_connect on # -P = permanente
semanage boolean -l | grep httpd_can_network   # con descrizione
```

### Etichette delle porte

```bash
sudo semanage port -a -t http_port_t -p tcp 8081   # permette ai server web di ascoltare sulla 8081
sudo semanage port -l | grep 8081
```

### Diagnosticare un blocco

I blocchi vengono registrati come messaggi **AVC** (*Access Vector Cache*) nel log di audit:

```bash
sudo ausearch -m AVC -ts recent                 # ultimi blocchi
sudo journalctl -t setroubleshoot --since "10 min ago"   # spiegazioni leggibili (pacchetto setroubleshoot-server)
sudo sealert -a /var/log/audit/audit.log        # analisi con soluzioni suggerite
```

`sealert` suggerisce spesso il comando corretto (`restorecon`, `setsebool` o `semanage port`). Valuta sempre il
suggerimento: l'ultima opzione proposta, creare un modulo di policy personalizzato con `audit2allow`, è da usare
solo quando hai capito che l'accesso è davvero legittimo.

### Il metodo in quattro domande

1. Il problema sparisce in modalità permissive (`setenforce 0`)? Se no, non è SELinux. **Rimetti subito `setenforce 1`**.
2. Cosa dice `ausearch -m AVC`: quale processo (tipo `scontext`), su quale oggetto (tipo `tcontext`), quale azione?
3. È un **file** con il contesto sbagliato (→ `restorecon`/`semanage fcontext`), una **porta** non etichettata
   (→ `semanage port`) o una **funzione** disattivata (→ `setsebool -P`)?
4. Applica la correzione e verifica in enforcing.

## AppArmor (percorso Debian/Ubuntu)

Svolgi questa parte sulla VM `lab`. AppArmor confina i programmi con **profili** in `/etc/apparmor.d/`, scritti in
base ai **percorsi**: per esempio "questo programma può leggere `/etc/nginx/**` e scrivere `/var/log/nginx/*`".

```bash
sudo aa-status                     # profili caricati e in quale modalità (enforce o complain)
ls /etc/apparmor.d/
sudo apt install -y apparmor-utils # strumenti aa-*
sudo aa-complain /etc/apparmor.d/usr.sbin.tcpdump   # modalità complain: registra senza bloccare
sudo aa-enforce /etc/apparmor.d/usr.sbin.tcpdump    # torna in enforce
sudo journalctl -k | grep -i apparmor               # messaggi DENIED nel log del kernel
```

Un caso tipico: `tcpdump` ha un profilo AppArmor che gli permette di scrivere le catture solo in alcune posizioni.
Se provi a salvare un file `.pcap` in una directory insolita, ricevi *Permission denied* anche come root, e nel log
del kernel trovi un messaggio `apparmor="DENIED"`.

Su Ubuntu molti servizi (per esempio Nginx installato da apt) **non** hanno un profilo AppArmor attivo per default,
a differenza di SELinux che confina tutti i servizi previsti dalla policy. Lo stesso concetto si ritrova nei
container: Docker applica un profilo AppArmor predefinito ai container su Ubuntu ([Fase 3](../fase-3/index.md)).

## Laboratorio

!!! laboratorio "Laboratorio · Un sito che SELinux blocca · 45 minuti"
    **Dove:** `rhel-lab`, con Nginx installato e la porta 80 aperta in firewalld.

    1. Crea nella tua home `index.html` con il testo `ciao da /srv/sito`, poi spostalo con `mv` in una nuova
       directory `/srv/sito`.
    2. Configura Nginx (`/etc/nginx/conf.d/sito.conf`) perché serva `/srv/sito` sulla porta **8081**.
    3. Riavvia Nginx: perché non parte? Correggi.
    4. Richiedi la pagina con `curl http://localhost:8081/`: perché ottieni `403`? Correggi in modo permanente.
    5. Verifica che la correzione sopravviva a `sudo restorecon -Rv /srv` e a un riavvio.

??? success "Soluzione"
    ```bash
    echo "ciao da /srv/sito" > ~/index.html
    sudo mkdir -p /srv/sito && sudo mv ~/index.html /srv/sito/
    sudo tee /etc/nginx/conf.d/sito.conf > /dev/null <<'EOF'
    server {
        listen 8081;
        root /srv/sito;
    }
    EOF
    sudo nginx -t && sudo systemctl restart nginx    # fallisce
    sudo journalctl -u nginx -n 5                     # bind() to 0.0.0.0:8081 failed (13: Permission denied)
    sudo semanage port -l | grep -w 8081              # la porta ha già un'etichetta di un altro servizio?
    sudo semanage port -a -t http_port_t -p tcp 8081 \
      || sudo semanage port -m -t http_port_t -p tcp 8081   # -a aggiunge; se è già definita serve -m (modify)
    sudo systemctl restart nginx                      # ora parte

    curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/   # 403
    ls -Z /srv/sito/index.html                        # unconfined_u:object_r:user_home_t:s0
    sudo ausearch -m AVC -ts recent | tail -3         # denied { read } ... tcontext=...user_home_t
    sudo semanage fcontext -a -t httpd_sys_content_t "/srv/sito(/.*)?"
    sudo restorecon -Rv /srv/sito
    curl http://localhost:8081/                       # ciao da /srv/sito

    sudo restorecon -Rv /srv && curl -s http://localhost:8081/   # funziona ancora: la regola è nella policy
    ```
    **Spiegazione.** Due blocchi diversi: la **porta** 8081 non era etichettata `http_port_t` (nella policy può
    essere già assegnata a un altro tipo di servizio: per questo si controlla prima con `semanage port -l`), e il **file**
    spostato con `mv` conservava il contesto `user_home_t` della home. `chcon` avrebbe risolto il secondo problema
    solo fino alla prossima rietichettatura; `semanage fcontext` + `restorecon` lo risolve per sempre. Nota che i
    permessi classici (`ls -l`) erano corretti: per questo un `403` su RHEL va sempre verificato anche con `ls -Z`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è la differenza principale tra DAC e MAC?","options":["Il MAC vale anche per root ed è definito da una policy di sistema","Il DAC è più sicuro","Il MAC sostituisce i permessi classici","Non c'è differenza"],"answer":[0],"explanation":"Il MAC si aggiunge ai permessi classici e confina anche i processi di root."},
 {"q":"Quale modalità di SELinux registra le violazioni senza bloccarle?","options":["enforcing","permissive","disabled","complain"],"answer":[1],"explanation":"complain è il termine equivalente di AppArmor."},
 {"q":"Hai spostato con mv un file nella document root e il server web risponde 403. Qual è la soluzione corretta e permanente?","options":["setenforce 0","chmod 777 sul file","restorecon (con una regola semanage fcontext se la directory non è standard)","chcon -t httpd_sys_content_t"],"answer":[2],"explanation":"chcon funziona solo fino alla prossima rietichettatura."},
 {"q":"Nginx su RHEL deve fare da reverse proxy verso la porta 8000 e restituisce 502 con Permission denied. Cosa attivi?","options":["setsebool -P httpd_can_network_connect on","semanage port -a -t http_port_t -p tcp 8000","setenforce 0","chmod +x nginx"],"answer":[0],"explanation":"Il booleano permette ai processi httpd_t di aprire connessioni di rete verso altre porte."},
 {"q":"Quale comando mostra i blocchi recenti di SELinux?","options":["ausearch -m AVC -ts recent","aa-status","getsebool -a","ls -Z"],"answer":[0],"explanation":"I blocchi sono messaggi AVC nel log di audit."},
 {"q":"Su quale distribuzione AppArmor è il sistema MAC predefinito?","options":["RHEL","Ubuntu","AlmaLinux","Fedora"],"answer":[1],"explanation":"Ubuntu, Debian e SUSE usano AppArmor; la famiglia Red Hat usa SELinux."}
]}
</script>

## Fonti

- [Red Hat: Using SELinux (RHEL 10)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/using_selinux/index)
- [SELinux Project](https://github.com/SELinuxProject)
- [semanage(8)](https://man7.org/linux/man-pages/man8/semanage.8.html), [restorecon(8)](https://man7.org/linux/man-pages/man8/restorecon.8.html)
- [Ubuntu Server docs: AppArmor](https://documentation.ubuntu.com/server/how-to/security/apparmor/)
- [AppArmor wiki](https://gitlab.com/apparmor/apparmor/-/wikis/home)
- [The Linux Kernel: Linux Security Module](https://docs.kernel.org/admin-guide/LSM/index.html)
