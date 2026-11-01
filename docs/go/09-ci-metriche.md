# G.9 Pipeline ed exporter Prometheus

!!! obiettivi "Obiettivi"
    - Scrivere una pipeline GitHub Actions per un progetto Go: formattazione, `go vet`, linter, test con race
      detector, copertura e immagine multi-architettura.
    - Cercare vulnerabilità note nelle dipendenze con `govulncheck`.
    - Esporre metriche Prometheus da un programma Go con `client_golang`: gauge, counter e istogrammi con label.
    - Trasformare racoon-probe in un **exporter** che Prometheus legge e su cui Grafana mostra gli alert.

!!! info "Quando studiarla"
    Nella [Fase 7](../fase-7/index.md), dopo la [lezione 7.6](../fase-7/06-prometheus.md). Hai già scritto la pipeline
    della Tasks API ([lezione 7.3](../fase-7/03-pipeline.md)) e l'hai strumentata in Python: qui fai lo stesso in Go.

## La pipeline di un progetto Go

La pipeline segue la stessa logica di quella della Tasks API: prima la qualità, poi l'immagine, che si pubblica solo
dai push su `main` e dai tag. Il file completo è
[`labs/go/probe/.github/workflows/ci.yml`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/probe/.github/workflows/ci.yml).

```yaml title=".github/workflows/ci.yml (job qualità)"
  qualita:
    name: Vet, lint e test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7

      - uses: actions/setup-go@v6
        with:
          go-version-file: go.mod     # (1)!

      - name: Formattazione
        run: test -z "$(gofmt -l .)" || { gofmt -l .; exit 1; }   # (2)!

      - name: go vet
        run: go vet ./...

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v9
        with:
          version: v2.15              # (3)!

      - name: Test con race detector e copertura
        run: go test -race -coverprofile=coverage.out ./...

      - name: Riepilogo copertura
        run: go tool cover -func=coverage.out | tail -1
```

1.  La versione di Go la decide `go.mod`, non la pipeline: una sola fonte di verità. `setup-go` mette anche in
    cache i moduli scaricati, in base a `go.sum`.
2.  `gofmt -l` non fallisce da solo se trova file non formattati: li elenca. Il test `-z` (stringa vuota) trasforma
    l'elenco in un errore.
3.  **golangci-lint** esegue decine di analizzatori insieme (errori non controllati, codice irraggiungibile,
    problemi di sicurezza). Fissare la versione evita che un aggiornamento del linter faccia fallire la pipeline
    senza che tu abbia cambiato nulla.

Il job `immagine` costruisce l'immagine della [lezione G.6](06-container.md) per `linux/amd64` e `linux/arm64` con
le azioni Docker già usate nella [lezione 7.3](../fase-7/03-pipeline.md), e la pubblica su GitHub Container Registry
con il permesso `packages: write` limitato a quel job.

### Vulnerabilità nelle dipendenze

`govulncheck` è lo strumento ufficiale del team Go. Confronta le dipendenze con il database delle vulnerabilità di
Go e, a differenza di molti scanner, segnala solo quelle nel codice che il tuo programma **chiama davvero**:

```console
$ go install golang.org/x/vuln/cmd/govulncheck@latest
$ govulncheck ./...
No vulnerabilities found.
```

Aggiungilo alla pipeline come passo del job di qualità:

```yaml
      - name: govulncheck
        run: go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## Metriche Prometheus da Go

Prometheus stesso è scritto in Go, e la libreria client ufficiale per Go è
[`client_golang`](https://github.com/prometheus/client_golang). I tipi di metrica sono quelli che conosci dalla
[lezione 7.6](../fase-7/06-prometheus.md#i-tipi-di-metrica):

| Tipo | In Go | In probe |
|---|---|---|
| Gauge | `prometheus.NewGaugeVec` | `probe_up`: 1 se l'ultimo controllo è riuscito |
| Counter | `prometheus.NewCounterVec` | `probe_checks_total`: controlli eseguiti, per esito |
| Histogram | `prometheus.NewHistogramVec` | `probe_duration_seconds`: quanto durano i controlli |

Il suffisso `Vec` indica una metrica con **label**: una serie per ogni combinazione di valori.

```go title="internal/server/metrics.go (estratto)"
up := prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "probe_up",
	Help: "1 se l'ultimo controllo del target è riuscito, 0 altrimenti.",
}, []string{"target", "type"})

reg := prometheus.NewRegistry()   // (1)!
reg.MustRegister(up, collectors.NewGoCollector())

up.WithLabelValues("tasks-api", "http").Set(1)

mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
```

1.  Un **registro** proprio invece di quello globale: i test possono creare più server senza conflitti di nomi.
    `NewGoCollector` aggiunge le metriche del runtime di Go (goroutine, memoria, garbage collector).

```console
$ curl -s localhost:9101/metrics | grep '^probe_'
probe_checks_total{result="success",target="tasks-api"} 12
probe_duration_seconds_bucket{target="tasks-api",type="http",le="0.005"} 11
...
probe_up{target="postgres",type="tcp"} 1
probe_up{target="tasks-api",type="http"} 1
```

!!! warning "Le label non sono illimitate"
    Ogni combinazione di valori delle label è una serie separata in Prometheus. Il nome del target va bene, perché
    i target sono pochi e decisi da te. Non mettere mai in una label valori senza limite, come un URL completo con
    parametri, l'IP di un client o un ID di richiesta: la memoria di Prometheus crescerebbe senza fine. È il
    problema della **cardinalità**.

### Un exporter che Prometheus conosce già

Il nome `probe_*` non è casuale: il **Blackbox Exporter** ufficiale di Prometheus fa lo stesso mestiere, in modo
molto più completo (HTTP, TCP, ICMP, DNS, TLS), e usa la stessa convenzione, come `probe_success`. racoon-probe è
un esercizio per capire come funziona un exporter dall'interno; in produzione, per questo caso, useresti il
Blackbox Exporter. Ma lo stesso schema ti serve ogni volta che un'applicazione o un sistema interno non espone
metriche da solo.

## Laboratorio

!!! laboratorio "Laboratorio · probe nello stack di osservabilità"
    1. Aggiungi al tuo `racoon-probe` le tre metriche con un registro dedicato, e l'endpoint `GET /metrics`.
       Estendi il test del server perché controlli le righe `probe_up` attese.
    2. Aggiungi la pipeline: pubblica il repository e controlla che il job di qualità sia verde. Crea il tag
       `v0.4.0` e verifica che l'immagine compaia nei *Packages* del repository.
    3. Avvia la Tasks API con lo stack di osservabilità della [Fase 7](../fase-7/index.md) e aggiungi il servizio
       `probe` sulle reti `backend` e `monitoraggio`.
    4. Aggiungi a `prometheus.yml` il job `probe` e controlla in *Status → Targets* che sia `UP`.
    5. Scrivi la regola di alert `TargetGiu`: `probe_up == 0` per 2 minuti. Ferma il database e guarda l'alert
       passare da *pending* a *firing* in Alertmanager.
    6. In Grafana crea un pannello con il 95° percentile della durata dei controlli per target.

??? success "Soluzione guidata"
    ```yaml title="compose.osservabilita.yaml (servizio da aggiungere)"
      probe:
        image: ghcr.io/TUONOME/racoon-probe:0.4.0
        volumes:
          - ./osservabilita/probe.yaml:/etc/probe/probe.yaml:ro
        networks: [backend, monitoraggio]
        restart: unless-stopped
    ```
    ```yaml title="osservabilita/prometheus.yml (job da aggiungere)"
      - job_name: probe
        static_configs:
          - targets: ["probe:9101"]
    ```
    ```yaml title="osservabilita/regole.yml (regola da aggiungere)"
      - alert: TargetGiu
        expr: probe_up == 0
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "{{ $labels.target }} non risponde"
          description: "Il controllo {{ $labels.type }} di {{ $labels.target }} fallisce da almeno 2 minuti."
    ```
    ```promql title="Pannello Grafana: 95° percentile"
    histogram_quantile(0.95, sum by (le, target) (rate(probe_duration_seconds_bucket[5m])))
    ```
    ```bash
    curl -X POST localhost:9090/-/reload      # ricarica la configurazione di Prometheus
    docker compose stop db                    # dopo 2 minuti l'alert TargetGiu per postgres va in firing
    ```
    **Spiegazione.** probe deve stare sulla rete `backend` per raggiungere `api` e `db`, e sulla rete
    `monitoraggio` per essere letto da Prometheus. Il tag Git `v0.4.0` diventa il tag d'immagine `0.4.0` grazie a
    `docker/metadata-action` con `type=semver`. Con `for: 2m` un singolo controllo fallito non sveglia nessuno: solo
    un problema che dura.

## Esercizi

!!! esercizio "Esercizio 1 · Quando è scaduto il certificato?"
    Aggiungi ai target HTTPS la metrica `probe_tls_cert_expiry_timestamp_seconds{target}`, con la data di scadenza
    del certificato del server in secondi Unix. Scrivi poi l'espressione PromQL che trova i certificati che scadono
    entro 14 giorni.

??? success "Soluzione"
    ```go
    // dopo resp, err := c.Client.Do(req), con resp != nil
    if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
    	scadenza := resp.TLS.PeerCertificates[0].NotAfter   // il primo è il certificato del server
    	m.certExpiry.WithLabelValues(t.Name).Set(float64(scadenza.Unix()))
    }
    ```
    ```promql
    probe_tls_cert_expiry_timestamp_seconds - time() < 14 * 24 * 3600
    ```
    **Spiegazione.** `resp.TLS` contiene i dettagli della connessione TLS (`nil` per HTTP in chiaro). Esporre una
    **data** e fare il calcolo in PromQL è meglio che esporre "giorni rimanenti": il valore non cambia a ogni
    lettura e l'alert si scrive con qualsiasi soglia. Il Blackbox Exporter espone una metrica analoga.

!!! esercizio "Esercizio 2 · Il linter ha sempre ragione?"
    Esegui `golangci-lint run` sul tuo progetto. Uno degli avvisi più comuni è `errcheck`: un errore restituito e non
    controllato, per esempio `w.Write([]byte("ok\n"))`. Correggi gli avvisi oppure, dove ignorare l'errore è una
    scelta consapevole, rendila esplicita.

??? success "Soluzione"
    ```go
    _, _ = w.Write([]byte("ok\n"))   // il client potrebbe essersi già disconnesso: non c'è niente da fare
    ```
    **Spiegazione.** Assegnare a `_` dice a chi legge (e al linter) che l'errore è stato ignorato **di proposito**.
    Per le scelte che valgono in tutto il progetto si usa un file `.golangci.yml` con la configurazione dei
    linter, versionato insieme al codice; per un singolo caso, un commento `//nolint:errcheck // motivo`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Perché nella pipeline si usa `go-version-file: go.mod`?","options":["È obbligatorio","Perché la versione di Go sia decisa in un solo posto, go.mod","Per scaricare la versione più recente","Per evitare la cache"],"answer":[1],"explanation":"Una sola fonte di verità per sviluppo e CI."},
 {"q":"Cosa distingue govulncheck da molti scanner di dipendenze?","options":["È a pagamento","Segnala solo le vulnerabilità nel codice che il programma chiama davvero","Controlla solo le immagini Docker","Corregge da solo le vulnerabilità"],"answer":[1],"explanation":"Analizza il grafo delle chiamate, quindi produce meno falsi allarmi."},
 {"q":"Quale tipo di metrica usi per il numero totale di controlli eseguiti?","options":["Gauge","Counter","Histogram","Summary"],"answer":[1],"explanation":"Un counter cresce sempre; con rate() ottieni la frequenza."},
 {"q":"Quale label NON dovresti mai aggiungere a probe_duration_seconds?","options":["target","type","l'ID univoco di ogni controllo","result"],"answer":[2],"explanation":"Valori senza limite creano una serie per ciascuno: problema di cardinalità."},
 {"q":"Perché probe usa un prometheus.NewRegistry() invece del registro globale?","options":["È più veloce","Per poter creare più server nei test senza registrare due volte le stesse metriche","Perché il registro globale è deprecato","Per cifrare le metriche"],"answer":[1],"explanation":"Il registro globale è condiviso da tutto il programma."},
 {"q":"Cosa fa `for: 2m` in una regola di alert?","options":["Ripete l'alert ogni 2 minuti","L'alert scatta solo se la condizione resta vera per 2 minuti","Silenzia l'alert per 2 minuti","Valuta la regola ogni 2 minuti"],"answer":[1],"explanation":"Prima di allora l'alert è pending."}
]}
</script>

## Fonti

- [GitHub Docs: Building and testing Go](https://docs.github.com/en/actions/use-cases-and-examples/building-and-testing/building-and-testing-go)
- [actions/setup-go](https://github.com/actions/setup-go) e [golangci-lint](https://golangci-lint.run/)
- [Go Vulnerability Management](https://go.dev/doc/security/vuln/) e [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [Prometheus: Instrumenting a Go application](https://prometheus.io/docs/guides/go-application/)
- [Prometheus: Writing exporters](https://prometheus.io/docs/instrumenting/writing_exporters/) e [Metric and label naming](https://prometheus.io/docs/practices/naming/)
- [prometheus/blackbox_exporter](https://github.com/prometheus/blackbox_exporter)
