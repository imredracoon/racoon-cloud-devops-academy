# 4.8 Database, analisi e AI

!!! obiettivi "Obiettivi"
    - Distinguere database relazionali e non relazionali e scegliere il tipo adatto ai dati.
    - Conoscere i database gestiti di AWS e Google Cloud e cosa significa "gestito".
    - Conoscere i servizi per l'analisi dei dati: data warehouse, elaborazione, streaming.
    - Orientarsi tra i servizi di intelligenza artificiale e machine learning.
    - Sapere quali strumenti aiutano a migrare database e applicazioni verso il cloud.

Questa lezione è una panoramica: molte domande degli esami Cloud Practitioner e Digital Leader chiedono di
**scegliere il servizio giusto** per un caso d'uso. Nella Fase 5 userai davvero un database gestito con la Tasks API.

## Relazionali e non relazionali

| | Relazionale (SQL) | Non relazionale (NoSQL) |
|---|---|---|
| Struttura | tabelle con righe e colonne, schema definito | documenti, chiave-valore, colonne larghe, grafi |
| Linguaggio | SQL | API specifiche |
| Relazioni e transazioni | join, transazioni ACID | limitate, dipende dal database |
| Scalabilità | soprattutto verticale (macchina più grande) e repliche di lettura | orizzontale, su molti nodi |
| Quando | dati strutturati con relazioni: ordini, clienti, fatture | grandi volumi, schema flessibile, latenza costante a qualsiasi scala |

## Cosa significa "gestito"

Installare PostgreSQL su una VM è possibile, ma tocca a te: installazione, aggiornamenti, backup, replica, failover,
monitoraggio, scalabilità. Un **database gestito** fa tutto questo al posto tuo; tu ti occupi di schema, query,
utenti, dimensionamento e configurazione di rete e sicurezza. È la responsabilità condivisa della lezione 4.1
applicata ai database.

## I servizi

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Relazionale gestito (PostgreSQL, MySQL, SQL Server...) | **Amazon RDS** | **Cloud SQL** |
| Relazionale ad alte prestazioni, compatibile PostgreSQL/MySQL | **Amazon Aurora** | **AlloyDB for PostgreSQL** |
| Relazionale distribuito globalmente, scalabile orizzontalmente | **Aurora DSQL** | **Spanner** |
| Chiave-valore / documenti serverless | **DynamoDB** | **Firestore** |
| Colonne larghe, enormi volumi, bassa latenza | **Keyspaces** (compatibile Cassandra) | **Bigtable** |
| Documenti compatibile MongoDB | **DocumentDB** | Firestore con compatibilità MongoDB |
| Cache in memoria (Redis, Valkey, Memcached) | **ElastiCache**, **MemoryDB** | **Memorystore** |
| Grafi | **Neptune** | Spanner Graph |

Funzioni ricorrenti dei database relazionali gestiti:

- **alta disponibilità** con una replica sincrona in standby in un'altra zona e failover automatico (RDS **Multi-AZ**,
  Cloud SQL **configurazione ad alta disponibilità**);
- **repliche di lettura** per distribuire le query di sola lettura;
- **backup automatici** con ripristino a un istante preciso (*point-in-time recovery*);
- **crittografia** a riposo e in transito.

## Analisi dei dati

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Data warehouse: SQL su enormi quantità di dati | **Redshift** | **BigQuery** (serverless) |
| SQL direttamente sui file nello storage a oggetti | **Athena** | BigQuery (tabelle esterne, BigLake) |
| Elaborazione distribuita (Spark, Hadoop) | **EMR** | **Dataproc** |
| ETL e catalogo dei dati | **Glue** | **Dataflow** (Apache Beam), Dataplex |
| Streaming di eventi | **Kinesis**, **MSK** (Kafka) | **Pub/Sub**, Managed Service for Apache Kafka |
| Dashboard e business intelligence | **QuickSight** | **Looker**, Looker Studio |
| Ricerca e analisi dei log | **OpenSearch Service** | — |

BigQuery è uno dei servizi più caratteristici di Google Cloud: un data warehouse **serverless**, in cui carichi i dati
e scrivi SQL senza gestire alcun server; si paga per i dati analizzati o per la capacità riservata.

## Intelligenza artificiale e machine learning

I servizi si dividono in tre livelli, dal più pronto al più personalizzabile:

| Livello | Cosa offre | AWS | Google Cloud |
|---|---|---|---|
| **Servizi AI pronti** | API per compiti specifici, senza competenze di ML | **Rekognition** (immagini), **Transcribe** (voce→testo), **Polly** (testo→voce), **Translate**, **Comprehend** (analisi del testo), **Textract** (documenti), **Lex** (chatbot) | **Vision AI**, **Speech-to-Text**, **Text-to-Speech**, **Translation AI**, **Natural Language AI**, **Document AI** |
| **AI generativa e modelli di base** | modelli linguistici e multimodali tramite API, assistenti | **Amazon Bedrock** (modelli di più fornitori), **Amazon Q** (assistente) | **Vertex AI** con i modelli **Gemini** e il Model Garden, **Gemini for Google Cloud** (assistente) |
| **Piattaforma di ML** | addestrare, distribuire e gestire modelli propri | **Amazon SageMaker AI** | **Vertex AI** |

Concetti da conoscere per gli esami di base:

- **intelligenza artificiale** è il campo generale; il **machine learning** è la parte in cui i sistemi imparano dai
  dati; il **deep learning** usa reti neurali profonde; l'**AI generativa** produce contenuti nuovi (testo, immagini,
  codice) a partire da un *prompt*;
- un **modello di base** (*foundation model*) è un grande modello pre-addestrato che si può usare così com'è, guidare
  con il prompt, arricchire con dati aziendali (*RAG*, *retrieval-augmented generation*) o adattare (*fine-tuning*);
- l'**AI responsabile**: equità, trasparenza, privacy, sicurezza e supervisione umana. I dati aziendali usati con i
  servizi gestiti restano del cliente: verifica sempre le condizioni d'uso del servizio.

## Migrare verso il cloud

Le strategie di migrazione sono spesso riassunte in "7 R": **retire** (dismettere), **retain** (lasciare dove sta),
**rehost** (*lift and shift*: spostare le VM così come sono), **relocate** (spostare la piattaforma di virtualizzazione),
**replatform** (piccoli adattamenti, per esempio passare a un database gestito), **repurchase** (passare a un SaaS),
**refactor** (riprogettare per il cloud).

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Migrare database, anche tra motori diversi | **Database Migration Service** (DMS), Schema Conversion | **Database Migration Service** |
| Migrare server e VM | **Application Migration Service** | **Migrate to Virtual Machines** |
| Pianificare e valutare | **Migration Hub**, **Migration Evaluator** | **Migration Center** |
| Quadro di riferimento per l'adozione | **Cloud Adoption Framework** (CAF) | **Cloud Adoption Framework** di Google Cloud |

## Esercizi

!!! esercizio "Esercizio · Scegli il servizio"
    Per ogni caso indica il servizio AWS e quello Google Cloud:

    1. il database PostgreSQL della Tasks API, con failover automatico in un'altra zona;
    2. il carrello di un e-commerce con milioni di utenti, che richiede latenza di pochi millisecondi a qualsiasi scala;
    3. analizzare con SQL 5 anni di dati di vendita (diversi terabyte), senza gestire server;
    4. ridurre il carico sul database tenendo in memoria i risultati delle query frequenti;
    5. estrarre automaticamente importi e date da migliaia di fatture scansionate;
    6. aggiungere a un'applicazione un assistente che risponde alle domande sui documenti aziendali usando un modello
       linguistico.

??? success "Soluzione"
    1. **RDS for PostgreSQL Multi-AZ** (o Aurora); **Cloud SQL for PostgreSQL** in alta disponibilità (o AlloyDB).
    2. **DynamoDB**; **Firestore** (o Bigtable per volumi enormi di dati a colonne).
    3. **Redshift** (o Athena sui file in S3); **BigQuery**.
    4. **ElastiCache**; **Memorystore**.
    5. **Textract**; **Document AI**.
    6. **Amazon Bedrock** con una knowledge base (RAG); **Vertex AI** con Gemini e la ricerca sui documenti (RAG).

    **Spiegazione.** Negli esami, le parole chiave guidano la scelta: "relazionale" e "transazioni" portano a RDS e
    Cloud SQL; "millisecondi a qualsiasi scala" e "chiave-valore" a DynamoDB; "data warehouse" e "analisi" a Redshift e
    BigQuery; "in memoria" a ElastiCache e Memorystore; "modelli di base" a Bedrock e Vertex AI.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale servizio è un data warehouse serverless di Google Cloud?","options":["Cloud SQL","BigQuery","Bigtable","Firestore"],"answer":[1],"explanation":"Su AWS il data warehouse è Redshift."},
 {"q":"Quale servizio AWS offre un database NoSQL chiave-valore serverless?","options":["RDS","Aurora","DynamoDB","Redshift"],"answer":[2],"explanation":"Latenza di pochi millisecondi a qualsiasi scala."},
 {"q":"Cosa fa RDS Multi-AZ?","options":["Distribuisce le letture su più regioni","Mantiene una replica sincrona in standby in un'altra AZ con failover automatico","Esegue backup su nastro","Rende il database pubblico"],"answer":[1],"explanation":"Per le letture si usano le repliche di lettura."},
 {"q":"Quali servizi permettono di usare modelli di AI generativa tramite API?","options":["Amazon Bedrock","Vertex AI","Amazon Redshift","Cloud DNS","AWS Glue"],"answer":[0,1],"explanation":"Bedrock offre modelli di più fornitori; Vertex AI include Gemini e il Model Garden."},
 {"q":"Spostare una VM nel cloud senza modificarla è la strategia di migrazione:","options":["Refactor","Rehost (lift and shift)","Repurchase","Retire"],"answer":[1],"explanation":"È la più rapida, ma non sfrutta i servizi gestiti."},
 {"q":"Cosa gestisce il fornitore in un database gestito come Cloud SQL?","options":["Lo schema delle tabelle","Le query","Patch, backup automatici e failover","Gli utenti dell'applicazione"],"answer":[2],"explanation":"Schema, query e dimensionamento restano al cliente."}
]}
</script>

## Fonti

- [AWS: Choosing an AWS database service](https://docs.aws.amazon.com/decision-guides/latest/databases-on-aws-how-to-choose/databases-on-aws-how-to-choose.html)
- [Amazon RDS: Multi-AZ deployments](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.MultiAZ.html)
- [Google Cloud: Google Cloud databases](https://cloud.google.com/products/databases)
- [Cloud SQL: About high availability](https://docs.cloud.google.com/sql/docs/postgres/high-availability)
- [BigQuery overview](https://docs.cloud.google.com/bigquery/docs/introduction)
- [AWS: Machine Learning and AI services](https://aws.amazon.com/ai/services/) e [Amazon Bedrock](https://aws.amazon.com/bedrock/)
- [Google Cloud: Vertex AI](https://cloud.google.com/vertex-ai)
- [AWS Prescriptive Guidance: Migration strategies (7 Rs)](https://docs.aws.amazon.com/prescriptive-guidance/latest/large-migration-guide/migration-strategies.html)
