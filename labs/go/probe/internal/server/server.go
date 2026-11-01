// Package server trasforma probe in un servizio: controlla i target a intervalli
// regolari ed espone i risultati in JSON e in formato Prometheus.
package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/check"
	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/config"
)

// Server tiene la configurazione, l'ultimo giro di risultati e le metriche.
type Server struct {
	cfg     config.Config
	checker *check.Checker
	log     *slog.Logger
	reg     *prometheus.Registry
	metrics *metrics

	mu      sync.RWMutex // protegge results: lo scrive il ciclo, lo leggono gli handler
	results []check.Result
}

// New prepara il server senza avviarlo.
func New(cfg config.Config, log *slog.Logger) *Server {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	return &Server{
		cfg:     cfg,
		checker: check.New(cfg.Timeout),
		log:     log,
		reg:     reg,
		metrics: newMetrics(reg),
	}
}

// RunOnce esegue un giro di controlli e aggiorna risultati e metriche.
func (s *Server) RunOnce(ctx context.Context) {
	res := s.checker.All(ctx, s.cfg.Targets, 8)
	for _, r := range res {
		s.metrics.record(r)
		if !r.OK {
			s.log.Warn("controllo fallito", "target", r.Target.Name, "errore", r.Error)
		}
	}
	s.mu.Lock()
	s.results = res
	s.mu.Unlock()
}

// Loop ripete i controlli finché il contesto non viene annullato.
func (s *Server) Loop(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()
	for {
		s.RunOnce(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// Handler restituisce le rotte HTTP del servizio.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok\n")) // se il client se n'è già andato non c'è altro da fare
	})
	mux.HandleFunc("GET /results", s.handleResults)
	mux.Handle("GET /metrics", promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{}))
	return mux
}

func (s *Server) handleResults(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	res := s.results
	s.mu.RUnlock()
	if res == nil {
		res = []check.Result{} // "[]" in JSON invece di "null" prima del primo giro
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error("risposta non inviata", "errore", err)
	}
}
