// Test di integrazione del primo progetto Terraform (lezione G.7 del percorso Go).
// Esegue davvero terraform init, apply e destroy: serve il binario terraform nel PATH.
//
//	go test -v -timeout 10m ./...
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
	opts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: dir,
		NoColor:      true,
	})

	// defer: il destroy avviene anche se un controllo più sotto fallisce.
	defer terraform.DestroyContext(t, ctx, opts)
	terraform.InitAndApplyContext(t, ctx, opts)

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

	// Un secondo plan non deve trovare differenze: la configurazione è idempotente.
	if code := terraform.PlanExitCodeContext(t, ctx, opts); code != 0 {
		t.Errorf("plan dopo apply ha exit code %d, voglio 0 (nessuna modifica)", code)
	}
}
