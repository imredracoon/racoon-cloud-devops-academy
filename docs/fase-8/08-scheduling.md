# 8.8 Scheduling e autoscaling

!!! obiettivi "Obiettivi"
    - Guidare lo scheduler con nodeSelector, affinità, taint e toleration, topology spread.
    - Limitare le risorse di un namespace con ResourceQuota e LimitRange.
    - Svuotare un nodo per manutenzione rispettando i PodDisruptionBudget.
    - Scalare in automatico i Pod (HPA) e i nodi (Cluster Autoscaler, Karpenter).

## Come decide lo scheduler

Per ogni Pod in attesa lo scheduler lavora in due passi:

1. **Filtro**: scarta i nodi non adatti (risorse richieste non disponibili, `nodeSelector` non soddisfatto, taint non
   tollerati, volume in un'altra zona).
2. **Punteggio**: ordina i nodi rimasti (spazio libero, distribuzione, affinità preferite) e sceglie il migliore.

Se nessun nodo passa il filtro, il Pod resta `Pending` e gli eventi spiegano perché, per esempio
`0/3 nodes are available: 1 node(s) had untolerated taint, 2 Insufficient memory`.

Lo scheduler considera le **richieste** (`requests`), non il consumo reale: un nodo "pieno" di richieste rifiuta nuovi
Pod anche se la CPU è ferma. Richieste realistiche sono quindi fondamentali.

## Scegliere i nodi

### nodeSelector e affinità per i nodi

```bash
kubectl label node academy-worker disco=ssd
```

```yaml
spec:
  nodeSelector:          # forma semplice: TUTTE le label devono corrispondere
    disco: ssd
```

L'**affinità per i nodi** è più espressiva: operatori (`In`, `NotIn`, `Exists`, `Gt`...), regole obbligatorie o
preferite.

```yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:     # obbligatoria
        nodeSelectorTerms:
          - matchExpressions:
              - key: zona
                operator: In
                values: ["a", "b"]
      preferredDuringSchedulingIgnoredDuringExecution:    # preferita, con un peso
        - weight: 50
          preference:
            matchExpressions:
              - key: disco
                operator: In
                values: ["ssd"]
```

*IgnoredDuringExecution* significa che, se le label del nodo cambiano dopo, i Pod già in esecuzione restano.

### Affinità e anti-affinità tra Pod

Avvicinano o allontanano Pod in base alle label **di altri Pod** e a una `topologyKey` (nodo, zona):

```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchLabels:
            app: web
        topologyKey: kubernetes.io/hostname   # mai due Pod "web" sullo stesso nodo
```

### Topology spread

Per distribuire le repliche in modo uniforme tra nodi o zone, la soluzione moderna è
`topologySpreadConstraints`, usata nel Deployment del laboratorio:

```yaml
topologySpreadConstraints:
  - maxSkew: 1                               # differenza massima tra due domini
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: DoNotSchedule          # oppure ScheduleAnyway (solo preferenza)
    labelSelector:
      matchLabels:
        app.kubernetes.io/name: tasks-api
```

### Taint e toleration

Un **taint** sul nodo **respinge** i Pod; una **toleration** sul Pod gli permette di ignorarlo. È il contrario
dell'affinità: serve a riservare nodi (GPU, nodi di sistema) o a tenerli fuori uso.

```bash
kubectl taint nodes academy-worker2 dedicato=gpu:NoSchedule
kubectl taint nodes academy-worker2 dedicato=gpu:NoSchedule-     # il "-" finale lo rimuove
kubectl describe node academy-control-plane | grep Taints
# node-role.kubernetes.io/control-plane:NoSchedule
```

```yaml
tolerations:
  - key: dedicato
    operator: Equal
    value: gpu
    effect: NoSchedule
```

Effetti: `NoSchedule` (nuovi Pod respinti), `PreferNoSchedule` (evitato se possibile), `NoExecute` (anche i Pod già
presenti vengono sfrattati). Una toleration **permette** di andare sul nodo ma non lo **impone**: per riservare
davvero un nodo si combina taint e affinità.

Il control plane ha un taint `NoSchedule`: per questo i tuoi Pod finiscono sui worker.

## Limiti per namespace

```yaml title="quota.yaml"
apiVersion: v1
kind: ResourceQuota
metadata:
  name: quota-team
  namespace: prove
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.memory: 16Gi
    pods: "20"
    persistentvolumeclaims: "5"
---
apiVersion: v1
kind: LimitRange
metadata:
  name: predefiniti
  namespace: prove
spec:
  limits:
    - type: Container
      default:           # limite applicato se il container non ne dichiara
        memory: 256Mi
      defaultRequest:    # richiesta applicata se il container non ne dichiara
        cpu: 100m
        memory: 128Mi
      max:
        memory: 1Gi
```

Con una ResourceQuota su CPU o memoria, ogni Pod **deve** dichiarare richieste e limiti, altrimenti viene rifiutato:
la LimitRange glieli assegna in automatico.

## Manutenzione dei nodi

```bash
kubectl cordon academy-worker        # nessun nuovo Pod sul nodo
kubectl drain academy-worker --ignore-daemonsets --delete-emptydir-data   # sfratta i Pod
kubectl uncordon academy-worker      # di nuovo disponibile
```

`drain` rispetta i **PodDisruptionBudget**: il PDB `tasks-api` del laboratorio (`minAvailable: 1`) impedisce di
sfrattare l'ultimo Pod disponibile della API. È la stessa protezione che usano gli aggiornamenti dei nodi su EKS e
GKE.

## Autoscaling

| Livello | Strumento | Cosa scala |
|---|---|---|
| Pod, in orizzontale | **HorizontalPodAutoscaler** (HPA) | il numero di repliche, in base a CPU, memoria o metriche personalizzate |
| Pod, in verticale | **VerticalPodAutoscaler** (VPA, progetto a parte) | richieste e limiti dei container |
| Nodi | **Cluster Autoscaler** o **Karpenter** (AWS), autoscaler di GKE e GKE Autopilot | il numero di nodi: li aggiunge se ci sono Pod `Pending` per mancanza di risorse, li toglie se sono poco usati |

L'HPA legge le metriche da **metrics-server** (lo stesso che alimenta `kubectl top`) e calcola:

```text
repliche desiderate = ceil( repliche attuali × utilizzo attuale / utilizzo obiettivo )
```

con l'utilizzo misurato **rispetto alle richieste**: un Pod che usa 140m con una richiesta di 100m è al 140%. Senza
`requests.cpu` l'HPA basato sulla CPU non funziona.

```yaml title="labs/k8s/base/hpa.yaml"
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: tasks-api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: tasks-api
  minReplicas: 2
  maxReplicas: 6
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

Quando un HPA gestisce un Deployment, **non** impostare `replicas` nel manifest applicato da GitOps: ogni
sincronizzazione riporterebbe il numero al valore di Git, in lotta con l'HPA.

## Laboratorio

!!! laboratorio "Laboratorio · Autoscaling della Tasks API"
    1. Installa metrics-server. In kind i certificati dei kubelet non sono firmati per il loro IP, quindi va
       aggiunto `--kubelet-insecure-tls` (**solo** in laboratorio):
       ```bash
       kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
       kubectl -n kube-system patch deployment metrics-server --type=json \
         -p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
       kubectl -n kube-system rollout status deployment metrics-server
       kubectl top nodes
       ```
    2. Controlla lo stato dell'HPA: `kubectl -n tasks get hpa tasks-api -w`.
    3. Le Network Policy dell'overlay `dev` accettano traffico verso la API solo dal Gateway ([lezione 8.10](10-networkpolicy.md)). Per
       questa prova sospendile, poi genera carico da un Pod:
       ```bash
       kubectl -n tasks delete networkpolicy api-da-gateway nega-tutto-in-ingresso
       kubectl -n tasks run carico --image=busybox:1.37 --restart=Never -- \
         sh -c 'while true; do wget -qO- http://tasks-api/tasks > /dev/null; done'
       ```
       In pochi minuti le repliche salgono. Elimina il Pod `carico` e osserva la discesa: è più lenta. Perché?
    4. Applica un taint `NoSchedule` a entrambi i worker, poi fai `kubectl rollout restart` della API. Cosa succede
       ai nuovi Pod? Rimuovi i taint.

??? success "Soluzione"
    3. Per evitare oscillazioni l'HPA aspetta una **finestra di stabilizzazione** di 5 minuti prima di ridurre le
       repliche (configurabile in `spec.behavior.scaleDown`), mentre per salire reagisce subito.

       Alla fine elimina il Pod `carico` e ripristina le Network Policy con `kubectl apply -k overlays/dev`.
    4. ```bash
       kubectl taint nodes academy-worker academy-worker2 manutenzione=si:NoSchedule
       kubectl -n tasks rollout restart deployment tasks-api
       kubectl -n tasks get pods
       kubectl -n tasks describe pod NOME-DEL-POD-PENDING | tail -5
       ```
       Il rollout crea un nuovo Pod che resta `Pending` (nessun nodo tollerato: il control plane ha il suo taint);
       grazie a `maxUnavailable: 0` i vecchi Pod restano in servizio. Rimuovi i taint con
       `kubectl taint nodes academy-worker academy-worker2 manutenzione=si:NoSchedule-`.

!!! esercizio "Esercizio · Il Pod che non parte"
    Un Pod resta `Pending` con l'evento `0/3 nodes are available: 1 node(s) had untolerated taint
    {node-role.kubernetes.io/control-plane: }, 2 node(s) didn't match Pod's node affinity/selector`. Cosa
    controlli e come risolvi?

??? success "Soluzione"
    Il control plane è escluso dal taint (corretto). I due worker non hanno le label richieste dal `nodeSelector`
    o dall'affinità del Pod. Confronta `kubectl get pod NOME -o yaml | grep -A5 nodeSelector` con
    `kubectl get nodes --show-labels`: correggi la label nel Pod (errore di battitura frequente) oppure etichetta il
    nodo giusto con `kubectl label node`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Su cosa si basa lo scheduler per verificare se un nodo ha spazio per un Pod?","options":["Sul consumo reale di CPU","Sulle richieste (requests) dei Pod","Sui limiti","Sul numero di Pod"],"answer":[1],"explanation":"Per questo le richieste vanno impostate in modo realistico."},
 {"q":"Cosa fa un taint con effetto NoExecute?","options":["Impedisce solo i nuovi Pod","Sfratta anche i Pod già presenti che non lo tollerano","Rallenta il nodo","Nulla, è solo informativo"],"answer":[1],"explanation":"NoSchedule agisce solo sui nuovi Pod."},
 {"q":"Una toleration garantisce che il Pod vada sul nodo con il taint?","options":["Sì","No, lo permette soltanto: serve anche un'affinità o un nodeSelector","Solo con NoExecute","Solo nel namespace kube-system"],"answer":[1],"explanation":"Per riservare nodi si combinano taint e affinità."},
 {"q":"Da dove legge le metriche di CPU l'HorizontalPodAutoscaler?","options":["Da Prometheus, sempre","Da metrics-server (Metrics API)","Da etcd","Dai log"],"answer":[1],"explanation":"Con un adattatore può usare anche metriche personalizzate."},
 {"q":"Quali comandi servono per una manutenzione ordinata di un nodo?","options":["kubectl cordon","kubectl drain --ignore-daemonsets","kubectl uncordon","kubectl delete node","kubectl taint nodes --all"],"answer":[0,1,2],"explanation":"drain esegue già il cordon, ma è utile conoscerli tutti."},
 {"q":"Cosa impedisce a kubectl drain di sfrattare l'ultimo Pod disponibile di un'applicazione?","options":["ResourceQuota","PodDisruptionBudget","LimitRange","PriorityClass"],"answer":[1],"explanation":"Il drain aspetta finché il PDB lo permette."},
 {"q":"Quale strumento aggiunge nodi quando ci sono Pod Pending per mancanza di risorse?","options":["HPA","VPA","Cluster Autoscaler o Karpenter","metrics-server"],"answer":[2],"explanation":"Su GKE Autopilot lo fa la piattaforma."}
]}
</script>

## Fonti

- [Kubernetes: Kubernetes Scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/) e [Assigning Pods to Nodes](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/)
- [Kubernetes: Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) e [Pod Topology Spread Constraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/)
- [Kubernetes: Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/) e [Limit Ranges](https://kubernetes.io/docs/concepts/policy/limit-range/)
- [Kubernetes: Safely Drain a Node](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/) e [Disruptions](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)
- [Kubernetes: Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) e [Node Autoscaling](https://kubernetes.io/docs/concepts/cluster-administration/node-autoscaling/)
- [metrics-server](https://github.com/kubernetes-sigs/metrics-server)
