# G.8 Testare Terraform con Terratest

!!! obiettivi "Obiettivi"
    - Capire la differenza tra i test nativi di Terraform e i test di integrazione con Terratest.
    - Scrivere un test Go che esegue `init`, `apply`, verifica il risultato e fa sempre `destroy`.
    - Leggere gli output di Terraform dal test e verificare l'infrastruttura dall'esterno, ad esempio con HTTP.
    - Tenere sotto controllo costi e conflitti quando i test creano risorse vere.

!!! info "Quando studiarla"
    Nella [Fase 6](../fase-6/index.md), dopo la [lezione 6.7](../fase-6/07-moduli.md). Il primo test non usa il cloud;
    il secondo crea risorse vere su AWS ed è facoltativo.

## Terraform è scritto in Go

Terraform, i suoi provider e molti strumenti attorno sono programmi Go. Non è un dettaglio: significa che puoi
leggere il codice di un provider quando la documentazione non basta, e che l'ecosistema di test più maturo per
Terraform è in Go.

## Due modi di testare

Nella [lezione 6.5](../fase-6/05-variabili.md) hai usato `terraform test` con i provider **simulati**
(`mock_provider`): veloce, gratuito e senza credenziali. Verifica che la configurazione **calcoli** le cose giuste.
Non può verificare che il server web risponda davvero, che il security group lasci passare la porta 80 o che lo
script di avvio funzioni.

| | `terraform test` con mock | Terratest |
|---|---|---|
| **Linguaggio** | HCL | Go |
| **Risorse vere** | no (con `mock_provider`) | sì: `apply` reale, poi `destroy` |
| **Cosa verifica** | valori calcolati, validazioni, piani | il comportamento dell'infrastruttura: HTTP, SSH, API del cloud |
| **Tempo e costo** | secondi, gratis | minuti, il costo delle risorse create |
| **Quando** | a ogni commit | prima di un rilascio del modulo, di notte, su richiesta |

Si usano **insieme**: molti test nativi veloci, pochi test Terratest che danno la certezza finale. È la cosiddetta
*piramide dei test*: tanti test rapidi alla base, pochi test lenti e realistici in cima.

## Il primo test, senza cloud

Il progetto della [lezione 6.2](../fase-6/02-primo-progetto.md) usa solo i provider `random` e `local`: un test
perfetto per iniziare. Il codice è in
[`labs/go/terratest`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/terratest).

```bash
mkdir -p ~/go-lab/terratest && cd ~/go-lab/terratest
go mod init terratest
go get github.com/gruntwork-io/terratest/modules/terraform
```

```go title="primo_progetto_test.go"
package terratest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestPrimoProgetto(t *testing.T) {
	t.Parallel()
	ctx := t.Context() // annullato alla fine del test

	dir := "../../terraform/01-primo-progetto"
	opts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{   // (1)!
		TerraformDir: dir,
		NoColor:      true,
	})

	defer terraform.DestroyContext(t, ctx, opts)   // (2)!
	terraform.InitAndApplyContext(t, ctx, opts)    // (3)!

	nome := terraform.OutputContext(t, ctx, opts, "nome_server")
	if parti := strings.Split(nome, "-"); len(parti) != 2 {
		t.Errorf("nome_server = %q, voglio due parole separate da un trattino", nome)
	}

	contenuto, err := os.ReadFile(filepath.Join(dir, "out", "benvenuto.txt"))
	if err != nil {
		t.Fatalf("il file non è stato creato: %v", err)
	}
	if !strings.Contains(string(contenuto), nome) {
		t.Errorf("il file non contiene il nome del server:\n%s", contenuto)
	}

	if code := terraform.PlanExitCodeContext(t, ctx, opts); code != 0 {   // (4)!
		t.Errorf("plan dopo apply ha exit code %d, voglio 0 (nessuna modifica)", code)
	}
}
```

1.  `WithDefaultRetryableErrors` ripete in automatico i comandi che falliscono per errori temporanei noti, come un
    download del provider interrotto.
2.  Il `defer` della [lezione G.2](02-collezioni-errori.md#leggere-file-e-defer) è fondamentale: il `destroy`
    avviene **anche se un controllo più sotto fallisce**. Senza, un test fallito lascerebbe risorse accese, e nel
    cloud risorse accese significa una fattura.
3.  Esegue `terraform init` e `terraform apply -auto-approve`. Se fallisce, il test si ferma (le funzioni di
    Terratest chiamano `t.Fatal` da sole); il `destroy` parte comunque.
4.  Un secondo `plan` subito dopo l'`apply` deve dire "nessuna modifica": se no, la configurazione non è
    **idempotente** e ogni esecuzione cambierebbe qualcosa.

```console
$ go test -v -timeout 10m ./...
=== RUN   TestPrimoProgetto
TestPrimoProgetto 2026-11-01T10:00:00Z logger.go:67: Running command terraform with args [init -upgrade=false]
...
TestPrimoProgetto 2026-11-01T10:00:09Z logger.go:67: Running command terraform with args [destroy -auto-approve -input=false]
--- PASS: TestPrimoProgetto (11.42s)
PASS
```

!!! warning "Il timeout di go test"
    Per impostazione predefinita `go test` interrompe un test dopo **10 minuti**. Un test che crea un cluster o un
    database può durarne di più: imposta sempre `-timeout` (ad esempio `-timeout 30m`). Un test interrotto dal
    timeout **non esegue i defer**, e le risorse restano accese.

## Un test con risorse vere

Il modulo `rete-web` della [Fase 6](../fase-6/index.md) crea una VPC e dei server web. Il test lo applica, aspetta che
il server risponda con il suo nome e poi distrugge tutto.

```go title="rete_web_aws_test.go"
func TestReteWebAWS(t *testing.T) {
	if os.Getenv("TERRATEST_AWS") != "1" {
		t.Skip("imposta TERRATEST_AWS=1 per eseguire il test su AWS")   // (1)!
	}
	t.Parallel()
	ctx := t.Context()

	nome := "tt-" + strings.ToLower(random.UniqueID())                  // (2)!

	opts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../terraform/aws/rete-web",
		Vars: map[string]any{
			"nome":          nome,
			"ambiente":      "test",
			"numero_server": 1,                                         // (3)!
		},
		NoColor: true,
	})
	defer terraform.DestroyContext(t, ctx, opts)
	terraform.InitAndApplyContext(t, ctx, opts)

	url := terraform.OutputListContext(t, ctx, opts, "indirizzi_web")[0]

	http_helper.HTTPGetWithRetryWithCustomValidationContext(t, ctx, url, nil, 30, 10*time.Second,   // (4)!
		func(stato int, corpo string) bool {
			return stato == 200 && strings.Contains(corpo, nome)
		})
}
```

1.  Senza la variabile, il test viene **saltato**: chi lancia `go test ./...` senza credenziali non crea nulla per
    sbaglio.
2.  Un nome unico per ogni esecuzione: due persone, o due pipeline, possono eseguire il test nello stesso account
    senza scontrarsi.
3.  Il minimo indispensabile per verificare il comportamento: un solo server, non due.
4.  Il server impiega un paio di minuti ad avviarsi e a eseguire lo script di cloud-init: il test riprova fino a
    30 volte a 10 secondi di distanza.

```bash
export AWS_PROFILE=sandbox            # un account o un profilo dedicato ai test
TERRATEST_AWS=1 go test -v -timeout 30m -run TestReteWebAWS ./...
```

!!! tip "Dove far girare questi test"
    In azienda i test Terratest girano in un **account separato** dedicato ai test, con un budget e un allarme di
    spesa ([lezione 4.9](../fase-4/09-costi.md)), e spesso con uno strumento che ogni notte cancella le risorse
    rimaste orfane. Mai nell'account di produzione.

## Laboratorio

!!! laboratorio "Laboratorio · Test di un tuo modulo"
    1. Esegui il test `TestPrimoProgetto` nel Dev Container e verifica che, alla fine, la cartella `out/` del
       progetto sia vuota o assente.
    2. Rompi apposta la configurazione (ad esempio cambia `length = 2` in `length = 3` nel `random_pet`): il test
       deve fallire con un messaggio chiaro, e il `destroy` deve partire comunque.
    3. Scrivi un test Terratest per uno dei moduli del tuo [Progetto 6](../fase-6/progetto.md). Se crea risorse vere,
       proteggilo con una variabile d'ambiente come `TestReteWebAWS`.

??? success "Soluzione guidata"
    Al punto 2 l'output contiene `nome_server = "allowed-sincere-cat", voglio due parole separate da un trattino` e,
    subito dopo, `Running command terraform with args [destroy ...]`. Per il punto 3, un modulo Google Cloud si
    testa allo stesso modo: Terratest ha anche il pacchetto `modules/gcp` per leggere le risorse create
    (ad esempio `gcp.FetchInstance(t, progetto, nome).GetPublicIP(t)`), e `modules/aws` fa lo stesso per AWS. Usa queste funzioni quando
    vuoi verificare una proprietà che non esce dagli `output` del modulo.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa può verificare Terratest che terraform test con mock_provider non può verificare?","options":["La sintassi HCL","Le validazioni delle variabili","Che il server creato risponda davvero via HTTP","I valori calcolati dei locals"],"answer":[2],"explanation":"Terratest crea risorse vere e le prova dall'esterno."},
 {"q":"Perché `defer terraform.DestroyContext(t, ctx, opts)` si scrive prima di InitAndApply?","options":["Per velocità","Perché il destroy deve avvenire anche se l'apply o un controllo successivo falliscono","Perché Terratest lo richiede per sintassi","Per eseguire il destroy prima dell'apply"],"answer":[1],"explanation":"Altrimenti un test fallito lascerebbe risorse accese."},
 {"q":"Cosa succede ai defer se go test interrompe un test per il timeout?","options":["Vengono eseguiti","Non vengono eseguiti: le risorse restano accese","Vengono eseguiti al test successivo","Terraform li esegue da solo"],"answer":[1],"explanation":"Per questo si imposta sempre un -timeout generoso."},
 {"q":"Cosa indica un PlanExitCode uguale a 2 subito dopo un apply?","options":["Tutto bene","Il plan vuole ancora modificare qualcosa: la configurazione non è idempotente","Errore di sintassi","Credenziali mancanti"],"answer":[1],"explanation":"Con -detailed-exitcode: 0 nessuna modifica, 1 errore, 2 modifiche."},
 {"q":"Perché il test AWS usa random.UniqueID() nel nome delle risorse?","options":["Per estetica","Per evitare conflitti tra esecuzioni parallele nello stesso account","Perché AWS lo richiede","Per ridurre i costi"],"answer":[1],"explanation":"Due pipeline possono girare insieme senza scontrarsi."}
]}
</script>

## Fonti

- [Terratest: documentazione](https://terratest.gruntwork.io/docs/) e [Quick start](https://terratest.gruntwork.io/docs/getting-started/quick-start/)
- [pkg.go.dev: terratest/modules/terraform](https://pkg.go.dev/github.com/gruntwork-io/terratest/modules/terraform)
- [Terraform: Tests](https://developer.hashicorp.com/terraform/language/tests) e [Mocks](https://developer.hashicorp.com/terraform/language/tests/mocking)
- [HashiCorp: Testing HashiCorp Terraform](https://www.hashicorp.com/blog/testing-hashicorp-terraform)
- [Command go: Testing flags](https://pkg.go.dev/cmd/go#hdr-Testing_flags), per `-timeout`
