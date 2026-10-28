# 8.13 Helm, Kustomize e operatori

!!! obiettivi "Obiettivi"
    - Installare, aggiornare e rimuovere applicazioni con Helm.
    - Scrivere un chart semplice con template e valori.
    - Gestire varianti per ambiente con Kustomize: base, overlay, generatori, patch.
    - Spiegare CRD e operatori e installarne uno.

## Il problema dei manifest copiati

La Tasks API in tre ambienti (sviluppo, collaudo, produzione) significa tre copie di quasi gli stessi file, che
divergono un po' alla volta. E installare un software di terzi (Envoy Gateway, Argo CD, Prometheus) significa
gestire centinaia di oggetti scritti da altri. Due strumenti, complementari, risolvono il problema:

| | Helm | Kustomize |
|---|---|---|
| Approccio | **template** con variabili, impacchettati in un *chart* | **sovrapposizioni**: YAML normale più patch |
| Punto di forza | distribuire software riusabile, con versioni e ritorno indietro | varianti dei propri manifest senza linguaggi di template |
| Dove | comando `helm` | integrato in `kubectl` (`kubectl apply -k`) |
| Uso tipico | installare componenti di terzi | gestire le proprie applicazioni per ambiente |

Si usano spesso insieme, e Argo CD ([lezione 8.15](15-gitops-cloud.md)) li supporta entrambi.

## Helm

**Helm** è il gestore di pacchetti di Kubernetes, progetto *graduated* della CNCF. La versione 4, uscita a novembre
2025, ha introdotto tra l'altro il supporto all'*apply* lato server e un nuovo sistema di plugin; i comandi di base
sono gli stessi di Helm 3, che riceve correzioni di sicurezza fino a novembre 2026.

| Concetto | Significato |
|---|---|
| **Chart** | il pacchetto: template dei manifest, valori predefiniti, metadati |
| **Release** | un'installazione di un chart nel cluster, con un nome; lo stesso chart si può installare più volte |
| **Revisione** | ogni `upgrade` crea una nuova revisione, a cui si può tornare |
| **Repository** | dove si pubblicano i chart: un registro **OCI** (come per le immagini) o un repository HTTP classico |

```bash
# Un repository classico
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm search repo prometheus-community/kube-prometheus-stack

# Installare, con valori personalizzati
helm show values oci://docker.io/envoyproxy/gateway-helm --version v1.9.3 | less
helm install eg oci://docker.io/envoyproxy/gateway-helm --version v1.9.3 \
  -n envoy-gateway-system --create-namespace -f miei-valori.yaml

# Gestire le release
helm list -A
helm status eg -n envoy-gateway-system
helm get values eg -n envoy-gateway-system
helm upgrade eg oci://docker.io/envoyproxy/gateway-helm --version v1.9.3 \
  -n envoy-gateway-system --reuse-values --set deployment.replicas=2
helm history eg -n envoy-gateway-system
helm rollback eg 1 -n envoy-gateway-system
helm uninstall eg -n envoy-gateway-system
```

`helm upgrade --install` installa se la release non esiste e aggiorna se esiste: è la forma usata nelle pipeline.
Lo stato delle release è salvato in Secret del namespace (`sh.helm.release.v1.NOME.vN`).

### Il chart della Tasks API

Nel laboratorio, `labs/k8s/chart/tasks-api` è la stessa API come chart:

```text
chart/tasks-api/
├── Chart.yaml          # nome, versione del chart, appVersion
├── values.yaml         # valori predefiniti
└── templates/
    ├── _helpers.tpl    # frammenti riusabili: nomi ed etichette
    ├── deployment.yaml
    ├── service.yaml
    ├── hpa.yaml        # solo se autoscaling.enabled
    ├── httproute.yaml  # solo se httpRoute.enabled
    └── NOTES.txt       # messaggio mostrato dopo l'installazione
```

I template usano il linguaggio di template di Go: `{{ .Values.replicaCount }}` legge un valore,
`{{ include "tasks-api.labels" . | nindent 4 }}` inserisce un frammento indentato, `{{- if ... }}` rende facoltativa
una risorsa.

```yaml title="templates/deployment.yaml (estratto)"
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  template:
    spec:
      containers:
        - name: api
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

```bash
cd labs/k8s
helm lint chart/tasks-api
helm template prova chart/tasks-api --set autoscaling.enabled=true | less   # vedi il risultato
helm install prova chart/tasks-api -n tasks --dry-run=server                # validato dal cluster
```

!!! tip "Prima di installare, guarda"
    `helm template` e `--dry-run` mostrano esattamente cosa verrà creato. Prima di installare un chart di terzi,
    leggi i valori disponibili (`helm show values`) e il risultato: un chart può creare ClusterRole molto potenti.

## Kustomize

**Kustomize** parte da YAML normale, valido così com'è, e applica trasformazioni dichiarate in un file
`kustomization.yaml`. La struttura del laboratorio:

```text
labs/k8s/
├── base/                     # i manifest comuni a tutti gli ambienti
│   ├── kustomization.yaml
│   ├── namespace.yaml  postgres.yaml  api.yaml  hpa.yaml  networkpolicy.yaml
├── rotta/                    # HTTPRoute, aggiunta dove serve
└── overlays/
    ├── dev/                  # 1 replica minima, Secret generato da un file locale
    ├── prod/                 # immagine dal registro, 3 repliche, distribuzione rigida
    └── gitops/               # per Argo CD, senza segreti
```

```yaml title="labs/k8s/base/kustomization.yaml"
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: tasks              # imposta il namespace su tutte le risorse
resources:
  - namespace.yaml
  - postgres.yaml
  - api.yaml
  - hpa.yaml
  - networkpolicy.yaml
configMapGenerator:           # genera la ConfigMap con suffisso hash
  - name: tasks-api
    literals:
      - PYTHONUNBUFFERED=1
labels:
  - pairs:
      app.kubernetes.io/managed-by: kustomize
```

```yaml title="labs/k8s/overlays/prod/kustomization.yaml"
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
  - ../../rotta
images:                        # cambia l'immagine senza toccare i manifest
  - name: tasks-api
    newName: ghcr.io/TUO-UTENTE/tasks-api
    newTag: "1.0.0"
configMapGenerator:
  - name: tasks-api
    namespace: tasks
    behavior: merge            # aggiunge chiavi alla ConfigMap della base
    literals:
      - AMBIENTE=prod
patches:                       # patch JSON 6902 mirate
  - target:
      kind: HorizontalPodAutoscaler
      name: tasks-api
    patch: |-
      - op: replace
        path: /spec/minReplicas
        value: 3
```

```bash
kubectl kustomize overlays/prod | less     # mostra il risultato
kubectl diff -k overlays/prod              # differenze con il cluster
kubectl apply -k overlays/prod
```

Oltre alle patch JSON esistono le patch *strategic merge*: un frammento di YAML con solo i campi da cambiare, che
Kustomize fonde con l'originale. Il generatore `secretGenerator` crea Secret da file o valori, con lo stesso
suffisso hash: nell'overlay `dev` legge la password da `db-password.txt`, escluso da Git.

## CRD e operatori

Una **CustomResourceDefinition** aggiunge un nuovo tipo all'API; un **operatore** è un controller che lo gestisce,
racchiudendo in codice il lavoro di un amministratore esperto. Per esempio, con l'operatore **CloudNativePG**:

```yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: tasks-db
spec:
  instances: 3              # un primario e due repliche, con failover automatico
  storage:
    size: 5Gi
```

L'operatore crea Pod, Service, PVC, configura la replica, gestisce backup e promozione di una replica se il primario
cade. Gli operatori si installano di solito con Helm o con manifest, e si trovano su **OperatorHub.io** e
**Artifact Hub**.

```bash
kubectl get crd
kubectl api-resources --api-group=gateway.networking.k8s.io
kubectl explain httproute.spec.rules
```

## Laboratorio

!!! laboratorio "Laboratorio · Helm e Kustomize a confronto"
    1. Installa il chart della Tasks API come release `bis` nel namespace `tasks`, con una sola replica e il nome
       del Secret del database generato da Kustomize:
       ```bash
       SECRET=$(kubectl -n tasks get secrets -o name | grep tasks-db | cut -d/ -f2)
       helm install bis chart/tasks-api -n tasks --set replicaCount=1 --set database.existingSecret=$SECRET
       ```
       Il Pod di `bis` diventa pronto? Il Service `tasks-api` della base gli manda traffico? Perché?
    2. Aggiorna la release a 2 repliche, poi torna alla revisione 1 con `helm rollback`. Guarda `helm history`.
    3. Confronta `kubectl kustomize overlays/dev` e `kubectl kustomize overlays/prod`: elenca tre differenze.
    4. Rimuovi la release `bis`.

??? success "Soluzione"
    1. Sì a entrambe le domande, ed è istruttivo. Il chart usa la label `app.kubernetes.io/name: tasks-api`, la
       stessa della base: la Network Policy `postgres-solo-da-api` lascia passare il Pod verso il database, che
       diventa pronto, **e** il selettore del Service `tasks-api` (che usa solo quella label) lo include tra i suoi
       endpoint. Verifica con `kubectl -n tasks get endpointslices -l kubernetes.io/service-name=tasks-api`. Il
       chart invece seleziona con due label, `name` e `instance`, così il suo Service non prende i Pod della base.
       Morale: le label sono un'interfaccia, e due installazioni con le stesse label si mescolano.
    2. ```bash
       helm upgrade bis chart/tasks-api -n tasks --reuse-values --set replicaCount=2
       helm history bis -n tasks
       helm rollback bis 1 -n tasks
       helm history bis -n tasks     # la revisione 3 è il rollback alla 1
       ```
       Il rollback crea una **nuova** revisione: la storia non si riscrive.
    3. Per esempio: l'immagine (`tasks-api:dev` contro `ghcr.io/TUO-UTENTE/tasks-api:1.0.0`), `minReplicas` (1
       contro 3), `whenUnsatisfiable` (`ScheduleAnyway` contro `DoNotSchedule`), la chiave `AMBIENTE` della
       ConfigMap, il Secret generato solo in `dev`, l'HTTPRoute solo in `prod`.
    4. `helm uninstall bis -n tasks`.

!!! esercizio "Esercizio · Un overlay di collaudo"
    Crea `overlays/collaudo` che parte dalla base, usa l'immagine `ghcr.io/TUO-UTENTE/tasks-api:1.1.0-rc1`,
    aggiunge `AMBIENTE=collaudo` alla ConfigMap e aggiunge a tutte le risorse la label `ambiente: collaudo`.

??? success "Soluzione"
    ```yaml title="overlays/collaudo/kustomization.yaml"
    apiVersion: kustomize.config.k8s.io/v1beta1
    kind: Kustomization
    resources:
      - ../../base
    images:
      - name: tasks-api
        newName: ghcr.io/TUO-UTENTE/tasks-api
        newTag: 1.1.0-rc1
    configMapGenerator:
      - name: tasks-api
        namespace: tasks
        behavior: merge
        literals:
          - AMBIENTE=collaudo
    labels:
      - pairs:
          ambiente: collaudo
    ```
    Con `labels` senza `includeSelectors: true` la label viene aggiunta ai metadati ma **non** ai selettori:
    è la scelta sicura, perché i selettori di un Deployment non si possono modificare dopo la creazione.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Come si chiama un'installazione di un chart Helm nel cluster?","options":["Package","Release","Deployment","Overlay"],"answer":[1],"explanation":"Ogni upgrade ne crea una nuova revisione."},
 {"q":"Quale comando Helm mostra i manifest generati senza installare nulla?","options":["helm template","helm show chart","helm list","helm get notes"],"answer":[0],"explanation":"Anche --dry-run su install e upgrade."},
 {"q":"Cosa fa helm rollback eg 1?","options":["Cancella le revisioni successive","Crea una nuova revisione con la configurazione della revisione 1","Reinstalla il chart da zero","Torna alla versione precedente di Helm"],"answer":[1],"explanation":"La cronologia resta."},
 {"q":"Con quale comando applichi un overlay Kustomize?","options":["kubectl apply -k overlays/prod","kubectl kustomize apply","helm install -k","kubectl apply -f kustomization.yaml"],"answer":[0],"explanation":"kubectl kustomize mostra solo il risultato."},
 {"q":"Quali sono funzioni di Kustomize?","options":["Generare ConfigMap e Secret con suffisso hash","Cambiare immagini e tag","Applicare patch per ambiente","Gestire revisioni e rollback delle release","Scaricare chart da repository OCI con versioni"],"answer":[0,1,2],"explanation":"Revisioni e repository sono funzioni di Helm."},
 {"q":"Cos'è un operatore Kubernetes?","options":["Un amministratore del cluster","Un controller che gestisce risorse personalizzate automatizzando compiti operativi","Un plugin di kubectl","Un tipo di nodo"],"answer":[1],"explanation":"Per esempio CloudNativePG per PostgreSQL."}
]}
</script>

## Fonti

- [Helm: Documentation](https://helm.sh/docs/) e [Using Helm](https://helm.sh/docs/intro/using_helm/)
- [Helm: Chart Template Guide](https://helm.sh/docs/chart_template_guide/)
- [Helm Blog: Helm 4 Released](https://helm.sh/blog/helm-4-released)
- [Kubernetes: Declarative Management of Kubernetes Objects Using Kustomize](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/)
- [Kustomize: documentazione di riferimento](https://kubectl.docs.kubernetes.io/references/kustomize/)
- [Kubernetes: Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) e [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [CloudNativePG](https://cloudnative-pg.io/documentation/current/)
