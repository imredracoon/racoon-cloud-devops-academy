# 6.5 Variabili, output ed espressioni

!!! obiettivi "Obiettivi"
    - Parametrizzare una configurazione con variabili tipizzate e validate.
    - Passare valori alle variabili e conoscere l'ordine di precedenza.
    - Esporre informazioni con gli output e semplificare con i valori locali.
    - Usare i tipi complessi, le espressioni, le funzioni e la console.
    - Creare più risorse con `count` e `for_each`, e blocchi ripetuti con `dynamic`.
    - Verificare la configurazione con condizioni personalizzate.

## Variabili di input

```hcl
variable "ambiente" {
  description = "Ambiente: dev, test o prod"
  type        = string
  default     = "dev"
  nullable    = false

  validation {
    condition     = contains(["dev", "test", "prod"], var.ambiente)
    error_message = "L'ambiente deve essere dev, test o prod."
  }
}
```

- Si usano con `var.ambiente`.
- Senza `default` la variabile è **obbligatoria**: Terraform la chiede in modo interattivo o fallisce in automazione.
- `sensitive = true` nasconde il valore nell'output della CLI ([lezione 6.8](08-segreti.md)).

### Tipi

| Tipo | Esempio di valore |
|---|---|
| `string`, `number`, `bool` | `"eu-south-1"`, `2`, `true` |
| `list(string)` | `["a", "b"]`: ordinata, ammette duplicati |
| `set(string)` | come una lista ma senza ordine né duplicati |
| `map(number)` | `{ web = 2, db = 1 }`: chiavi stringa, valori dello stesso tipo |
| `object({ nome = string, porta = number })` | struttura con attributi di tipi diversi |
| `tuple([string, number])` | sequenza di lunghezza fissa con tipi per posizione |

I tipi si combinano. Gli attributi di un oggetto possono essere **facoltativi**, con un valore predefinito:

```hcl
variable "server" {
  type = map(object({
    tipo  = string
    disco = optional(number, 20)   # se omesso vale 20
  }))
  default = {
    web = { tipo = "t3.micro" }
    db  = { tipo = "t3.small", disco = 50 }
  }
}
```

### Assegnare i valori

In ordine di **precedenza crescente** (l'ultimo vince):

1. il `default` nella dichiarazione;
2. le variabili d'ambiente `TF_VAR_nome` (per esempio `export TF_VAR_ambiente=prod`);
3. il file `terraform.tfvars`, poi `terraform.tfvars.json`;
4. i file `*.auto.tfvars` e `*.auto.tfvars.json`, in ordine alfabetico;
5. le opzioni `-var` e `-var-file` sulla riga di comando, nell'ordine in cui le scrivi.

Un file `.tfvars` contiene solo assegnazioni:

```hcl
# prod.tfvars
ambiente      = "prod"
numero_server = 4
```

```bash
terraform plan -var-file=prod.tfvars
```

## Output

```hcl
output "indirizzi_web" {
  description = "URL dei server web"
  value       = [for i in aws_instance.web : "http://${i.public_ip}"]
}
```

Gli output si leggono con `terraform output` (`-raw nome` per usarli negli script, `-json` per tutti) e servono anche
a passare valori da un modulo a chi lo usa ([lezione 6.7](07-moduli.md)).

## Valori locali

I **locals** danno un nome a un'espressione usata più volte:

```hcl
locals {
  prefisso = "${var.nome}-${var.ambiente}"
  etichette_comuni = {
    progetto = var.nome
    ambiente = var.ambiente
  }
}

resource "aws_vpc" "principale" {
  cidr_block = var.cidr_vpc
  tags       = merge(local.etichette_comuni, { Name = "${local.prefisso}-vpc" })
}
```

Differenza: una **variabile** è un input che chi usa la configurazione può cambiare; un **local** è un calcolo interno.

## Espressioni

| Espressione | Esempio | Risultato |
|---|---|---|
| interpolazione | `"${var.nome}-vpc"` | `"academy-vpc"` |
| condizionale | `var.ambiente == "prod" ? "t3.large" : "t3.micro"` | `"t3.micro"` in dev |
| `for` su lista | `[for p in [80, 443] : "porta-${p}"]` | `["porta-80", "porta-443"]` |
| `for` su mappa | `{ for nome, s in var.server : nome => s.disco }` | `{ db = 50, web = 20 }` |
| `for` con filtro | `[for s in values(var.server) : s.tipo if s.disco > 30]` | `["t3.small"]` |
| splat | `aws_instance.web[*].id` | lista degli ID di tutte le istanze |
| indice | `var.server["web"].disco` | `20` |

### Funzioni

Terraform ha decine di funzioni predefinite (non puoi scriverne di tue). Alcune delle più usate:

| Categoria | Funzioni |
|---|---|
| stringhe | `format("%s-%03d", "web", 7)` → `"web-007"`, `join`, `split`, `upper`, `lower`, `replace`, `trimspace` |
| collezioni | `length`, `merge`, `concat`, `contains`, `keys`, `values`, `lookup(mappa, chiave, predefinito)`, `flatten`, `toset`, `slice`, `element` |
| rete | `cidrsubnet("10.20.0.0/16", 8, 3)` → `"10.20.3.0/24"`, `cidrhost` |
| file e codifica | `file`, `templatefile`, `jsonencode`, `yamlencode`, `base64encode` |
| tipi e valori | `coalesce("", "x")` → `"x"`, `try`, `can`, `tostring`, `tonumber` |

Prova le espressioni nella **console**:

```bash
terraform console
> cidrsubnet("10.20.0.0/16", 8, 1)
"10.20.1.0/24"
> [for p in [80, 443] : "porta-${p}"]
```

La console legge anche variabili, locals e (dopo un `apply`) gli attributi delle risorse dello state.

## Più risorse: count e for_each

| | `count` | `for_each` |
|---|---|---|
| Valore | un numero | una mappa o un insieme di stringhe |
| Indirizzi | `aws_instance.web[0]`, `[1]`... | `aws_subnet.pubblica["eu-south-1a"]` |
| Nel blocco | `count.index` | `each.key`, `each.value` |
| Se togli un elemento in mezzo | gli indici successivi scalano: Terraform modifica o ricrea risorse che non volevi toccare | si distrugge solo quell'elemento |
| Uso tipico | N copie identiche; `count = var.crea ? 1 : 0` per una risorsa facoltativa | elementi distinti con un nome stabile |

Nel laboratorio `rete-web` le subnet usano `for_each` su una mappa zona → CIDR, costruita con un'espressione `for`:

```hcl
locals {
  zone = slice(data.aws_availability_zones.disponibili.names, 0, 2)
  subnet_pubbliche = {
    for i, zona in local.zone : zona => cidrsubnet(var.cidr_vpc, 8, i)
  }
}

resource "aws_subnet" "pubblica" {
  for_each = local.subnet_pubbliche

  vpc_id            = aws_vpc.principale.id
  availability_zone = each.key
  cidr_block        = each.value
}
```

## Blocchi ripetuti: dynamic

Alcune risorse hanno blocchi annidati ripetibili (regole di un security group, dischi di una VM). Un blocco
`dynamic` li genera da una collezione:

```hcl
variable "porte_aperte" {
  type    = list(number)
  default = [80, 443]
}

resource "aws_security_group" "web" {
  name = "web"

  dynamic "ingress" {
    for_each = var.porte_aperte
    content {
      from_port   = ingress.value
      to_port     = ingress.value
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }
}
```

Usali con moderazione: rendono la configurazione meno leggibile.

## Condizioni personalizzate

Oltre alla `validation` delle variabili, Terraform permette di dichiarare **ipotesi** e **garanzie**:

```hcl
data "aws_ec2_instance_type" "scelto" {
  instance_type = var.tipo_istanza
}

resource "aws_instance" "web" {
  # ...
  lifecycle {
    precondition {                     # controllata prima di creare o modificare
      condition     = data.aws_ec2_instance_type.scelto.free_tier_eligible
      error_message = "Per il laboratorio usa un tipo di istanza del piano gratuito."
    }
    postcondition {                    # controllata dopo, sull'oggetto reale (self)
      condition     = self.public_ip != ""
      error_message = "L'istanza deve avere un indirizzo IP pubblico."
    }
  }
}
```

I blocchi **`check`** verificano proprietà dell'infrastruttura a ogni `plan` e `apply` **senza bloccare**: producono
un avviso. Esempio, con il provider `hashicorp/http`:

```hcl
check "sito_raggiungibile" {
  data "http" "home" {
    url = "https://www.esempio.it"
  }

  assert {
    condition     = data.http.home.status_code == 200
    error_message = "Il sito non risponde con 200."
  }
}
```

| Meccanismo | Dove | Se fallisce |
|---|---|---|
| `validation` | blocco `variable` | errore: la variabile non è valida |
| `precondition` / `postcondition` | `lifecycle` di risorse e data source, blocchi `output` | errore: l'operazione si ferma |
| `check` | blocco autonomo | **avviso**: l'operazione continua |

## Laboratorio

!!! laboratorio "Laboratorio · Parametrizza la rete"
    Nella cartella `labs/terraform/aws/rete-web` (o `gcp/rete-web`):

    1. Leggi `variables.tf` e prova `terraform plan -var="ambiente=collaudo"`: la validazione fallisce.
    2. Crea `dev.tfvars` con `numero_server = 1` e `prod.tfvars` con `numero_server = 3` e `ambiente = "prod"`.
       Confronta `terraform plan -var-file=dev.tfvars` e `-var-file=prod.tfvars` (senza applicare).
    3. Imposta `export TF_VAR_numero_server=2` e lancia `terraform plan -var-file=dev.tfvars`: quanti server? Perché?
    4. Apri `terraform console` e calcola le subnet: `local.subnet_pubbliche`.
    5. Lancia i test senza cloud: `terraform test`. Leggi `tests/rete.tftest.hcl`: il provider è **simulato**
       (`mock_provider`), quindi i test girano anche senza credenziali.

??? success "Spiegazione"
    3. **Uno**: `-var-file` sulla riga di comando ha precedenza sulle variabili d'ambiente `TF_VAR_`.
    4. Con le zone simulate o reali ottieni una mappa come `{ "eu-south-1a" = "10.20.0.0/24", "eu-south-1b" =
       "10.20.1.0/24" }`.
    5. `terraform test` esegue i blocchi `run` (qui con `command = plan`) e verifica le `assert`. Il test
       `ambiente_non_valido` **si aspetta** il fallimento della validazione (`expect_failures`).

## Esercizi

!!! esercizio "Esercizio · count o for_each?"
    Hai `count = length(var.utenti)` con `var.utenti = ["anna", "bruno", "carla"]` per creare tre utenti IAM. Togli
    `"bruno"` dalla lista. Cosa propone Terraform? Come riscrivi la risorsa per evitarlo?

??? success "Soluzione"
    Con `count`, `carla` passa dall'indice 2 all'indice 1: Terraform propone di **rinominare** `[1]` (da bruno a
    carla) e di distruggere `[2]`. Con `for_each` ogni utente ha un indirizzo stabile:

    ```hcl
    resource "aws_iam_user" "utenti" {
      for_each = toset(var.utenti)
      name     = each.value
    }
    ```

    Togliendo `bruno` viene distrutto solo `aws_iam_user.utenti["bruno"]`.

!!! esercizio "Esercizio · Scrivi le espressioni"
    Data `var.server` della sezione sui tipi:

    1. Una lista con i nomi dei server (`["db", "web"]`).
    2. Una mappa nome → tipo, ma solo per i server con disco di almeno 30 GB.
    3. Il tipo di istanza da usare: `"t3.large"` se l'ambiente è `prod`, altrimenti `"t3.micro"`.

??? success "Soluzione"
    1. `keys(var.server)`
    2. `{ for nome, s in var.server : nome => s.tipo if s.disco >= 30 }` → `{ db = "t3.small" }`
    3. `var.ambiente == "prod" ? "t3.large" : "t3.micro"`

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale fonte di valori per una variabile ha la precedenza più alta?","options":["Il default","Una variabile d'ambiente TF_VAR_","Il file terraform.tfvars","L'opzione -var sulla riga di comando"],"answer":[3],"explanation":"-var e -var-file vincono su tutto, nell'ordine in cui compaiono."},
 {"q":"Quale file viene caricato automaticamente senza opzioni?","options":["prod.tfvars","terraform.tfvars","variables.tf.json","vars.hcl"],"answer":[1],"explanation":"Anche i file *.auto.tfvars sono caricati automaticamente."},
 {"q":"Perché preferire for_each a count per risorse con un'identità propria?","options":["È più veloce","Gli indirizzi sono basati su chiavi stabili: togliere un elemento non sposta gli altri","count non esiste più","for_each non richiede variabili"],"answer":[1],"explanation":"Con count gli indici scalano."},
 {"q":"Cosa produce cidrsubnet(\"10.0.0.0/16\", 8, 2)?","options":["10.0.2.0/24","10.0.0.2/24","10.2.0.0/16","10.0.8.0/24"],"answer":[0],"explanation":"Aggiunge 8 bit alla maschera (/24) e prende la terza subnet (indice 2)."},
 {"q":"Quali meccanismi bloccano l'operazione se la condizione non è rispettata?","options":["validation nelle variabili","precondition","check","Un commento # TODO","Un output"],"answer":[0,1],"explanation":"Un blocco check produce solo un avviso."},
 {"q":"Qual è la differenza tra una variabile e un local?","options":["Nessuna","La variabile è un input impostabile dall'esterno; il local è un valore calcolato internamente","I local sono più lenti","Le variabili non hanno tipo"],"answer":[1],"explanation":"I locals evitano di ripetere espressioni."}
]}
</script>

## Fonti

- [Terraform: Input variables](https://developer.hashicorp.com/terraform/language/values/variables) e [Variable definition precedence](https://developer.hashicorp.com/terraform/language/values/variables#variable-definition-precedence)
- [Terraform: Output values](https://developer.hashicorp.com/terraform/language/values/outputs) e [Local values](https://developer.hashicorp.com/terraform/language/values/locals)
- [Terraform: Type constraints](https://developer.hashicorp.com/terraform/language/expressions/type-constraints)
- [Terraform: for expressions](https://developer.hashicorp.com/terraform/language/expressions/for), [dynamic blocks](https://developer.hashicorp.com/terraform/language/expressions/dynamic-blocks), [Functions](https://developer.hashicorp.com/terraform/language/functions)
- [Terraform: count](https://developer.hashicorp.com/terraform/language/meta-arguments/count) e [for_each](https://developer.hashicorp.com/terraform/language/meta-arguments/for_each)
- [Terraform: Custom conditions](https://developer.hashicorp.com/terraform/language/expressions/custom-conditions) e [Checks](https://developer.hashicorp.com/terraform/language/checks)
- [Terraform: Tests](https://developer.hashicorp.com/terraform/language/tests)
