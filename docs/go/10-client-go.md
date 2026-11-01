# G.10 Parlare con Kubernetes: client-go

!!! obiettivi "Obiettivi"
    - Capire come un programma parla con l'API server di Kubernetes, esattamente come `kubectl`.
    - Usare `client-go` per leggere risorse, dal tuo PC (kubeconfig) e da dentro un Pod (service account).
    - Dare a un programma solo i permessi RBAC che gli servono.
    - Capire cosa sono *watch* e *informer*, la base degli operator.

!!! info "Quando studiarla"
    Nella [Fase 8](../fase-8/index.md), dopo la [lezione 8.11](../fase-8/11-rbac.md). Serve il cluster kind della
    [lezione 8.2](../fase-8/02-kubectl.md).

## Tutto passa dall'API server

Nella [lezione 8.1](../fase-8/01-architettura.md) hai visto che ogni componente di Kubernetes (scheduler,
controller, kubelet) e ogni client (`kubectl`, Argo CD, Helm) fa una cosa sola: chiama l'**API server** con
richieste HTTPS. Puoi vederlo con `kubectl` stesso:

```console
$ kubectl get pods -n tasks -v=6
... GET https://127.0.0.1:41235/api/v1/namespaces/tasks/pods?limit=500 200 OK in 9 milliseconds
```

**client-go** è la libreria Go con cui sono scritti `kubectl`, i controller di Kubernetes e quasi tutti gli
operator. Ti dà un client **tipizzato**: invece di costruire richieste e leggere JSON, lavori con struct Go come
`corev1.Pod` e `appsv1.Deployment`.

```bash
go get k8s.io/client-go@latest k8s.io/api@latest k8s.io/apimachinery@latest
```

| Modulo | Contenuto |
|---|---|
| `k8s.io/api` | le struct delle risorse: `core/v1` (Pod, Service...), `apps/v1` (Deployment...) |
| `k8s.io/apimachinery` | i tipi comuni: `metav1.ObjectMeta`, `ListOptions`, errori dell'API |
| `k8s.io/client-go` | il client vero e proprio, la configurazione, i *watch* e gli *informer* |

!!! info "Le versioni"
    Le versioni di questi moduli seguono quelle di Kubernetes: `client-go v0.37.x` corrisponde a Kubernetes 1.37.
    Un client è compatibile con le versioni del cluster vicine alla sua; tieni allineati i tre moduli tra loro.

## Connettersi al cluster

```go
func connetti(kubeconfig string) (*kubernetes.Clientset, error) {
	cfg, err := rest.InClusterConfig()                                  // (1)!
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)        // (2)!
		if err != nil {
			return nil, fmt.Errorf("configurazione del cluster: %w", err)
		}
	}
	return kubernetes.NewForConfig(cfg)
}
```

1.  **Dentro un Pod**: Kubernetes monta nel container il token del service account e il certificato della CA in
    `/var/run/secrets/kubernetes.io/serviceaccount/`, e l'indirizzo dell'API server nelle variabili d'ambiente.
    `InClusterConfig` li legge.
2.  **Sul tuo PC**: lo stesso file `~/.kube/config` che usa `kubectl`, con il contesto corrente.

## Leggere le risorse

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

pods, err := clientset.CoreV1().Pods("tasks").List(ctx, metav1.ListOptions{   // "" = tutti i namespace
	LabelSelector: "app.kubernetes.io/name=tasks-api",                          // come kubectl -l
})
if err != nil {
	return err
}
for _, p := range pods.Items {
	fmt.Println(p.Name, p.Status.Phase, p.Spec.NodeName)
}

dep, err := clientset.AppsV1().Deployments("tasks").Get(ctx, "tasks-api", metav1.GetOptions{})
if apierrors.IsNotFound(err) {                       // k8s.io/apimachinery/pkg/api/errors
	fmt.Println("il Deployment non esiste")
}
```

La struttura delle chiamate rispecchia quella dell'API: gruppo e versione (`CoreV1`, `AppsV1`), risorsa (`Pods`,
`Deployments`), namespace, verbo (`List`, `Get`, `Create`, `Update`, `Delete`, `Watch`). Gli stessi verbi che
scrivi nelle regole RBAC.

## podcheck

Il tool di questa lezione, `podcheck`, elenca i Pod con problemi: in attesa (`Pending`), con container bloccati
(`CrashLoopBackOff`, `ImagePullBackOff`), riavviati troppe volte o non pronti. È quello che fai a mano nella
[lezione 8.14](../fase-8/14-troubleshooting.md) con `kubectl get pods` e `kubectl describe`. Il codice è in
[`labs/go/podcheck`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/podcheck).

```go title="main.go (estratto)"
func problemi(p corev1.Pod, soglia int32) []string {
	var out []string
	if p.Status.Phase == corev1.PodPending || p.Status.Phase == corev1.PodFailed {
		out = append(out, "fase "+string(p.Status.Phase))
	}
	for _, c := range p.Status.ContainerStatuses {
		if w := c.State.Waiting; w != nil && w.Reason != "ContainerCreating" {
			out = append(out, fmt.Sprintf("container %s in attesa: %s", c.Name, w.Reason))
		}
		if c.RestartCount >= soglia {
			out = append(out, fmt.Sprintf("container %s riavviato %d volte", c.Name, c.RestartCount))
		}
		if !c.Ready && p.Status.Phase == corev1.PodRunning {
			out = append(out, fmt.Sprintf("container %s non pronto", c.Name))
		}
	}
	return out
}
```

La funzione `problemi` riceve un `corev1.Pod` e non chiama l'API: si testa costruendo Pod finti in memoria, senza
cluster (vedi `main_test.go`). È la stessa separazione tra logica e I/O delle lezioni precedenti.

```console
$ kubectl -n tasks set image deployment/tasks-api api=ghcr.io/esempio/non-esiste:1   # un guasto voluto
$ go run . -n tasks
NAMESPACE  POD                         PROBLEMA
tasks      tasks-api-6d9f7c8b5d-x2k4q  container api in attesa: ImagePullBackOff

3 Pod controllati, 1 problemi trovati
$ kubectl -n tasks rollout undo deployment/tasks-api
```

## Dentro il cluster, con il minimo dei permessi

Un programma che gira nel cluster usa il **service account** del suo Pod. Per `podcheck` bastano `get` e `list`
sui Pod: niente scrittura, niente Secret. Il file
[`labs/go/podcheck/deploy.yaml`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/podcheck/deploy.yaml)
contiene namespace, ServiceAccount, ClusterRole, ClusterRoleBinding e un **CronJob** che lo esegue ogni 10 minuti:

```yaml title="deploy.yaml (estratto)"
kind: ClusterRole
metadata:
  name: podcheck-lettura-pod
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list"]
```

```bash
docker build -t podcheck:dev .
kind load docker-image podcheck:dev --name academy      # copia l'immagine nei nodi di kind, senza registry
kubectl apply -f deploy.yaml
kubectl -n podcheck create job --from=cronjob/podcheck prova   # esegue subito, senza aspettare l'orario
kubectl -n podcheck logs job/prova
```

Verifica che il permesso sia davvero minimo con `kubectl auth can-i`, come nella
[lezione 8.11](../fase-8/11-rbac.md):

```console
$ kubectl auth can-i list pods -A --as=system:serviceaccount:podcheck:podcheck
yes
$ kubectl auth can-i list secrets -A --as=system:serviceaccount:podcheck:podcheck
no
```

## Watch e informer

`List` fotografa la situazione in un istante. Per **reagire** ai cambiamenti, l'API offre il **watch**: una
connessione che resta aperta e su cui l'API server invia ogni evento (aggiunto, modificato, eliminato).

```go
w, err := clientset.CoreV1().Pods("tasks").Watch(ctx, metav1.ListOptions{})
if err != nil {
	return err
}
for ev := range w.ResultChan() {          // un channel, come nella lezione G.5
	if p, ok := ev.Object.(*corev1.Pod); ok {
		fmt.Println(ev.Type, p.Name, p.Status.Phase)
	}
}
```

Un watch "nudo" ha però dei limiti: si interrompe e va riaperto, e ogni programma che chiede di continuo le stesse
informazioni appesantisce l'API server. Gli **informer** di client-go risolvono entrambi i problemi: tengono una
**cache locale** sempre aggiornata, si riconnettono da soli e ti chiamano quando qualcosa cambia. Sono il motore di
tutti i controller di Kubernetes, e li userai, senza vederli direttamente, nella
[lezione G.11](11-operator.md).

## Esercizi

!!! esercizio "Esercizio 1 · Scalare da Go"
    Scrivi `scala NAMESPACE DEPLOYMENT REPLICHE`, che cambia il numero di repliche di un Deployment come
    `kubectl scale`. Che permesso RBAC servirebbe se girasse in un Pod?

??? success "Soluzione"
    ```go
    s, err := clientset.AppsV1().Deployments(ns).GetScale(ctx, nome, metav1.GetOptions{})
    if err != nil {
    	return err
    }
    s.Spec.Replicas = int32(repliche)
    _, err = clientset.AppsV1().Deployments(ns).UpdateScale(ctx, nome, s, metav1.UpdateOptions{})
    return err
    ```
    ```yaml
    rules:
      - apiGroups: ["apps"]
        resources: ["deployments/scale"]
        verbs: ["get", "update"]
    ```
    **Spiegazione.** `scale` è una **sottorisorsa** del Deployment: si può dare il permesso di scalare senza dare
    quello di modificare l'intero Deployment (immagine, variabili, comandi). È lo stesso meccanismo che usa
    l'HorizontalPodAutoscaler della [lezione 8.8](../fase-8/08-scheduling.md).

!!! esercizio "Esercizio 2 · Seguire i Pod in tempo reale"
    Aggiungi a `podcheck` l'opzione `-watch`: invece di un solo elenco, resta in ascolto e stampa un problema appena
    compare. Provalo cancellando un Pod della Tasks API in un altro terminale.

??? success "Soluzione"
    ```go
    w, err := clientset.CoreV1().Pods(*namespace).Watch(ctx, metav1.ListOptions{})
    if err != nil {
    	return err
    }
    defer w.Stop()
    for ev := range w.ResultChan() {
    	p, ok := ev.Object.(*corev1.Pod)
    	if !ok || ev.Type == watch.Deleted {
    		continue
    	}
    	for _, problema := range problemi(*p, limite) {
    		fmt.Printf("%s %s/%s: %s\n", time.Now().Format(time.TimeOnly), p.Namespace, p.Name, problema)
    	}
    }
    ```
    **Spiegazione.** Quando cancelli un Pod, il ReplicaSet ne crea uno nuovo: vedi il nuovo Pod passare da
    `Pending` a `Running`, prima segnalato come non pronto e poi sparire dagli avvisi quando la *readiness probe* ha
    successo. Il `ctx` con timeout di 15 secondi dell'elenco non va bene qui: usa quello di
    `signal.NotifyContext` della [lezione G.5](05-concorrenza-http.md#context-timeout-e-annullamenti).

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Con chi parla un programma che usa client-go?","options":["Con etcd direttamente","Con il kubelet di ogni nodo","Con l'API server, come kubectl","Con lo scheduler"],"answer":[2],"explanation":"Tutto passa dall'API server."},
 {"q":"Dove trova le credenziali un programma che usa rest.InClusterConfig()?","options":["In ~/.kube/config","Nel token del service account montato nel Pod","In una variabile scritta nel codice","Nel Secret kube-system/admin"],"answer":[1],"explanation":"Kubernetes lo monta in /var/run/secrets/kubernetes.io/serviceaccount/."},
 {"q":"Quali verbi RBAC servono a podcheck sui Pod?","options":["get e list","create e delete","*","watch e patch"],"answer":[0],"explanation":"Minimo privilegio: solo lettura (più watch se usi -watch)."},
 {"q":"Come carichi in kind un'immagine costruita in locale senza un registry?","options":["docker push","kind load docker-image","kubectl cp","Non si può"],"answer":[1],"explanation":"La copia nei nodi del cluster kind."},
 {"q":"Cosa aggiunge un informer rispetto a un watch semplice?","options":["Niente","Una cache locale aggiornata, riconnessione automatica e callback sui cambiamenti","Permessi di scrittura","La cifratura del traffico"],"answer":[1],"explanation":"È il motore dei controller."}
]}
</script>

## Fonti

- [client-go](https://github.com/kubernetes/client-go) e i suoi [esempi ufficiali](https://github.com/kubernetes/client-go/tree/master/examples)
- [pkg.go.dev: k8s.io/client-go/kubernetes](https://pkg.go.dev/k8s.io/client-go/kubernetes)
- [Kubernetes: Access the API from inside a Pod](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/)
- [Kubernetes: Efficient detection of changes (watch)](https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes)
- [Kubernetes: Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
- [kind: Loading an image into your cluster](https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster)
