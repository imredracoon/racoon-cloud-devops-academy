# 7.8 Log con Loki

!!! obiettivi "Obiettivi"
    - Descrivere le sfide dei log in sistemi distribuiti e l'approccio di Loki.
    - Raccogliere i log dei container con Grafana Alloy.
    - Scrivere query LogQL: selettori, filtri, parser, metriche dai log.
    - Collegare metriche e log durante un'indagine.
    - Scrivere log utili: livelli, formato strutturato, niente dati sensibili.

## Perché centralizzare i log

Con un solo server i log si leggono con `journalctl` o `tail -f`. Con dieci container che nascono e muoiono su più
macchine, i log devono essere **raccolti** in un unico posto, **conservati** anche dopo la fine dei container e
**interrogati** insieme.

Gli strumenti tradizionali (come Elasticsearch) **indicizzano tutto il testo** di ogni riga: ricerche potenti, ma
costi di memoria e disco elevati. **Grafana Loki** fa una scelta diversa:

- indicizza **solo le etichette** (`service="api"`, `container="..."`), come Prometheus;
- conserva il testo compresso, a blocchi, anche su storage a oggetti economico (S3, Cloud Storage);
- filtra il contenuto al momento della query.

Il risultato è molto più economico da gestire, a patto di scegliere bene le etichette: **poche**, con **pochi
valori**. Come per Prometheus, mai ID utente, IP o URL come etichette.

## Raccogliere i log: Grafana Alloy

**Grafana Alloy** è l'agente di raccolta di Grafana Labs (distribuzione di OpenTelemetry Collector): legge log, metriche
e tracce e li invia alle destinazioni. Ha sostituito **Promtail**, il vecchio agente di Loki, ora deprecato.

La configurazione del laboratorio raccoglie i log di tutti i container dal socket di Docker:

```alloy title="osservabilita/config.alloy"
// 1. Trova i container in esecuzione tramite il socket di Docker.
discovery.docker "contenitori" {
  host = "unix:///var/run/docker.sock"
}

// 2. Trasforma i metadati di Docker in etichette leggibili.
discovery.relabel "contenitori" {
  targets = []

  rule {
    source_labels = ["__meta_docker_container_name"]
    regex         = "/(.*)"
    target_label  = "container"
  }

  rule {
    source_labels = ["__meta_docker_container_label_com_docker_compose_service"]
    target_label  = "service"
  }
}

// 3. Legge i log dei container trovati.
loki.source.docker "contenitori" {
  host          = "unix:///var/run/docker.sock"
  targets       = discovery.docker.contenitori.targets
  relabel_rules = discovery.relabel.contenitori.rules
  forward_to    = [loki.write.locale.receiver]
}

// 4. Li invia a Loki.
loki.write "locale" {
  endpoint {
    url = "http://loki:3100/loki/api/v1/push"
  }
}
```

La configurazione è una catena di **componenti**: ognuno ha un tipo (`discovery.docker`), un nome (`"contenitori"`)
ed esporta valori usati dai successivi. L'interfaccia di Alloy (`http://localhost:12345`) mostra il grafo dei
componenti e il loro stato.

Fuori da Docker, Alloy legge i file di log (`local.file_match` e `loki.source.file`), il journal di systemd
(`loki.source.journal`) o i log dei pod Kubernetes.

## LogQL

**LogQL** è il linguaggio di Loki, ispirato a PromQL. Si prova in Grafana → **Explore** → sorgente Loki.

### Selezionare e filtrare

| Query | Significato |
|---|---|
| `{service="api"}` | tutte le righe del servizio `api` (il **selettore** tra graffe è obbligatorio) |
| <code>{service=~"api&#124;proxy"}</code> | più servizi, con un'espressione regolare |
| <code>{service="api"} &#124;= "POST"</code> | righe che **contengono** `POST` |
| `{service="api"} != "/health"` | righe che **non** contengono `/health` |
| <code>{service="api"} &#124;~ "HTTP/1.1\" 5[0-9]{2}"</code> | righe che corrispondono a un'espressione regolare (errori 5xx) |

### Estrarre campi con i parser

Le righe di accesso di Gunicorn hanno questo formato:

```text
172.20.0.3 - - [28/Oct/2026:13:41:48 +0000] "GET /tasks/9999 HTTP/1.1" 404 38 "-" "curl/8.5.0"
```

Il parser **`pattern`** estrae i campi descrivendo la forma della riga (`<_>` ignora una parte):

```text
{service="api"}
  | pattern `<ip> - - [<_>] "<metodo> <percorso> <_>" <stato> <_>`
  | stato >= 400
```

Dopo il parser `ip`, `metodo`, `percorso` e `stato` sono **etichette estratte**, utilizzabili nei filtri. Per i log in
JSON si usa `| json`, per il formato `chiave=valore` `| logfmt`; `| line_format` riscrive la riga visualizzata.

### Metriche dai log

```text
sum by (stato) (
  count_over_time(
    {service="api"} | pattern `<ip> - - [<_>] "<metodo> <percorso> <_>" <stato> <_>` | stato != "" [5m]
  )
)
```

Conta le richieste per codice di stato negli ultimi 5 minuti, partendo dai log. Funzioni utili: `count_over_time`,
`rate`, `bytes_over_time`, `sum by`, `topk`. Sono comode quando un'applicazione non espone metriche, ma costano più
delle metriche di Prometheus: per i numeri che guardi spesso, strumenta il codice.

## Dalla metrica al log

Il flusso tipico di un'indagine in Grafana:

1. un **alert** o un pannello mostra un picco di errori 5xx alle 14:05;
2. in **Explore**, con **Split**, affianchi la query Prometheus degli errori e la query Loki
   `{service="api"} | pattern ... | stato >= 500` sullo **stesso intervallo**;
3. leggi le righe di log di quei minuti e trovi la causa (per esempio il database irraggiungibile:
   la Tasks API risponde 503 su `/health`).

Il passo successivo, con le **tracce** (Grafana Tempo, X-Ray, Cloud Trace) e gli identificativi di richiesta nei log,
porta dalla riga di log al percorso completo della richiesta.

## Scrivere log utili

- **Livelli**: `DEBUG` (solo in sviluppo), `INFO` (eventi normali), `WARNING` (qualcosa di anomalo ma gestito),
  `ERROR` (un'operazione è fallita), `CRITICAL`.
- **Formato strutturato** (JSON o `chiave=valore`): più facile da interrogare di una frase libera.
- **Contesto**: un identificativo della richiesta, l'utente (come ID, non come email), l'operazione.
- **Mai** password, token, numeri di carte, dati personali non necessari: i log finiscono in molti posti e sono
  conservati a lungo. Il GDPR si applica anche ai log.
- Nei container, scrivi su **stdout/stderr**: raccolta e rotazione sono compito della piattaforma.

## Laboratorio

!!! laboratorio "Laboratorio · Indagine con metriche e log"
    Con lo stack attivo:

    1. Apri `http://localhost:12345` e osserva il grafo dei componenti di Alloy: tutti sani?
    2. In Grafana → **Explore** → Loki, prova `{service="api"}`, poi aggiungi `!= "/health"`.
    3. Con il parser `pattern`, mostra solo le richieste con stato `404`.
    4. Crea la query metrica che conta le richieste per stato e confrontala con la query Prometheus
       `sum by (status) (increase(http_requests_total[5m]))`.
    5. **Simula un incidente**: ferma il database
       (`docker compose -f compose.yaml -f compose.osservabilita.yaml stop db`), genera traffico su
       `localhost:8080/tasks` e `localhost:8080/health` per un minuto, poi riavvia il database.
    6. Indaga come se non sapessi nulla: nella dashboard **Tasks API** individua il momento del problema e i codici di
       errore; in Explore, con Split, trova le righe di log corrispondenti.

??? success "Cosa dovresti trovare"
    3. La query:

        ```text
        {service="api"} | pattern `<ip> - - [<_>] "<metodo> <percorso> <_>" <stato> <_>` | stato = "404"
        ```

    4. I due numeri sono simili ma non identici: le metriche contano le richieste nell'applicazione, i log le righe
       scritte da Gunicorn (comprese quelle di `/metrics`, escluse dalle metriche), e gli intervalli non coincidono al
       secondo.
    6. Nel pannello degli errori compaiono risposte `500` su `/tasks` e `503` su `/health`; l'healthcheck di Docker
       segna l'API come `unhealthy`. Nei log dell'API trovi le righe di accesso con quei codici e, nel servizio `db`,
       le righe di arresto e riavvio di PostgreSQL. La sequenza ricostruisce l'incidente: database giù, API in errore,
       ripristino al riavvio del database, senza dover riavviare l'API.

    **Pulizia.** `docker compose -f compose.yaml -f compose.osservabilita.yaml down` (con `-v` cancelli anche i
    volumi con dati, metriche e log).

## Esercizi

!!! esercizio "Esercizio · Etichette o contenuto?"
    Per ciascun campo decidi se è una buona **etichetta** di Loki o se deve restare nel **contenuto** della riga (ed
    essere estratto al momento della query): nome del servizio, ambiente, ID della richiesta, indirizzo IP del
    client, livello di log, percorso dell'URL.

??? success "Soluzione"
    **Etichette**: servizio, ambiente, livello di log (pochi valori, usati sempre per filtrare). **Contenuto**: ID
    della richiesta, IP del client, percorso dell'URL (valori praticamente infiniti: come etichette moltiplicherebbero
    gli stream e degraderebbero Loki). Si estraggono con `pattern`, `json` o `logfmt` quando servono.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa indicizza Loki?","options":["Tutto il testo di ogni riga","Solo le etichette degli stream","Nulla","Solo i timestamp"],"answer":[1],"explanation":"Il contenuto viene filtrato al momento della query: è più economico."},
 {"q":"Quale agente raccoglie i log nel laboratorio, al posto del deprecato Promtail?","options":["Grafana Alloy","Fluent Bit","Logstash","rsyslog"],"answer":[0],"explanation":"Alloy raccoglie log, metriche e tracce."},
 {"q":"Quale query LogQL mostra le righe del servizio api che contengono ERROR?","options":["{service=\"api\"} |= \"ERROR\"","service=api AND ERROR","SELECT * FROM api WHERE ERROR","{ERROR}"],"answer":[0],"explanation":"Il selettore tra graffe è sempre obbligatorio."},
 {"q":"A cosa serve | pattern in LogQL?","options":["A estrarre campi da righe con un formato fisso","A cifrare i log","A cancellare le righe","A creare etichette indicizzate"],"answer":[0],"explanation":"I campi estratti si usano nei filtri e nelle aggregazioni, senza indicizzarli."},
 {"q":"Quali dati non devono mai finire nei log?","options":["Password e token","Numeri di carte di pagamento","Il codice di stato HTTP","Il nome del servizio","La durata della richiesta"],"answer":[0,1],"explanation":"Anche i dati personali vanno limitati al necessario (GDPR)."},
 {"q":"Come trasformi i log in una metrica in LogQL?","options":["Con funzioni come count_over_time o rate su un intervallo","Non si può","Con un export in CSV","Con Alertmanager"],"answer":[0],"explanation":"Per esempio sum by (stato) (count_over_time({...}[5m]))."}
]}
</script>

## Fonti

- [Grafana Loki: Overview](https://grafana.com/docs/loki/latest/get-started/overview/) e [Labels best practices](https://grafana.com/docs/loki/latest/get-started/labels/bp-labels/)
- [Grafana Loki: LogQL](https://grafana.com/docs/loki/latest/query/), [Log queries](https://grafana.com/docs/loki/latest/query/log_queries/), [Metric queries](https://grafana.com/docs/loki/latest/query/metric_queries/)
- [Grafana Alloy: Introduction](https://grafana.com/docs/alloy/latest/introduction/)
- [Grafana Alloy: loki.source.docker](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.source.docker/) e [discovery.docker](https://grafana.com/docs/alloy/latest/reference/components/discovery/discovery.docker/)
- [Grafana Loki: Promtail deprecation and migration to Alloy](https://grafana.com/docs/loki/latest/send-data/promtail/)
- [Grafana: Explore, logs and metrics side by side](https://grafana.com/docs/grafana/latest/explore/)
- [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)
