package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1alpha1 "github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe-operator/api/v1alpha1"
)

// ProbeReconciler riconcilia gli oggetti Probe.
type ProbeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	HTTP   *http.Client
}

// +kubebuilder:rbac:groups=monitoring.academy.racoon.dev,resources=probes,verbs=get;list;watch
// +kubebuilder:rbac:groups=monitoring.academy.racoon.dev,resources=probes/status,verbs=get;update;patch

func (r *ProbeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Leggi l'oggetto. Se è stato cancellato non c'è niente da fare.
	var p monitoringv1alpha1.Probe
	if err := r.Get(ctx, req.NamespacedName, &p); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Osserva il mondo: esegui il controllo.
	inizio := time.Now()
	errCheck := r.controlla(ctx, p.Spec.URL, int(p.Spec.ExpectStatus))
	durata := time.Since(inizio)

	// 3. Scrivi ciò che hai osservato nello status.
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
	meta.SetStatusCondition(&p.Status.Conditions, cond)
	adesso := metav1.Now()
	p.Status.LastCheckTime = &adesso
	p.Status.LastDuration = durata.Round(time.Millisecond).String()
	if err := r.Status().Update(ctx, &p); err != nil {
		return ctrl.Result{}, err // un errore fa ripetere la riconciliazione, con attese crescenti
	}

	// 4. Torna a controllare dopo l'intervallo.
	return ctrl.Result{RequeueAfter: p.Spec.Interval.Duration}, nil
}

func (r *ProbeReconciler) controlla(ctx context.Context, url string, atteso int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := r.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != atteso {
		return fmt.Errorf("stato HTTP %d, atteso %d", resp.StatusCode, atteso)
	}
	return nil
}

// SetupWithManager registra il controller: guarda gli oggetti Probe.
func (r *ProbeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Probe{}).
		Named("probe").
		Complete(r)
}
