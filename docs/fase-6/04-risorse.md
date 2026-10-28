# 6.4 Risorse, data source e dipendenze

!!! obiettivi "Obiettivi"
    - Distinguere i blocchi `resource` e `data`.
    - Collegare le risorse con i riferimenti ai loro attributi.
    - Capire dipendenze implicite ed esplicite e il grafo delle risorse.
    - Usare i meta-argomenti `lifecycle` per proteggere e controllare le risorse.
    - Costruire una rete completa con server web, su AWS o Google Cloud.

## Resource e data

| Blocco | Cosa fa | Esempio |
|---|---|---|
| `resource` | **crea e gestisce** un oggetto: Terraform ne è responsabile per tutto il ciclo di vita | `resource "aws_vpc" "principale"` |
| `data` | **legge** informazioni su un oggetto esistente, gestito da altri o da AWS/Google | `data "aws_availability_zones" "disponibili"` |

Un **data source** si usa per non scrivere valori fissi che cambiano nel tempo o tra regioni: l'immagine più recente
di un sistema operativo, le zone disponibili, l'ID di una rete creata da un altro team.

=== "AWS"

    ```hcl
    data "aws_availability_zones" "disponibili" {
      state = "available"
    }

    # L'immagine Amazon Linux 2023 più recente, pubblicata da AWS come parametro pubblico.
    data "aws_ssm_parameter" "al2023" {
      name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
    }
    ```

    Si usano con il prefisso `data.`: `data.aws_availability_zones.disponibili.names`.

=== "Google Cloud"

    ```hcl
    data "google_compute_zones" "disponibili" {
      region = "europe-west8"
    }

    data "google_compute_image" "debian" {
      family  = "debian-12"
      project = "debian-cloud"
    }
    ```

    Si usano con il prefisso `data.`: `data.google_compute_image.debian.self_link`.

## Riferimenti e dipendenze

Quando un argomento usa un attributo di un'altra risorsa, Terraform crea una **dipendenza implicita**:

```hcl
resource "aws_subnet" "pubblica" {
  vpc_id     = aws_vpc.principale.id      # prima la VPC, poi la subnet
  cidr_block = "10.20.0.0/24"
}
```

Terraform costruisce un **grafo** di tutte le dipendenze: crea in parallelo le risorse indipendenti e in ordine
quelle collegate; in distruzione segue l'ordine inverso. Puoi vederlo con `terraform graph` (formato DOT, da
visualizzare con Graphviz).

Quando la dipendenza esiste ma non passa da un attributo, la dichiari con **`depends_on`** (dipendenza
**esplicita**). Esempio: un'istanza il cui script di avvio scarica file da un bucket deve aspettare che la policy che
le dà accesso esista.

```hcl
resource "aws_instance" "app" {
  # ...
  depends_on = [aws_iam_role_policy.lettura_bucket]
}
```

Usa `depends_on` solo quando serve: rende i piani più conservativi e la configurazione più difficile da leggere.

## Il blocco lifecycle

| Argomento | Effetto | Uso tipico |
|---|---|---|
| `create_before_destroy = true` | nelle sostituzioni crea la nuova risorsa prima di distruggere la vecchia | ridurre i tempi di indisponibilità |
| `prevent_destroy = true` | qualunque piano che distruggerebbe la risorsa fallisce | database, bucket con dati, chiavi KMS |
| `ignore_changes = [tags]` | ignora le differenze su certi attributi | attributi modificati da altri processi (per esempio la capacità desiderata di un ASG modificata dall'autoscaling) |
| `replace_triggered_by = [...]` | sostituisce la risorsa quando cambia un'altra | ricreare le istanze quando cambia il loro template |

```hcl
resource "aws_db_instance" "tasks" {
  # ...
  lifecycle {
    prevent_destroy = true
  }
}
```

!!! warning "prevent_destroy non è un backup"
    Protegge solo dagli errori **in Terraform**: se togli il blocco o la risorsa dalla configurazione, la protezione
    sparisce. Per i dati importanti attiva anche la protezione del servizio (`deletion_protection` su RDS e Cloud SQL)
    e i backup.

## Laboratorio

!!! laboratorio "Laboratorio · Rete e server web"
    Usa la cartella `labs/terraform/aws/rete-web` oppure `labs/terraform/gcp/rete-web`. Prima di eseguire, **leggi**
    i file: `versions.tf`, `variables.tf`, `main.tf`, `outputs.tf` e il modello `web.sh.tftpl`.

    === "AWS"

        La configurazione crea: VPC, Internet gateway, **due subnet pubbliche** in due zone (con `for_each`),
        tabella di routing e associazioni, un security group con regole separate per HTTP in ingresso e traffico in
        uscita, e **due istanze** (con `count`) con nginx installato dallo script di avvio, IMDSv2 obbligatorio.

        ```bash
        cd labs/terraform/aws/rete-web
        terraform init
        terraform plan          # conta le risorse: quante sono?
        terraform apply
        terraform output indirizzi_web
        ```

    === "Google Cloud"

        La configurazione crea: rete VPC in modalità personalizzata, una subnet, una regola firewall per HTTP verso le
        VM con il tag `web`, e **due VM** (con `count`) in due zone diverse, con nginx installato dallo script di
        avvio e Secure Boot attivo.

        ```bash
        cd labs/terraform/gcp/rete-web
        cp terraform.tfvars.example terraform.tfvars   # e scrivi l'ID del tuo progetto
        terraform init
        terraform plan
        terraform apply
        terraform output indirizzi_web
        ```

    Dopo un paio di minuti apri gli indirizzi nel browser: ogni server mostra il suo nome.

    Poi:

    1. `terraform graph | head -30`: individua le dipendenze tra rete, subnet e istanze.
    2. `terraform state list`: confronta con le risorse viste nella console.
    3. Cambia il testo nel file `web.sh.tftpl` e lancia `terraform plan`. Che cosa succede ai server, e perché?
    4. Aggiungi `lifecycle { prevent_destroy = true }` alla rete e prova `terraform destroy`. Poi toglilo.
    5. `terraform destroy`.

??? success "Spiegazione"
    1. Le istanze dipendono dalle subnet (e dal security group o dalla regola firewall), le subnet dalla rete.
    3. Su AWS lo script di avvio è eseguito solo al primo avvio: con `user_data_replace_on_change = true` Terraform
       **sostituisce** le istanze, così il nuovo script viene eseguito. Su Google Cloud cambia il metadato
       `startup-script` e la VM viene sostituita (`metadata_startup_script` forza la sostituzione). È l'approccio
       **immutabile**.
    4. `Error: Instance cannot be destroyed`: il piano fallisce prima di toccare qualunque risorsa.

## Esercizi

!!! esercizio "Esercizio · resource o data?"
    1. La VPC condivisa dell'azienda, gestita da un altro team, in cui devi creare le tue subnet.
    2. Il bucket per i log della tua applicazione.
    3. L'ID dell'account AWS in cui stai lavorando.
    4. L'ultima immagine Ubuntu 24.04 pubblicata da Canonical.

??? success "Soluzione"
    1. `data` (per esempio `data "aws_vpc"` filtrato per etichetta): non è tua, non devi gestirla.
    2. `resource`.
    3. `data "aws_caller_identity" "corrente" {}` e poi `data.aws_caller_identity.corrente.account_id`.
    4. `data "aws_ami"` con `owners = ["099720109477"]` (Canonical) e un filtro sul nome, oppure il parametro SSM
       pubblico di Canonical; su Google Cloud `data "google_compute_image"` con `family = "ubuntu-2404-lts-amd64"` e
       `project = "ubuntu-os-cloud"`.

!!! esercizio "Esercizio · Zero downtime"
    Un certificato TLS collegato a un load balancer deve essere sostituito. Con il comportamento predefinito il
    vecchio viene distrutto prima di creare il nuovo, ma il load balancer non può restare senza certificato. Cosa
    aggiungi?

??? success "Soluzione"
    ```hcl
    lifecycle {
      create_before_destroy = true
    }
    ```

    Il nuovo certificato viene creato, il load balancer aggiornato, e solo dopo il vecchio è distrutto. È il motivo per
    cui il piano mostra `+/-` invece di `-/+`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è la differenza tra un blocco resource e un blocco data?","options":["Nessuna","resource crea e gestisce un oggetto; data legge informazioni su un oggetto esistente","data crea risorse più velocemente","resource legge, data scrive"],"answer":[1],"explanation":"Un data source non viene mai distrutto da Terraform."},
 {"q":"Come crei una dipendenza implicita tra due risorse?","options":["Con depends_on","Usando in un argomento un attributo dell'altra risorsa","Mettendole nello stesso file","Nominandole in ordine alfabetico"],"answer":[1],"explanation":"depends_on crea invece una dipendenza esplicita."},
 {"q":"Quale argomento di lifecycle fa fallire ogni piano che distruggerebbe la risorsa?","options":["prevent_destroy","ignore_changes","create_before_destroy","deletion_protection"],"answer":[0],"explanation":"deletion_protection è un argomento di alcune risorse (RDS, Cloud SQL), non di lifecycle."},
 {"q":"Un ASG ha desired_capacity in Terraform, ma l'autoscaling la cambia di continuo e ogni piano vuole ripristinarla. Cosa usi?","options":["lifecycle { ignore_changes = [desired_capacity] }","depends_on","prevent_destroy","-target"],"answer":[0],"explanation":"Terraform ignorerà le differenze su quell'attributo."},
 {"q":"Quando serve depends_on?","options":["Sempre, per ogni risorsa","Quando una dipendenza esiste ma non è espressa da un riferimento a un attributo","Per i data source soltanto","Per ordinare gli output"],"answer":[1],"explanation":"Va usato con parsimonia."}
]}
</script>

## Fonti

- [Terraform: Resources](https://developer.hashicorp.com/terraform/language/resources/syntax)
- [Terraform: Data sources](https://developer.hashicorp.com/terraform/language/data-sources)
- [Terraform: depends_on](https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on)
- [Terraform: The lifecycle meta-argument](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle)
- [Terraform: terraform graph](https://developer.hashicorp.com/terraform/cli/commands/graph)
- [AWS Systems Manager: Query for the latest Amazon Linux AMI IDs](https://docs.aws.amazon.com/systems-manager/latest/userguide/parameter-store-public-parameters-ami.html)
