# 6.6 Lo state

!!! obiettivi "Obiettivi"
    - Spiegare a cosa serve lo state e perché va protetto.
    - Distinguere il backend locale da un backend remoto, e configurare S3 o Cloud Storage.
    - Capire il blocco dello state (*state locking*).
    - Rilevare e gestire le differenze tra codice e realtà (*drift*).
    - Ispezionare e modificare lo state con la CLI e con i blocchi `moved`, `removed`, `import`.
    - Usare i workspace della CLI.

## A cosa serve lo state

Lo **state** è un file JSON in cui Terraform registra, per ogni risorsa della configurazione, **quale oggetto reale**
le corrisponde (l'ID nel cloud) e i suoi attributi. Serve a:

- **collegare** il codice agli oggetti reali: `aws_instance.web[0]` è `i-0abc...`;
- **calcolare i piani**: confronta configurazione, state e realtà;
- conoscere le **dipendenze** anche delle risorse tolte dalla configurazione, per distruggerle nell'ordine giusto;
- **prestazioni**: per configurazioni grandi evita di interrogare tutto ogni volta (con `-refresh=false`).

!!! danger "Lo state contiene segreti"
    Nello state finiscono **in chiaro** tutti gli attributi delle risorse, comprese password generate e chiavi
    (lezione 6.8). Mai in Git, sempre in un backend cifrato con accesso limitato.

## Backend locale e remoto

Per impostazione predefinita Terraform usa il **backend locale**: il file `terraform.tfstate` nella cartella di
lavoro (con una copia precedente in `terraform.tfstate.backup`). Va bene per studiare da soli, non per lavorare in
team: ognuno avrebbe il suo state, e se perdi il file Terraform "dimentica" le risorse che ha creato.

Un **backend remoto** conserva lo state in un servizio condiviso, durevole e cifrato.

=== "AWS"

    ```hcl
    terraform {
      backend "s3" {
        bucket       = "academy-tfstate-123456789012"
        key          = "rete-web/dev/terraform.tfstate"
        region       = "eu-south-1"
        encrypt      = true
        use_lockfile = true      # blocco nativo con un file .tflock nel bucket
      }
    }
    ```

    Il bucket deve esistere prima, con **versioning** attivo (per recuperare versioni precedenti dello state) e
    accesso pubblico bloccato. Il blocco con una tabella DynamoDB (`dynamodb_table`) è il metodo precedente, ancora
    diffuso nei progetti esistenti ma **deprecato** a favore di `use_lockfile`.

=== "Google Cloud"

    ```hcl
    terraform {
      backend "gcs" {
        bucket = "academy-tfstate-ID-PROGETTO"
        prefix = "rete-web/dev"
      }
    }
    ```

    Il backend `gcs` supporta il blocco in modo nativo. Il bucket deve esistere prima, con **versioning** attivo e
    accesso uniforme.

Altri backend: `azurerm`, `consul`, `pg`, `kubernetes`, `http`, e **HCP Terraform** (lezione 6.9), che oltre allo
state esegue anche le operazioni.

!!! note "Il problema dell'uovo e della gallina"
    Il bucket per lo state non può essere creato dalla configurazione che lo usa come backend. Di solito si crea con
    una piccola configurazione separata (con state locale, oppure migrata subito dopo) o con la CLI.

Le **variabili non sono ammesse** nel blocco `backend`: i valori devono essere scritti per esteso o passati a `init`
con `-backend-config=file.hcl` (configurazione parziale), utile per usare lo stesso codice con bucket diversi.

## Il blocco dello state

Se due persone lanciano `apply` insieme sullo stesso state, le modifiche si sovrascrivono e lo state si corrompe. Il
**blocco** (*locking*) lo impedisce: chi inizia un'operazione di scrittura acquisisce il blocco, gli altri ricevono
`Error acquiring the state lock` e aspettano.

- `-lock-timeout=5m`: attende fino a 5 minuti invece di fallire subito.
- `terraform force-unlock ID`: rimuove un blocco rimasto appeso (per esempio dopo un crash). Usalo **solo** se sei
  certo che nessuno stia lavorando.

Il backend locale blocca il file sul tuo computer; non tutti i backend remoti supportano il blocco: verifica nella
documentazione.

## Drift: quando la realtà cambia

Il **drift** è una differenza tra lo state e la realtà, causata da modifiche fatte fuori da Terraform (dalla console,
da uno script, da un collega).

- `terraform plan` aggiorna prima lo state in memoria leggendo la realtà (*refresh*) e propone di **riportare** la
  realtà alla configurazione.
- `terraform plan -refresh-only` mostra **solo** le differenze trovate; `terraform apply -refresh-only` aggiorna lo
  state accettando la realtà, senza modificare le risorse. Poi dovrai aggiornare il codice, se la modifica va tenuta.

La regola d'oro: **ciò che è gestito da Terraform si modifica solo con Terraform**.

## Ispezionare e modificare lo state

| Comando | Uso |
|---|---|
| `terraform state list` | elenca le risorse nello state |
| `terraform state show aws_instance.web[0]` | mostra gli attributi di una risorsa |
| `terraform show` | mostra tutto lo state (o un piano salvato) in forma leggibile |
| `terraform output` | mostra gli output |
| `terraform state mv A B` | rinomina una risorsa nello state senza ricrearla |
| `terraform state rm A` | toglie una risorsa dallo state senza distruggerla: Terraform smette di gestirla |
| `terraform state pull` / `push` | scarica o carica lo state remoto (con estrema cautela) |
| `terraform state replace-provider` | cambia l'indirizzo del provider (per esempio dopo un fork) |

Dalla versione 1.1 in poi è meglio esprimere molte di queste operazioni **nel codice**, così passano dalla revisione
come ogni altra modifica:

```hcl
# Rinominare una risorsa (o spostarla in un modulo) senza ricrearla
moved {
  from = random_pet.nome
  to   = random_pet.nome_server
}

# Smettere di gestire una risorsa senza distruggerla
removed {
  from = local_file.benvenuto

  lifecycle {
    destroy = false
  }
}
```

## Importare risorse esistenti

Quando una risorsa esiste già (creata a mano) e vuoi che Terraform la gestisca, la **importi**. Il metodo moderno è il
blocco `import`, che Terraform esegue durante `plan` e `apply`:

=== "AWS"

    ```hcl
    import {
      to = aws_s3_bucket.vecchio
      id = "nome-del-bucket-creato-a-mano"
    }
    ```

=== "Google Cloud"

    ```hcl
    import {
      to = google_storage_bucket.vecchio
      id = "nome-del-bucket-creato-a-mano"
    }
    ```

Non hai ancora scritto la risorsa? Fattela proporre:

```bash
terraform plan -generate-config-out=importati.tf
```

Terraform scrive `importati.tf` con gli attributi letti dal cloud: rivedilo, pulisci i valori superflui, poi `apply`.
Il formato dell'`id` cambia per ogni tipo di risorsa: lo trovi in fondo alla sua pagina nel registry (sezione
*Import*). Il vecchio comando `terraform import INDIRIZZO ID` funziona ancora, ma importa una risorsa alla volta,
senza revisione e richiede di aver già scritto il blocco `resource`.

## Workspace della CLI

Un **workspace** è uno state separato per la stessa configurazione:

```bash
terraform workspace list          # * default
terraform workspace new test
terraform workspace select test
terraform apply                   # crea risorse separate, con uno state separato
```

Nel codice il nome è disponibile come `terraform.workspace` (per esempio per i nomi: `"academy-${terraform.workspace}"`).
I workspace della CLI vanno bene per copie temporanee e simili; per ambienti davvero diversi (dev e prod con account,
permessi e backend separati) si preferiscono **cartelle o configurazioni root separate** che usano gli stessi moduli.
I **workspace di HCP Terraform** sono un concetto diverso (lezione 6.9).

## Condividere dati tra configurazioni

Una configurazione può leggere gli **output** di un'altra con il data source `terraform_remote_state`:

```hcl
data "terraform_remote_state" "rete" {
  backend = "s3"
  config = {
    bucket = "academy-tfstate-123456789012"
    key    = "rete/terraform.tfstate"
    region = "eu-south-1"
  }
}

# data.terraform_remote_state.rete.outputs.vpc_id
```

Richiede l'accesso in lettura a tutto lo state dell'altra configurazione: in alternativa si usano data source mirati
(per esempio `data "aws_vpc"` filtrato per etichetta).

## Laboratorio

!!! laboratorio "Laboratorio · Dallo state locale a quello remoto"

    === "AWS"

        1. Crea il bucket per lo state (una volta sola):

            ```bash
            BUCKET=academy-tfstate-$(aws sts get-caller-identity --query Account --output text)
            aws s3api create-bucket --bucket "$BUCKET" --region eu-south-1 \
              --create-bucket-configuration LocationConstraint=eu-south-1
            aws s3api put-bucket-versioning --bucket "$BUCKET" --versioning-configuration Status=Enabled
            aws s3api put-public-access-block --bucket "$BUCKET" --public-access-block-configuration \
              BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true
            echo "$BUCKET"
            ```

        2. In `labs/terraform/01-primo-progetto` lancia `terraform apply` (state locale).
        3. Aggiungi al blocco `terraform` un `backend "s3"` con il tuo bucket e `key = "primo-progetto/terraform.tfstate"`.
        4. `terraform init -migrate-state` e rispondi `yes`: lo state viene copiato nel bucket.
        5. Verifica: `aws s3 ls "s3://$BUCKET/primo-progetto/"`. I file locali `terraform.tfstate` e
           `terraform.tfstate.backup` non servono più: cancellali, così nessuno li usa per errore.

    === "Google Cloud"

        1. Crea il bucket per lo state (una volta sola):

            ```bash
            BUCKET=academy-tfstate-$(gcloud config get-value project)
            gcloud storage buckets create "gs://$BUCKET" --location=europe-west8 \
              --uniform-bucket-level-access --public-access-prevention
            gcloud storage buckets update "gs://$BUCKET" --versioning
            echo "$BUCKET"
            ```

        2. In `labs/terraform/01-primo-progetto` lancia `terraform apply` (state locale).
        3. Aggiungi al blocco `terraform` un `backend "gcs"` con il tuo bucket e `prefix = "primo-progetto"`.
        4. `terraform init -migrate-state` e rispondi `yes`.
        5. Verifica: `gcloud storage ls "gs://$BUCKET/primo-progetto/"` (il file si chiama `default.tfstate`, dal nome
           del workspace).

    Poi, nella stessa cartella:

    6. Apri due terminali e lancia `terraform apply` in entrambi senza confermare nel primo: cosa dice il secondo?
    7. Rinomina `random_pet.nome` in `random_pet.nome_server` **con un blocco `moved`** e verifica che il piano
       mostri `has moved` e `0 to add, 0 to change, 0 to destroy`.
    8. Togli `local_file.benvenuto` con un blocco `removed` e `destroy = false`: dopo l'`apply` il file esiste ancora,
       ma non è più nello state (`terraform state list`).
    9. `terraform destroy`. Tieni il bucket dello state: ti servirà nel progetto.

??? success "Spiegazione"
    6. Il secondo terminale riceve `Error acquiring the state lock` con le informazioni su chi detiene il blocco:
       è il meccanismo che protegge lo state nei team. Annulla il primo apply e il blocco viene rilasciato.
    7. Senza `moved` Terraform distruggerebbe `random_pet.nome` e creerebbe `random_pet.nome_server`: con un database
       o una VM sarebbe un disastro.
    8. È l'equivalente dichiarativo di `terraform state rm`, ma passa dalla revisione del codice.

## Esercizi

!!! esercizio "Esercizio · Cosa fai?"
    1. Un collega ha aumentato a mano dalla console la dimensione di un'istanza gestita da Terraform. Vuoi tenere la
       modifica.
    2. Un `apply` si è interrotto per un crash e ora tutti ricevono `Error acquiring the state lock`.
    3. Vuoi spostare tre risorse in un nuovo modulo senza ricrearle.
    4. Un bucket creato a mano due anni fa deve passare sotto Terraform.
    5. Un collega propone di mettere `terraform.tfstate` su Git "così lo condividiamo".

??? success "Soluzione"
    1. Aggiorna il codice con il nuovo tipo di istanza; `terraform plan` deve risultare senza modifiche (eventualmente
       `apply -refresh-only` per allineare lo state).
    2. Verifica che nessuno stia davvero lavorando, poi `terraform force-unlock ID-DEL-BLOCCO`.
    3. Tre blocchi `moved` (da `aws_x.y` a `module.nome.aws_x.y`).
    4. Un blocco `import` (con `-generate-config-out` se non hai ancora scritto la risorsa), poi `plan` senza
       modifiche.
    5. No: lo state contiene segreti in chiaro e Git non offre il blocco, quindi due modifiche contemporanee lo
       corromperebbero. Serve un backend remoto.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Dove conserva lo state Terraform se non configuri un backend?","options":["In HCP Terraform","Nel file terraform.tfstate della cartella di lavoro","In memoria, e lo perde a ogni esecuzione","Nel file .terraform.lock.hcl"],"answer":[1],"explanation":"È il backend locale."},
 {"q":"Perché lo state non va messo in un repository Git?","options":["È troppo grande","Può contenere segreti in chiaro e Git non offre il blocco delle scritture","Git non accetta JSON","Terraform non lo ritroverebbe"],"answer":[1],"explanation":"Si usa un backend remoto cifrato, con blocco."},
 {"q":"Cosa impedisce il blocco dello state?","options":["La lettura dello state","Che due operazioni scrivano lo state contemporaneamente","La creazione di nuove risorse","L'uso dei moduli"],"answer":[1],"explanation":"Il secondo processo riceve un errore o attende con -lock-timeout."},
 {"q":"Come aggiorni lo state con le modifiche fatte fuori da Terraform, senza cambiare le risorse?","options":["terraform apply -refresh-only","terraform destroy","terraform state rm","terraform init -upgrade"],"answer":[0],"explanation":"Poi aggiorna il codice se vuoi tenere la modifica."},
 {"q":"Vuoi rinominare una risorsa nel codice senza che venga ricreata. Cosa usi?","options":["Un blocco moved","Un blocco removed","-replace","Un blocco import"],"answer":[0],"explanation":"In alternativa terraform state mv, ma il blocco moved passa dalla revisione."},
 {"q":"Quali affermazioni sul backend S3 sono corrette?","options":["use_lockfile = true attiva il blocco nativo nel bucket","Il bucket deve esistere prima di terraform init","Le variabili sono ammesse nel blocco backend","Il versioning è sconsigliato","Il blocco con DynamoDB è l'unico metodo supportato"],"answer":[0,1],"explanation":"Le variabili non sono ammesse nel blocco backend; il versioning è raccomandato; DynamoDB è il metodo precedente, deprecato."},
 {"q":"Quale comando genera la configurazione delle risorse indicate nei blocchi import?","options":["terraform import -all","terraform plan -generate-config-out=importati.tf","terraform show -json","terraform state pull"],"answer":[1],"explanation":"Il file generato va rivisto prima dell'apply."}
]}
</script>

## Fonti

- [Terraform: State](https://developer.hashicorp.com/terraform/language/state) e [Sensitive data in state](https://developer.hashicorp.com/terraform/language/state/sensitive-data)
- [Terraform: Backend configuration](https://developer.hashicorp.com/terraform/language/backend) e [State locking](https://developer.hashicorp.com/terraform/language/state/locking)
- [Terraform: S3 backend](https://developer.hashicorp.com/terraform/language/backend/s3) e [GCS backend](https://developer.hashicorp.com/terraform/language/backend/gcs)
- [Terraform: Manage resource drift](https://developer.hashicorp.com/terraform/tutorials/state/resource-drift)
- [Terraform: Refactoring with moved](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring) e [removed](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources)
- [Terraform: Import](https://developer.hashicorp.com/terraform/language/import) e [Generating configuration](https://developer.hashicorp.com/terraform/language/import/generating-configuration)
- [Terraform: CLI workspaces](https://developer.hashicorp.com/terraform/cli/workspaces)
- [Terraform: terraform_remote_state](https://developer.hashicorp.com/terraform/language/state/remote-state-data)
