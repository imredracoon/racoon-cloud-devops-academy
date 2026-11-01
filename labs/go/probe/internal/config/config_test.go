package config

import (
	"strings"
	"testing"
	"time"
)

func TestParseValoriPredefiniti(t *testing.T) {
	c, err := Parse([]byte(`
targets:
  - name: api
    type: http
    url: http://localhost:8000/healthz
`))
	if err != nil {
		t.Fatalf("errore inatteso: %v", err)
	}
	if c.Interval != 30*time.Second || c.Timeout != 5*time.Second {
		t.Errorf("predefiniti sbagliati: interval=%s timeout=%s", c.Interval, c.Timeout)
	}
	if c.Targets[0].ExpectStatus != 200 {
		t.Errorf("expect_status = %d, voglio 200", c.Targets[0].ExpectStatus)
	}
}

func TestParseErrori(t *testing.T) {
	casi := []struct {
		nome     string
		yaml     string
		contiene string
	}{
		{"nessun target", "interval: 10s", "nessun target"},
		{"tipo sconosciuto", "targets: [{name: x, type: udp}]", "sconosciuto"},
		{"url mancante", "targets: [{name: x, type: http}]", "url mancante"},
		{"address mancante", "targets: [{name: x, type: tcp}]", "address mancante"},
		{"duplicato", "targets: [{name: x, type: tcp, address: a:1}, {name: x, type: tcp, address: b:1}]", "duplicato"},
		{"timeout troppo lungo", "interval: 5s\ntimeout: 10s\ntargets: [{name: x, type: tcp, address: a:1}]", "minore di interval"},
		{"yaml rotto", "targets: [", "YAML non valido"},
	}
	for _, tc := range casi {
		t.Run(tc.nome, func(t *testing.T) {
			_, err := Parse([]byte(tc.yaml))
			if err == nil || !strings.Contains(err.Error(), tc.contiene) {
				t.Fatalf("errore = %v, voglio che contenga %q", err, tc.contiene)
			}
		})
	}
}
