# 5.9 Container e serverless

!!! obiettivi "Obiettivi"
    - Eseguire container gestiti senza amministrare server: Cloud Run, ECS con Fargate.
    - Pubblicare immagini nei registry dei cloud.
    - Scrivere e collegare funzioni serverless agli eventi.
    - Esporre API con un gateway gestito.
    - Scegliere tra VM, container gestiti, Kubernetes e funzioni.

## I container gestiti

Nella Fase 3 hai containerizzato la Tasks API; nelle lezioni precedenti l'hai eseguita su VM con Docker installato
da uno script. I servizi di **container gestiti** eliminano le VM: dai l'immagine, la configurazione e le risorse, e il
fornitore la esegue, la scala e la sostituisce se si guasta.

=== "AWS"

    **Amazon ECS** (*Elastic Container Service*) è l'orchestratore di container di AWS:

    - una **task definition** descrive uno o più container: immagine, CPU, memoria, porte, variabili, segreti (da
      Secrets Manager o Parameter Store), log verso CloudWatch, **ruolo del task** (i permessi dell'applicazione) e
      **ruolo di esecuzione** (i permessi di ECS per scaricare l'immagine e leggere i segreti);
    - un **servizio** mantiene in esecuzione N copie del task, dietro un load balancer, con autoscaling e rilasci
      progressivi;
    - il **tipo di avvio**: **Fargate** (serverless: nessuna istanza da gestire, paghi vCPU e memoria al secondo) o
      **EC2** (gestisci tu le istanze del cluster).

    Le immagini private stanno in **Amazon ECR**, con scansione delle vulnerabilità, regole di ciclo di vita e tag
    immutabili.

    Per Kubernetes gestito c'è **EKS** (Fase 8); per un'esperienza più semplice, **App Runner**.

=== "Google Cloud"

    **Cloud Run** esegue container che rispondono a richieste HTTP (o eventi), ed è il servizio più semplice per
    pubblicare un'API containerizzata:

    - scala automaticamente in base alle richieste, **anche a zero** (nessun costo senza traffico);
    - ogni istanza gestisce più richieste in parallelo (**concorrenza**, predefinita 80);
    - HTTPS con un URL `*.run.app` incluso; dominio personalizzato tramite load balancer;
    - variabili d'ambiente e **segreti da Secret Manager**, **service account** dedicato;
    - collegamento privato alla VPC (**Direct VPC egress**) per raggiungere Cloud SQL o altri servizi privati;
    - i **job** di Cloud Run eseguono container fino al completamento (elaborazioni batch);
    - **revisioni** immutabili, con suddivisione del traffico tra revisioni per i rilasci graduali.

    Le immagini private stanno in **Artifact Registry**, con scansione delle vulnerabilità.

    Per Kubernetes gestito c'è **GKE**, anche nella modalità **Autopilot** in cui Google gestisce i nodi (Fase 8).

## Le funzioni

Una **funzione** è un frammento di codice che il fornitore esegue in risposta a un **evento**:

=== "AWS"

    **AWS Lambda**:

    - linguaggi supportati (Python, Node.js, Java, .NET, Go, Ruby) o immagini container;
    - memoria da 128 MB a 10 GB (la CPU cresce con la memoria), **durata massima 15 minuti**;
    - **eventi**: richieste HTTP (API Gateway, URL della funzione), oggetti S3, messaggi SQS, flussi DynamoDB e
      Kinesis, regole EventBridge, pianificazioni;
    - **ruolo di esecuzione** per i permessi;
    - si paga per numero di richieste e per durata × memoria;
    - **avvio a freddo** (*cold start*): la prima richiesta su un'istanza nuova è più lenta; la **concorrenza
      riservata** e la **concorrenza predisposta** controllano capacità e latenza.

=== "Google Cloud"

    **Cloud Run functions** (prima Cloud Functions): scrivi una funzione, Google la impacchetta come container e la
    esegue su Cloud Run. Eventi HTTP o tramite **Eventarc** (Cloud Storage, Pub/Sub, audit log e molti altri).

## API gestite

**Amazon API Gateway** (API REST, HTTP e WebSocket) e **API Gateway** / **Apigee** su Google Cloud offrono davanti a
funzioni e servizi: autenticazione (anche con Cognito o JWT), limitazione delle richieste (*throttling*), quote per
cliente, cache, trasformazioni, versioni e ambienti.

## Come scegliere

| Esigenza | Scelta tipica |
|---|---|
| Controllo completo del sistema operativo, software legacy | VM (EC2, Compute Engine) |
| Un'API o un sito in container, team piccolo, traffico variabile | **Cloud Run**, **App Runner**, **ECS con Fargate** |
| Molti servizi, team con competenze Kubernetes, portabilità tra cloud | **EKS**, **GKE** |
| Reazione a eventi, elaborazioni brevi, picchi improvvisi | **Lambda**, **Cloud Run functions** |
| Elaborazioni batch lunghe | AWS Batch, Cloud Run jobs, Batch |

## Laboratorio

!!! laboratorio "Laboratorio · La Tasks API senza server"
    === "Google Cloud"

        Cloud Run è il modo più rapido. Usa il database Cloud SQL e il segreto delle lezioni 5.4 e 5.6.

        ```bash
        gcloud services enable run.googleapis.com artifactregistry.googleapis.com
        gcloud artifacts repositories create tasks --repository-format=docker --location=REGIONE
        gcloud auth configure-docker REGIONE-docker.pkg.dev

        IMG=REGIONE-docker.pkg.dev/ID-PROGETTO/tasks/tasks-api:1.1.0
        docker pull ghcr.io/TUO-UTENTE/tasks-api:1.1.0
        docker tag ghcr.io/TUO-UTENTE/tasks-api:1.1.0 "$IMG"
        docker push "$IMG"

        gcloud run deploy tasks-api --image "$IMG" --region REGIONE \
          --service-account tasks-app@ID-PROGETTO.iam.gserviceaccount.com \
          --set-secrets DATABASE_URL=tasks-database-url:latest \
          --network tasks --subnet app --vpc-egress private-ranges-only \
          --port 8000 --allow-unauthenticated --max-instances 3
        ```

        Apri l'URL restituito, aggiungendo `/health`. Poi:

        ```bash
        gcloud run services describe tasks-api --region REGIONE --format='value(status.url)'
        gcloud run services logs read tasks-api --region REGIONE --limit 20
        ```

        **Pulizia.** `gcloud run services delete tasks-api --region REGIONE` e il repository di Artifact Registry.

    === "AWS"

        Con **ECS e Fargate**, nella rete della lezione 5.1:

        1. **ECR** → crea il repository `tasks-api`. Nella pagina, **View push commands**: esegui login, tag e push
           dell'immagine `ghcr.io/TUO-UTENTE/tasks-api:1.1.0`.
        2. **ECS** → **Clusters** → **Create cluster** `tasks` (solo **Fargate**).
        3. **Task definitions** → **Create**: Fargate, 0,25 vCPU, 0,5 GB; container `api` con l'immagine ECR, porta 8000;
           variabile d'ambiente `DATABASE_URL` di tipo **ValueFrom** con l'ARN del segreto `tasks/database-url`
           seguito da `:DATABASE_URL::` (estrae la chiave JSON); log su CloudWatch. Il **ruolo di esecuzione** deve poter
           leggere il segreto (aggiungi la policy della lezione 5.6 a `ecsTaskExecutionRole`).
        4. **Create service**: 2 task, subnet `app-a` e `app-b`, security group `app-sg`, senza IP pubblico, con un
           **nuovo Application Load Balancer** (o il target group esistente, di tipo **IP**), health check `/health`.
        5. Apri `http://DNS-DEL-ALB/health`.

        **Pulizia.** Imposta i task desiderati a 0, elimina il servizio, il cluster, il load balancer e il target group,
        il repository ECR (con le immagini).

??? success "Spiegazione"
    Lo stesso container della Fase 3 gira ora senza alcuna VM da gestire: niente patch del sistema operativo, niente
    Docker da installare, scalabilità e sostituzione automatiche. Su Cloud Run, senza traffico, le istanze scendono a
    zero. La configurazione (variabili, segreti, rete, identità) è separata dall'immagine: lo stesso artefatto va in
    test e in produzione con configurazioni diverse.

## Esercizi

!!! esercizio "Esercizio · Progetta la soluzione serverless"
    Un'agenzia riceve foto dagli utenti tramite un'app. Ogni foto deve essere ridimensionata in tre formati, i metadati
    salvati in un database e un'email inviata all'utente al termine. Il traffico è molto variabile: zero di notte,
    migliaia di foto all'ora durante gli eventi. Disegna l'architettura su AWS o Google Cloud.

??? success "Soluzione"
    === "AWS"

        L'app carica la foto su **S3** con un **URL prefirmato** (niente server in mezzo). L'evento di S3 va in una coda
        **SQS** (per assorbire i picchi e ritentare) che attiva una **Lambda** di ridimensionamento; questa salva le tre
        versioni su S3 e i metadati su **DynamoDB**, poi pubblica su **SNS** (o tramite **SES** per l'email). Una DLQ
        raccoglie le foto che falliscono.

    === "Google Cloud"

        L'app carica su **Cloud Storage** con un **URL firmato**. L'evento passa da **Eventarc** (o una notifica
        **Pub/Sub**) a una **Cloud Run function** che ridimensiona, salva su Cloud Storage e scrive i metadati in
        **Firestore**, poi pubblica un messaggio per il servizio che invia l'email. Un dead-letter topic raccoglie gli
        errori.

    **Spiegazione.** Nessun componente costa quando non ci sono foto, tutto scala con il traffico, e la coda protegge
    dai picchi e dagli errori temporanei.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale tipo di avvio di ECS evita di gestire istanze EC2?","options":["EC2","Fargate","Outposts","Spot"],"answer":[1],"explanation":"Fargate è il motore serverless per i container."},
 {"q":"Qual è la durata massima di un'esecuzione AWS Lambda?","options":["5 minuti","15 minuti","1 ora","Illimitata"],"answer":[1],"explanation":"Per lavori più lunghi si usano Fargate, Batch o Step Functions."},
 {"q":"Quale servizio Google Cloud esegue container HTTP scalando anche a zero istanze?","options":["Compute Engine","Cloud Run","GKE Standard","Bigtable"],"answer":[1],"explanation":"Si paga solo quando gestisce richieste (con la fatturazione predefinita)."},
 {"q":"In una task definition di ECS, a cosa serve il ruolo del task?","options":["A scaricare l'immagine","A dare i permessi AWS all'applicazione nel container","A creare il cluster","A gestire il DNS"],"answer":[1],"explanation":"Il ruolo di esecuzione serve invece a ECS per scaricare immagini e leggere i segreti."},
 {"q":"Dove conservi le immagini container private su AWS e Google Cloud?","options":["Amazon ECR","Artifact Registry","Amazon EBS","Cloud SQL","Amazon SQS"],"answer":[0,1],"explanation":"Entrambi offrono la scansione delle vulnerabilità."},
 {"q":"Quale servizio aggiunge autenticazione, throttling e quote davanti a funzioni Lambda?","options":["Amazon API Gateway","Amazon SQS","AWS Config","Amazon EFS"],"answer":[0],"explanation":"Su Google Cloud: API Gateway o Apigee."}
]}
</script>

## Fonti

- [Amazon ECS Developer Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html)
- [Amazon ECS: Pass Secrets Manager secrets](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/secrets-envvar-secrets-manager.html)
- [AWS Fargate](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html), [Amazon ECR](https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html)
- [AWS Lambda Developer Guide](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html) e [quote](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-limits.html)
- [Amazon API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html)
- [Cloud Run documentation](https://docs.cloud.google.com/run/docs) e [Secrets](https://docs.cloud.google.com/run/docs/configuring/services/secrets)
- [Cloud Run: Direct VPC egress](https://docs.cloud.google.com/run/docs/configuring/vpc-direct-vpc)
- [Artifact Registry](https://docs.cloud.google.com/artifact-registry/docs/overview)
