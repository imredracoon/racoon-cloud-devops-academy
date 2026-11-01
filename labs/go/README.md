# Laboratori del percorso Go

Le versioni di riferimento dei programmi del [percorso Go](https://imredracoon.github.io/racoon-cloud-devops-academy/go/).
Consultale **dopo** aver scritto la tua versione: sono qui per confrontare, non per copiare.

| Cartella | Lezioni | Cosa contiene |
|---|---|---|
| [`probe`](probe/) | G.4, G.5, G.6, G.9 | **racoon-probe**: CLI e servizio che controlla target HTTP e TCP, con metriche Prometheus, Dockerfile e pipeline |
| [`inventario`](inventario/) | G.7 | elenco di istanze e bucket di AWS e Google Cloud dietro un'unica interfaccia |
| [`terratest`](terratest/) | G.8 | test di integrazione dei progetti Terraform in `labs/terraform` |
| [`podcheck`](podcheck/) | G.10 | trova i Pod con problemi usando client-go; CronJob con RBAC minimo |
| [`probe-operator`](probe-operator/) | G.11 | operator Kubernetes con la CRD `Probe`, scritto con controller-runtime |

Ogni cartella è un **modulo Go indipendente** (ha il suo `go.mod`). Dal Dev Container:

```bash
cd labs/go/probe
go test -race ./...
go run ./cmd/probe check -config probe.example.yaml
```

Tutti i moduli richiedono Go 1.27 (già installato nel Dev Container). I test non richiedono credenziali cloud né un
cluster, tranne `terratest`: `TestPrimoProgetto` richiede il binario `terraform`, `TestReteWebAWS` crea risorse vere
su AWS e parte solo con `TERRATEST_AWS=1`.

La pipeline di `probe` (`probe/.github/workflows/ci.yml`) funziona quando la cartella è la radice di un repository
GitHub, come il tuo `racoon-probe`.
