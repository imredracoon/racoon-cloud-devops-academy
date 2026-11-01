package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	monitoringv1alpha1 "github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe-operator/api/v1alpha1"
)

// Il client "fake" di controller-runtime simula l'API server in memoria.
func TestReconcile(t *testing.T) {
	bersaglio := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rotto" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer bersaglio.Close()

	casi := []struct {
		nome     string
		percorso string
		stato    metav1.ConditionStatus
	}{
		{"risponde", "/", metav1.ConditionTrue},
		{"rotto", "/rotto", metav1.ConditionFalse},
	}
	for _, c := range casi {
		t.Run(c.nome, func(t *testing.T) {
			scheme := runtime.NewScheme()
			if err := monitoringv1alpha1.AddToScheme(scheme); err != nil {
				t.Fatal(err)
			}
			p := &monitoringv1alpha1.Probe{
				ObjectMeta: metav1.ObjectMeta{Name: "web", Namespace: "default"},
				Spec: monitoringv1alpha1.ProbeSpec{
					URL:          bersaglio.URL + c.percorso,
					ExpectStatus: 200,
					Interval:     metav1.Duration{Duration: 30 * time.Second},
				},
			}
			cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(p).
				WithStatusSubresource(p).Build()
			r := &ProbeReconciler{Client: cl, Scheme: scheme, HTTP: bersaglio.Client()}

			chiave := types.NamespacedName{Name: "web", Namespace: "default"}
			res, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: chiave})
			if err != nil {
				t.Fatal(err)
			}
			if res.RequeueAfter != 30*time.Second {
				t.Errorf("RequeueAfter = %s, voglio 30s", res.RequeueAfter)
			}

			var aggiornato monitoringv1alpha1.Probe
			if err := cl.Get(context.Background(), chiave, &aggiornato); err != nil {
				t.Fatal(err)
			}
			cond := meta.FindStatusCondition(aggiornato.Status.Conditions, "Available")
			if cond == nil || cond.Status != c.stato {
				t.Errorf("condizione Available = %+v, voglio %s", cond, c.stato)
			}
		})
	}
}
