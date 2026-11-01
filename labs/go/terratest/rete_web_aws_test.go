package terratest

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestReteWebAWS crea davvero VPC e server web su AWS, li prova e li distrugge.
// Costa qualche centesimo e richiede credenziali: parte solo con TERRATEST_AWS=1.
//
//	TERRATEST_AWS=1 go test -v -timeout 30m -run TestReteWebAWS ./...
func TestReteWebAWS(t *testing.T) {
	if os.Getenv("TERRATEST_AWS") != "1" {
		t.Skip("imposta TERRATEST_AWS=1 per eseguire il test su AWS")
	}
	t.Parallel()
	ctx := t.Context()

	// Un nome unico evita conflitti se due test girano insieme nello stesso account.
	nome := "tt-" + strings.ToLower(random.UniqueID())

	opts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../terraform/aws/rete-web",
		Vars: map[string]any{
			"nome":          nome,
			"ambiente":      "test",
			"numero_server": 1,
		},
		NoColor: true,
	})
	defer terraform.DestroyContext(t, ctx, opts)
	terraform.InitAndApplyContext(t, ctx, opts)

	url := terraform.OutputListContext(t, ctx, opts, "indirizzi_web")[0]

	// Il server impiega qualche minuto ad avviarsi: riprova per 5 minuti.
	http_helper.HTTPGetWithRetryWithCustomValidationContext(t, ctx, url, nil, 30, 10*time.Second,
		func(stato int, corpo string) bool {
			return stato == 200 && strings.Contains(corpo, nome)
		})
	fmt.Println("server web raggiungibile su", url)
}
