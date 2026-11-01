# G.7 Gli SDK di AWS e Google Cloud

!!! obiettivi "Obiettivi"
    - Capire cos'è un SDK e come trova le credenziali senza che tu le scriva nel codice.
    - Usare l'AWS SDK for Go v2 per leggere istanze EC2 e bucket S3.
    - Usare le librerie client di Google Cloud per leggere istanze Compute Engine e bucket Cloud Storage.
    - Gestire la paginazione delle API.
    - Nascondere i cloud dietro un'interfaccia per poter testare il codice senza un account.

!!! info "Quando studiarla"
    Dopo la [lezione 4.10](../fase-4/10-cli.md), quando hai configurato le CLI `aws` e `gcloud`, oppure durante la
    [Fase 5](../fase-5/index.md). Per **scrivere e testare** il codice non serve un account cloud; per eseguirlo
    davvero sì: segui prima le regole di [Account cloud in sicurezza](../ambiente/account-cloud.md).

## Dalla CLI al codice

La CLI `aws ec2 describe-instances` e la console web fanno la stessa cosa: chiamano le **API** del cloud, cioè
richieste HTTPS firmate. Un **SDK** (*Software Development Kit*) è una libreria che fa queste chiamate per te, nel
tuo linguaggio: costruisce la richiesta, la firma con le credenziali, ripete in caso di errori temporanei e ti
restituisce struct già pronte.

Scriverai un tool `inventario` che elenca istanze e bucket, con la stessa interfaccia per i due cloud. La versione
di riferimento è in
[`labs/go/inventario`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/inventario).

## Le credenziali: mai nel codice

Gli SDK trovano le credenziali da soli, seguendo una **catena** di posti in cui cercarle. È la stessa catena delle
CLI, quindi se `aws sts get-caller-identity` o `gcloud auth list` funzionano, funziona anche il tuo programma.

=== "AWS"

    In ordine: variabili d'ambiente (`AWS_ACCESS_KEY_ID`...), file `~/.aws/config` e `~/.aws/credentials` (incluso
    l'accesso con IAM Identity Center, `aws sso login`), e infine il **ruolo** dell'istanza EC2, del task ECS o
    del Pod. Nel cloud non si usano chiavi: si usano i ruoli ([lezione 4.4](../fase-4/04-iam.md)).

=== "Google Cloud"

    Le **Application Default Credentials** (ADC): la variabile `GOOGLE_APPLICATION_CREDENTIALS`, poi le
    credenziali create con `gcloud auth application-default login`, e infine il **service account** collegato alla
    VM, a Cloud Run o al Pod GKE. Anche qui, nel cloud niente file di chiavi.

!!! warning "Sola lettura"
    Per questa lezione usa un utente o un ruolo con permessi di **sola lettura** (ad esempio la policy gestita
    `ReadOnlyAccess` su AWS, o il ruolo `roles/viewer` su Google Cloud). Un programma che sbaglia un ciclo può
    fare molte chiamate in pochi secondi: meglio che non possa creare nulla.

## AWS SDK for Go v2

```bash
go get github.com/aws/aws-sdk-go-v2/config github.com/aws/aws-sdk-go-v2/service/ec2 github.com/aws/aws-sdk-go-v2/service/s3
```

Ogni servizio AWS ha il suo modulo (`service/ec2`, `service/s3`, ...): scarichi solo quelli che usi.

```go
cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-south-1"))   // (1)!
if err != nil {
	return err
}
client := ec2.NewFromConfig(cfg)

p := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})   // (2)!
for p.HasMorePages() {
	pagina, err := p.NextPage(ctx)
	if err != nil {
		return err
	}
	for _, res := range pagina.Reservations {
		for _, i := range res.Instances {
			fmt.Println(aws.ToString(i.InstanceId), i.State.Name)   // (3)!
		}
	}
}
```

1.  Carica regione e credenziali con la catena descritta sopra.
2.  Le API restituiscono i risultati **a pagine** (ad esempio 1000 alla volta). Il *paginatore* chiede le pagine
    successive per te: senza, vedresti solo la prima.
3.  Molti campi delle struct dell'SDK sono **puntatori** (`*string`), perché l'API distingue "vuoto" da "assente".
    `aws.ToString` restituisce la stringa, o `""` se il puntatore è `nil`, senza rischio di panic.

Gli errori delle API si riconoscono con `errors.As`:

```go
var apiErr smithy.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.ErrorCode())   // ad esempio "UnauthorizedOperation"
}
```

## Librerie client di Google Cloud

```bash
go get cloud.google.com/go/compute cloud.google.com/go/storage google.golang.org/api
```

```go
client, err := storage.NewClient(ctx)          // usa le Application Default Credentials
if err != nil {
	return err
}
defer client.Close()

it := client.Buckets(ctx, "mio-progetto")      // un iteratore: anche lui gestisce le pagine
for {
	b, err := it.Next()
	if errors.Is(err, iterator.Done) {         // fine dell'elenco: non è un errore vero
		break
	}
	if err != nil {
		return err
	}
	fmt.Println(b.Name, b.Location)
}
```

Le istanze Compute Engine sono **zonali**: `AggregatedList` le restituisce tutte, raggruppate per zona. Il campo
`Zone` è un URL completo (`.../zones/europe-west8-a`): `path.Base` ne prende l'ultima parte.

## Un'interfaccia per tutti i cloud

`main` non deve sapere quale cloud sta usando. Definisci il comportamento che ti serve con un'interfaccia
([lezione G.3](03-pacchetti-rete.md#le-interfacce)) e scrivi un tipo per ogni cloud:

```go
type Risorsa struct {
	Tipo, Nome, Zona, Stato string
}

type Fornitore interface {
	Istanze(ctx context.Context) ([]Risorsa, error)
	Bucket(ctx context.Context) ([]Risorsa, error)
}

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
```

Il vantaggio più grande è nei **test**: un fornitore finto soddisfa la stessa interfaccia e non fa nessuna chiamata
di rete. Così la logica si verifica nella pipeline, gratis e in millisecondi.

```go title="main_test.go"
type finto struct {
	istanze []Risorsa
	err     error
}

func (f finto) Istanze(context.Context) ([]Risorsa, error) { return f.istanze, f.err }
func (f finto) Bucket(context.Context) ([]Risorsa, error) {
	return []Risorsa{{Tipo: "bucket", Nome: "log"}}, nil
}

func TestInventarioErrore(t *testing.T) {
	errAPI := errors.New("AccessDenied")
	_, err := Inventario(context.Background(), finto{err: errAPI})
	if !errors.Is(err, errAPI) {
		t.Errorf("errore = %v, voglio che avvolga %v", err, errAPI)
	}
}
```

## Laboratorio

!!! laboratorio "Laboratorio · inventario"
    1. Crea il modulo `inventario` con i file `main.go` (opzioni `-cloud`, `-region`, `-project`, stampa a tabella),
       `aws.go` e `gcp.go`, ognuno con un tipo che soddisfa `Fornitore`.
    2. Scrivi i test con un fornitore finto ed eseguili: devono passare **senza** credenziali.
    3. Se hai già un account cloud configurato, esegui `go run . -cloud aws` e `go run . -cloud gcp -project
       TUO_PROGETTO` e confronta il risultato con `aws ec2 describe-instances` e `gcloud compute instances list`.
    4. Prova a eseguire il tool con un profilo senza permessi (`AWS_PROFILE=nessuno`): l'errore è comprensibile?

??? success "Soluzione guidata"
    La soluzione completa è in
    [`labs/go/inventario`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/inventario).
    ```console
    $ go test ./...
    ok  	inventario	0.019s
    $ go run . -cloud aws -region eu-south-1
    TIPO     NOME               ZONA           STATO
    istanza  web-a              eu-south-1a    running
    istanza  web-b              eu-south-1b    stopped
    bucket   racoon-sito-xyz    eu-south-1     -
    ```
    **Spiegazione.** `ListBuckets` di S3 è un'operazione **globale**: elenca i bucket di tutte le regioni, mentre
    `DescribeInstances` vede solo la regione configurata. Su Google Cloud, invece, sia istanze sia bucket si
    leggono per progetto. Sono le stesse differenze che hai visto con le CLI nella
    [lezione 4.10](../fase-4/10-cli.md). Gli errori arrivano arricchiti (`istanze: operation error EC2:
    DescribeInstances, ... UnauthorizedOperation`): chi legge capisce subito quale chiamata è fallita.

## Esercizi

!!! esercizio "Esercizio 1 · Solo le istanze accese"
    Aggiungi l'opzione `-stato running` che mostra solo le risorse con quello stato. Implementa il filtro **una
    volta sola**, in `main`, non nei due fornitori, e scrivi il test con il fornitore finto.

??? success "Soluzione"
    ```go
    func filtra(r []Risorsa, stato string) []Risorsa {
    	if stato == "" {
    		return r
    	}
    	var out []Risorsa
    	for _, x := range r {
    		if strings.EqualFold(x.Stato, stato) {   // AWS dice "running", Google "RUNNING"
    			out = append(out, x)
    		}
    	}
    	return out
    }
    ```
    **Spiegazione.** `EqualFold` confronta ignorando maiuscole e minuscole, così lo stesso filtro funziona con i due
    cloud. Mettere la logica sopra l'interfaccia evita di scriverla, e testarla, due volte. Per account con
    migliaia di istanze conviene invece filtrare **lato API** (`Filters` in EC2, `Filter` in Compute Engine), per
    scaricare meno dati.

!!! esercizio "Esercizio 2 · Due regioni insieme"
    Modifica il fornitore AWS perché accetti più regioni (`-region eu-south-1,eu-west-1`) e le interroghi **in
    parallelo**, con gli strumenti della [lezione G.5](05-concorrenza-http.md).

??? success "Soluzione"
    ```go
    func (a *AWSMulti) Istanze(ctx context.Context) ([]Risorsa, error) {
    	risultati := make([][]Risorsa, len(a.regioni))
    	errori := make([]error, len(a.regioni))
    	var wg sync.WaitGroup
    	for i, r := range a.regioni {
    		wg.Go(func() { risultati[i], errori[i] = r.Istanze(ctx) })   // r è un *AWS per regione
    	}
    	wg.Wait()
    	if err := errors.Join(errori...); err != nil {
    		return nil, err
    	}
    	return slices.Concat(risultati...), nil
    }
    ```
    **Spiegazione.** Ogni goroutine scrive in una posizione diversa delle slice, quindi non servono mutex.
    `errors.Join` ignora i `nil` e restituisce `nil` se tutte le regioni hanno risposto. Un client per regione è
    necessario perché la regione fa parte della configurazione del client.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Dove deve trovare le credenziali un programma Go che gira su un'istanza EC2?","options":["In una costante nel codice","In un file .env nel repository","Nel ruolo IAM dell'istanza, tramite la catena predefinita dell'SDK","In una variabile globale"],"answer":[2],"explanation":"Nel cloud si usano ruoli e service account, non chiavi."},
 {"q":"Perché si usa il paginatore con DescribeInstances?","options":["Per velocizzare la chiamata","Perché l'API restituisce i risultati a pagine e senza paginatore vedresti solo la prima","Per filtrare per stato","Per cambiare regione"],"answer":[1],"explanation":"Il paginatore chiede le pagine successive finché ce ne sono."},
 {"q":"A cosa serve aws.ToString(i.InstanceId)?","options":["A convertire un numero","A leggere un *string restituendo \"\" se è nil, senza panic","A cifrare l'ID","A tradurre il testo"],"answer":[1],"explanation":"Molti campi dell'SDK sono puntatori."},
 {"q":"Come capisci che un iteratore delle librerie Google Cloud è finito?","options":["Restituisce nil, nil","Restituisce l'errore iterator.Done","Va in panic","Restituisce un elemento vuoto"],"answer":[1],"explanation":"iterator.Done non è un errore vero: segnala la fine."},
 {"q":"Qual è il vantaggio principale dell'interfaccia Fornitore?","options":["Rende il programma più veloce","Permette di testare la logica con un fornitore finto, senza account e senza rete","Riduce i costi del cloud","È obbligatoria negli SDK"],"answer":[1],"explanation":"La logica si verifica in pipeline, gratis e in millisecondi."}
]}
</script>

## Fonti

- [AWS SDK for Go v2: Developer Guide](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/welcome.html)
- [AWS SDK for Go v2: Configure the SDK](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-gosdk.html) e [Using paginators](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/using.html)
- [Google Cloud Client Libraries for Go](https://cloud.google.com/go/docs/reference) e [Cloud Storage client library](https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go)
- [How Application Default Credentials works](https://cloud.google.com/docs/authentication/application-default-credentials)
- [pkg.go.dev: cloud.google.com/go/compute/apiv1](https://pkg.go.dev/cloud.google.com/go/compute/apiv1)
