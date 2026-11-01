package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWS usa due client del SDK: EC2 per le istanze, S3 per i bucket.
type AWS struct {
	ec2 *ec2.Client
	s3  *s3.Client
}

// NuovoAWS carica le credenziali con la stessa catena della AWS CLI:
// variabili d'ambiente, ~/.aws/config e credentials (anche SSO), ruolo dell'istanza.
func NuovoAWS(ctx context.Context, regione string) (*AWS, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(regione))
	if err != nil {
		return nil, err
	}
	return &AWS{ec2: ec2.NewFromConfig(cfg), s3: s3.NewFromConfig(cfg)}, nil
}

func (a *AWS) Istanze(ctx context.Context) ([]Risorsa, error) {
	var out []Risorsa
	// Il paginatore chiede le pagine successive finché ce ne sono.
	p := ec2.NewDescribeInstancesPaginator(a.ec2, &ec2.DescribeInstancesInput{})
	for p.HasMorePages() {
		pagina, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, res := range pagina.Reservations {
			for _, i := range res.Instances {
				nome := aws.ToString(i.InstanceId)
				for _, t := range i.Tags {
					if aws.ToString(t.Key) == "Name" {
						nome = aws.ToString(t.Value)
					}
				}
				out = append(out, Risorsa{
					Tipo:  "istanza",
					Nome:  nome,
					Zona:  aws.ToString(i.Placement.AvailabilityZone),
					Stato: string(i.State.Name),
				})
			}
		}
	}
	return out, nil
}

func (a *AWS) Bucket(ctx context.Context) ([]Risorsa, error) {
	var out []Risorsa
	p := s3.NewListBucketsPaginator(a.s3, &s3.ListBucketsInput{})
	for p.HasMorePages() {
		pagina, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, b := range pagina.Buckets {
			out = append(out, Risorsa{Tipo: "bucket", Nome: aws.ToString(b.Name), Zona: aws.ToString(b.BucketRegion), Stato: "-"})
		}
	}
	return out, nil
}
