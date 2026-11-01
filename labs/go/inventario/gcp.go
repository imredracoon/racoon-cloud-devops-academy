package main

import (
	"context"
	"errors"
	"path"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// GCP usa i client di Compute Engine e Cloud Storage.
type GCP struct {
	progetto string
	istanze  *compute.InstancesClient
	storage  *storage.Client
}

// NuovoGCP usa le Application Default Credentials:
// gcloud auth application-default login sul PC, il service account nel cloud.
func NuovoGCP(ctx context.Context, progetto string) (*GCP, error) {
	if progetto == "" {
		return nil, errors.New("serve -project")
	}
	ic, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	sc, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCP{progetto: progetto, istanze: ic, storage: sc}, nil
}

func (g *GCP) Istanze(ctx context.Context) ([]Risorsa, error) {
	var out []Risorsa
	// AggregatedList restituisce le istanze di tutte le zone, raggruppate per zona.
	it := g.istanze.AggregatedList(ctx, &computepb.AggregatedListInstancesRequest{Project: g.progetto})
	for {
		coppia, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, i := range coppia.Value.GetInstances() {
			out = append(out, Risorsa{
				Tipo:  "istanza",
				Nome:  i.GetName(),
				Zona:  path.Base(i.GetZone()), // l'API restituisce un URL: tengo solo l'ultima parte
				Stato: i.GetStatus(),
			})
		}
	}
	return out, nil
}

func (g *GCP) Bucket(ctx context.Context) ([]Risorsa, error) {
	var out []Risorsa
	it := g.storage.Buckets(ctx, g.progetto)
	for {
		b, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, Risorsa{Tipo: "bucket", Nome: b.Name, Zona: b.Location, Stato: "-"})
	}
	return out, nil
}
