// Package check esegue i controlli sui target, anche in parallelo.
package check

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/config"
)

// Result è l'esito di un controllo.
type Result struct {
	Target    config.Target `json:"target"`
	OK        bool          `json:"ok"`
	Duration  time.Duration `json:"duration_ns"`
	Error     string        `json:"error,omitempty"`
	CheckedAt time.Time     `json:"checked_at"`
}

// Checker raccoglie ciò che serve per controllare: il client HTTP e il timeout.
type Checker struct {
	Client  *http.Client
	Timeout time.Duration
}

// New crea un Checker con un client HTTP che non segue i redirect:
// un 301 inatteso è un'informazione utile, non va nascosto.
func New(timeout time.Duration) *Checker {
	return &Checker{
		Timeout: timeout,
		Client: &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
}

// One controlla un singolo target rispettando il timeout.
func (c *Checker) One(ctx context.Context, t config.Target) Result {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	start := time.Now()
	var err error
	switch t.Type {
	case "http":
		err = c.http(ctx, t)
	case "tcp":
		err = tcp(ctx, t)
	default:
		err = fmt.Errorf("tipo %q non supportato", t.Type)
	}
	r := Result{Target: t, OK: err == nil, Duration: time.Since(start), CheckedAt: start}
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

func (c *Checker) http(ctx context.Context, t config.Target) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "racoon-probe")
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != t.ExpectStatus {
		return fmt.Errorf("stato HTTP %d, atteso %d", resp.StatusCode, t.ExpectStatus)
	}
	return nil
}

func tcp(ctx context.Context, t config.Target) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", t.Address)
	if err != nil {
		return err
	}
	return conn.Close()
}

// All controlla tutti i target in parallelo, al massimo `parallel` alla volta,
// e restituisce i risultati nello stesso ordine dei target.
func (c *Checker) All(ctx context.Context, targets []config.Target, parallel int) []Result {
	if parallel < 1 {
		parallel = 1
	}
	results := make([]Result, len(targets))
	sem := make(chan struct{}, parallel) // semaforo: canale con capacità limitata
	var wg sync.WaitGroup
	for i, t := range targets {
		wg.Go(func() {
			sem <- struct{}{}        // occupa un posto
			defer func() { <-sem }() // e lo libera alla fine
			results[i] = c.One(ctx, t)
		})
	}
	wg.Wait()
	return results
}
