// probe controlla che i servizi rispondano. Due modi d'uso:
//
//	probe check -config probe.yaml   un giro di controlli, esce con 1 se qualcosa non va
//	probe serve -config probe.yaml   servizio con /healthz, /results e /metrics
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/check"
	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/config"
	"github.com/imredracoon/racoon-cloud-devops-academy/labs/go/probe/internal/server"
)

// version viene sostituita in fase di build: go build -ldflags "-X main.version=v1.0.0"
var version = "dev"

const uso = `uso: probe <comando> [opzioni]

comandi:
  check     esegue un giro di controlli e stampa i risultati
  serve     avvia il servizio HTTP con /healthz, /results e /metrics
  version   stampa la versione
`

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}
}

// run è separata da main per poterla testare: riceve gli argomenti e dove scrivere.
func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		fmt.Fprint(out, uso)
		return errors.New("comando mancante")
	}
	switch args[0] {
	case "check":
		return runCheck(args[1:], out)
	case "serve":
		return runServe(args[1:])
	case "version":
		fmt.Fprintln(out, version)
		return nil
	default:
		fmt.Fprint(out, uso)
		return fmt.Errorf("comando %q sconosciuto", args[0])
	}
}

var errFalliti = errors.New("almeno un controllo è fallito")

func runCheck(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	path := fs.String("config", "probe.yaml", "file di configurazione")
	formato := fs.String("o", "table", "formato di uscita: table o json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := config.Load(*path)
	if err != nil {
		return err
	}

	res := check.New(cfg.Timeout).All(context.Background(), cfg.Targets, 8)

	switch *formato {
	case "json":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res); err != nil {
			return err
		}
	case "table":
		tw := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "TARGET\tTIPO\tESITO\tDURATA\tERRORE")
		for _, r := range res {
			esito := "OK"
			if !r.OK {
				esito = "KO"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", r.Target.Name, r.Target.Type, esito,
				r.Duration.Round(time.Millisecond), r.Error)
		}
		if err := tw.Flush(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("formato %q sconosciuto", *formato)
	}

	for _, r := range res {
		if !r.OK {
			return errFalliti
		}
	}
	return nil
}

func runServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	path := fs.String("config", "probe.yaml", "file di configurazione")
	addr := fs.String("addr", ":9101", "indirizzo su cui ascoltare")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := config.Load(*path)
	if err != nil {
		return err
	}

	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	s := server.New(cfg, log)

	// Il contesto viene annullato con Ctrl+C (SIGINT) o con SIGTERM, il segnale
	// che Docker e Kubernetes inviano per fermare un container.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go s.Loop(ctx)

	httpSrv := &http.Server{Addr: *addr, Handler: s.Handler(), ReadHeaderTimeout: 5 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- httpSrv.ListenAndServe() }()
	log.Info("probe avviato", "addr", *addr, "versione", version, "target", len(cfg.Targets))

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	log.Info("arresto in corso")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return httpSrv.Shutdown(shutdownCtx) // termina le richieste in corso, poi chiude
}
