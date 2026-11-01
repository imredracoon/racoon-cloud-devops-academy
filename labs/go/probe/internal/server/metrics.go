package server

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/check"
)

// metrics contiene le metriche esposte su /metrics. Ogni Server ha il suo registro,
// così i test possono creare più server senza conflitti.
type metrics struct {
	up       *prometheus.GaugeVec
	duration *prometheus.HistogramVec
	checks   *prometheus.CounterVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		up: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_up",
			Help: "1 se l'ultimo controllo del target è riuscito, 0 altrimenti.",
		}, []string{"target", "type"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "probe_duration_seconds",
			Help:    "Durata dei controlli, in secondi.",
			Buckets: prometheus.DefBuckets,
		}, []string{"target", "type"}),
		checks: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "probe_checks_total",
			Help: "Numero di controlli eseguiti, per esito.",
		}, []string{"target", "result"}),
	}
	reg.MustRegister(m.up, m.duration, m.checks)
	return m
}

func (m *metrics) record(r check.Result) {
	up, esito := 0.0, "failure"
	if r.OK {
		up, esito = 1, "success"
	}
	m.up.WithLabelValues(r.Target.Name, r.Target.Type).Set(up)
	m.duration.WithLabelValues(r.Target.Name, r.Target.Type).Observe(r.Duration.Seconds())
	m.checks.WithLabelValues(r.Target.Name, esito).Inc()
}
