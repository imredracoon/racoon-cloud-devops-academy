# G.11 Un operator con Kubebuilder

!!! obiettivi "Obiettivi"
    - Capire il *pattern* operator: risorsa personalizzata + controller che riconcilia.
    - Definire una **CustomResourceDefinition** a partire da una struct Go con i marcatori di Kubebuilder.
    - Scrivere la funzione `Reconcile` con controller-runtime: leggere, osservare, aggiornare lo status, ripianificare.
    - Testare il controller senza cluster con il client finto, e poi eseguirlo in kind.

!!! info "Quando studiarla"
    Nella [Fase 8](../fase-8/index.md), dopo la [lezione 8.13](../fase-8/13-helm-kustomize.md), dove hai visto cosa fa
    un operator dal punto di vista di chi lo usa. Qui ne scrivi uno. Serve la [lezione G.10](10-client-go.md).

## Il pattern operator

Kubernetes funziona per **riconciliazione**: tu dichiari lo stato desiderato (`replicas: 3`), e un **controller**
osserva continuamente lo stato reale e agisce per avvicinarlo a quello desiderato. Il controller dei ReplicaSet fa
questo con i Pod, quello dei Deployment con i ReplicaSet.

Un **operator** applica la stessa idea a qualcosa che Kubernetes non conosce:

1. una **CustomResourceDefinition** (CRD) aggiunge all'API un nuovo tipo di risorsa, con il suo schema;
2. un **controller** scritto da te osserva gli oggetti di quel tipo e fa ciò che serve.

In questa lezione racoon-probe diventa un operator: invece di un file YAML, i controlli sono **risorse Kubernetes**.

```yaml
apiVersion: monitoring.academy.racoon.dev/v1alpha1
kind: Probe
metadata:
  name: tasks-api
  namespace: tasks
spec:
  url: http://tasks-api.tasks.svc.cluster.local/health
  interval: 15s
```

```console
$ kubectl get probes -A
NAMESPACE   NAME         URL                                               DISPONIBILE   DURATA   ETÀ
tasks       tasks-api    http://tasks-api.tasks.svc.cluster.local/health   True          4ms      2m
default     sito-corso   https://imredracoon.github.io/racoon-cloud-...    True          212ms    2m
```

Ogni team può dichiarare i controlli dei propri servizi nel proprio namespace, con i permessi RBAC e il flusso GitOps
che già usa ([lezione 8.15](../fase-8/15-gitops-cloud.md)).

## Gli strumenti

| Strumento | A cosa serve |
|---|---|
| **controller-runtime** | la libreria che gestisce informer, cache, code di lavoro e riconnessioni: tu scrivi solo `Reconcile` |
| **controller-tools** (`controller-gen`) | genera la CRD, il ClusterRole e il codice di copia profonda dai **marcatori** nei commenti |
| **Kubebuilder** | crea lo scheletro del progetto (Makefile, Dockerfile, manifest, test) e i file per ogni nuova API |

Kubebuilder è un progetto ufficiale del gruppo Kubernetes SIG API Machinery. L'alternativa più diffusa,
**Operator SDK**, usa Kubebuilder al suo interno per i progetti in Go.

```bash
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
kubebuilder version
```

## Lo scheletro

```bash
mkdir -p ~/go-lab/probe-operator && cd ~/go-lab/probe-operator
kubebuilder init --domain academy.racoon.dev --repo github.com/TUONOME/probe-operator
kubebuilder create api --group monitoring --version v1alpha1 --kind Probe --resource --controller
```

I file su cui lavori sono due; il resto è generato:

```text
api/v1alpha1/probe_types.go             # la struct: diventa lo schema della CRD
internal/controller/probe_controller.go # la funzione Reconcile
config/crd/, config/rbac/               # generati da "make manifests"
cmd/main.go                             # avvia il manager (lo genera Kubebuilder)
```

La versione di riferimento è in
[`labs/go/probe-operator`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/probe-operator):
contiene solo i file essenziali, con un `cmd/main.go` ridotto al minimo e un piccolo Makefile, per leggerla in
pochi minuti. Il progetto generato da Kubebuilder è più ricco (metriche protette, webhook, test end-to-end), ma fa
le stesse cose.

## Il tipo: da struct a CRD

```go title="api/v1alpha1/probe_types.go"
type ProbeSpec struct {
	// URL da controllare con una richiesta GET.
	// +kubebuilder:validation:Pattern=`^https?://`
	URL string `json:"url"`

	// Codice HTTP atteso.
	// +kubebuilder:default=200
	// +kubebuilder:validation:Minimum=100
	// +kubebuilder:validation:Maximum=599
	// +optional
	ExpectStatus int32 `json:"expectStatus,omitempty"`

	// Ogni quanto ripetere il controllo.
	// +kubebuilder:default="30s"
	// +optional
	Interval metav1.Duration `json:"interval,omitempty"`
}

type ProbeStatus struct {
	// +optional
	LastCheckTime *metav1.Time `json:"lastCheckTime,omitempty"`
	// +optional
	LastDuration string `json:"lastDuration,omitempty"`
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.spec.url`
// +kubebuilder:printcolumn:name="Disponibile",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].status`
// +kubebuilder:printcolumn:name="Durata",type=string,JSONPath=`.status.lastDuration`
// +kubebuilder:printcolumn:name="Età",type=date,JSONPath=`.metadata.creationTimestamp`

type Probe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbeSpec   `json:"spec"`
	Status ProbeStatus `json:"status,omitempty"`
}
```

I commenti che iniziano con `+kubebuilder:` sono **marcatori**: `controller-gen` li legge e li trasforma nello
schema OpenAPI della CRD. L'API server rifiuterà così un `url: ftp://...` o un `expectStatus: 42` **prima** che
arrivino al tuo controller, e riempirà da solo i valori predefiniti.

**Spec** e **status** sono separati come in tutte le risorse di Kubernetes: lo spec lo scrive l'utente, lo status lo
scrive il controller. `+kubebuilder:subresource:status` crea la sottorisorsa `/status`, così i due non si
sovrascrivono a vicenda e il permesso di scrivere lo status si dà a parte. Le **condizioni** (`metav1.Condition`)
sono il formato standard per dire "come sta" una risorsa: le stesse che vedi in `kubectl describe deployment`.

```bash
make manifests generate    # CRD in config/crd/, ClusterRole in config/rbac/, zz_generated.deepcopy.go
```

## Reconcile

```go title="internal/controller/probe_controller.go"
// +kubebuilder:rbac:groups=monitoring.academy.racoon.dev,resources=probes,verbs=get;list;watch
// +kubebuilder:rbac:groups=monitoring.academy.racoon.dev,resources=probes/status,verbs=get;update;patch

func (r *ProbeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var p monitoringv1alpha1.Probe
	if err := r.Get(ctx, req.NamespacedName, &p); err != nil {      // (1)!
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	inizio := time.Now()
	errCheck := r.controlla(ctx, p.Spec.URL, int(p.Spec.ExpectStatus))   // (2)!
	durata := time.Since(inizio)

	cond := metav1.Condition{
		Type:               "Available",
		Status:             metav1.ConditionTrue,
		Reason:             "CheckSucceeded",
		Message:            fmt.Sprintf("risposta %d in %s", p.Spec.ExpectStatus, durata.Round(time.Millisecond)),
		ObservedGeneration: p.Generation,
	}
	if errCheck != nil {
		cond.Status, cond.Reason, cond.Message = metav1.ConditionFalse, "CheckFailed", errCheck.Error()
		log.Info("controllo fallito", "probe", req.NamespacedName, "errore", errCheck)
	}
	meta.SetStatusCondition(&p.Status.Conditions, cond)              // (3)!
	adesso := metav1.Now()
	p.Status.LastCheckTime = &adesso
	p.Status.LastDuration = durata.Round(time.Millisecond).String()
	if err := r.Status().Update(ctx, &p); err != nil {
		return ctrl.Result{}, err                                    // (4)!
	}

	return ctrl.Result{RequeueAfter: p.Spec.Interval.Duration}, nil  // (5)!
}
```

1.  `Reconcile` riceve solo **nome e namespace**, non l'oggetto né "cosa è cambiato". Rilegge sempre lo stato
    attuale dalla cache. Se l'oggetto è stato cancellato, `IgnoreNotFound` trasforma l'errore in "niente da fare".
2.  La stessa logica di controllo HTTP della [lezione G.4](04-cli-test.md), con un timeout.
3.  `SetStatusCondition` aggiunge la condizione o aggiorna quella con lo stesso `Type`, e cambia
    `lastTransitionTime` solo quando lo stato passa davvero da `True` a `False` o viceversa.
4.  Se l'aggiornamento fallisce (ad esempio perché qualcuno ha modificato l'oggetto nel frattempo), restituire
    l'errore basta: controller-runtime rimette la richiesta in coda e riprova, con attese crescenti.
5.  Nessun ciclo con `time.Sleep`: chiedi di essere richiamato tra `interval`. Se nel frattempo l'utente modifica
    il Probe, il controller viene chiamato subito.

!!! tip "Reconcile deve essere idempotente"
    La funzione può essere chiamata molte volte per lo stesso oggetto, per un cambiamento, per un riavvio
    dell'operator o per una risincronizzazione periodica. Deve dare lo stesso risultato ogni volta, come un
    `terraform apply` senza modifiche ([lezione G.8](08-terratest.md)). Per questo non "reagisce a un evento": guarda
    lo stato e lo corregge.

## Test senza cluster

controller-runtime offre un **client finto** che simula l'API server in memoria. Il test crea un Probe, chiama
`Reconcile` e controlla lo status, contro un server HTTP finto di `httptest`:

```go title="internal/controller/probe_controller_test.go (estratto)"
cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(p).WithStatusSubresource(p).Build()
r := &ProbeReconciler{Client: cl, Scheme: scheme, HTTP: bersaglio.Client()}

res, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: chiave})
// ...
cond := meta.FindStatusCondition(aggiornato.Status.Conditions, "Available")
if cond == nil || cond.Status != c.stato {
	t.Errorf("condizione Available = %+v, voglio %s", cond, c.stato)
}
```

Il progetto generato da Kubebuilder include anche test con **envtest**, che avvia un vero API server ed etcd senza
nodi: più lenti, ma verificano anche lo schema della CRD.

## Laboratorio

!!! laboratorio "Laboratorio · L'operator Probe in kind"
    1. Crea il progetto con Kubebuilder, scrivi `ProbeSpec`, `ProbeStatus` e i marcatori, poi `make manifests`.
       Apri la CRD generata e trova lo schema di `url` con il `pattern`.
    2. Scrivi `Reconcile` e il test con il client finto; `make test` deve essere verde.
    3. Esegui l'operator **fuori dal cluster** con `make run` (usa il tuo kubeconfig) e applica un Probe verso un
       sito pubblico. Guarda `kubectl get probes -w`.
    4. Prova la validazione: applica un Probe con `url: ftp://esempio` e leggi l'errore dell'API server.
    5. Costruisci l'immagine, caricala in kind ed esegui l'operator **nel cluster**. Applica il Probe verso la
       Tasks API (`config/samples/probes.yaml`).
    6. Rompi la Tasks API (scala il database a zero repliche) e guarda il Probe diventare `False`, con il motivo
       nel messaggio. Poi ripristina.

??? success "Soluzione guidata"
    Con la versione di riferimento i comandi sono:

    ```bash
    cd labs/go/probe-operator
    make test
    make kind-deploy                                  # CRD, ClusterRole, Deployment dell'operator
    kubectl -n probe-operator logs deploy/probe-operator -f &
    kubectl apply -f config/samples/probes.yaml
    kubectl get probes -A -w
    ```
    ```console
    $ kubectl apply -f - <<'EOF'
    apiVersion: monitoring.academy.racoon.dev/v1alpha1
    kind: Probe
    metadata: { name: sbagliato }
    spec: { url: "ftp://esempio" }
    EOF
    The Probe "sbagliato" is invalid: spec.url: Invalid value: "ftp://esempio": spec.url in body should match '^https?://'
    ```
    ```bash
    kubectl -n tasks scale statefulset/postgres --replicas=0   # /health della Tasks API risponde 503
    kubectl -n tasks describe probe tasks-api                   # Reason: CheckFailed, stato HTTP 503, atteso 200
    kubectl -n tasks scale statefulset/postgres --replicas=1
    ```
    **Spiegazione.** Fuori dal cluster l'operator non risolve i nomi `*.svc.cluster.local`, che esistono solo nel DNS
    interno del cluster: per questo il Probe verso la Tasks API funziona solo con l'operator nel cluster. La
    validazione è fatta dall'**API server** con lo schema della CRD: il tuo codice non vede mai l'oggetto non
    valido. Con il database fermo, `/health` restituisce 503 e la condizione passa a `False` al controllo
    successivo. Il nome dello StatefulSet dipende dai tuoi manifest: controlla con `kubectl -n tasks get sts`.

## Esercizi

!!! esercizio "Esercizio 1 · Un evento quando cambia"
    Fai emettere all'operator un **Event** Kubernetes quando un Probe passa da disponibile a non disponibile, così
    compare in `kubectl describe probe` e in `kubectl get events`.

??? success "Soluzione"
    ```go
    // nel ProbeReconciler: Recorder recorder.EventRecorder (sigs.k8s.io/controller-runtime/pkg/recorder),
    // impostato in main con mgr.GetEventRecorder("probe-operator")
    prima := meta.IsStatusConditionTrue(p.Status.Conditions, "Available")
    meta.SetStatusCondition(&p.Status.Conditions, cond)
    if prima && errCheck != nil {
    	r.Recorder.Eventf(&p, nil, corev1.EventTypeWarning, "CheckFailed", "Controllo", "%s non risponde: %v", p.Spec.URL, errCheck)
    }
    ```
    ```go
    // +kubebuilder:rbac:groups=events.k8s.io,resources=events,verbs=create;patch
    ```
    **Spiegazione.** L'evento va emesso solo **al cambio di stato**, non a ogni controllo fallito, altrimenti in
    pochi minuti riempi il cluster di eventi uguali. Il marcatore RBAC aggiunge al ClusterRole il permesso di
    creare eventi; senza, l'operator non può scriverli. Nel codice meno recente troverai `GetEventRecorderFor`, che
    usa la vecchia API degli eventi ed è deprecata.

!!! esercizio "Esercizio 2 · Metriche dall'operator"
    Esponi dall'operator la metrica `probe_up{namespace,name}` della [lezione G.9](09-ci-metriche.md). Suggerimento:
    controller-runtime ha già un registro Prometheus e un endpoint `/metrics`.

??? success "Soluzione"
    ```go
    import "sigs.k8s.io/controller-runtime/pkg/metrics"

    var probeUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    	Name: "probe_up", Help: "1 se l'ultimo controllo è riuscito.",
    }, []string{"namespace", "name"})

    func init() { metrics.Registry.MustRegister(probeUp) }

    // in Reconcile, dopo il controllo:
    probeUp.WithLabelValues(p.Namespace, p.Name).Set(map[bool]float64{true: 1, false: 0}[errCheck == nil])
    ```
    **Spiegazione.** `metrics.Registry` è il registro che il manager espone già su `/metrics` (nel progetto di
    Kubebuilder con autenticazione, sulla porta 8443). Quando un Probe viene cancellato la serie resta: per
    toglierla usa `probeUp.DeleteLabelValues(ns, nome)` nel ramo `IsNotFound`, oppure un *finalizer* se devi fare
    pulizia anche fuori da Kubernetes.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Da cosa è composto un operator?","options":["Un Helm chart","Una CustomResourceDefinition e un controller che riconcilia gli oggetti di quel tipo","Un DaemonSet","Uno script cron nel cluster"],"answer":[1],"explanation":"Estende Kubernetes con un nuovo tipo e la logica per gestirlo."},
 {"q":"Cosa riceve la funzione Reconcile?","options":["L'oggetto modificato e la differenza","Solo nome e namespace dell'oggetto da riconciliare","Tutti gli oggetti del cluster","L'evento HTTP originale"],"answer":[1],"explanation":"Rilegge lo stato attuale ogni volta: per questo è idempotente."},
 {"q":"Come si chiede a controller-runtime di richiamare Reconcile tra 30 secondi?","options":["time.Sleep(30*time.Second)","return ctrl.Result{RequeueAfter: 30*time.Second}, nil","Un goroutine con un ticker","Un CronJob"],"answer":[1],"explanation":"Il controller resta libero di lavorare su altri oggetti."},
 {"q":"Chi rifiuta un Probe con url: ftp://esempio?","options":["Il controller","L'API server, grazie allo schema della CRD generato dai marcatori","kubectl","etcd"],"answer":[1],"explanation":"La validazione avviene prima che l'oggetto venga salvato."},
 {"q":"Perché spec e status sono separati e status è una sottorisorsa?","options":["Per estetica","Perché l'utente scrive lo spec, il controller lo status, senza sovrascriversi, con permessi distinti","Perché etcd lo richiede","Per ridurre lo spazio"],"answer":[1],"explanation":"È lo stesso modello di tutte le risorse di Kubernetes."},
 {"q":"Cosa genera `make manifests` in un progetto Kubebuilder?","options":["L'immagine Docker","La CRD e le regole RBAC a partire dai marcatori +kubebuilder","Il codice di Reconcile","I test"],"answer":[1],"explanation":"Lo fa controller-gen; make generate produce il codice di copia profonda."}
]}
</script>

## Fonti

- [The Kubebuilder Book](https://book.kubebuilder.io/), in particolare [Tutorial: Building CronJob](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial)
- [Kubebuilder: CRD generation e marcatori](https://book.kubebuilder.io/reference/markers/crd-validation)
- [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) e [pkg.go.dev](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [Kubernetes: Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Kubernetes: Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Kubernetes API conventions: spec, status e conditions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#spec-and-status)
- [Operator SDK](https://sdk.operatorframework.io/)
