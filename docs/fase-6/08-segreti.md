# 6.8 Dati sensibili e Vault

!!! obiettivi "Obiettivi"
    - Capire dove possono finire i segreti in un progetto Terraform: codice, variabili, piani, state, log.
    - Usare `sensitive`, i valori **effimeri** e gli argomenti **write-only**.
    - Leggere i segreti da un gestore (AWS Secrets Manager, Google Secret Manager, HashiCorp Vault).
    - Conoscere i segreti dinamici di Vault.

## Dove si nascondono i segreti

| Luogo | Rischio | Contromisura |
|---|---|---|
| file `.tf` | password scritte nel codice finiscono in Git | mai segreti nel codice |
| file `.tfvars` | spesso contengono password | escluderli da Git, o usare variabili d'ambiente `TF_VAR_` o un gestore di segreti |
| output della CLI e log della pipeline | i valori compaiono a schermo | `sensitive = true` |
| **state** e **piani salvati** | contengono gli attributi delle risorse **in chiaro** | backend cifrato con accesso limitato; valori effimeri e argomenti write-only |
| credenziali dei provider | chiavi di lunga durata rubabili | credenziali temporanee: SSO in locale, OIDC nelle pipeline ([Fase 7](../fase-7/index.md)), ruoli |

## Valori sensitive

```hcl
variable "password_db" {
  type      = string
  sensitive = true
}

output "stringa_connessione" {
  value     = "postgresql://tasks:${var.password_db}@${aws_db_instance.tasks.address}/tasks"
  sensitive = true
}
```

Terraform mostra `(sensitive value)` al posto del valore in `plan`, `apply` e `output`; un output che usa un valore
sensitive deve essere dichiarato a sua volta `sensitive`. Ma attenzione: **`sensitive` non cifra nulla**. Il valore
resta leggibile nello state e con `terraform output -raw stringa_connessione` o `-json`. Molti attributi delle
risorse (come `password` di un database) sono già marcati sensitive dal provider.

## Valori effimeri e argomenti write-only

Da Terraform 1.10 esistono i valori **effimeri** (*ephemeral*): esistono solo durante l'esecuzione e **non vengono
mai salvati** né nello state né nei piani.

- Variabili e output dei moduli figli con `ephemeral = true`.
- **Risorse effimere** (blocco `ephemeral`): generano o leggono un valore al momento, per esempio una password casuale
  o un segreto letto da un gestore.
- Da Terraform 1.11 gli argomenti **write-only** (con suffisso `_wo`) accettano valori effimeri: il provider li invia
  all'API ma non li salva. Si accompagnano a un argomento `_wo_version` da incrementare quando vuoi cambiare il valore.

=== "AWS"

    ```hcl
    # Una password generata al momento: non finisce nello state.
    ephemeral "random_password" "db" {
      length           = 24
      override_special = "!#$%&*()-_=+[]{}<>:?"
    }

    resource "aws_secretsmanager_secret" "db" {
      name = "academy/db-password"
    }

    # La password è scritta nel segreto senza essere salvata nello state.
    resource "aws_secretsmanager_secret_version" "db" {
      secret_id                = aws_secretsmanager_secret.db.id
      secret_string_wo         = ephemeral.random_password.db.result
      secret_string_wo_version = 1
    }

    # Rilettura effimera del segreto, per passarlo al database.
    ephemeral "aws_secretsmanager_secret_version" "db" {
      secret_id = aws_secretsmanager_secret_version.db.secret_id
    }

    resource "aws_db_instance" "tasks" {
      identifier          = "academy-tasks"
      engine              = "postgres"
      instance_class      = "db.t4g.micro"
      allocated_storage   = 20
      username            = "tasks"
      password_wo         = ephemeral.aws_secretsmanager_secret_version.db.secret_string
      password_wo_version = 1
      skip_final_snapshot = true
    }
    ```

    In alternativa RDS può generare e ruotare da sé la password principale in Secrets Manager con
    `manage_master_user_password = true`: Terraform non la vede mai.

=== "Google Cloud"

    ```hcl
    ephemeral "random_password" "db" {
      length = 24
    }

    resource "google_secret_manager_secret" "db" {
      secret_id = "academy-db-password"
      replication {
        auto {}
      }
    }

    resource "google_secret_manager_secret_version" "db" {
      secret                 = google_secret_manager_secret.db.id
      secret_data_wo         = ephemeral.random_password.db.result
      secret_data_wo_version = 1
    }

    ephemeral "google_secret_manager_secret_version" "db" {
      secret = google_secret_manager_secret_version.db.secret
    }

    resource "google_sql_user" "tasks" {
      name                = "tasks"
      instance            = google_sql_database_instance.tasks.name
      password_wo         = ephemeral.google_secret_manager_secret_version.db.secret_data
      password_wo_version = 1
    }
    ```

Non tutti gli argomenti hanno una versione write-only: controlla la documentazione della risorsa nel registry. Dove
non esiste, il valore finirà nello state, che quindi va protetto in ogni caso.

## HashiCorp Vault

**Vault** è il gestore di segreti di HashiCorp, indipendente dal cloud: conserva segreti statici (chiave/valore),
genera **segreti dinamici** (credenziali create al momento con una scadenza: per un database, per AWS, per Google
Cloud), cifra dati per conto delle applicazioni e gestisce certificati. Si usa installato in proprio o come servizio
gestito (**HCP Vault Dedicated**).

Il provider `hashicorp/vault` permette a Terraform di **leggere** segreti da Vault e di **configurare** Vault stesso:

```hcl
provider "vault" {
  # indirizzo e token da VAULT_ADDR e VAULT_TOKEN, mai nel codice
}

# Lettura effimera di un segreto chiave/valore (motore KV versione 2)
ephemeral "vault_kv_secret_v2" "db" {
  mount = "secret"
  name  = "tasks/db"
}

# ephemeral.vault_kv_secret_v2.db.data["password"]
```

Con i **segreti dinamici** il vantaggio è ancora maggiore: Terraform (o l'applicazione) chiede a Vault credenziali AWS
valide per pochi minuti, invece di usare chiavi di lunga durata. Se vengono rubate, scadono da sole.

!!! note "Per l'esame"
    L'obiettivo 4h dell'esame Terraform Associate chiede di conoscere le buone pratiche per i dati sensibili, Vault
    compreso: sensitive non cifra, lo state contiene segreti, Vault fornisce segreti (anche dinamici) a Terraform, i
    valori effimeri non vengono salvati.

## Laboratorio

!!! laboratorio "Laboratorio · Cosa finisce nello state?"
    Senza cloud, nella cartella `labs/terraform/01-primo-progetto`:

    1. Aggiungi:

        ```hcl
        resource "random_password" "normale" {
          length = 16
        }

        ephemeral "random_password" "effimera" {
          length = 16
        }

        output "password" {
          value     = random_password.normale.result
          sensitive = true
        }
        ```

    2. `terraform init` e `terraform apply`. Che cosa mostra la CLI per l'output?
    3. `terraform output -raw password`: è leggibile?
    4. Cerca le password nello state: `grep -A3 '"result"' terraform.tfstate`. Trovi quella normale? E quella effimera?
    5. Prova ad aggiungere un output con `value = ephemeral.random_password.effimera.result`: cosa succede?
    6. Togli le aggiunte e lancia `terraform apply`.

??? success "Spiegazione"
    2. `password = <sensitive>`.
    3. Sì: `sensitive` nasconde il valore nella CLI, non lo protegge.
    4. La password della risorsa `random_password.normale` è nello state **in chiaro**; quella effimera non c'è.
    5. Errore: un valore effimero non può finire in un output del modulo root, perché verrebbe salvato.

## Esercizi

!!! esercizio "Esercizio · Revisione di sicurezza"
    In una pull request trovi:

    ```hcl
    provider "aws" {
      region     = "eu-south-1"
      access_key = "AKIA..."
      secret_key = "wJalr..."
    }

    variable "db_password" {
      default = "Password123!"
    }
    ```

    Elenca i problemi e la correzione.

??? success "Soluzione"
    1. **Credenziali nel codice**: vanno revocate subito (sono ormai nella cronologia di Git) e sostituite con un
       profilo SSO in locale o con OIDC nella pipeline.
    2. **Password come valore predefinito**: chiunque legge il repository la conosce. Va tolta, la variabile marcata
       `sensitive = true` (o meglio `ephemeral = true`) e il valore fornito da un gestore di segreti.
    3. Meglio ancora: generare la password con `ephemeral "random_password"` e scriverla con un argomento write-only,
       o lasciarla gestire al servizio (`manage_master_user_password`).
    4. Attiva la scansione dei segreti nel repository (GitHub secret scanning) per intercettare questi errori.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa fa sensitive = true su una variabile?","options":["Cifra il valore nello state","Nasconde il valore nell'output della CLI, ma resta in chiaro nello state","Impedisce di usare la variabile","Salva il valore in Vault"],"answer":[1],"explanation":"Per non salvarlo serve un valore effimero."},
 {"q":"Quale caratteristica di Terraform evita che un valore venga salvato nello state e nei piani?","options":["I valori effimeri (ephemeral) con gli argomenti write-only","sensitive = true","Il file .gitignore","terraform fmt"],"answer":[0],"explanation":"Disponibili da Terraform 1.10 e 1.11."},
 {"q":"Cosa sono i segreti dinamici di Vault?","options":["Segreti che cambiano nome","Credenziali generate al momento, con una durata limitata, revocabili","Segreti cifrati due volte","Variabili d'ambiente"],"answer":[1],"explanation":"Riducono il rischio delle credenziali di lunga durata."},
 {"q":"Quali sono buone pratiche per i dati sensibili in Terraform?","options":["Proteggere lo state con un backend cifrato e accessi limitati","Leggere i segreti da un gestore come Vault o Secrets Manager","Scrivere le password come default delle variabili","Versionare il file terraform.tfstate","Stampare i segreti nei log per il debug"],"answer":[0,1],"explanation":"Lo state contiene spesso segreti: va trattato come un segreto."},
 {"q":"Perché un argomento write-only (_wo) ha anche un argomento _wo_version?","options":["Per indicare la versione del provider","Perché Terraform non salva il valore e non può accorgersi che è cambiato: incrementando la versione forzi l'aggiornamento","Per cifrare il valore","Non serve"],"answer":[1],"explanation":"Il valore non è nello state, quindi non c'è nulla con cui confrontarlo."}
]}
</script>

## Fonti

- [Terraform: Sensitive data in state](https://developer.hashicorp.com/terraform/language/state/sensitive-data)
- [Terraform: Manage sensitive data](https://developer.hashicorp.com/terraform/language/manage-sensitive-data)
- [Terraform: Ephemeral values](https://developer.hashicorp.com/terraform/language/manage-sensitive-data/ephemeral) e [Write-only arguments](https://developer.hashicorp.com/terraform/language/manage-sensitive-data/write-only)
- [Terraform: Ephemeral resources](https://developer.hashicorp.com/terraform/language/resources/ephemeral)
- [AWS provider: aws_db_instance](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance) e [aws_secretsmanager_secret_version](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret_version)
- [Google provider: google_secret_manager_secret_version](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_version)
- [HashiCorp Vault: What is Vault?](https://developer.hashicorp.com/vault/docs/about-vault/what-is-vault) e [Vault provider for Terraform](https://registry.terraform.io/providers/hashicorp/vault/latest/docs)
