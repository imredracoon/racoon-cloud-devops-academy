# Laboratori Terraform

Configurazioni usate nella [Fase 6](../../docs/fase-6/index.md) del corso.

| Cartella | Lezione | Costo |
|---|---|---|
| `01-primo-progetto/` | 6.2: provider `local` e `random`, tutto sul tuo computer | nessuno |
| `aws/bucket/`, `gcp/bucket/` | 6.2: il primo bucket nel cloud | trascurabile |
| `aws/rete-web/`, `gcp/rete-web/` | 6.4 e 6.5: rete, firewall e server web | qualche centesimo l'ora: distruggi alla fine |

Ogni configurazione si usa così:

```bash
cd labs/terraform/aws/rete-web
terraform init
terraform plan
terraform apply
terraform destroy
```

Le cartelle `rete-web` contengono anche **test senza cloud** (`tests/*.tftest.hcl`) che simulano il provider:
`terraform test` funziona anche senza credenziali.
