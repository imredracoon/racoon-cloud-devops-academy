# 7.7 Dashboard e alert con Grafana

!!! obiettivi "Obiettivi"
    - Collegare Grafana a sorgenti dati diverse.
    - Costruire dashboard con pannelli, unità, soglie e variabili.
    - Gestire dashboard e sorgenti dati come codice (*provisioning*).
    - Creare regole di alert in Grafana e instradarle verso i punti di contatto.

## Cos'è Grafana

**Grafana** è lo strumento open source più diffuso per visualizzare dati: non conserva le metriche, ma interroga le
**sorgenti dati** (*data source*) e ne mostra i risultati in **dashboard**. Le sorgenti possono essere Prometheus,
Loki, CloudWatch, Cloud Monitoring, database SQL, Elasticsearch e decine di altre: una sola interfaccia per tutti i
segnali.

Nel laboratorio Grafana è già configurato: apri `http://localhost:3000`, accedi con `admin` / `admin` e scegli una
nuova password.

## Provisioning: Grafana come codice

Sorgenti dati e dashboard create a mano nell'interfaccia si perdono se il volume viene cancellato e non passano dalla
revisione. Con il **provisioning** Grafana le legge all'avvio da file versionati:

```yaml title="osservabilita/grafana/provisioning/datasources/datasources.yml"
apiVersion: 1
datasources:
  - name: Prometheus
    uid: prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
  - name: Loki
    uid: loki
    type: loki
    access: proxy
    url: http://loki:3100
```

```yaml title="osservabilita/grafana/provisioning/dashboards/dashboards.yml"
apiVersion: 1
providers:
  - name: academy
    folder: Academy
    type: file
    options:
      path: /var/lib/grafana/dashboards
```

Ogni file JSON nella cartella delle dashboard diventa una dashboard. Il JSON si ottiene da qualunque dashboard con
**Export → Export as JSON**. `access: proxy` significa che è il server di Grafana a interrogare la sorgente, non il
browser: per questo l'indirizzo è il nome del servizio sulla rete Docker.

## Anatomia di una dashboard

La dashboard **Academy → Tasks API** del laboratorio contiene:

| Pannello | Tipo | Query |
|---|---|---|
| API raggiungibile | Stat, con mappatura 0 → GIÙ (rosso), 1 → OK (verde) | `up{job="tasks-api"}` |
| Richieste al secondo | Stat | `sum(rate(http_requests_total{job="tasks-api"}[5m]))` |
| Errori 5xx (%) | Stat, soglia rossa al 5% | rapporto tra richieste 5xx e totali |
| Latenza p95 | Stat, unità secondi | `histogram_quantile(0.95, ...)` |
| Richieste per endpoint e stato | Time series | `sum by (method, endpoint, status) (rate(...))` |
| Latenza per endpoint | Time series, due query (p50 e p95) | `histogram_quantile` per endpoint |
| Log della Tasks API | Logs, da Loki | `{service="api"}` |

Alcune impostazioni da conoscere:

- **Unità** (*Standard options → Unit*): `reqps`, `s`, `percentunit`, `bytes`. Grafana formatta i valori di
  conseguenza (`0.05` diventa `5%`).
- **Legend**: `{{endpoint}} {{status}}` usa le etichette della serie come nome.
- **Soglie** (*Thresholds*) e **mappature dei valori** (*Value mappings*) colorano i pannelli.
- **Intervallo di tempo** e **aggiornamento automatico** in alto a destra.

### Le variabili

Una **variabile** rende la dashboard interattiva: un menu a tendina in alto che filtra tutti i pannelli. In
*Dashboard settings → Variables → New variable*:

- tipo **Query**, sorgente Prometheus, query `label_values(http_requests_total, endpoint)`, nome `endpoint`,
  *Multi-value* e *Include All option* attivi;
- nei pannelli usa il filtro `endpoint=~"$endpoint"`.

### Dashboard della comunità

Su [grafana.com/grafana/dashboards](https://grafana.com/grafana/dashboards/) trovi migliaia di dashboard pronte, da
importare con il loro numero (*Dashboards → New → Import*). Per esempio la popolare **Node Exporter Full** (ID
`1860`) mostra tutte le metriche di `node_exporter`. Come per i moduli Terraform, controllane la qualità prima di
adottarle.

## Alert in Grafana

Grafana ha un proprio sistema di alert (*Grafana Alerting*), utile quando vuoi allertare su più sorgenti dati
(compresi i log di Loki) da un unico posto. Gli elementi:

- **regola di alert**: una query, una condizione (per esempio "ultimo valore sopra 0,05"), un intervallo di
  valutazione e una durata (*pending period*);
- **punto di contatto** (*contact point*): dove inviare le notifiche (email, Slack, Microsoft Teams, webhook,
  PagerDuty...);
- **politica di notifica** (*notification policy*): quale punto di contatto riceve quali alert, in base alle
  etichette, con raggruppamento e tempi di ripetizione;
- **silenzi** e **finestre di manutenzione**.

Nel laboratorio usiamo **sia** le regole di Prometheus con Alertmanager ([lezione 7.6](06-prometheus.md)), lo standard negli ambienti
Kubernetes, **sia** Grafana Alerting per vederne le differenze. In un progetto reale scegli un sistema principale per
non duplicare gli alert.

## Laboratorio

!!! laboratorio "Laboratorio · La tua dashboard"
    Con lo stack della [lezione 7.6](06-prometheus.md) attivo e un po' di traffico generato:

    1. Apri **Dashboards → Academy → Tasks API** ed esplora i pannelli; prova gli intervalli `Last 5 minutes` e
       `Last 1 hour`.
    2. Crea una nuova dashboard **Sistema** con tre pannelli da `node_exporter`: CPU usata (%), memoria disponibile
       (byte), disco radice usato (%), con unità e soglie sensate.
    3. Aggiungi alla dashboard **Tasks API** (prima fai **Save as** per crearne una copia modificabile: quella
       caricata dal provisioning non si può salvare) la variabile `endpoint` e usala nel pannello delle richieste.
    4. Esporta la dashboard **Sistema** in JSON, salvala in `osservabilita/grafana/dashboards/sistema.json`, riavvia
       Grafana (`docker compose ... restart grafana`) e verifica che compaia nella cartella **Academy**.
    5. Crea un **punto di contatto** di tipo *Webhook* verso `http://alertmanager:9093/api/v2/alerts` oppure, più
       semplicemente, di tipo email (non verrà inviata senza un server SMTP, ma vedrai lo stato). Crea una **regola**
       "Errori 5xx sopra il 5%" sulla query della percentuale di errori, valutata ogni minuto, con pending period di 5
       minuti.
    6. Importa la dashboard **1860** (*Node Exporter Full*) e confrontala con la tua.

??? success "Soluzione: le query della dashboard Sistema"
    ```text
    CPU usata (unità percent (0-100)):
    100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])))

    Memoria disponibile (unità bytes(IEC)):
    node_memory_MemAvailable_bytes

    Disco radice usato (unità percentunit, soglia arancione 0.8 e rossa 0.9):
    1 - node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}
    ```

    **Spiegazione.** Nel Dev Container `node_exporter` gira dentro un container: vede il kernel della macchina che
    esegue Docker (in Codespaces la VM di GitHub, altrimenti il tuo computer), ma i dischi e i file system che vede
    sono quelli del container. I numeri sono reali, ma vanno letti sapendo da dove vengono. Su un server Linux vedresti il
    server stesso. Al punto 4 la dashboard ricompare dopo il riavvio perché ora è **codice** nel repository.

## Esercizi

!!! esercizio "Esercizio · Il pannello giusto"
    Quale tipo di pannello e quale unità useresti per:

    1. l'andamento delle richieste al secondo nelle ultime 24 ore;
    2. il valore attuale dello spazio libero su disco, con colore verde, giallo o rosso;
    3. le ultime righe di log con errori;
    4. la distribuzione delle latenze nel tempo, per vedere se ci sono "code" lente.

??? success "Soluzione"
    1. **Time series**, unità `reqps`.
    2. **Gauge** o **Stat**, unità `percentunit` o `bytes`, con soglie.
    3. **Logs**, sorgente Loki.
    4. **Heatmap** sui bucket dell'histogram (`sum by (le) (rate(..._bucket[5m]))`, formato *Heatmap*).

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Grafana conserva le metriche?","options":["Sì, nel suo database","No: interroga le sorgenti dati come Prometheus e Loki","Solo i log","Solo per 24 ore"],"answer":[1],"explanation":"Il database di Grafana contiene utenti, dashboard e impostazioni."},
 {"q":"A cosa serve il provisioning di Grafana?","options":["A creare sorgenti dati e dashboard da file versionati all'avvio","A installare plugin a pagamento","A fare backup delle metriche","A scalare Grafana"],"answer":[0],"explanation":"Grafana come codice, riproducibile."},
 {"q":"A cosa serve una variabile di dashboard?","options":["A filtrare tutti i pannelli con un menu, per esempio per endpoint o istanza","A salvare la password","A cambiare la sorgente dei log","A creare alert"],"answer":[0],"explanation":"Per esempio label_values(http_requests_total, endpoint)."},
 {"q":"In Grafana Alerting, dove si definisce chi riceve le notifiche?","options":["Nei punti di contatto e nelle politiche di notifica","Nel pannello","Nella sorgente dati","Nel file prometheus.yml"],"answer":[0],"explanation":"Le politiche instradano gli alert in base alle etichette."},
 {"q":"Quale unità mostra 0.05 come 5%?","options":["percent (0-100)","percentunit (0.0-1.0)","short","reqps"],"answer":[1],"explanation":"percent si aspetta valori già moltiplicati per 100."}
]}
</script>

## Fonti

- [Grafana documentation: Introduction](https://grafana.com/docs/grafana/latest/introduction/)
- [Grafana: Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/)
- [Grafana: Dashboards, panels and visualizations](https://grafana.com/docs/grafana/latest/visualizations/)
- [Grafana: Variables](https://grafana.com/docs/grafana/latest/dashboards/variables/) e [Prometheus template variables](https://grafana.com/docs/grafana/latest/datasources/prometheus/template-variables/)
- [Grafana Alerting](https://grafana.com/docs/grafana/latest/alerting/)
- [Grafana dashboards: Node Exporter Full](https://grafana.com/grafana/dashboards/1860-node-exporter-full/)
