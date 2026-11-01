# Progetto Go · racoon-probe

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai scrivere in Go uno strumento infrastrutturale completo: una CLI con test, un servizio concorrente che
    si ferma in modo pulito, un'immagine minima multi-architettura, una pipeline, metriche Prometheus e un operator
    Kubernetes. È il tipo di codice che si scrive davvero nei team Platform e SRE, e risponde in modo concreto agli
    annunci che chiedono "Go" tra i requisiti.

| | |
|---|---|
| **Durata** | si costruisce lungo il corso, nei laboratori delle lezioni G.4-G.11; 6-10 ore di rifinitura finale |
| **Consegna** | Due repository pubblici: `racoon-probe` e `probe-operator` |
| **Ambiente** | Dev Container, lo stack Compose della [Fase 7](../fase-7/index.md), il cluster kind della [Fase 8](../fase-8/index.md) |

## Le tappe

Il progetto non si fa tutto alla fine: ogni tappa chiude il laboratorio di una lezione. Segna un **tag Git** a ogni
tappa, così la storia del repository racconta la crescita del tool.

| Tappa | Lezione | Cosa aggiunge | Tag |
|---|---|---|---|
| 1 | [G.4](04-cli-test.md#laboratorio-racoon-probe-primo-passo) | `probe check`, configurazione YAML validata, test a tabella | `v0.1.0` |
| 2 | [G.5](05-concorrenza-http.md#laboratorio-probe-serve) | controlli in parallelo, `probe serve`, `/healthz` e `/results`, arresto pulito, unità systemd | `v0.2.0` |
| 3 | [G.6](06-container.md#laboratorio) | Dockerfile multi-stage su distroless, immagine amd64 + arm64 | `v0.3.0` |
| 4 | [G.9](09-ci-metriche.md#laboratorio) | `/metrics`, pipeline GitHub Actions, alert `TargetGiu` e pannello Grafana | `v0.4.0` |
| 5 | [G.11](11-operator.md#laboratorio) | repository `probe-operator`: CRD `Probe` e controller | `v0.1.0` |

## Il compito

### Repository `racoon-probe`

1. Struttura `cmd/probe` + `internal/config`, `internal/check`, `internal/server`.
2. Sottocomandi `check` (uscita a tabella o JSON, codice 1 se un controllo fallisce), `serve` e `version`.
3. Target `http` (codice atteso) e `tcp`; timeout per controllo; controlli in parallelo con un limite.
4. Metriche `probe_up`, `probe_duration_seconds` e `probe_checks_total` su un registro dedicato.
5. Log JSON con `slog`; arresto pulito su SIGTERM.
6. Test per ogni pacchetto, con `httptest` e senza accessi a Internet; `go test -race` verde.
7. `Dockerfile` multi-stage su `distroless/static:nonroot`, versione incorporata con `-ldflags`.
8. Pipeline: `gofmt`, `go vet`, golangci-lint, `govulncheck`, test con race detector, immagine multi-architettura
   pubblicata su GHCR dai tag.
9. Una funzionalità in più a tua scelta, con i suoi test. Per esempio: la metrica di scadenza dei certificati
   ([lezione G.9](09-ci-metriche.md#esercizi)), un target `dns` che risolve un nome, oppure la ricarica della
   configurazione con `SIGHUP`.
10. `README.md` con: a cosa serve, esempio di `probe.yaml`, come installarlo (binario, systemd, container),
    l'elenco delle metriche, come eseguire i test.

### Repository `probe-operator`

1. La CRD `Probe` con validazione, valori predefiniti e colonne per `kubectl get`.
2. Il controller con condizione `Available`, `lastCheckTime`, `lastDuration` e ripianificazione su `interval`.
3. Test del controller con il client finto (o con envtest).
4. Permessi RBAC minimi, generati dai marcatori.
5. `README.md` con l'installazione su kind e un esempio di `Probe`.

### Nei tre livelli

=== "Livello 1 · Un PC"

    `probe` gira nello stack Compose del Dev Container e controlla la Tasks API; l'operator gira nel cluster kind.

=== "Livello 2 · Workstation + server"

    In più, `probe` è installato su **`srv1`** come servizio systemd (tappa 2) e controlla le VM e i servizi di
    casa; il Prometheus della [Fase 7](../fase-7/index.md#nei-tre-livelli) lo legge su `srv1:9101`.

=== "Livello 3 · Mini data center"

    `probe` gira su `nodo1` e controlla tutti i nodi (SSH, NFS, la Tasks API su `nodo2` e `nodo3`); l'operator
    gira nel cluster kubeadm e i `Probe` dei servizi del cluster sono versionati nel repository GitOps di Argo CD.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `go test -race ./...` in entrambi i repository | tutti i test passano |
| `golangci-lint run` e `govulncheck ./...` | nessun problema |
| `probe check -config probe.yaml` con un target spento | tabella con `KO`, codice di uscita 1 |
| `probe serve`, poi `kill -TERM` del processo | nei log `arresto in corso`, uscita con codice 0 |
| `docker buildx imagetools inspect ghcr.io/TUONOME/racoon-probe:0.4.0` | due piattaforme: `linux/amd64` e `linux/arm64` |
| `docker run --rm ghcr.io/TUONOME/racoon-probe:0.4.0 version` | `v0.4.0` |
| Prometheus, *Status → Targets* | job `probe` in stato `UP` |
| fermato il database della Tasks API per 2 minuti | alert `TargetGiu` in *firing* |
| `kubectl apply` di un `Probe` con `url: ftp://x` | rifiutato dall'API server |
| `kubectl get probes -A` | colonne URL, DISPONIBILE e DURATA aggiornate |

## Criteri di valutazione

- [ ] Tutte le verifiche della tabella passano.
- [ ] `main.go` è sottile: la logica è nei pacchetti `internal` ed è coperta da test.
- [ ] Ogni operazione di rete ha un timeout o un `context` con scadenza.
- [ ] Gli errori sono arricchiti con `%w` e i messaggi dicono cosa è fallito e dove.
- [ ] Nessun dato condiviso tra goroutine senza mutex, atomic o channel.
- [ ] I README permettono a un estraneo di installare e usare i tool senza chiederti nulla.
- [ ] Ogni tappa è un tag Git e ogni funzionalità è arrivata in `main` con una pull request.

## Come raccontarlo

Nel CV e su LinkedIn questo progetto si descrive in una riga che un selezionatore capisce subito, per esempio:
*"racoon-probe: exporter Prometheus e operator Kubernetes in Go per il controllo di servizi HTTP e TCP; immagine
distroless multi-architettura, CI con GitHub Actions"*. In un colloquio tecnico preparati a spiegare tre scelte:
perché un registro Prometheus dedicato, perché `Reconcile` restituisce `RequeueAfter` invece di dormire, e come
hai evitato le race condition nel server.

## Fonti

- [Effective Go](https://go.dev/doc/effective_go) e [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Organizing a Go module](https://go.dev/doc/modules/layout)
- [Prometheus: Writing exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)
- [The Kubebuilder Book](https://book.kubebuilder.io/)
