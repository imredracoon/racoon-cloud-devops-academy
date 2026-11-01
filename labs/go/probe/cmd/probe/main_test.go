package main

import (
	"bytes"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	dir := t.TempDir()
	scrivi := func(nome, contenuto string) string {
		p := filepath.Join(dir, nome)
		if err := os.WriteFile(p, []byte(contenuto), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	buono := scrivi("buono.yaml", "targets: [{name: locale, type: tcp, address: "+ln.Addr().String()+"}]")
	cattivo := scrivi("cattivo.yaml", "targets: [{name: spento, type: tcp, address: 127.0.0.1:1}]")

	var out bytes.Buffer
	if err := run([]string{"check", "-config", buono}, &out); err != nil {
		t.Fatalf("errore inatteso: %v", err)
	}
	if !strings.Contains(out.String(), "locale") || !strings.Contains(out.String(), "OK") {
		t.Errorf("uscita inattesa:\n%s", out.String())
	}

	out.Reset()
	if err := run([]string{"check", "-config", cattivo, "-o", "json"}, &out); !errors.Is(err, errFalliti) {
		t.Errorf("errore = %v, voglio errFalliti", err)
	}
	if !strings.Contains(out.String(), `"ok": false`) {
		t.Errorf("JSON inatteso:\n%s", out.String())
	}
}

func TestComandoSconosciuto(t *testing.T) {
	if err := run([]string{"boh"}, &bytes.Buffer{}); err == nil {
		t.Error("mi aspettavo un errore")
	}
}
