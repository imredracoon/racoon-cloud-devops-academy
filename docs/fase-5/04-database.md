# 5.4 Database gestiti

!!! obiettivi "Obiettivi"
    - Configurare un database relazionale gestito con alta disponibilità, backup e repliche di lettura.
    - Distinguere alta disponibilità (standby) e scalabilità delle letture (repliche).
    - Collegare in modo privato e sicuro un'applicazione al database.
    - Conoscere i database ad alte prestazioni e distribuiti globalmente.
    - Progettare con un database NoSQL: chiavi, capacità, repliche globali.
    - Ridurre il carico sul database con una cache.

## Il database relazionale gestito

=== "AWS"

    **Amazon RDS** gestisce PostgreSQL, MySQL, MariaDB, Oracle, SQL Server e Db2. Le scelte principali alla creazione:

    | Scelta | Opzioni |
    |---|---|
    | Classe dell'istanza | `db.t4g.micro` per le prove, `db.m7g`/`db.r7g` in produzione |
    | Storage | gp3 o io2, con **autoscaling dello storage** |
    | Disponibilità | **Single-AZ**; **Multi-AZ instance** (uno standby sincrono non leggibile in un'altra AZ, failover automatico in uno-due minuti); **Multi-AZ DB cluster** (due standby leggibili, failover più rapido) |
    | Rete | un **DB subnet group** con subnet private di almeno due AZ; accesso pubblico **no**; security group |
    | Crittografia | KMS, da scegliere **alla creazione** (un database non cifrato si cifra solo passando da uno snapshot) |
    | Backup | automatici, conservati da 1 a 35 giorni, con **ripristino a un istante** (*point-in-time recovery*); gli snapshot manuali restano finché non li cancelli |

    Altre funzioni:

    - **Repliche di lettura** (*read replica*): copie **asincrone**, anche in altre regioni, per distribuire le letture o
      per il disaster recovery (una replica si può **promuovere** a database autonomo);
    - **RDS Proxy**: un pool di connessioni gestito tra applicazione e database, utilissimo con Lambda e con molte
      istanze dell'applicazione; rende anche il failover più rapido per l'applicazione;
    - **autenticazione IAM**: l'applicazione si autentica con un token IAM invece che con una password;
    - **parameter group**: la configurazione del motore (l'equivalente di `postgresql.conf`).

=== "Google Cloud"

    **Cloud SQL** gestisce PostgreSQL, MySQL e SQL Server.

    | Scelta | Opzioni |
    |---|---|
    | Edizione | **Enterprise** o **Enterprise Plus** (prestazioni e disponibilità superiori) |
    | Tipo di macchina | da shared-core per le prove a dedicate in produzione |
    | Disponibilità | **zonale** oppure **alta disponibilità (regionale)**: un'istanza principale e uno standby in un'altra zona, con disco replicato in sincrono e failover automatico |
    | Rete | **IP privato** tramite *private services access* (consigliato) e/o IP pubblico con reti autorizzate |
    | Crittografia | sempre attiva; CMEK opzionale |
    | Backup | automatici e on demand, **ripristino a un istante** con i log delle transazioni |

    Altre funzioni:

    - **Repliche di lettura**, anche tra regioni, promuovibili;
    - **Cloud SQL Auth Proxy** e **connettori** per linguaggio: connessioni cifrate e autorizzate con IAM, senza gestire
      certificati né elenchi di IP;
    - **autenticazione IAM** dei database per utenti e service account;
    - **flag del database** per la configurazione del motore.

### Alta disponibilità contro repliche di lettura

È una distinzione chiesta in quasi tutti gli esami:

| | Standby (Multi-AZ / HA) | Replica di lettura |
|---|---|---|
| Scopo | **disponibilità**: sopravvivere al guasto | **prestazioni**: distribuire le letture |
| Replica | sincrona | asincrona (può essere indietro di qualche secondo) |
| Leggibile | no (salvo Multi-AZ DB cluster su RDS) | sì |
| Failover | automatico, stesso nome DNS | manuale (promozione) |
| In un'altra regione | no | sì, per il disaster recovery |

## Database ad alte prestazioni e globali

=== "AWS"

    **Amazon Aurora** è un database relazionale compatibile con PostgreSQL e MySQL, progettato per il cloud:

    - lo storage conserva **sei copie dei dati su tre AZ** e cresce automaticamente;
    - fino a **15 repliche di lettura** con ritardo molto basso, e un *reader endpoint* che le bilancia;
    - **Aurora Serverless v2** adatta la capacità al carico in frazioni di secondo;
    - **Aurora Global Database**: repliche in altre regioni con ritardo tipicamente inferiore a un secondo, per il
      disaster recovery tra regioni;
    - **Aurora DSQL**: relazionale distribuito e serverless, attivo-attivo su più regioni.

=== "Google Cloud"

    - **AlloyDB for PostgreSQL**: compatibile con PostgreSQL, ad alte prestazioni, con un motore colonnare per le query
      analitiche e alta disponibilità integrata.
    - **Spanner**: relazionale **distribuito globalmente**, con transazioni coerenti tra regioni e disponibilità fino al
      99,999% nelle configurazioni multi-regione. Scala orizzontalmente aggiungendo capacità di calcolo.

## NoSQL: progettare per l'accesso

=== "AWS"

    **DynamoDB** è un database chiave-valore e documentale serverless:

    - ogni elemento ha una **chiave primaria**: una **chiave di partizione** e, facoltativa, una **chiave di
      ordinamento**. Si progetta partendo dalle **query** che servono, non dalla struttura dei dati;
    - **indici secondari** (globali e locali) per interrogare per altri attributi;
    - **modalità di capacità**: *on-demand* (paghi per richiesta, nessuna pianificazione) o *provisioned* (capacità di
      lettura e scrittura definita, con autoscaling, più economica per carichi stabili);
    - **DAX**: una cache in memoria specifica per DynamoDB, con latenze di microsecondi;
    - **DynamoDB Streams**: il flusso delle modifiche, per attivare funzioni Lambda;
    - **Global Tables**: replica multi-regione attiva-attiva;
    - **TTL** per scadenza automatica degli elementi, **PITR** fino a 35 giorni.

=== "Google Cloud"

    - **Firestore**: documenti organizzati in collezioni, serverless, con aggiornamenti in tempo reale verso le app
      mobili e web, transazioni e configurazioni multi-regione.
    - **Bigtable**: colonne larghe, per enormi volumi (serie temporali, IoT, analisi) con latenza di pochi millisecondi;
      compatibile con l'API di HBase.

## Le cache

Molte richieste leggono gli stessi dati. Una **cache in memoria** davanti al database riduce latenza e carico:

- **ElastiCache** (Valkey, Redis OSS, Memcached) e **MemoryDB** (durevole) su AWS;
- **Memorystore** (Valkey, Redis, Memcached) su Google Cloud.

La strategia più comune è il **cache-aside** (*lazy loading*): l'applicazione cerca nella cache; se non trova, legge
dal database e salva il risultato in cache con una **scadenza** (TTL). Il problema da gestire è l'**invalidazione**:
quando un dato cambia, la copia in cache va aggiornata o cancellata.

## Laboratorio

!!! laboratorio "Laboratorio · Un database condiviso per la Tasks API"
    Riprendi l'architettura del laboratorio [5.2](02-ha-scaling.md#laboratorio) e risolvi il problema dei dati diversi su ogni istanza.

    === "AWS"

        1. **RDS** → **Subnet groups** → crea `tasks-db-subnets` con le subnet `dati-a` e `dati-b`.
        2. **Create database**: PostgreSQL 17, modello **Free tier** o **Dev/Test**, `db.t4g.micro`, identificativo
           `tasks-db`, nome iniziale del database `tasks`, utente `tasks`, password generata e annotata; subnet group
           `tasks-db-subnets`, **Public access: No**, security group `db-sg`; backup 1 giorno. (Multi-AZ: provalo solo se
           i crediti lo consentono: raddoppia il costo.)
        3. Quando è disponibile, copia l'**endpoint**.
        4. Crea una nuova versione del launch template `tasks-lt` con lo script di avvio modificato:

            ```bash
            docker run -d --restart unless-stopped -p 8000:8000 --name api \
              -e DATABASE_URL='postgresql+psycopg://tasks:PASSWORD@ENDPOINT:5432/tasks' \
              ghcr.io/TUO-UTENTE/tasks-api:1.1.0
            ```

        5. Imposta la nuova versione come predefinita nell'ASG e avvia un **instance refresh**.
        6. Ripeti la prova della [lezione 5.2](02-ha-scaling.md): ora `POST` e `GET` danno risultati coerenti, qualunque istanza risponda.

        **Pulizia.** Elimina il database (senza snapshot finale, per un laboratorio), il subnet group e le risorse della
        [lezione 5.2](02-ha-scaling.md).

    === "Google Cloud"

        1. **Rete VPC** `tasks` → **Accesso privato ai servizi**: alloca un intervallo IP (per esempio `/20`) e crea la
           connessione privata.
        2. **Cloud SQL** → **Crea istanza** PostgreSQL 17, edizione **Enterprise**, preimpostazione **Sandbox**,
           disponibilità **zona singola**, regione del corso; **IP privato** sulla rete `tasks`, IP pubblico
           disattivato. Crea il database `tasks` e l'utente `tasks`.
        3. Annota l'**IP privato** dell'istanza.
        4. Non serve una nuova regola firewall: le regole implicite permettono il traffico in uscita, e Cloud SQL è
           raggiungibile attraverso la connessione privata.
        5. Crea un nuovo modello di istanza con lo script di avvio che passa
           `-e DATABASE_URL='postgresql+psycopg://tasks:PASSWORD@IP-PRIVATO:5432/tasks'`, e aggiorna il MIG al nuovo
           modello (**Aggiorna le VM** → sostituzione progressiva).
        6. Ripeti la prova della [lezione 5.2](02-ha-scaling.md): i dati sono coerenti.

        **Pulizia.** Elimina l'istanza Cloud SQL (disattiva prima la protezione dall'eliminazione), la connessione
        privata e l'intervallo allocato, e le risorse della [lezione 5.2](02-ha-scaling.md).

??? success "Spiegazione e il prossimo problema"
    Ora tutte le istanze condividono lo stesso database: l'applicazione è davvero **senza stato** e può scalare. Ma
    guarda lo script di avvio: la **password del database è scritta in chiaro** nel launch template o nel modello di
    istanza, visibile a chiunque possa leggerli. Nella [lezione 5.6](06-sicurezza.md) la sposterai in un gestore di segreti.

## Esercizi

!!! esercizio "Esercizio · Requisiti e scelte"
    1. Un e-commerce ha molte più letture che scritture e il database è al limite della CPU.
    2. Un'applicazione bancaria deve continuare a funzionare se un'intera regione AWS diventa irraggiungibile, perdendo
       al massimo un secondo di dati.
    3. Mille funzioni Lambda aprono connessioni al database PostgreSQL e lo saturano.
    4. Un'app mobile deve sincronizzare in tempo reale i dati tra i dispositivi degli utenti.
    5. Il database di produzione non è cifrato e una nuova norma lo richiede.

??? success "Soluzione"
    1. **Repliche di lettura** (e, se le stesse letture si ripetono, una **cache**).
    2. **Aurora Global Database** (RPO tipicamente sotto il secondo, promozione della regione secondaria); su Google
       Cloud, **Spanner** multi-regione.
    3. **RDS Proxy**, che riusa un numero limitato di connessioni.
    4. **Firestore** (aggiornamenti in tempo reale); su AWS si userebbe AWS AppSync con DynamoDB.
    5. Crea uno **snapshot**, **copialo cifrandolo** con una chiave KMS, ripristina un nuovo database dallo snapshot
       cifrato e sposta l'applicazione. Un'istanza RDS esistente non si può cifrare "al volo".

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è lo scopo principale di RDS Multi-AZ (istanza con standby)?","options":["Aumentare le prestazioni in lettura","Alta disponibilità con failover automatico","Ridurre i costi","Replicare in un'altra regione"],"answer":[1],"explanation":"Per le letture si usano le repliche di lettura."},
 {"q":"Come si replicano i dati verso una replica di lettura RDS?","options":["In modo sincrono","In modo asincrono","Con backup notturni","Non si replicano"],"answer":[1],"explanation":"Per questo può essere leggermente indietro rispetto al primario."},
 {"q":"Quante copie dei dati conserva lo storage di Aurora?","options":["2 in 2 AZ","3 in 1 AZ","6 in 3 AZ","1 con backup"],"answer":[2],"explanation":"Sei copie su tre zone di disponibilità."},
 {"q":"Quale servizio gestisce un pool di connessioni tra molte funzioni Lambda e un database RDS?","options":["DAX","RDS Proxy","ElastiCache","Read replica"],"answer":[1],"explanation":"Evita di saturare il database con connessioni."},
 {"q":"Quali servizi offrono un database relazionale coerente distribuito su più regioni?","options":["Spanner","Aurora DSQL","Cloud SQL zonale","ElastiCache","DynamoDB DAX"],"answer":[0,1],"explanation":"Aurora Global Database è invece un primario con repliche secondarie."},
 {"q":"In DynamoDB, da cosa si parte per progettare la chiave primaria?","options":["Dalla struttura delle tabelle relazionali","Dalle query che l'applicazione deve eseguire","Dal numero di colonne","Dalla regione"],"answer":[1],"explanation":"Il modello NoSQL si progetta per i modelli di accesso."},
 {"q":"Come colleghi in modo privato un'applicazione su Compute Engine a Cloud SQL?","options":["IP pubblico con 0.0.0.0/0 autorizzato","IP privato tramite accesso privato ai servizi","Con Cloud CDN","Con un bucket"],"answer":[1],"explanation":"In alternativa, Cloud SQL Auth Proxy o i connettori."}
]}
</script>

## Fonti

- [Amazon RDS User Guide: Multi-AZ deployments](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.MultiAZ.html)
- [Amazon RDS: Working with read replicas](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ReadRepl.html)
- [Amazon RDS: Backing up and restoring](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_CommonTasks.BackupRestore.html)
- [Amazon RDS Proxy](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html)
- [Amazon Aurora: Overview](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/CHAP_AuroraOverview.html)
- [Amazon DynamoDB Developer Guide: Core components](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html)
- [Cloud SQL for PostgreSQL: High availability](https://docs.cloud.google.com/sql/docs/postgres/high-availability) e [Private IP](https://docs.cloud.google.com/sql/docs/postgres/private-ip)
- [Spanner overview](https://docs.cloud.google.com/spanner/docs/overview), [AlloyDB overview](https://docs.cloud.google.com/alloydb/docs/overview)
- [AWS Whitepaper: Database caching strategies using Redis](https://docs.aws.amazon.com/whitepapers/latest/database-caching-strategies-using-redis/caching-patterns.html)
