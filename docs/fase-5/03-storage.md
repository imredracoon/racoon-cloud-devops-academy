# 5.3 Storage avanzato

!!! obiettivi "Obiettivi"
    - Scegliere la crittografia giusta per lo storage a oggetti.
    - Replicare i dati tra regioni e proteggerli da cancellazioni con il blocco degli oggetti.
    - Gestire le prestazioni dello storage a oggetti per file grandi e molte richieste.
    - Automatizzare snapshot e backup di dischi, file system e database con un servizio centralizzato.
    - Conoscere gli storage per carichi ibridi e ad alte prestazioni.

## Crittografia dello storage a oggetti

=== "AWS"

    Dal 2023 S3 cifra **sempre** gli oggetti nuovi. Le opzioni:

    | Opzione | Chi gestisce la chiave | Quando |
    |---|---|---|
    | **SSE-S3** | AWS, in modo trasparente | predefinita, nessun costo |
    | **SSE-KMS** | una chiave in **AWS KMS**, con permessi e audit su CloudTrail | requisiti di conformità, controllo su chi può decifrare |
    | **DSSE-KMS** | doppia cifratura con KMS | normative che richiedono due livelli |
    | **SSE-C** | il cliente fornisce la chiave a ogni richiesta | casi particolari: AWS non conserva la chiave |
    | **Cifratura lato client** | il cliente cifra prima di caricare | il fornitore non vede mai i dati in chiaro |

    Con SSE-KMS, attiva le **S3 Bucket Keys** per ridurre le chiamate (e i costi) verso KMS.

=== "Google Cloud"

    Cloud Storage cifra **sempre** i dati. Le opzioni:

    | Opzione | Chi gestisce la chiave |
    |---|---|
    | **Chiavi gestite da Google** | Google (predefinita) |
    | **CMEK** (*customer-managed encryption keys*) | una chiave in **Cloud KMS**, gestita dal cliente |
    | **CSEK** (*customer-supplied encryption keys*) | il cliente fornisce la chiave a ogni richiesta |
    | **Cifratura lato client** | il cliente cifra prima di caricare |

La gestione delle chiavi con KMS è trattata nella [lezione 5.6](06-sicurezza.md).

## Replica e protezione

- **Replica tra regioni** per il disaster recovery o la vicinanza agli utenti: S3 **Cross-Region Replication** (CRR) e
  **Same-Region Replication** (SRR, per esempio verso un account separato per i backup); su Cloud Storage i bucket
  **dual-region** e **multi-region**, con la **replica turbo** opzionale (RPO di 15 minuti per i bucket dual-region).
- **Blocco degli oggetti**: S3 **Object Lock** in modalità *governance* (alcuni utenti autorizzati possono ancora
  rimuovere il blocco) o *compliance* (nessuno, nemmeno root, può cancellare prima della scadenza), più il *legal hold*;
  su Cloud Storage la **retention policy** del bucket, bloccabile in modo permanente con **Bucket Lock**, e le
  conservazioni sui singoli oggetti.
- **MFA Delete** su S3: richiede l'MFA per cancellare versioni o disattivare il versioning.
- **Accesso tra account**: le bucket policy (S3) o IAM sul bucket (Cloud Storage); su S3 gli **Access Points**
  semplificano la gestione di molti consumatori diversi dello stesso bucket.

## Prestazioni dello storage a oggetti

- **Caricamento in più parti** (*multipart upload*): per file grandi (consigliato oltre 100 MB su S3), il file viene
  diviso in parti caricate in parallelo e riprese in caso di errore. Le CLI lo fanno in automatico. Una regola del
  ciclo di vita dovrebbe eliminare i caricamenti incompleti dopo qualche giorno.
- **S3 Transfer Acceleration**: caricamenti da lontano più veloci, sfruttando i punti di presenza di CloudFront.
- **Prefissi**: S3 scala automaticamente a migliaia di richieste al secondo per prefisso; distribuire gli oggetti su
  più prefissi aumenta il totale.
- **Notifiche di eventi**: un oggetto creato o cancellato può attivare una funzione, una coda o EventBridge ([lezione 5.5](05-messaggistica.md));
  su Cloud Storage tramite **Pub/Sub** o **Eventarc**.
- **S3 Select** e query in loco con **Athena** o le tabelle esterne di **BigQuery** evitano di scaricare interi file.

## Backup centralizzati

=== "AWS"

    - Gli **snapshot EBS** sono incrementali; **Data Lifecycle Manager** li crea e li elimina secondo una politica.
    - Si possono **copiare tra regioni** per il disaster recovery, e un **AMI** è sostanzialmente uno snapshot avviabile.
    - **AWS Backup** centralizza i backup di EBS, EC2, RDS, Aurora, DynamoDB, EFS, FSx, S3 e altri: **piani di backup**
      con frequenza e conservazione, assegnati alle risorse tramite **etichette**, copie in altre regioni o account,
      **Vault Lock** per backup immutabili, report di conformità.

=== "Google Cloud"

    - Gli **snapshot** dei dischi sono incrementali e possono essere programmati con **pianificazioni di snapshot**
      (*snapshot schedules*) collegate ai dischi.
    - I **dischi regionali** replicano in sincrono su due zone, per il failover di una VM in un'altra zona.
    - **Backup and DR Service** centralizza i backup di VM di Compute Engine, database e altri carichi, con **backup
      vault** immutabili.

Una regola che vale ovunque: **un backup che non hai mai provato a ripristinare non è un backup**. Pianifica prove di
ripristino periodiche.

## File system e storage ibrido

=== "AWS"

    - **EFS**: modalità di prestazioni e di throughput (*Elastic* è la scelta predefinita), classi di storage
      (Standard, Infrequent Access, Archive) con ciclo di vita, opzione *One Zone* più economica.
    - **FSx for Windows File Server** (SMB, Active Directory), **FSx for Lustre** (calcolo ad alte prestazioni,
      collegabile a S3), **FSx for NetApp ONTAP**, **FSx for OpenZFS**.
    - **Storage Gateway** per collegare la sede: **S3 File Gateway** (NFS/SMB davanti a S3), **Volume Gateway**
      (volumi iSCSI con copie su AWS), **Tape Gateway** (librerie a nastro virtuali per i software di backup).

=== "Google Cloud"

    - **Filestore** con diversi livelli (zonale, regionale per l'alta disponibilità, enterprise).
    - **Google Cloud NetApp Volumes** per carichi aziendali NFS/SMB.
    - **Parallelstore** e **Managed Lustre** per il calcolo ad alte prestazioni.
    - **Cloud Storage FUSE** per montare un bucket come file system (con limiti: non è un vero file system POSIX).

## Laboratorio

!!! laboratorio "Laboratorio · Condivisione temporanea e replica"
    === "AWS"

        1. Crea due bucket con versioning: `academy-orig-...` nella regione del corso e `academy-copia-...` in
           `eu-central-1`.
        2. Sul primo, **Management** → **Replication rules** → regola su tutti gli oggetti verso il secondo bucket,
           lasciando che la console crei il ruolo IAM.
        3. Carica un file e verifica dopo qualche istante che compaia nel bucket di destinazione con
           `Replication status: REPLICA`.
        4. Genera un URL prefirmato valido 5 minuti e aprilo da una finestra anonima del browser:

            ```bash
            aws s3 presign s3://academy-orig-NOME/file.txt --expires-in 300
            ```

            Dopo 5 minuti l'URL restituisce `AccessDenied`.

        **Pulizia.** Elimina la regola di replica, svuota ed elimina i bucket, elimina il ruolo IAM creato.

    === "Google Cloud"

        1. Crea un bucket **dual-region** (per esempio la configurazione predefinita `EUR4` o una coppia di regioni
           europee) e caricaci un file.
        2. Genera un URL firmato valido 5 minuti. Serve un service account: impersonalo con il tuo utente (ruolo
           **Service Account Token Creator** sul service account), poi:

            ```bash
            gcloud storage sign-url gs://NOME-BUCKET/file.txt --duration=5m \
              --impersonate-service-account=lab-lettore@ID-PROGETTO.iam.gserviceaccount.com
            ```

            Il service account deve poter leggere l'oggetto (`roles/storage.objectViewer`, [lezione 4.4](../fase-4/04-iam.md)).
        3. Apri l'URL da una finestra anonima; dopo 5 minuti non funziona più.

        **Pulizia.** Elimina bucket e oggetti; rimuovi il ruolo concesso al tuo utente.

??? success "Spiegazione"
    La replica e la configurazione dual-region proteggono dalla perdita di una regione: il bucket di destinazione è il
    punto di partenza del tuo piano di disaster recovery. L'URL firmato concede l'accesso a **un solo oggetto**, per un
    tempo limitato, con i permessi di chi lo ha firmato: è il modo corretto di condividere file privati (per esempio una
    fattura scaricabile dal sito) senza rendere pubblico nulla.

## Esercizi

!!! esercizio "Esercizio · Requisiti e soluzioni"
    1. Un'assicurazione deve conservare i documenti per 10 anni in modo che **nessuno**, nemmeno un amministratore,
       possa cancellarli prima.
    2. Un'azienda vuole cifrare i dati di un bucket e poter dimostrare a un revisore ogni accesso alla chiave.
    3. Un ufficio vuole continuare a usare una cartella di rete SMB ma salvare i dati su S3.
    4. Un team vuole un backup giornaliero di tutte le istanze, i volumi e i database etichettati `backup=si`, con copia
       in un'altra regione.

??? success "Soluzione"
    1. **S3 Object Lock in modalità compliance** con conservazione di 10 anni; su Google Cloud una **retention policy
       bloccata con Bucket Lock**.
    2. **SSE-KMS** con una chiave gestita dal cliente: ogni uso della chiave è registrato in **CloudTrail**; su Google
       Cloud **CMEK** con gli **audit log** di Cloud KMS.
    3. **Storage Gateway in modalità S3 File Gateway** (SMB).
    4. **AWS Backup** con un piano giornaliero assegnato tramite l'etichetta `backup=si` e una regola di copia verso
       un'altra regione; su Google Cloud **Backup and DR Service**.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale opzione di cifratura S3 registra ogni uso della chiave in CloudTrail e permette di controllare chi può decifrare?","options":["SSE-S3","SSE-KMS","SSE-C","Nessuna cifratura"],"answer":[1],"explanation":"La chiave è in AWS KMS, con i suoi permessi."},
 {"q":"In quale modalità di S3 Object Lock nemmeno l'utente root può cancellare gli oggetti prima della scadenza?","options":["Governance","Compliance","Legal hold disattivato","Versioning"],"answer":[1],"explanation":"In governance gli utenti con un permesso speciale possono rimuovere il blocco."},
 {"q":"Quale servizio centralizza i backup di EBS, RDS, DynamoDB, EFS e altri con piani assegnati tramite etichette?","options":["Data Lifecycle Manager","AWS Backup","S3 Glacier","Storage Gateway"],"answer":[1],"explanation":"Data Lifecycle Manager gestisce solo snapshot EBS e AMI."},
 {"q":"Come caricare in modo affidabile un file da 20 GB su S3?","options":["Con un'unica richiesta PUT","Con il caricamento in più parti (multipart upload)","Via email","Con Storage Gateway Tape"],"answer":[1],"explanation":"Le parti sono caricate in parallelo e riprese in caso di errore."},
 {"q":"Quali servizi collegano applicazioni in sede allo storage cloud tramite protocolli di file?","options":["S3 File Gateway","FSx for Windows File Server raggiunto via VPN","Amazon Athena","AWS KMS","CloudFront"],"answer":[0,1],"explanation":"Athena interroga i dati, KMS gestisce chiavi, CloudFront è una CDN."},
 {"q":"Come condividi per un'ora un file privato di Cloud Storage con un cliente esterno?","options":["Rendendo pubblico il bucket","Con un URL firmato","Inviando la chiave del service account","Con una VPN"],"answer":[1],"explanation":"L'URL scade e vale solo per quell'oggetto."}
]}
</script>

## Fonti

- [Amazon S3: Protecting data with encryption](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingEncryption.html)
- [Amazon S3: Replicating objects](https://docs.aws.amazon.com/AmazonS3/latest/userguide/replication.html)
- [Amazon S3: Object Lock](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock.html)
- [Amazon S3: Uploading and copying objects using multipart upload](https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html)
- [AWS Backup Developer Guide](https://docs.aws.amazon.com/aws-backup/latest/devguide/whatisbackup.html)
- [AWS Storage Gateway](https://docs.aws.amazon.com/storagegateway/latest/userguide/WhatIsStorageGateway.html)
- [Cloud Storage: Encryption options](https://docs.cloud.google.com/storage/docs/encryption) e [Bucket Lock](https://docs.cloud.google.com/storage/docs/bucket-lock)
- [Cloud Storage: Turbo replication](https://docs.cloud.google.com/storage/docs/availability-durability#turbo-replication)
- [Cloud Storage: Signed URLs](https://docs.cloud.google.com/storage/docs/access-control/signed-urls)
- [Google Cloud: Backup and DR Service](https://docs.cloud.google.com/backup-disaster-recovery/docs/concepts/backup-dr)
