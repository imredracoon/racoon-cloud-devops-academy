# 4.6 Storage

!!! obiettivi "Obiettivi"
    - Distinguere storage a blocchi, a file e a oggetti e sapere quando usare ciascuno.
    - Usare lo storage a oggetti: bucket, oggetti, classi di archiviazione, versioning, cicli di vita.
    - Proteggere i dati: accesso pubblico bloccato, crittografia, permessi.
    - Conoscere i dischi delle VM, gli snapshot e i file system condivisi.
    - Conoscere gli strumenti per trasferire grandi quantità di dati nel cloud.

## Tre tipi di storage

| Tipo | Come si usa | Esempio nella vita reale | AWS | Google Cloud |
|---|---|---|---|---|
| **A blocchi** | un disco collegato a **una** VM, con un filesystem ([lezione 1.6](../fase-1/06-storage.md)) | il disco del portatile | **EBS** | **Persistent Disk**, **Hyperdisk** |
| **A file** | una cartella condivisa via rete tra **più** macchine (NFS, SMB) | il server di file dell'ufficio | **EFS**, **FSx** | **Filestore** |
| **A oggetti** | file (oggetti) in contenitori (bucket), letti e scritti via **API HTTP** | Google Drive, Dropbox | **S3** | **Cloud Storage** |

Lo storage a oggetti è il cuore del cloud: economico, praticamente illimitato, estremamente durevole, accessibile
da ovunque. Ci si salvano backup, log, immagini, siti statici, dati per l'analisi.

## Lo storage a oggetti

=== "AWS"

    In **Amazon S3** (*Simple Storage Service*):

    - un **bucket** è il contenitore; il suo nome è **univoco in tutto il mondo** e il bucket vive in una regione;
    - un **oggetto** è un file con i suoi metadati, identificato da una **chiave** (`foto/2026/mare.jpg`); le
      "cartelle" sono solo prefissi della chiave;
    - un singolo oggetto può arrivare a 5 TB;
    - S3 è progettato per una **durabilità** del 99,999999999% ("undici nove"), replicando i dati su almeno tre AZ
      (tranne le classi One Zone).

=== "Google Cloud"

    In **Cloud Storage**:

    - un **bucket** ha un nome univoco a livello mondiale e una **posizione**: **regione** (una regione),
      **doppia regione** (due regioni specifiche) o **multi-regione** (per esempio `EU`, un'area geografica ampia);
    - un **oggetto** è un file con i metadati; anche qui le cartelle sono prefissi (o *cartelle gestite*);
    - un singolo oggetto può arrivare a 5 TiB;
    - la durabilità annuale dichiarata è di almeno 99,999999999%.

### Le classi di archiviazione

Si paga per i GB conservati, per le operazioni e per i dati in uscita. Le classi "fredde" costano meno per GB ma di
più per la lettura, e hanno una durata minima di conservazione.

=== "AWS"

    | Classe | Uso | Note |
    |---|---|---|
    | **S3 Standard** | dati usati spesso | la predefinita |
    | **S3 Intelligent-Tiering** | accesso imprevedibile | sposta gli oggetti tra livelli in automatico in base all'uso |
    | **S3 Standard-IA** | accesso raro, ma lettura immediata | durata minima 30 giorni, costo per GB letto |
    | **S3 One Zone-IA** | come IA, dati ricreabili | **una sola AZ**: costa meno, meno resiliente |
    | **S3 Express One Zone** | latenza bassissima | una sola AZ, per carichi ad alte prestazioni |
    | **S3 Glacier Instant Retrieval** | archivio letto raramente, in millisecondi | minimo 90 giorni |
    | **S3 Glacier Flexible Retrieval** | archivio, recupero in minuti o ore | minimo 90 giorni |
    | **S3 Glacier Deep Archive** | archivio a lunghissimo termine, recupero entro 12-48 ore | il più economico, minimo 180 giorni |

=== "Google Cloud"

    | Classe | Uso | Durata minima |
    |---|---|---|
    | **Standard** | dati usati spesso | nessuna |
    | **Nearline** | accesso circa una volta al mese | 30 giorni |
    | **Coldline** | accesso circa una volta al trimestre | 90 giorni |
    | **Archive** | accesso meno di una volta all'anno, backup e conformità | 365 giorni |

    Tutte le classi offrono lettura in millisecondi: cambia il costo, non la velocità. **Autoclass** sposta gli oggetti
    tra le classi in automatico in base all'uso.

### Versioning e cicli di vita

- Il **versioning** conserva ogni versione di un oggetto: una sovrascrittura o una cancellazione per errore si possono
  annullare. È la prima difesa contro errori e ransomware.
- Le **regole del ciclo di vita** (*lifecycle*) automatizzano: "sposta in una classe fredda dopo 30 giorni",
  "cancella dopo un anno", "cancella le versioni non correnti dopo 7 giorni".
- Il **blocco degli oggetti** (S3 **Object Lock**, Cloud Storage **Bucket Lock** e blocco degli oggetti) impedisce
  modifiche e cancellazioni per un periodo stabilito: modello *WORM* (*write once, read many*), richiesto da alcune
  normative.
- La **replica** copia automaticamente gli oggetti in un altro bucket, anche in un'altra regione (S3 Cross-Region
  Replication; con Cloud Storage si sceglie una posizione doppia o multi-regione).

### Proteggere i dati

!!! danger "Il bucket pubblico per errore"
    Molte fughe di dati della storia del cloud nascono da un bucket reso pubblico per sbaglio. Per questo oggi i
    bucket nascono **privati** e con l'accesso pubblico **bloccato**:

    - AWS: **S3 Block Public Access** attivo per impostazione predefinita su ogni nuovo bucket (e si può imporre su
      tutto l'account);
    - Google Cloud: **prevenzione dell'accesso pubblico** (*public access prevention*) e **accesso uniforme a livello di
      bucket**, che disattiva i permessi sui singoli oggetti a favore di IAM.

    Rendi pubblico qualcosa solo consapevolmente, come nel progetto di questa fase, e mai un bucket con dati.

- **Crittografia a riposo**: entrambi i cloud **cifrano sempre** i dati salvati, con chiavi gestite dal fornitore.
  Per maggiore controllo si usano chiavi gestite dal cliente in **AWS KMS** o **Cloud KMS** ([Fase 5](../fase-5/index.md)).
- **Crittografia in transito**: le API sono HTTPS; si può imporre di rifiutare le richieste non cifrate.
- **Permessi**: IAM e, su S3, le **bucket policy**.
- **URL prefirmati** (*presigned URL* su S3, *signed URL* su Cloud Storage): danno accesso temporaneo a un oggetto
  privato, senza renderlo pubblico.

## I dischi delle VM

=== "AWS"

    Un volume **EBS** è un disco di rete collegato a un'istanza nella **stessa AZ**. Sopravvive all'istanza (se non
    impostato per essere cancellato con lei).

    | Tipo | Uso |
    |---|---|
    | **gp3** (SSD uso generale) | la scelta predefinita: prestazioni configurabili indipendentemente dalla dimensione |
    | **io2 Block Express** (SSD ad alte prestazioni) | database critici |
    | **st1** (HDD ottimizzato per il throughput) | big data, log, letture sequenziali |
    | **sc1** (HDD cold) | dati poco usati, al costo più basso |

    Gli **snapshot** EBS sono copie incrementali salvate (in modo trasparente) su S3: servono per backup, per creare
    nuovi volumi e per copiare dati in un'altra AZ o regione.

    L'**instance store** è invece un disco fisicamente collegato al server: velocissimo ma **temporaneo**, i dati si
    perdono se l'istanza si ferma.

=== "Google Cloud"

    I dischi di Compute Engine sono **Persistent Disk** (standard, bilanciato, SSD) e la generazione più recente
    **Hyperdisk**, con prestazioni configurabili. Possono essere **zonali** o **regionali** (replicati in sincrono
    in due zone, per l'alta disponibilità).

    Gli **snapshot** sono incrementali e, per impostazione predefinita, conservati in più regioni. Le **SSD locali**
    sono dischi fisici temporanei, come l'instance store di AWS.

## File system condivisi

- **Amazon EFS**: NFS gestito per Linux, cresce e si riduce in automatico, accessibile da più AZ.
- **Amazon FSx**: file system gestiti specializzati: Windows File Server (SMB), Lustre (calcolo ad alte prestazioni),
  NetApp ONTAP, OpenZFS.
- **Filestore** (Google Cloud): NFS gestito.

## Portare i dati nel cloud

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Trasferimenti online tra sistemi in sede e cloud | **DataSync** | **Storage Transfer Service** |
| Grandi volumi offline, con un dispositivo fisico spedito | **Snowball Edge** | **Transfer Appliance** |
| Collegare applicazioni in sede allo storage cloud | **Storage Gateway** | — |

## Laboratorio

!!! laboratorio "Laboratorio · Un bucket con versioning e ciclo di vita"
    === "AWS"

        1. Console → **S3** → **Create bucket**, nella regione del corso, nome `academy-NOME-AAAAMMGG` (minuscolo,
           univoco). Lascia **Block all public access** attivo.
        2. Nelle opzioni attiva **Bucket Versioning**.
        3. Carica un file `nota.txt` con il testo `versione 1`. Modificalo in locale (`versione 2`) e ricaricalo.
        4. Attiva **Show versions**: vedi due versioni. Scarica la prima.
        5. Cancella l'oggetto: compare un **delete marker**. Elimina il delete marker: l'oggetto torna.
        6. **Management** → **Create lifecycle rule**: sposta gli oggetti in **Standard-IA** dopo 30 giorni e
           cancella le versioni non correnti dopo 7 giorni.

    === "Google Cloud"

        1. Console → **Cloud Storage** → **Bucket** → **Crea**, nome univoco, posizione **Regione** del corso, classe
           **Standard**, **prevenzione dell'accesso pubblico** attiva, controllo dell'accesso **uniforme**.
        2. Nella sezione sulla protezione dei dati attiva il **controllo delle versioni degli oggetti**.
        3. Carica `nota.txt` (`versione 1`), poi la versione modificata con lo stesso nome.
        4. Nei dettagli dell'oggetto apri la scheda **Cronologia delle versioni**: vedi le due versioni.
        5. Cancella l'oggetto e ripristina la versione precedente dalla cronologia.
        6. Scheda **Ciclo di vita** → **Aggiungi regola**: passa a **Nearline** dopo 30 giorni; elimina le versioni non
           correnti dopo 7 giorni (condizione sul numero di giorni da quando l'oggetto è diventato non corrente).

    **Pulizia.** Svuota ed elimina il bucket (con il versioning attivo, vanno eliminate anche tutte le versioni: la
    console lo fa quando svuoti il bucket).

??? success "Spiegazione"
    Con il versioning, una cancellazione su S3 non elimina i dati: aggiunge un *delete marker* che nasconde l'oggetto.
    Su Cloud Storage l'oggetto diventa una versione **non corrente**. In entrambi i casi puoi tornare indietro. Il
    costo è che paghi lo spazio di tutte le versioni: per questo si abbina sempre una regola che elimina le versioni
    vecchie.

## Esercizi

!!! esercizio "Esercizio · Scegli lo storage"
    1. I dischi di sistema di un database PostgreSQL installato su una VM.
    2. Le fatture degli ultimi dieci anni, che per legge non devono essere modificate e si consultano quasi mai.
    3. Una cartella condivisa tra dieci server web Linux che servono gli stessi file caricati dagli utenti.
    4. I log dell'applicazione, consultati spesso nella prima settimana e quasi mai dopo, da cancellare dopo un anno.
    5. 80 TB di archivio video da trasferire da un data center con una connessione lenta.

??? success "Soluzione"
    1. **A blocchi**: EBS gp3 o io2; Persistent Disk SSD o Hyperdisk. (Meglio ancora: un database gestito, [lezione 4.8](08-database-ai.md).)
    2. **A oggetti**, classe d'archivio (S3 Glacier Deep Archive, Cloud Storage Archive) con **Object Lock / Bucket
       Lock** per la non modificabilità.
    3. **A file**: EFS o Filestore. (In alternativa, cambiare l'applicazione perché salvi i file su S3/Cloud Storage.)
    4. **A oggetti** con una **regola di ciclo di vita**: Standard, poi una classe fredda dopo 7-30 giorni, cancellazione
       dopo 365 giorni.
    5. **Trasferimento offline**: Snowball Edge o Transfer Appliance.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale tipo di storage si usa via API HTTP, organizzato in bucket?","options":["A blocchi","A file","A oggetti","Instance store"],"answer":[2],"explanation":"S3 e Cloud Storage."},
 {"q":"Un volume EBS può essere collegato a un'istanza in un'altra AZ?","options":["Sì, sempre","No, deve essere nella stessa AZ; si può passare da uno snapshot","Solo con gp3","Solo con io2"],"answer":[1],"explanation":"EBS è un servizio zonale."},
 {"q":"Quale funzione protegge da cancellazioni e sovrascritture accidentali in un bucket?","options":["Il ciclo di vita","Il versioning","La classe Standard","La CDN"],"answer":[1],"explanation":"Conserva ogni versione degli oggetti."},
 {"q":"Quale classe di Cloud Storage scegli per dati letti meno di una volta all'anno?","options":["Standard","Nearline","Coldline","Archive"],"answer":[3],"explanation":"Durata minima 365 giorni, costo per GB più basso."},
 {"q":"Come dai accesso temporaneo a un oggetto privato a un utente esterno?","options":["Rendendo pubblico il bucket","Con un URL prefirmato o firmato","Inviando le tue credenziali","Copiandolo su EBS"],"answer":[1],"explanation":"L'URL scade dopo il tempo stabilito."},
 {"q":"Quali servizi offrono file system condivisi via rete?","options":["Amazon EFS","Filestore","Amazon S3 Glacier","Persistent Disk","Instance store"],"answer":[0,1],"explanation":"EFS e Filestore sono NFS gestiti."}
]}
</script>

## Fonti

- [Amazon S3 User Guide](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html) e [Storage classes](https://aws.amazon.com/s3/storage-classes/)
- [Amazon S3: Blocking public access](https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-block-public-access.html)
- [Amazon S3: Versioning](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Versioning.html) e [Lifecycle](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lifecycle-mgmt.html)
- [Amazon EBS volume types](https://docs.aws.amazon.com/ebs/latest/userguide/ebs-volume-types.html)
- [Cloud Storage: Storage classes](https://docs.cloud.google.com/storage/docs/storage-classes) e [Bucket locations](https://docs.cloud.google.com/storage/docs/locations)
- [Cloud Storage: Object Versioning](https://docs.cloud.google.com/storage/docs/object-versioning) e [Lifecycle](https://docs.cloud.google.com/storage/docs/lifecycle)
- [Cloud Storage: Public access prevention](https://docs.cloud.google.com/storage/docs/public-access-prevention)
- [Compute Engine: Storage options](https://docs.cloud.google.com/compute/docs/disks)
