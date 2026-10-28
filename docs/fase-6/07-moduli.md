# 6.7 Moduli

!!! obiettivi "Obiettivi"
    - Capire cos'è un modulo e la differenza tra modulo root e moduli figli.
    - Scrivere un modulo con input, output e documentazione.
    - Usare moduli locali, del registry e da Git, fissandone le versioni.
    - Capire la visibilità delle variabili tra moduli e il passaggio dei provider.
    - Rifattorizzare una configurazione in moduli senza ricreare le risorse.

## Cos'è un modulo

Ogni cartella con file `.tf` è un **modulo**. Quella in cui lanci `terraform` è il **modulo root**; i moduli che
richiama con un blocco `module` sono **moduli figli**. Un modulo è come una funzione: riceve **input** (le sue
variabili), crea risorse e restituisce **output**.

```text
progetto/
├── main.tf              # modulo root: usa i moduli
├── versions.tf
└── modules/
    └── rete/            # modulo figlio
        ├── main.tf
        ├── variables.tf
        ├── outputs.tf
        └── README.md
```

I moduli servono a:

- **riusare** lo stesso schema (una rete, un server web, un database) in progetti e ambienti diversi;
- **incapsulare** le scelte (cifratura, etichette, regole di sicurezza) in un unico punto;
- **ridurre** la configurazione root a poche righe leggibili.

## Usare un modulo

```hcl
module "rete_dev" {
  source = "./modules/rete"         # da dove prendere il codice

  nome = "academy-dev"              # le variabili del modulo
  cidr = "10.20.0.0/16"
  subnet_pubbliche = {
    "eu-south-1a" = "10.20.0.0/24"
    "eu-south-1b" = "10.20.1.0/24"
  }
}

module "rete_test" {                # lo stesso modulo, una seconda volta
  source = "./modules/rete"

  nome             = "academy-test"
  cidr             = "10.30.0.0/16"
  subnet_pubbliche = { "eu-south-1a" = "10.30.0.0/24" }
}

output "vpc_dev" {
  value = module.rete_dev.vpc_id     # gli output del modulo
}
```

Le risorse dentro il modulo hanno indirizzi come `module.rete_dev.aws_vpc.this`. Dopo aver aggiunto o cambiato un
modulo serve `terraform init` (o `terraform get`).

## Le sorgenti dei moduli

| Sorgente | Esempio di `source` | Versione |
|---|---|---|
| cartella locale | `"./modules/rete"` | segue il codice del repository |
| Terraform Registry pubblico | `"terraform-aws-modules/vpc/aws"` | argomento `version` |
| registry privato (HCP Terraform) | `"app.terraform.io/ORGANIZZAZIONE/rete/aws"` | argomento `version` |
| GitHub o Git | `"git::https://github.com/org/moduli.git//rete?ref=v1.2.0"` | `ref` = tag, ramo o commit |
| bucket S3 o GCS, URL HTTP | `"s3::https://..."`, `"gcs::https://..."` | nel percorso |

L'argomento **`version`** vale solo per i moduli dei registry. Per Git si usa `?ref=` con un **tag**. Fissa sempre la
versione: un modulo aggiornato senza controllo può cambiare le tue risorse.

### Moduli del registry

Il [Terraform Registry](https://registry.terraform.io/browse/modules) ospita migliaia di moduli. Quelli dei gruppi
`terraform-aws-modules` (comunità) e `terraform-google-modules` (mantenuti da Google) sono molto usati:

```hcl
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"   # controlla l'ultima versione nella pagina del modulo

  name = "academy"
  cidr = "10.20.0.0/16"

  azs             = ["eu-south-1a", "eu-south-1b"]
  public_subnets  = ["10.20.0.0/24", "10.20.1.0/24"]
  private_subnets = ["10.20.10.0/24", "10.20.11.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true
}
```

Prima di usare un modulo pubblico leggi il codice e la documentazione: stai delegando a terzi la definizione della
tua infrastruttura. I moduli **verificati** (*partner*) sono pubblicati da aziende partner di HashiCorp.

## Visibilità delle variabili

- Un modulo vede **solo** le sue variabili, i suoi locals e le sue risorse. Non può leggere `var.ambiente` del root né
  le risorse di un altro modulo.
- Tutto ciò che il modulo deve ricevere passa dagli **argomenti** del blocco `module`; tutto ciò che deve restituire
  passa dai suoi **output**.
- Per collegare due moduli, passa l'output dell'uno come input dell'altro:

```hcl
module "web" {
  source     = "./modules/web"
  subnet_ids = module.rete_dev.subnet_ids
}
```

## Provider e moduli

Un modulo figlio **eredita** la configurazione predefinita dei provider dal root. Un modulo riutilizzabile non
dovrebbe contenere blocchi `provider`, ma solo dichiarare in `required_providers` quali provider e versioni gli
servono. Per passargli una configurazione con alias:

```hcl
module "certificato" {
  source = "./modules/certificato"
  providers = {
    aws = aws.virginia
  }
}
```

## Meta-argomenti dei moduli

Anche i moduli accettano `count`, `for_each` e `depends_on`:

```hcl
module "rete" {
  for_each = {
    dev  = "10.20.0.0/16"
    test = "10.30.0.0/16"
  }

  source           = "./modules/rete"
  nome             = "academy-${each.key}"
  cidr             = each.value
  subnet_pubbliche = { "eu-south-1a" = cidrsubnet(each.value, 8, 0) }
}
```

## Buone pratiche

- Un modulo deve fare **una cosa** (la rete, il database) e farla bene; non avvolgere una sola risorsa in un modulo
  senza aggiungere valore.
- Variabili con `description`, `type` e, dove ha senso, `validation`; output con `description`.
- Un `README.md` con un esempio d'uso. Lo strumento della comunità `terraform-docs` genera le tabelle di input e
  output.
- Versioni semantiche (`v1.2.0`) con tag Git, e un registro delle modifiche.
- Test con `terraform test` dentro la cartella del modulo.

## Laboratorio

!!! laboratorio "Laboratorio · Estrai il modulo della rete"
    Parti da una copia di `labs/terraform/aws/rete-web` (o `gcp/rete-web`) **con le risorse create** (`apply`).

    1. Crea `modules/rete/` con `main.tf`, `variables.tf`, `outputs.tf`.
    2. **Sposta** nel modulo la rete e tutto ciò che la compone (su AWS: VPC, Internet gateway, subnet, tabella di
       routing e associazioni; su Google Cloud: rete e subnet). Il modulo riceve nome, CIDR e zone, e restituisce gli
       ID della rete e delle subnet.
    3. Nel root sostituisci le risorse con un blocco `module "rete"` e usa i suoi output per le istanze.
    4. `terraform init`, poi `terraform plan`: Terraform vuole distruggere e ricreare la rete. **Non applicare.**
    5. Aggiungi un blocco `moved` per ogni risorsa spostata e rilancia `terraform plan` fino a ottenere
       `0 to add, 0 to change, 0 to destroy`.
    6. `terraform destroy` alla fine.

??? success "Soluzione (AWS)"
    `modules/rete/variables.tf`:

    ```hcl
    variable "nome" {
      description = "Prefisso dei nomi"
      type        = string
    }

    variable "cidr" {
      description = "Intervallo della VPC"
      type        = string
    }

    variable "subnet_pubbliche" {
      description = "Mappa zona => CIDR delle subnet pubbliche"
      type        = map(string)
    }
    ```

    `modules/rete/main.tf` contiene le risorse spostate, rinominate in modo neutro (`aws_vpc.this`,
    `aws_internet_gateway.this`, `aws_subnet.pubblica`, `aws_route_table.pubblica`,
    `aws_route_table_association.pubblica`), con `var.cidr` e `for_each = var.subnet_pubbliche`.

    `modules/rete/outputs.tf`:

    ```hcl
    output "vpc_id" {
      description = "ID della VPC"
      value       = aws_vpc.this.id
    }

    output "subnet_ids" {
      description = "ID delle subnet pubbliche, nell'ordine delle zone"
      value       = [for s in aws_subnet.pubblica : s.id]
    }
    ```

    Nel root:

    ```hcl
    module "rete" {
      source           = "./modules/rete"
      nome             = local.prefisso
      cidr             = var.cidr_vpc
      subnet_pubbliche = local.subnet_pubbliche
    }

    moved {
      from = aws_vpc.principale
      to   = module.rete.aws_vpc.this
    }

    moved {
      from = aws_internet_gateway.principale
      to   = module.rete.aws_internet_gateway.this
    }

    moved {
      from = aws_subnet.pubblica
      to   = module.rete.aws_subnet.pubblica
    }

    moved {
      from = aws_route_table.pubblica
      to   = module.rete.aws_route_table.pubblica
    }

    moved {
      from = aws_route_table_association.pubblica
      to   = module.rete.aws_route_table_association.pubblica
    }
    ```

    Il security group e le istanze restano nel root: `vpc_id = module.rete.vpc_id` e
    `subnet_id = module.rete.subnet_ids[count.index % length(module.rete.subnet_ids)]`.

    **Spiegazione.** Un blocco `moved` su una risorsa con `for_each` sposta tutte le sue istanze insieme. Se nel
    modulo hai cambiato anche gli argomenti (per esempio le etichette `Name`), il piano mostrerà delle modifiche sul
    posto: sono accettabili, purché non ci siano sostituzioni. I blocchi `moved` si possono rimuovere quando tutti gli
    state che usano la configurazione sono stati aggiornati.

## Esercizi

!!! esercizio "Esercizio · Sorgenti e versioni"
    Scrivi l'argomento `source` (e `version` se serve) per:

    1. il modulo nella cartella `../moduli/database`;
    2. il tag `v2.1.0` del modulo nella sottocartella `rete` del repository `https://github.com/acme/infra-moduli`;
    3. il modulo `terraform-aws-modules/s3-bucket/aws`, qualunque versione 5.x.

??? success "Soluzione"
    1. `source = "../moduli/database"`
    2. `source = "git::https://github.com/acme/infra-moduli.git//rete?ref=v2.1.0"` (la doppia barra separa il
       repository dalla sottocartella)
    3. `source = "terraform-aws-modules/s3-bucket/aws"` e `version = "~> 5.0"`

!!! esercizio "Esercizio · Perché non funziona?"
    Nel modulo `modules/web` scrivi `subnet_id = aws_subnet.pubblica["eu-south-1a"].id`, ma la subnet è definita nel
    modulo root. `terraform validate` dà errore. Perché, e come si risolve?

??? success "Soluzione"
    Un modulo non vede le risorse del root (né di altri moduli). Si aggiunge al modulo una variabile `subnet_id` e il
    root gliela passa: `subnet_id = aws_subnet.pubblica["eu-south-1a"].id` nel blocco `module "web"`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cos'è il modulo root?","options":["Un modulo del registry ufficiale","La cartella in cui esegui i comandi di Terraform","Un modulo che richiede i permessi di root","Il primo modulo in ordine alfabetico"],"answer":[1],"explanation":"I moduli richiamati con un blocco module sono figli."},
 {"q":"Per quali sorgenti si usa l'argomento version?","options":["Per tutti i moduli","Solo per i moduli di un registry (pubblico o privato)","Solo per i moduli locali","Solo per i moduli Git"],"answer":[1],"explanation":"Per Git si usa ?ref= nel source."},
 {"q":"Un modulo figlio può leggere direttamente le variabili del modulo root?","options":["Sì, sempre","No: riceve solo i valori passati come argomenti nel blocco module","Solo se sono sensitive","Solo con depends_on"],"answer":[1],"explanation":"I valori tornano al root tramite gli output del modulo."},
 {"q":"Come leggi l'output vpc_id del modulo rete?","options":["var.rete.vpc_id","module.rete.vpc_id","output.rete.vpc_id","rete.outputs.vpc_id"],"answer":[1],"explanation":"module.NOME.OUTPUT."},
 {"q":"Dopo aver aggiunto un nuovo blocco module, quale comando serve prima del plan?","options":["terraform init","terraform fmt","terraform refresh","terraform output"],"answer":[0],"explanation":"init scarica o collega i moduli (anche terraform get)."},
 {"q":"Quali sono buone pratiche per i moduli riutilizzabili?","options":["Dichiarare required_providers senza blocchi provider propri","Fissare le versioni dei moduli usati","Leggere le variabili del root con var.","Usare sempre -target","Mettere le credenziali nel modulo"],"answer":[0,1],"explanation":"I provider sono configurati dal root ed ereditati o passati con providers."}
]}
</script>

## Fonti

- [Terraform: Modules overview](https://developer.hashicorp.com/terraform/language/modules)
- [Terraform: Module blocks](https://developer.hashicorp.com/terraform/language/modules/syntax) e [Module sources](https://developer.hashicorp.com/terraform/language/modules/sources)
- [Terraform: Providers within modules](https://developer.hashicorp.com/terraform/language/modules/develop/providers)
- [Terraform: Refactoring modules with moved](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring)
- [Terraform: Module creation, recommended pattern](https://developer.hashicorp.com/terraform/tutorials/modules/pattern-module-creation)
- [Terraform Registry: terraform-aws-modules/vpc](https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/latest) e [terraform-google-modules/network](https://registry.terraform.io/modules/terraform-google-modules/network/google/latest)
