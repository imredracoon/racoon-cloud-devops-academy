package check

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/config"
)

func TestHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rotto" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	c := New(time.Second)
	casi := []struct {
		nome   string
		target config.Target
		ok     bool
	}{
		{"risponde 200", config.Target{Name: "a", Type: "http", URL: srv.URL + "/", ExpectStatus: 200}, true},
		{"risponde 500", config.Target{Name: "b", Type: "http", URL: srv.URL + "/rotto", ExpectStatus: 200}, false},
		{"500 atteso", config.Target{Name: "c", Type: "http", URL: srv.URL + "/rotto", ExpectStatus: 500}, true},
	}
	for _, tc := range casi {
		t.Run(tc.nome, func(t *testing.T) {
			r := c.One(context.Background(), tc.target)
			if r.OK != tc.ok {
				t.Errorf("OK = %v, voglio %v (errore: %s)", r.OK, tc.ok, r.Error)
			}
		})
	}
}

func TestHTTPTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done():
		}
	}))
	defer srv.Close()

	r := New(100*time.Millisecond).One(context.Background(),
		config.Target{Name: "lento", Type: "http", URL: srv.URL, ExpectStatus: 200})
	if r.OK || r.Duration > time.Second {
		t.Fatalf("il timeout non ha funzionato: ok=%v durata=%s", r.OK, r.Duration)
	}
}

func TestTCP(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0") // porta libera scelta dal sistema
	if err != nil {
		t.Fatal(err)
	}
	indirizzo := ln.Addr().String()

	c := New(time.Second)
	if r := c.One(context.Background(), config.Target{Name: "aperta", Type: "tcp", Address: indirizzo}); !r.OK {
		t.Errorf("porta aperta segnalata come chiusa: %s", r.Error)
	}
	_ = ln.Close()
	if r := c.One(context.Background(), config.Target{Name: "chiusa", Type: "tcp", Address: indirizzo}); r.OK {
		t.Error("porta chiusa segnalata come aperta")
	}
}

func TestAllMantieneOrdine(t *testing.T) {
	targets := []config.Target{
		{Name: "uno", Type: "tcp", Address: "127.0.0.1:1"},
		{Name: "due", Type: "tcp", Address: "127.0.0.1:1"},
		{Name: "tre", Type: "tcp", Address: "127.0.0.1:1"},
	}
	res := New(time.Second).All(context.Background(), targets, 2)
	for i, r := range res {
		if r.Target.Name != targets[i].Name {
			t.Errorf("posizione %d: %s, voglio %s", i, r.Target.Name, targets[i].Name)
		}
	}
}
