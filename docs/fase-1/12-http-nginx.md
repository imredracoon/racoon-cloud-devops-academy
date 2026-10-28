# 1.12 HTTP, TLS e Nginx

!!! obiettivi "Obiettivi"
    - Leggere richieste e risposte HTTP: metodi, codici di stato, intestazioni.
    - Capire come TLS protegge HTTPS: certificati, autorità di certificazione, catena di fiducia.
    - Configurare Nginx come server web e come reverse proxy.
    - Attivare HTTPS con un certificato e reindirizzare HTTP su HTTPS.

## HTTP

**HTTP** è un protocollo testuale di richiesta e risposta. Ecco una richiesta e la sua risposta, così come
viaggiano (le vedi con `curl -v`):

```http
GET /api/tasks HTTP/1.1
Host: app.example.com
User-Agent: curl/8.5.0
Accept: application/json
```

```http
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 27
Cache-Control: no-store

[{"id": 1, "title": "Ciao"}]
```

### Metodi

| Metodo | Uso |
|---|---|
| `GET` | Leggere una risorsa. Non deve modificare nulla. |
| `POST` | Creare una risorsa o inviare dati. |
| `PUT` / `PATCH` | Sostituire / modificare in parte una risorsa. |
| `DELETE` | Cancellare una risorsa. |
| `HEAD` | Come GET, ma solo intestazioni (`curl -I`). |

### Codici di stato

| Classe | Significato | Esempi da ricordare |
|---|---|---|
| **2xx** | Successo | `200 OK`, `201 Created`, `204 No Content` |
| **3xx** | Reindirizzamento | `301 Moved Permanently`, `302 Found`, `304 Not Modified` |
| **4xx** | Errore del **client** | `400 Bad Request`, `401 Unauthorized` (non autenticato), `403 Forbidden` (autenticato ma senza permesso), `404 Not Found`, `429 Too Many Requests` |
| **5xx** | Errore del **server** | `500 Internal Server Error`, `502 Bad Gateway`, `503 Service Unavailable`, `504 Gateway Timeout` |

!!! tip "502 e 504 parlano del backend"
    Dietro un reverse proxy, `502` significa "il proxy ha ricevuto una risposta non valida, o nessuna connessione,
    dal servizio dietro di lui" (spesso il servizio è fermo) e `504` "il servizio dietro di me non ha risposto in
    tempo". Il problema quasi mai è nel proxy: guarda l'applicazione.

### Intestazioni importanti

`Host` (quale sito si vuole: permette di ospitare molti siti sullo stesso IP), `Content-Type`, `Location` (nei
reindirizzamenti), `Cache-Control`, `Set-Cookie`/`Cookie`, `Authorization`, `X-Forwarded-For` (aggiunta dai proxy
con l'IP originale del client).

## TLS e HTTPS

**HTTPS** è HTTP trasportato dentro **TLS** (*Transport Layer Security*). TLS garantisce:

- **riservatezza**: il traffico è cifrato;
- **integrità**: non può essere modificato senza che ce ne si accorga;
- **autenticazione**: il client verifica di parlare con il server giusto.

L'autenticazione si basa sui **certificati X.509**: un documento che lega un nome (per esempio `app.example.com`)
a una chiave pubblica, **firmato** da un'**autorità di certificazione** (CA). Il browser e il sistema operativo
contengono l'elenco delle CA radice fidate; un certificato è valido se esiste una **catena** di firme che porta
a una di esse:

```text
CA radice (nel sistema) → CA intermedia → certificato di app.example.com
```

Il server deve inviare il proprio certificato **più gli intermedi** (*fullchain*). Una catena incompleta è un
errore comune: alcuni client funzionano, altri no.

**Let's Encrypt** è una CA gratuita e automatizzata: con lo strumento `certbot` ottieni e rinnovi certificati
validi per un dominio che controlli. Senza un dominio pubblico (come nel laboratorio) si usa un certificato
**autofirmato**: cifra correttamente, ma nessuna CA lo garantisce, quindi i client mostrano un avviso.

```bash
openssl s_client -connect example.com:443 -servername example.com </dev/null | openssl x509 -noout -subject -issuer -dates
```

Questo comando mostra a chi è intestato il certificato di un sito, chi l'ha firmato e quando scade. I certificati
scaduti sono una delle cause più frequenti di disservizio: si monitorano ([Fase 7](../fase-7/index.md)).

## Nginx

**Nginx** è un server web e **reverse proxy** ad alte prestazioni, tra i più usati al mondo.

```bash
sudo apt install -y nginx
systemctl status nginx
curl -I http://localhost
```

### La struttura della configurazione (Debian/Ubuntu)

| Percorso | Contenuto |
|---|---|
| `/etc/nginx/nginx.conf` | Configurazione principale. |
| `/etc/nginx/sites-available/` | Un file per ogni sito. |
| `/etc/nginx/sites-enabled/` | **Link simbolici** ai siti attivi. |
| `/var/www/` | File dei siti. |
| `/var/log/nginx/access.log`, `error.log` | Log degli accessi e degli errori. |

=== "Debian/Ubuntu"

    La struttura è quella della tabella: un file per sito in `sites-available`, attivato con un link simbolico in
    `sites-enabled`. Nginx parte da solo dopo l'installazione e il firewall va aperto a mano (ufw o nftables).

=== "RHEL (Rocky/AlmaLinux)"

    - Non esistono `sites-available` e `sites-enabled`: ogni sito è un file **`/etc/nginx/conf.d/NOME.conf`**.
    - Dopo l'installazione il servizio va avviato e abilitato: `sudo systemctl enable --now nginx`.
    - Il firewall va aperto con firewalld: `sudo firewall-cmd --permanent --add-service={http,https} && sudo firewall-cmd --reload`.
    - **SELinux** ([lezione 1.14](14-selinux-apparmor.md)) controlla cosa può fare Nginx. Due casi da ricordare:
        - come reverse proxy, Nginx deve potersi collegare a un'altra porta: `sudo setsebool -P httpd_can_network_connect on`.
          Senza, ottieni un `502` e nel log `connect() ... failed (13: Permission denied)`;
        - i file del sito fuori da `/usr/share/nginx/html` o `/var/www` devono avere il contesto SELinux corretto
          (`httpd_sys_content_t`), altrimenti Nginx risponde `403`.
    - Il server web "di casa" Red Hat è **Apache httpd** (pacchetto `httpd`), che troverai spesso nelle aziende e
      negli esercizi RHCSA. I concetti (virtual host, reverse proxy, TLS) sono gli stessi.

```bash
sudo nginx -t                 # controlla la configurazione: fallo SEMPRE prima di ricaricare
sudo systemctl reload nginx   # applica senza interrompere le connessioni
```

### Un sito statico

```nginx title="/etc/nginx/sites-available/sito.lab"
server {
    listen 80;
    server_name sito.lab;

    root /var/www/sito.lab;
    index index.html;

    location / {
        try_files $uri $uri/ =404;
    }
}
```

```bash
sudo mkdir -p /var/www/sito.lab
echo "<h1>Ciao da sito.lab</h1>" | sudo tee /var/www/sito.lab/index.html
sudo ln -s /etc/nginx/sites-available/sito.lab /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
curl -H "Host: sito.lab" http://localhost/
```

Nginx sceglie il blocco `server` confrontando l'intestazione `Host` con `server_name`: è così che un solo server
ospita molti siti (*virtual host*).

### Reverse proxy

Un **reverse proxy** riceve le richieste dei client e le inoltra a uno o più servizi interni (*backend* o
*upstream*). Vantaggi: un solo punto d'ingresso, HTTPS gestito in un solo posto, bilanciamento del carico,
applicazioni che ascoltano solo su `127.0.0.1`.

```nginx title="/etc/nginx/sites-available/app.lab"
upstream app_backend {
    server 127.0.0.1:8000;
}

server {
    listen 80;
    server_name app.lab;

    location / {
        proxy_pass http://app_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Le intestazioni `proxy_set_header` dicono all'applicazione chi è davvero il client e se la richiesta originale era
in HTTPS: senza, l'app vedrebbe sempre `127.0.0.1` come client.

### HTTPS con certificato autofirmato

```bash
sudo mkdir -p /etc/nginx/tls
sudo openssl req -x509 -newkey rsa:2048 -nodes -days 365 \
  -keyout /etc/nginx/tls/app.lab.key -out /etc/nginx/tls/app.lab.crt \
  -subj "/CN=app.lab" -addext "subjectAltName=DNS:app.lab"
sudo chmod 600 /etc/nginx/tls/app.lab.key
```

```nginx title="/etc/nginx/sites-available/app.lab (versione HTTPS)"
server {
    listen 80;
    server_name app.lab;
    return 301 https://$host$request_uri;      # tutto il traffico HTTP va su HTTPS
}

server {
    listen 443 ssl;
    server_name app.lab;

    ssl_certificate     /etc/nginx/tls/app.lab.crt;
    ssl_certificate_key /etc/nginx/tls/app.lab.key;
    ssl_protocols       TLSv1.2 TLSv1.3;

    add_header Strict-Transport-Security "max-age=31536000" always;

    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
curl -k --resolve app.lab:443:127.0.0.1 https://app.lab/   # -k accetta il certificato autofirmato
```

!!! warning "-k solo in laboratorio"
    `curl -k` disattiva la verifica del certificato. In laboratorio, con un certificato autofirmato, è accettabile;
    in produzione mai: annullerebbe la protezione di TLS.

## Laboratorio

!!! laboratorio "Laboratorio · Reverse proxy con HTTPS · 1,5 ore"
    **Dove:** VM `lab` (con il firewall della [lezione 1.11](11-firewall.md)) e VM `client`.

    1. Su `lab` avvia un'applicazione di prova che ascolta solo su `127.0.0.1:8000`:
       `python3 -m http.server 8000 --bind 127.0.0.1 --directory /etc` (mostra il contenuto di `/etc`: va bene per
       una prova, mai in produzione).
    2. Configura Nginx come reverse proxy per `app.lab` con HTTPS e reindirizzamento da HTTP.
    3. Su `client` aggiungi a `/etc/hosts` la riga che fa risolvere `app.lab` nell'IP di `lab`.
    4. Da `client` verifica: `http://app.lab` risponde con `301`; `https://app.lab` risponde `200`; la porta 8000
       di `lab` non è raggiungibile da `client`.
    5. Ferma l'applicazione e ripeti la richiesta HTTPS: che codice ottieni? In quale log di `lab` trovi la causa?

    **Percorso RHEL:** ripeti su `rhel-lab` con il file `/etc/nginx/conf.d/app.lab.conf`. Prima di attivare il
    boolean SELinux prova la richiesta: che codice ottieni e cosa dice il log degli errori?

??? success "Soluzione"
    ```bash
    # su lab
    python3 -m http.server 8000 --bind 127.0.0.1 --directory /etc &
    # crea certificato e file di configurazione come negli esempi, poi:
    sudo ln -s /etc/nginx/sites-available/app.lab /etc/nginx/sites-enabled/
    sudo nginx -t && sudo systemctl reload nginx

    # su client
    echo "IP-LAB app.lab" | sudo tee -a /etc/hosts
    curl -sI http://app.lab | head -3          # HTTP/1.1 301 Moved Permanently, Location: https://app.lab/
    curl -skI https://app.lab | head -1        # HTTP/1.1 200 OK
    curl -m 5 http://IP-LAB:8000/              # timeout: il firewall scarta e comunque l'app ascolta solo su 127.0.0.1

    # su lab: ferma l'app
    kill %1
    # su client
    curl -skI https://app.lab | head -1        # HTTP/1.1 502 Bad Gateway
    # su lab
    sudo tail -n 3 /var/log/nginx/error.log    # connect() failed (111: Connection refused) while connecting to upstream
    ```
    **Spiegazione.** Il `502` viene da Nginx, che non riesce a collegarsi al backend: l'errore 111 è proprio il
    *Connection refused* della [lezione 1.9](09-diagnostica.md). Questo è il ragionamento a strati che userai nel troubleshooting.

    **Percorso RHEL.** Con l'applicazione attiva ma senza il boolean ottieni comunque `502`, e in
    `/var/log/nginx/error.log` trovi `connect() to 127.0.0.1:8000 failed (13: Permission denied)`. L'errore 13 non
    viene dai permessi dei file ma da **SELinux**: `sudo ausearch -m avc -ts recent` mostra il blocco. Si risolve
    con `sudo setsebool -P httpd_can_network_connect on`. Stesso sintomo, causa diversa: per questo si leggono
    sempre i log.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Un utente non autenticato prova ad accedere a una risorsa protetta. Quale codice è corretto?","options":["200","401","403","500"],"answer":[1],"explanation":"401 significa che serve autenticazione; 403 che sei autenticato ma non autorizzato."},
 {"q":"Nginx restituisce 502 Bad Gateway. Dove cerchi prima il problema?","options":["Nel DNS del client","Nel servizio backend dietro Nginx","Nel browser","Nel certificato"],"answer":[1],"explanation":"502 indica che il proxy non ha ottenuto una risposta valida dal backend."},
 {"q":"Cosa permette a Nginx di ospitare più siti sullo stesso IP e porta?","options":["L'intestazione Host confrontata con server_name","Il record MX","La porta del client","Il codice di stato"],"answer":[0],"explanation":"È il meccanismo dei virtual host basati sul nome."},
 {"q":"Prima di ricaricare Nginx dopo una modifica cosa esegui?","options":["nginx -t","systemctl stop nginx","rm sites-enabled/*","curl -k"],"answer":[0],"explanation":"nginx -t verifica la sintassi e i riferimenti ai file."},
 {"q":"Perché un certificato autofirmato genera un avviso nel browser?","options":["Perché non cifra","Perché nessuna CA fidata lo ha firmato","Perché usa la porta 443","Perché è scaduto per forza"],"answer":[1],"explanation":"La cifratura funziona, ma l'identità non è garantita da una CA."},
 {"q":"Quali sono garanzie di TLS?","options":["Riservatezza","Integrità","Autenticazione del server","Disponibilità del servizio","Velocità maggiore"],"answer":[0,1,2],"explanation":"TLS non garantisce che il servizio sia disponibile."}
]}
</script>

## Fonti

- [MDN: HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP) e [codici di stato](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status)
- [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 8446: TLS 1.3](https://www.rfc-editor.org/rfc/rfc8446)
- [Nginx: documentazione ufficiale](https://nginx.org/en/docs/), [Beginner's Guide](https://nginx.org/en/docs/beginners_guide.html), [ngx_http_proxy_module](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
- [Let's Encrypt: How It Works](https://letsencrypt.org/how-it-works/) e [Certbot](https://certbot.eff.org/)
- [Mozilla: Server Side TLS](https://wiki.mozilla.org/Security/Server_Side_TLS)
- [Red Hat: Deploying web servers and reverse proxies (RHEL 10)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/deploying_web_servers_and_reverse_proxies/index)
