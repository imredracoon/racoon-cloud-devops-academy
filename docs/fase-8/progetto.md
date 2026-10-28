# Progetto 8 · Cluster con GitOps

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai gestire una piattaforma Kubernetes come si fa in azienda: tutto dichiarato in Git, applicato da un agente
    GitOps, sicuro per impostazione predefinita, osservabile, e con una procedura provata per i disastri. È il
    progetto che chiude il percorso e unisce tutte le fasi: Linux, container, CI/CD, IaC, osservabilità.

| | |
|---|---|
| **Durata** | 30-40 ore |
| **Consegna** | un repository `tasks-platform` (o una cartella del tuo fork) con README, diagramma e documentazione |
| **Costo** | zero con kind; con EKS o GKE solo per le ore di prova (cancella il cluster alla fine) |

## Il compito

### Parte A · La piattaforma

1. Un cluster **kind** con tre nodi (o, se vuoi dimostrare il cloud, **EKS** o **GKE** creato con **Terraform**).
2. **Argo CD** installato, e da quel momento **tutto il resto** installato da Argo CD con lo schema
   *app of apps*: un'Application `radice` che punta alla cartella `piattaforma/applicazioni/`, dove ogni file è
   un'Application:
    - **Envoy Gateway** (chart Helm) e il Gateway `academy`;
    - **metrics-server**;
    - **kube-prometheus-stack** (Prometheus Operator, Prometheus, Alertmanager, Grafana) nel namespace
      `monitoraggio`;
    - la **Tasks API** in due ambienti, `tasks-staging` e `tasks-prod`, dagli overlay Kustomize.

### Parte B · L'applicazione, sicura per impostazione predefinita

3. Overlay `staging` e `prod` con immagini da GHCR con tag immutabile (`sha-...`), repliche e risorse diverse.
4. Namespace con Pod Security `restricted` e **ResourceQuota** e **LimitRange**.
5. **Network Policy** in ingresso *e* in uscita: la API riceve solo dal Gateway e da Prometheus, parla solo con DNS
   e PostgreSQL; PostgreSQL riceve solo dalla API.
6. HPA, PodDisruptionBudget, sonde corrette (readiness con il database, liveness senza), `topologySpreadConstraints`.
7. Nessun segreto in Git: Secret creati fuori banda e documentati, oppure (facoltativo, consigliato)
   **External Secrets Operator** o **Sealed Secrets**.
8. Un service account `ci` con il minimo privilegio, se la pipeline deve leggere lo stato del cluster.

### Parte C · Il flusso di rilascio

9. La pipeline della [Fase 7](../fase-7/index.md), a ogni merge su `main`, costruisce l'immagine e **apre una pull request** che aggiorna
   il tag nell'overlay `staging` ([lezione 8.15](15-gitops-cloud.md)). Il passaggio a `prod` è una seconda pull request, approvata a mano.
10. Un **ServiceMonitor** per la Tasks API, la dashboard SLO della [Fase 7](../fase-7/index.md) importata in Grafana come ConfigMap, e
    almeno una regola di alert come `PrometheusRule`.

### Parte D · Il giorno del disastro

11. Uno script `scripts/backup-etcd.sh` che esegue lo snapshot di etcd ([lezione 8.12](12-kubeadm.md)) e un **runbook** di
    ripristino, provato: elimina un namespace, ripristina, documenta i tempi.
12. Un documento `docs/game-day.md` con tre guasti provocati di proposito (un nodo fermo, un rilascio rotto, il
    database eliminato) e, per ciascuno: sintomi visti in Grafana, diagnosi, riparazione, tempo impiegato,
    miglioramento proposto.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `kubectl get applications -n argocd` | tutte `Synced` e `Healthy` |
| Modifica manuale di un Deployment | annullata da Argo CD entro pochi minuti |
| Merge di una PR sul codice | nuova immagine, PR automatica con il nuovo tag, rilascio in staging dopo il merge |
| `kubectl -n tasks-prod run t --image=busybox:1.37 -- wget -qO- -T 3 tasks-api` | rifiutato (Pod Security) o in timeout (Network Policy) |
| `kubectl drain` di un worker | la API resta raggiungibile dal Gateway per tutto il tempo |
| Prometheus, pagina *Targets* | la Tasks API di entrambi gli ambienti in stato UP |
| Ripristino di etcd dopo l'eliminazione di un namespace | namespace e risorse di nuovo presenti |

## Suggerimenti

??? tip "La struttura del repository"
    ```text
    tasks-platform/
    ├── README.md                     # cosa c'è, diagramma, come si avvia
    ├── bootstrap/
    │   ├── kind-cluster.yaml
    │   └── radice.yaml               # l'unica Application applicata a mano
    ├── piattaforma/
    │   └── applicazioni/             # un'Application per componente
    │       ├── envoy-gateway.yaml
    │       ├── metrics-server.yaml
    │       ├── monitoraggio.yaml
    │       ├── tasks-staging.yaml
    │       └── tasks-prod.yaml
    ├── tasks/
    │   ├── base/
    │   └── overlays/{staging,prod}/
    ├── scripts/backup-etcd.sh
    └── docs/{runbook-etcd.md,game-day.md}
    ```

??? tip "Le Network Policy e Prometheus"
    Prometheus deve poter leggere `/metrics` sulla porta del container. Senza una policy dedicata, il target
    risulta `DOWN` con un timeout: è un buon esempio da raccontare nel `game-day.md`.

??? success "Una possibile soluzione: Application radice, un componente Helm e il ServiceMonitor"
    ```yaml title="bootstrap/radice.yaml"
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: radice
      namespace: argocd
    spec:
      project: default
      source:
        repoURL: https://github.com/TUO-UTENTE/tasks-platform.git
        targetRevision: main
        path: piattaforma/applicazioni
      destination:
        server: https://kubernetes.default.svc
        namespace: argocd
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
    ```

    ```yaml title="piattaforma/applicazioni/envoy-gateway.yaml"
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: envoy-gateway
      namespace: argocd
    spec:
      project: default
      source:
        repoURL: registry-1.docker.io/envoyproxy   # registro OCI, senza il prefisso oci://
        chart: gateway-helm
        targetRevision: v1.9.3
      destination:
        server: https://kubernetes.default.svc
        namespace: envoy-gateway-system
      syncPolicy:
        automated: {}
        syncOptions:
          - CreateNamespace=true
          - ServerSideApply=true                   # CRD grandi: come con kubectl --server-side
    ```

    ```yaml title="tasks/base/servicemonitor.yaml"
    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      name: tasks-api
      labels:
        release: monitoraggio      # il kube-prometheus-stack seleziona i ServiceMonitor con l'etichetta della release
    spec:
      selector:
        matchLabels:
          app.kubernetes.io/name: tasks-api
      endpoints:
        - port: http               # il NOME della porta del Service
          path: /metrics
          interval: 30s
    ```

    ```yaml title="tasks/base/networkpolicy-prometheus.yaml"
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: api-da-prometheus
    spec:
      podSelector:
        matchLabels:
          app.kubernetes.io/name: tasks-api
      policyTypes: ["Ingress"]
      ingress:
        - from:
            - namespaceSelector:
                matchLabels:
                  kubernetes.io/metadata.name: monitoraggio
              podSelector:
                matchLabels:
                  app.kubernetes.io/name: prometheus
          ports:
            - protocol: TCP
              port: 8000
    ```

    Il ServiceMonitor raggiunge i Pod passando dal Service, ma Prometheus si collega **direttamente** agli IP dei Pod
    sulla porta del container: per questo la policy usa 8000. Tutti i manifest di questa soluzione sono stati
    validati con un API server Kubernetes 1.37 con le CRD di Argo CD e Prometheus Operator.

## Come presentarlo

- Un **diagramma** nel README: cluster, namespace, flusso GitOps, flusso delle richieste e delle metriche.
- Una GIF o un breve video: modifica in Git → Argo CD sincronizza → la dashboard mostra la nuova versione.
- Il `game-day.md` è la parte che colpisce di più in un colloquio: dimostra che sai cosa fare quando le cose vanno
  male. Preparati a raccontarlo in due minuti.

## Fonti

- [Argo CD: Cluster Bootstrapping (app of apps)](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/)
- [Argo CD: Helm (repository OCI)](https://argo-cd.readthedocs.io/en/stable/user-guide/helm/)
- [Argo CD: Sync Options (ServerSideApply)](https://argo-cd.readthedocs.io/en/stable/user-guide/sync-options/)
- [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
- [Prometheus Operator: ServiceMonitor](https://prometheus-operator.dev/docs/developer/getting-started/)
- [External Secrets Operator](https://external-secrets.io/latest/) · [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets)
