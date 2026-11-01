package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/check"
	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/config"
)

func TestServer(t *testing.T) {
	bersaglio := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer bersaglio.Close()

	cfg := config.Config{
		Interval: time.Minute,
		Timeout:  time.Second,
		Targets: []config.Target{
			{Name: "web", Type: "http", URL: bersaglio.URL, ExpectStatus: 200},
			{Name: "spento", Type: "tcp", Address: "127.0.0.1:1"},
		},
	}
	s := New(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	s.RunOnce(context.Background())

	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	t.Run("healthz", func(t *testing.T) {
		if got := leggi(t, srv.URL+"/healthz"); got != "ok\n" {
			t.Errorf("healthz = %q", got)
		}
	})

	t.Run("results", func(t *testing.T) {
		var res []check.Result
		if err := json.Unmarshal([]byte(leggi(t, srv.URL+"/results")), &res); err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 || !res[0].OK || res[1].OK {
			t.Errorf("risultati inattesi: %+v", res)
		}
	})

	t.Run("metrics", func(t *testing.T) {
		body := leggi(t, srv.URL+"/metrics")
		for _, riga := range []string{
			`probe_up{target="web",type="http"} 1`,
			`probe_up{target="spento",type="tcp"} 0`,
			`probe_checks_total{result="failure",target="spento"} 1`,
		} {
			if !strings.Contains(body, riga) {
				t.Errorf("manca %q in /metrics", riga)
			}
		}
	})

	t.Run("metodo non permesso", func(t *testing.T) {
		resp, err := http.Post(srv.URL+"/results", "text/plain", nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("POST /results = %d, voglio 405", resp.StatusCode)
		}
	})
}

func leggi(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
