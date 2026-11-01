// Package config legge e controlla il file di configurazione di probe.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

// Target è un singolo servizio da controllare.
type Target struct {
	Name         string `yaml:"name" json:"name"`
	Type         string `yaml:"type" json:"type"`                                       // "http" oppure "tcp"
	URL          string `yaml:"url,omitempty" json:"url,omitempty"`                     // solo per http
	Address      string `yaml:"address,omitempty" json:"address,omitempty"`             // solo per tcp, host:porta
	ExpectStatus int    `yaml:"expect_status,omitempty" json:"expect_status,omitempty"` // solo per http, 200 se assente
}

// Config è l'intero file probe.yaml.
type Config struct {
	Interval time.Duration `yaml:"interval"` // ogni quanto ripetere i controlli in modalità serve
	Timeout  time.Duration `yaml:"timeout"`  // tempo massimo per un singolo controllo
	Targets  []Target      `yaml:"targets"`
}

// Load legge il file indicato, applica i valori predefiniti e lo valida.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // il percorso lo sceglie chi avvia probe: è voluto
	if err != nil {
		return Config{}, fmt.Errorf("lettura della configurazione: %w", err)
	}
	return Parse(data)
}

// Parse interpreta il contenuto YAML. È separata da Load per poterla testare senza file.
func Parse(data []byte) (Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("YAML non valido: %w", err)
	}
	if c.Interval == 0 {
		c.Interval = 30 * time.Second
	}
	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
	}
	for i := range c.Targets {
		if c.Targets[i].Type == "http" && c.Targets[i].ExpectStatus == 0 {
			c.Targets[i].ExpectStatus = 200
		}
	}
	return c, c.Validate()
}

// Validate restituisce tutti gli errori trovati, non solo il primo.
func (c Config) Validate() error {
	var errs []error
	if len(c.Targets) == 0 {
		errs = append(errs, errors.New("nessun target definito"))
	}
	if c.Timeout >= c.Interval {
		errs = append(errs, fmt.Errorf("timeout (%s) deve essere minore di interval (%s)", c.Timeout, c.Interval))
	}
	visti := map[string]bool{}
	for i, t := range c.Targets {
		switch {
		case t.Name == "":
			errs = append(errs, fmt.Errorf("target %d: name mancante", i+1))
		case visti[t.Name]:
			errs = append(errs, fmt.Errorf("target %q: nome duplicato", t.Name))
		}
		visti[t.Name] = true
		switch t.Type {
		case "http":
			if t.URL == "" {
				errs = append(errs, fmt.Errorf("target %q: url mancante", t.Name))
			}
		case "tcp":
			if t.Address == "" {
				errs = append(errs, fmt.Errorf("target %q: address mancante", t.Name))
			}
		default:
			errs = append(errs, fmt.Errorf("target %q: type %q sconosciuto (usa http o tcp)", t.Name, t.Type))
		}
	}
	return errors.Join(errs...)
}
