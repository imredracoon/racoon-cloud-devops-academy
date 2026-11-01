// inventario elenca istanze e bucket di un account AWS o di un progetto Google Cloud
// (lezione G.7 del percorso Go). Usa le credenziali già configurate per le CLI.
//
//	go run . -cloud aws -region eu-south-1
//	go run . -cloud gcp -project mio-progetto
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

// Risorsa è una riga dell'inventario, uguale per tutti i cloud.
type Risorsa struct {
	Tipo  string // "istanza" o "bucket"
	Nome  string
	Zona  string
	Stato string
}

// Fornitore è ciò che ogni cloud deve saper fare. main non sa quale ha davanti.
type Fornitore interface {
	Istanze(ctx context.Context) ([]Risorsa, error)
	Bucket(ctx context.Context) ([]Risorsa, error)
}

func main() {
	cloud := flag.String("cloud", "aws", "aws o gcp")
	region := flag.String("region", "eu-south-1", "regione AWS")
	project := flag.String("project", "", "progetto Google Cloud")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var f Fornitore
	var err error
	switch *cloud {
	case "aws":
		f, err = NuovoAWS(ctx, *region)
	case "gcp":
		f, err = NuovoGCP(ctx, *project)
	default:
		err = fmt.Errorf("cloud %q sconosciuto", *cloud)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}

	risorse, err := Inventario(ctx, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "TIPO\tNOME\tZONA\tSTATO")
	for _, r := range risorse {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Tipo, r.Nome, r.Zona, r.Stato)
	}
	if err := tw.Flush(); err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}
}

// Inventario raccoglie istanze e bucket da un qualsiasi Fornitore.
func Inventario(ctx context.Context, f Fornitore) ([]Risorsa, error) {
	istanze, err := f.Istanze(ctx)
	if err != nil {
		return nil, fmt.Errorf("istanze: %w", err)
	}
	bucket, err := f.Bucket(ctx)
	if err != nil {
		return nil, fmt.Errorf("bucket: %w", err)
	}
	return append(istanze, bucket...), nil
}
