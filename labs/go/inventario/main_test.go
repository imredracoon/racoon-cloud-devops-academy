package main

import (
	"context"
	"errors"
	"testing"
)

// finto è un Fornitore che non parla con nessun cloud: serve ai test.
type finto struct {
	istanze []Risorsa
	err     error
}

func (f finto) Istanze(context.Context) ([]Risorsa, error) { return f.istanze, f.err }
func (f finto) Bucket(context.Context) ([]Risorsa, error) {
	return []Risorsa{{Tipo: "bucket", Nome: "log"}}, nil
}

func TestInventario(t *testing.T) {
	f := finto{istanze: []Risorsa{{Tipo: "istanza", Nome: "web01", Stato: "running"}}}
	r, err := Inventario(context.Background(), f)
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 2 || r[0].Nome != "web01" || r[1].Nome != "log" {
		t.Errorf("inventario inatteso: %+v", r)
	}
}

func TestInventarioErrore(t *testing.T) {
	errAPI := errors.New("AccessDenied")
	_, err := Inventario(context.Background(), finto{err: errAPI})
	if !errors.Is(err, errAPI) {
		t.Errorf("errore = %v, voglio che avvolga %v", err, errAPI)
	}
}
