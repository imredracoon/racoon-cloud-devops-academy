// Punto di avvio dell'operator (lezione G.11 del percorso Go). Kubebuilder genera un
// main.go più ricco (metriche protette, webhook, leader election configurabile):
// questa è la versione minima, utile per capire cosa fa.
package main

import (
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	monitoringv1alpha1 "github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe-operator/api/v1alpha1"
	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe-operator/internal/controller"
)

func main() {
	ctrl.SetLogger(zap.New())
	log := ctrl.Log.WithName("avvio")

	// Lo scheme dice al client quali tipi conosce: quelli di Kubernetes più il nostro Probe.
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = monitoringv1alpha1.AddToScheme(scheme)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: ":8081",
		LeaderElection:         os.Getenv("LEADER_ELECTION") == "true",
		LeaderElectionID:       "probe-operator.academy.racoon.dev",
	})
	if err != nil {
		log.Error(err, "manager non creato")
		os.Exit(1)
	}

	if err := (&controller.ProbeReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		HTTP:   &http.Client{Timeout: 10 * time.Second},
	}).SetupWithManager(mgr); err != nil {
		log.Error(err, "controller non registrato")
		os.Exit(1)
	}

	_ = mgr.AddHealthzCheck("healthz", healthz.Ping)
	_ = mgr.AddReadyzCheck("readyz", healthz.Ping)

	log.Info("operator avviato")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil { // si ferma con SIGTERM
		log.Error(err, "arresto con errore")
		os.Exit(1)
	}
}
