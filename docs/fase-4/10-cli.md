# 4.10 Le CLI: aws e gcloud

!!! obiettivi "Obiettivi"
    - Configurare l'accesso della CLI con credenziali temporanee: `aws configure sso` e `gcloud auth login`.
    - Usare profili e configurazioni per lavorare con più account e progetti.
    - Filtrare e formattare l'output (`--query`, `--filter`, `--format`) per usarlo negli script.
    - Creare, elencare e distruggere risorse dalla riga di comando.

## Perché la riga di comando

La console web è ottima per imparare ed esplorare. Per il lavoro quotidiano e l'automazione si usa la **CLI**: i
comandi si possono salvare, ripetere, mettere in uno script, revisionare in una pull request. Dietro le quinte, la
console, la CLI, gli SDK (come `boto3` per Python) e Terraform chiamano tutti le stesse **API** del fornitore.

Nel Dev Container trovi già installate la **AWS CLI v2** (`aws`) e la **Google Cloud CLI** (`gcloud`).

```bash
aws --version
gcloud --version
```

## Configurare l'accesso

=== "AWS"

    Con IAM Identity Center (lezione 4.4) la CLI ottiene **credenziali temporanee** senza chiavi permanenti:

    ```bash
    aws configure sso
    # SSO session name: academy
    # SSO start URL: https://d-xxxxxxxxxx.awsapps.com/start   (il tuo AWS access portal URL)
    # SSO region: la regione in cui hai attivato Identity Center
    # SSO registration scopes: lascia il predefinito
    # → si apre il browser (o ti dà un codice da inserire): accedi e autorizza
    # scegli account e ruolo (AdministratorAccess)
    # Default client Region: eu-south-1
    # Default output format: json
    # Profile name: academy
    ```

    La configurazione finisce in `~/.aws/config`. Da quel momento:

    ```bash
    export AWS_PROFILE=academy          # oppure --profile academy a ogni comando
    aws sts get-caller-identity         # chi sono? account, ruolo, ARN
    aws sso login                       # quando la sessione scade
    ```

    !!! warning "Evita aws configure con le access key"
        `aws configure` chiede *Access Key ID* e *Secret Access Key*: sono credenziali **permanenti** salvate in chiaro
        in `~/.aws/credentials`. Usale solo se non hai alternative, mai per l'utente root, e non pubblicarle mai.

=== "Google Cloud"

    ```bash
    gcloud init                          # guida: accesso, scelta del progetto, regione e zona predefinite
    # oppure, passo per passo:
    gcloud auth login                    # si apre il browser
    gcloud config set project ID-PROGETTO
    gcloud config set compute/region europe-west8
    gcloud config set compute/zone europe-west8-a
    gcloud config list
    gcloud auth list                     # gli account autenticati
    ```

    Le **configurazioni con nome** permettono di passare tra progetti diversi:

    ```bash
    gcloud config configurations create academy
    gcloud config configurations activate academy
    gcloud config configurations list
    ```

    Nota: `gcloud auth login` autentica **la CLI**. Le **librerie** (Python, Terraform) usano le *Application Default
    Credentials*, che si impostano a parte con `gcloud auth application-default login`. Lo userai nella Fase 6.

    In GitHub Codespaces e nel Dev Container, se il browser non si apre, usa `gcloud auth login --no-launch-browser`
    e segui le istruzioni.

## La struttura dei comandi

=== "AWS"

    ```text
    aws <servizio> <operazione> [--parametri]
    aws ec2 describe-instances --region eu-south-1
    aws s3api create-bucket ...
    ```

    Per S3 esistono due insiemi di comandi: `aws s3` (alto livello, simile ai comandi di Linux: `ls`, `cp`, `sync`,
    `rm`) e `aws s3api` (accesso diretto a tutte le API).

    L'aiuto: `aws ec2 help`, `aws ec2 describe-instances help`.

=== "Google Cloud"

    ```text
    gcloud <gruppo> [<sottogruppo>] <comando> [--flag]
    gcloud compute instances list
    gcloud storage buckets create ...
    ```

    Per lo storage si usa `gcloud storage` (il vecchio strumento `gsutil` è ancora disponibile ma non più consigliato).

    L'aiuto: `gcloud compute instances --help`, `gcloud help compute instances create`.

## Filtrare e formattare l'output

Negli script ti serve estrarre un valore preciso (un ID, un indirizzo IP) dall'output.

=== "AWS"

    `--output` sceglie il formato (`json`, `table`, `text`, `yaml`); `--query` filtra con il linguaggio **JMESPath**:

    ```bash
    aws ec2 describe-instances \
      --query 'Reservations[].Instances[].[InstanceId, State.Name, PublicIpAddress]' \
      --output table

    # solo le istanze in esecuzione, solo gli ID, come testo
    aws ec2 describe-instances \
      --filters Name=instance-state-name,Values=running \
      --query 'Reservations[].Instances[].InstanceId' --output text
    ```

    `--filters` filtra **sul server** (meno dati trasferiti); `--query` filtra **sul client**, dopo la risposta.

=== "Google Cloud"

    `--format` sceglie formato e campi; `--filter` filtra:

    ```bash
    gcloud compute instances list \
      --format="table(name, zone.basename(), status, networkInterfaces[0].accessConfigs[0].natIP)"

    # solo le VM in esecuzione, solo i nomi
    gcloud compute instances list --filter="status=RUNNING" --format="value(name)"

    gcloud compute instances describe lab-web --zone europe-west8-a --format=json
    ```

    `value(...)` stampa i valori senza intestazioni: perfetto per gli script.

In alternativa a entrambe, puoi chiedere JSON e filtrarlo con **jq**, già presente nel Dev Container:
`aws ec2 describe-instances | jq -r '.Reservations[].Instances[].InstanceId'`.

## Esempi di uso comune

| Operazione | AWS | Google Cloud |
|---|---|---|
| Chi sono | `aws sts get-caller-identity` | `gcloud auth list`, `gcloud config list` |
| Elenco delle regioni | `aws ec2 describe-regions --output table` | `gcloud compute regions list` |
| Creare un bucket | `aws s3 mb s3://NOME` | `gcloud storage buckets create gs://NOME --location=REGIONE` |
| Copiare un file | `aws s3 cp file.txt s3://NOME/` | `gcloud storage cp file.txt gs://NOME/` |
| Sincronizzare una cartella | `aws s3 sync ./sito s3://NOME/` | `gcloud storage rsync ./sito gs://NOME --recursive` |
| Elencare | `aws s3 ls s3://NOME/` | `gcloud storage ls gs://NOME/` |
| Eliminare un bucket e il contenuto | `aws s3 rb s3://NOME --force` | `gcloud storage rm --recursive gs://NOME` |
| Elencare le VM | `aws ec2 describe-instances` | `gcloud compute instances list` |
| Collegarsi in SSH | `aws ec2-instance-connect ssh --instance-id ID` | `gcloud compute ssh NOME` |
| Abilitare un servizio | — (i servizi sono già attivi) | `gcloud services enable compute.googleapis.com` |

Su Google Cloud molti servizi vanno **abilitati** nel progetto prima dell'uso (`gcloud services enable ...`): è una
differenza frequente rispetto ad AWS.

## Laboratorio

!!! laboratorio "Laboratorio · Un ciclo completo da riga di comando"
    Scrivi lo script `ciclo-bucket.sh` che, con la CLI del tuo cloud:

    1. crea un bucket con un nome univoco (per esempio con la data e un numero casuale) nella regione del corso,
       con un'etichetta `progetto=academy`;
    2. carica tre file di prova;
    3. elenca il contenuto;
    4. elimina bucket e contenuto, anche se uno dei passi precedenti fallisce (usa `trap`, lezione 2.4).

??? success "Soluzione"
    === "AWS"

        ```bash
        #!/usr/bin/env bash
        set -euo pipefail
        REGIONE=${REGIONE:-eu-south-1}
        BUCKET="academy-lab-$(date +%Y%m%d)-$RANDOM"

        pulizia() { aws s3 rb "s3://$BUCKET" --force > /dev/null 2>&1 || true; echo "Eliminato $BUCKET"; }
        trap pulizia EXIT

        aws s3api create-bucket --bucket "$BUCKET" --region "$REGIONE" \
          --create-bucket-configuration LocationConstraint="$REGIONE"
        aws s3api put-bucket-tagging --bucket "$BUCKET" \
          --tagging 'TagSet=[{Key=progetto,Value=academy}]'

        tmp=$(mktemp -d)
        for i in 1 2 3; do echo "file $i" > "$tmp/file$i.txt"; done
        aws s3 sync "$tmp" "s3://$BUCKET/"
        aws s3 ls "s3://$BUCKET/"
        rm -rf "$tmp"
        ```
        **Spiegazione.** Fuori da `us-east-1`, `create-bucket` richiede `LocationConstraint` con la regione. Il `trap`
        garantisce che il bucket venga eliminato anche se un comando fallisce: nessuna risorsa dimenticata.

    === "Google Cloud"

        ```bash
        #!/usr/bin/env bash
        set -euo pipefail
        REGIONE=${REGIONE:-europe-west8}
        BUCKET="academy-lab-$(date +%Y%m%d)-$RANDOM"

        pulizia() { gcloud storage rm --recursive "gs://$BUCKET" --quiet > /dev/null 2>&1 || true; echo "Eliminato $BUCKET"; }
        trap pulizia EXIT

        gcloud storage buckets create "gs://$BUCKET" --location="$REGIONE" \
          --uniform-bucket-level-access --public-access-prevention
        gcloud storage buckets update "gs://$BUCKET" --update-labels=progetto=academy

        tmp=$(mktemp -d)
        for i in 1 2 3; do echo "file $i" > "$tmp/file$i.txt"; done
        gcloud storage rsync "$tmp" "gs://$BUCKET" --recursive
        gcloud storage ls "gs://$BUCKET/"
        rm -rf "$tmp"
        ```
        **Spiegazione.** `gcloud storage rm --recursive gs://NOME` elimina gli oggetti e poi il bucket stesso.
        Le opzioni di creazione impongono subito le protezioni della lezione 4.6.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è il modo consigliato per autenticare la AWS CLI di una persona?","options":["aws configure con le access key dell'utente root","aws configure sso con IAM Identity Center","Copiare le credenziali di un collega","Salvare la password nel codice"],"answer":[1],"explanation":"Credenziali temporanee, nessuna chiave permanente."},
 {"q":"Quale comando mostra con quale identità sta lavorando la AWS CLI?","options":["aws whoami","aws sts get-caller-identity","aws iam get-user --me","aws configure list-users"],"answer":[1],"explanation":"Restituisce account, ID e ARN."},
 {"q":"Come estrai solo i nomi delle VM in esecuzione con gcloud?","options":["gcloud compute instances list --filter=\"status=RUNNING\" --format=\"value(name)\"","gcloud compute instances list | grep RUNNING","gcloud vm names","gcloud compute list --running"],"answer":[0],"explanation":"--filter seleziona, --format sceglie i campi."},
 {"q":"Quali comandi autenticano strumenti diversi su Google Cloud?","options":["gcloud auth login, per la CLI","gcloud auth application-default login, per librerie e strumenti come Terraform","gcloud config set project, per l'autenticazione","gsutil login","gcloud services enable"],"answer":[0,1],"explanation":"config set project sceglie il progetto, non autentica."},
 {"q":"Su AWS, che differenza c'è tra --filters e --query?","options":["Nessuna","--filters filtra sul server, --query sul client con JMESPath","--query filtra sul server","--filters cambia il formato"],"answer":[1],"explanation":"Filtrare sul server riduce i dati trasferiti."}
]}
</script>

## Fonti

- [AWS CLI User Guide: Configuring IAM Identity Center authentication](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html)
- [AWS CLI: Filtering output (--query)](https://docs.aws.amazon.com/cli/latest/userguide/cli-usage-filter.html)
- [AWS CLI Command Reference](https://docs.aws.amazon.com/cli/latest/reference/)
- [JMESPath tutorial](https://jmespath.org/tutorial.html)
- [Google Cloud CLI: Initializing](https://docs.cloud.google.com/sdk/docs/initializing) e [Managing configurations](https://docs.cloud.google.com/sdk/docs/configurations)
- [gcloud topic formats](https://docs.cloud.google.com/sdk/gcloud/reference/topic/formats) e [filters](https://docs.cloud.google.com/sdk/gcloud/reference/topic/filters)
- [gcloud storage](https://docs.cloud.google.com/sdk/gcloud/reference/storage)
- [Google Cloud: Application Default Credentials](https://docs.cloud.google.com/docs/authentication/application-default-credentials)
