# 1.13 Metodo di troubleshooting

!!! obiettivi "Obiettivi"
    - Applicare un metodo strutturato alla risoluzione dei problemi.
    - Diagnosticare per strati, dal basso verso l'alto.
    - Usare un elenco di controlli rapidi per lo stato di un server.
    - Scrivere un runbook e un resoconto di incidente.

## Perché serve un metodo

Davanti a un guasto la tentazione è provare cose a caso: riavviare, cambiare configurazioni, cercare l'errore su
Internet e incollare il primo comando trovato. Funziona a volte, ma spesso peggiora le cose e non insegna nulla.
I professionisti seguono un **metodo**: più lento nei primi minuti, molto più veloce nel complesso.

## Il ciclo in sei passi

1. **Definisci il problema** con precisione. "Il sito non va" non basta: chi è colpito? Da quando? Cosa succede
   esattamente (errore, lentezza, pagina bianca)? Cosa è cambiato di recente?
2. **Raccogli i fatti**: log, metriche, output dei comandi. Annota l'ora di ogni osservazione.
3. **Formula un'ipotesi** che spieghi **tutti** i fatti.
4. **Verificala** con un test che possa smentirla, cambiando **una cosa alla volta**.
5. **Correggi** e verifica che il sintomo originale sia sparito.
6. **Documenta**: cosa è successo, perché, come l'hai risolto e come evitare che si ripeta.

!!! tip "Cosa è cambiato?"
    La maggior parte dei guasti segue un **cambiamento**: un rilascio, un aggiornamento, una modifica di
    configurazione, un certificato scaduto, un disco che si è riempito. Chiedere "cosa è cambiato e quando?" è la
    domanda più produttiva in assoluto.

## Diagnosticare per strati

Per un servizio di rete si risale la pila, dal basso verso l'alto. Ogni domanda ha il suo comando:

| # | Domanda | Comando |
|---|---|---|
| 1 | La macchina è accesa e ha risorse? | `uptime`, `free -h`, `df -h`, `top` |
| 2 | L'interfaccia di rete è attiva, con IP e rotta? | `ip -br a`, `ip route` |
| 3 | Il nome si risolve nell'IP giusto? | `dig +short nome`, `getent hosts nome` |
| 4 | L'host è raggiungibile? | `ping`, `tracepath` |
| 5 | La porta è aperta e raggiungibile? | dal client `nc -zv host porta`; sul server `ss -tlpn` |
| 6 | Un firewall blocca? | `sudo nft list ruleset`, regole cloud |
| 7 | Il servizio è attivo? | `systemctl status`, `journalctl -u` |
| 8 | Il servizio risponde correttamente? | `curl -v`, log dell'applicazione |

Puoi anche partire dall'alto e scendere, ma l'importante è **non saltare** passaggi e annotare cosa hai escluso.

## I primi 60 secondi su un server

Quando entri su un server con un problema, questi comandi danno in un minuto un quadro generale:

```bash
uptime                          # da quanto è acceso, carico medio
dmesg -T | tail -20             # ultimi messaggi del kernel: errori di disco, memoria esaurita (OOM)
journalctl -p err -b --no-pager | tail -20   # errori dall'avvio
systemctl --failed              # servizi in errore
free -h                         # memoria e swap
df -h; df -i                    # spazio e inode
top -b -n1 | head -20           # processi che consumano di più
ss -tlpn                        # cosa ascolta
ip -br a; ip route              # rete
```

Il messaggio `Out of memory: Killed process` in `dmesg` significa che il kernel ha terminato un processo per
mancanza di RAM (l'**OOM killer**): spiega molti servizi "spariti senza motivo".

## Runbook e resoconti

Un **runbook** è una procedura scritta per affrontare un problema noto: chi lo segue, anche alle 3 di notte, sa
cosa controllare e in che ordine. Un buon runbook ha: sintomo, impatto, controlli in ordine, azioni correttive,
quando e a chi passare il problema (**escalation**).

Dopo un incidente si scrive un **resoconto** (*post-mortem*): cronologia, causa principale, impatto, azioni
correttive. La cultura **blameless** (senza colpevoli) descritta nel libro SRE di Google si concentra sui
processi e sugli strumenti, non sulle persone: così le persone raccontano ciò che è successo davvero e il sistema
migliora.

## Laboratorio: tre guasti

!!! laboratorio "Laboratorio · Guasti da risolvere · 2 ore"
    **Dove:** VM `lab` con la configurazione della [lezione 1.12](12-http-nginx.md) funzionante (Nginx in HTTPS davanti all'app sulla
    porta 8000) e VM `client`. Fai uno **snapshot** prima di iniziare.

    Per ogni scenario: esegui su `lab` il comando che introduce il guasto **senza leggerlo con attenzione** (meglio
    ancora: fallo eseguire a un'altra persona), poi da `client` osserva il sintomo e trova la causa usando il metodo.
    Scrivi i passaggi seguiti prima di aprire la soluzione.

    **Scenario A**
    ```bash
    sudo sed -i 's/127.0.0.1:8000/127.0.0.1:8001/' /etc/nginx/sites-available/app.lab && sudo systemctl reload nginx
    ```

    **Scenario B**
    ```bash
    sudo nft insert rule inet filtro input tcp dport 443 drop
    ```

    **Scenario C**
    ```bash
    sudo sed -i 's/IP-LAB/10.255.255.1/' /etc/hosts   # da eseguire su client, sostituendo IP-LAB con l'IP reale
    ```

??? success "Soluzione: scenario A"
    **Sintomo:** `curl -sk https://app.lab` → `502 Bad Gateway`. **Strati 1-6 ok** (la risposta arriva da Nginx).
    **Strato 7-8:** `sudo tail /var/log/nginx/error.log` → `connect() failed (111: Connection refused) while
    connecting to upstream, upstream: "http://127.0.0.1:8001/"`. `ss -tlpn` mostra l'app sulla 8000.
    **Causa:** porta del backend sbagliata nella configurazione. **Correzione:** ripristina `8000`, `nginx -t`, reload.

??? success "Soluzione: scenario B"
    **Sintomo:** `curl -k https://app.lab` resta appeso e va in timeout; `http://app.lab` restituisce il 301.
    **Indizi:** il DNS risolve, `ping` risponde, `nc -zv app.lab 443` va in timeout ma la porta 80 risponde.
    Sul server `ss -tlpn` mostra Nginx in ascolto sulla 443, quindi il servizio c'è. Un timeout con servizio in
    ascolto indica un **firewall**: `sudo nft -a list chain inet filtro input` mostra la regola `drop` in cima.
    **Correzione:** `sudo nft delete rule inet filtro input handle N`.

??? success "Soluzione: scenario C"
    **Sintomo:** da `client` qualsiasi richiesta ad `app.lab` va in timeout, mentre `curl -k https://IP-LAB` con
    l'IP funziona. **Indizio decisivo:** `getent hosts app.lab` restituisce `10.255.255.1`, un indirizzo che non è
    quello di `lab`. `dig app.lab` non aiuta, perché non legge `/etc/hosts` ([lezione 1.8](08-dns-routing.md)).
    **Causa:** risoluzione del nome sbagliata. **Correzione:** ripristina l'IP corretto in `/etc/hosts`.

    **Lezione generale:** quando con l'IP funziona e con il nome no, il problema è nella risoluzione dei nomi.

## Esercizio finale

!!! esercizio "Esercizio · Il tuo primo runbook"
    Scrivi in Markdown il runbook "Il sito app.lab non risponde", basandoti sui tre scenari. Deve contenere sintomi,
    controlli in ordine con i comandi, interpretazione dei risultati e azioni correttive. Salvalo nel repository
    `academy-lab`: lo riutilizzerai nel progetto della fase.

??? success "Soluzione (traccia)"
    ```markdown
    # Runbook: app.lab non risponde

    ## Sintomi
    Errore nel browser, timeout o codice 5xx.

    ## Controlli
    1. Dal client: `getent hosts app.lab` → l'IP deve essere quello del server. Se no: DNS o /etc/hosts.
    2. `curl -skv https://app.lab -o /dev/null`
       - timeout → passo 3
       - 502/504 → passo 5
    3. `nc -zv app.lab 443`. Timeout: firewall (passo 4). Refused: Nginx fermo (`systemctl status nginx`).
    4. Sul server: `sudo nft list ruleset`, cerca regole drop sulla 443. Controlla anche il firewall cloud.
    5. Sul server: `sudo tail /var/log/nginx/error.log`, `ss -tlpn | grep 8000`, `systemctl status app`.

    ## Escalation
    Se l'applicazione è in errore dopo un rilascio, contatta il team di sviluppo con l'estratto dei log.
    ```

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è la domanda più utile all'inizio di un'indagine?","options":["Chi è il colpevole?","Cosa è cambiato di recente?","Quale comando ho trovato su Internet?","Posso riavviare tutto?"],"answer":[1],"explanation":"La maggior parte dei guasti segue un cambiamento."},
 {"q":"Perché si cambia una sola cosa alla volta?","options":["Per risparmiare tempo","Per sapere quale modifica ha prodotto l'effetto","Perché è obbligatorio per legge","Per non scrivere documentazione"],"answer":[1],"explanation":"Più modifiche insieme rendono impossibile capire la causa."},
 {"q":"Con l'IP il sito funziona, con il nome no. Dove si trova il problema?","options":["Nel firewall","Nella risoluzione dei nomi","Nel certificato","Nel disco"],"answer":[1],"explanation":"La differenza tra i due casi è solo la risoluzione del nome."},
 {"q":"Quale comando ti mostra se il kernel ha terminato processi per mancanza di memoria?","options":["dmesg","ss","dig","ip route"],"answer":[0],"explanation":"I messaggi dell'OOM killer compaiono nel log del kernel."},
 {"q":"Cos'è un post-mortem blameless?","options":["Un resoconto che individua chi punire","Un resoconto dell'incidente centrato su cause e miglioramenti, non sulle colpe","Un log automatico","Un riavvio del server"],"answer":[1],"explanation":"La cultura blameless favorisce la trasparenza e il miglioramento."}
]}
</script>

## Fonti

- [Google SRE Book: Effective Troubleshooting](https://sre.google/sre-book/effective-troubleshooting/)
- [Google SRE Book: Postmortem Culture](https://sre.google/sre-book/postmortem-culture/)
- [dmesg(1)](https://man7.org/linux/man-pages/man1/dmesg.1.html)
- [Nginx: Debugging and logging](https://nginx.org/en/docs/debugging_log.html)
