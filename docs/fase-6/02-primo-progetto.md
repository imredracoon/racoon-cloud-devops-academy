# 6.2 Provider e primo progetto

!!! obiettivi "Obiettivi"
    - Leggere e scrivere la sintassi di base di HCL: blocchi, argomenti, riferimenti.
    - Dichiarare i provider e le loro versioni.
    - Creare, osservare e distruggere le prime risorse, prima senza cloud e poi su AWS o Google Cloud.
    - Usare più provider e più configurazioni dello stesso provider.

## Installare Terraform

Nel Dev Container del corso Terraform è già installato: verifica con `terraform version`. Per installarlo su una tua
macchina usa il repository ufficiale di HashiCorp:

=== "Debian/Ubuntu"

    ```bash
    wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
      | sudo tee /etc/apt/sources.list.d/hashicorp.list
    sudo apt update && sudo apt install terraform
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y dnf-plugins-core
    sudo dnf config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo
    sudo dnf install -y terraform
    ```

Terraform è un singolo file eseguibile: in alternativa ai repository puoi scaricare lo zip dalla pagina ufficiale e copiarlo in `/usr/local/bin`.

## La sintassi di HCL

Un file `.tf` è fatto di **blocchi**. Ogni blocco ha un tipo, zero o più etichette e un corpo con **argomenti**:

```hcl
# tipo     etichetta 1     etichetta 2
resource "aws_s3_bucket" "dati" {
  bucket = "academy-dati-1234"   # argomento = espressione
}
```

- Un blocco `resource` ha due etichette: il **tipo di risorsa** (`aws_s3_bucket`, il prefisso `aws_` indica il
  provider) e il **nome locale** (`dati`), che usi per riferirti a quella risorsa nel resto della configurazione.
- L'**indirizzo** della risorsa è `aws_s3_bucket.dati`; i suoi attributi si leggono con `aws_s3_bucket.dati.arn`.
- Commenti con `#` (o `//`, e `/* ... */` su più righe).
- Terraform legge **tutti** i file `.tf` della cartella come un'unica configurazione: la divisione in `main.tf`,
  `variables.tf`, `outputs.tf`, `versions.tf` è una convenzione per la leggibilità.

## Dichiarare i provider

```hcl
terraform {
  required_version = ">= 1.12"          # versione minima della CLI

  required_providers {
    aws = {
      source  = "hashicorp/aws"         # indirizzo nel registry: registry.terraform.io/hashicorp/aws
      version = "~> 6.0"                # vincolo di versione
    }
  }
}

provider "aws" {
  region = "eu-south-1"                 # configurazione del provider
}
```

I vincoli di versione:

| Vincolo | Significato |
|---|---|
| `= 6.10.0` | esattamente questa |
| `>= 6.0` | questa o successive |
| `~> 6.0` | dalla 6.0 in su, ma sotto la 7.0 (cambia solo l'ultima cifra indicata) |
| `~> 6.10.0` | dalla 6.10.0 in su, ma sotto la 6.11 |
| `>= 6.0, < 6.20` | combinazione |

`terraform init` scarica i provider nella cartella `.terraform/` e registra le versioni scelte nel file
`.terraform.lock.hcl` ([lezione 6.3](03-workflow.md)).

## Laboratorio 1: senza cloud

Il primo progetto usa due provider che non creano nulla nel cloud: `random` (genera valori casuali) e `local`
(scrive file sul tuo disco). Puoi sperimentare senza costi e senza credenziali.

!!! laboratorio "Laboratorio · Il primo apply"
    Apri la cartella `labs/terraform/01-primo-progetto` del corso. Contiene un solo file, `main.tf`:

    ```hcl
    terraform {
      required_version = ">= 1.12"

      required_providers {
        local = {
          source  = "hashicorp/local"
          version = "~> 2.9"
        }
        random = {
          source  = "hashicorp/random"
          version = "~> 3.9"
        }
      }
    }

    resource "random_pet" "nome" {
      length = 2
    }

    resource "local_file" "benvenuto" {
      filename = "${path.module}/out/benvenuto.txt"
      content  = "Ciao! Questo server si chiama ${random_pet.nome.id}.\n"
    }

    output "nome_server" {
      value = random_pet.nome.id
    }
    ```

    1. `terraform init`: osserva i provider scaricati e la cartella `.terraform/`.
    2. `terraform plan`: legge la configurazione e mostra cosa farebbe. Cerca `Plan: 2 to add`.
    3. `terraform apply`: rimostra il piano e chiede conferma; rispondi `yes`.
    4. `cat out/benvenuto.txt` e `terraform output`.
    5. `terraform apply` di nuovo: `No changes`. È l'**idempotenza**.
    6. Modifica a mano il file `out/benvenuto.txt` e rilancia `terraform plan`. Cosa propone Terraform?
    7. Cambia `length = 2` in `length = 3` e lancia `terraform plan`. Cerca `forces replacement`.
    8. `terraform destroy`.

??? success "Spiegazione"
    - Al punto 2 la risorsa `local_file` usa `random_pet.nome.id`: Terraform capisce da solo che deve creare prima il
      nome e poi il file (**dipendenza implicita**, [lezione 6.4](04-risorse.md)). Il valore è `(known after apply)` perché si conosce
      solo dopo la creazione.
    - Al punto 6 Terraform rileva che il file reale non corrisponde più alla configurazione (**drift**) e propone di
      ricrearlo: la configurazione è la fonte di verità.
    - Al punto 7 `random_pet` non può cambiare lunghezza sul posto: va **sostituito** (`-/+`), e di conseguenza
      cambia anche il contenuto del file. Nel cloud lo stesso accade, per esempio, cambiando l'immagine di una VM.
    - Dopo `apply` compare il file `terraform.tfstate`: è lo **state**, che studierai nella [lezione 6.6](06-state.md). Non va mai
      modificato a mano né messo in Git.

## Autenticarsi nel cloud

I provider leggono le credenziali dagli stessi posti delle CLI che hai configurato nella [lezione 4.10](../fase-4/10-cli.md). **Non**
scrivere mai credenziali nei file `.tf`.

=== "AWS"

    ```bash
    aws sso login --profile academy        # accesso con IAM Identity Center
    export AWS_PROFILE=academy             # il provider userà questo profilo
    aws sts get-caller-identity            # verifica
    ```

=== "Google Cloud"

    ```bash
    gcloud auth application-default login  # credenziali per le librerie e per Terraform (ADC)
    gcloud config set project ID-PROGETTO
    ```

    Le *Application Default Credentials* sono diverse dal login della CLI (`gcloud auth login`): Terraform usa le
    prime.

## Laboratorio 2: il primo bucket

!!! laboratorio "Laboratorio · Un bucket gestito da Terraform"
    Usa la cartella `labs/terraform/aws/bucket` oppure `labs/terraform/gcp/bucket`. Leggi il file `main.tf`: crea un
    suffisso casuale con `random_id` (i nomi dei bucket sono unici al mondo) e un bucket privato con versioning.

    === "AWS"

        ```bash
        cd labs/terraform/aws/bucket
        terraform init
        terraform apply
        aws s3api get-bucket-tagging --bucket "$(terraform output -raw nome_bucket)"
        ```

        Le etichette `progetto` e `gestito` non sono scritte sul bucket: arrivano dal blocco `default_tags` del
        provider.

    === "Google Cloud"

        ```bash
        cd labs/terraform/gcp/bucket
        terraform init
        terraform apply -var="progetto=ID-PROGETTO"
        gcloud storage buckets describe "gs://$(terraform output -raw nome_bucket)" --format="value(labels)"
        ```

        Le etichette arrivano dal blocco `default_labels` del provider.

    Ora prova a cancellare il bucket dalla **console** e rilancia `terraform plan`: Terraform si accorge che manca e
    propone di ricrearlo. Infine `terraform destroy`.

## Più provider nella stessa configurazione

Una configurazione può usare provider diversi (come `random` e `aws` nel laboratorio) e anche **più configurazioni
dello stesso provider**, distinte da un **alias**. Il caso tipico è lavorare in due regioni:

```hcl
provider "aws" {
  region = "eu-south-1"            # configurazione predefinita
}

provider "aws" {
  alias  = "virginia"
  region = "us-east-1"
}

# I certificati per CloudFront devono stare in us-east-1 (lezione 5.7).
resource "aws_acm_certificate" "cdn" {
  provider          = aws.virginia
  domain_name       = "www.esempio.it"
  validation_method = "DNS"
}
```

Su Google Cloud si fa lo stesso con `provider "google" { alias = "..." project = "altro-progetto" }`.

## Esercizi

!!! esercizio "Esercizio · Leggi la configurazione"
    Data la risorsa `resource "google_compute_network" "principale" { ... }`:

    1. Qual è il suo indirizzo?
    2. Come scrivi un riferimento al suo attributo `id` in un'altra risorsa?
    3. Quale provider la gestisce?
    4. Il vincolo `version = "~> 8.3"` accetta la versione 8.9.0? E la 9.0.0?

??? success "Soluzione"
    1. `google_compute_network.principale`.
    2. `network = google_compute_network.principale.id`.
    3. Il provider `google` (prefisso del tipo di risorsa), cioè `hashicorp/google`.
    4. Sì la 8.9.0 (`~> 8.3` significa `>= 8.3, < 9.0`); no la 9.0.0.

!!! esercizio "Esercizio · Due regioni"
    Aggiungi alla configurazione `aws/bucket` un secondo bucket, `aws_s3_bucket.copia`, nella regione `eu-central-1`,
    usando un provider con alias `francoforte`.

??? success "Soluzione"
    ```hcl
    provider "aws" {
      alias  = "francoforte"
      region = "eu-central-1"
    }

    resource "aws_s3_bucket" "copia" {
      provider = aws.francoforte
      bucket   = "academy-copia-${random_id.suffisso.hex}"
    }
    ```

    `default_tags` va ripetuto nel nuovo blocco `provider` se vuoi le stesse etichette: ogni configurazione del
    provider è indipendente.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale comando scarica i provider indicati nella configurazione?","options":["terraform plan","terraform init","terraform get-providers","terraform apply"],"answer":[1],"explanation":"init prepara la cartella di lavoro: provider, moduli e backend."},
 {"q":"Cosa significa version = \"~> 6.0\" per un provider?","options":["Esattamente la 6.0","Qualunque versione dalla 6.0 in su, ma inferiore alla 7.0","Qualunque versione","Solo versioni inferiori alla 6.0"],"answer":[1],"explanation":"L'operatore ~> lascia variare solo l'ultima cifra indicata."},
 {"q":"Come usi una seconda configurazione dello stesso provider per un'altra regione?","options":["Con un secondo blocco provider con alias, e l'argomento provider nella risorsa","Con due cartelle separate, obbligatoriamente","Non è possibile","Con una variabile d'ambiente"],"answer":[0],"explanation":"Per esempio provider = aws.virginia."},
 {"q":"Dove devono stare le credenziali AWS usate da Terraform?","options":["Nel blocco provider, in chiaro","Fuori dalla configurazione: profilo della CLI, variabili d'ambiente, ruoli","In un output","Nel file .terraform.lock.hcl"],"answer":[1],"explanation":"Mai credenziali nei file .tf versionati."},
 {"q":"In resource \"aws_instance\" \"web\" { ... }, cos'è \"web\"?","options":["Il nome dell'istanza su AWS","Il nome locale, usato per riferirsi alla risorsa nella configurazione","Il tipo di risorsa","Il provider"],"answer":[1],"explanation":"Il nome visibile su AWS si imposta con le etichette (tag Name)."}
]}
</script>

## Fonti

- [Install Terraform](https://developer.hashicorp.com/terraform/install)
- [Terraform language: Configuration syntax](https://developer.hashicorp.com/terraform/language/syntax/configuration)
- [Terraform: Provider requirements](https://developer.hashicorp.com/terraform/language/providers/requirements) e [Version constraints](https://developer.hashicorp.com/terraform/language/expressions/version-constraints)
- [Terraform: Provider configuration and alias](https://developer.hashicorp.com/terraform/language/providers/configuration)
- [AWS provider: Authentication and configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration)
- [Google provider: Authentication](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#authentication)
