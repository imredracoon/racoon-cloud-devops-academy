# 4.4 Identità e accessi (IAM)

!!! obiettivi "Obiettivi"
    - Spiegare autenticazione e autorizzazione e il principio del privilegio minimo.
    - Conoscere gli elementi di IAM su AWS: utenti, gruppi, ruoli, policy, IAM Identity Center.
    - Conoscere gli elementi di IAM su Google Cloud: principal, ruoli, policy di autorizzazione, service account.
    - Leggere e scrivere una policy semplice.
    - Creare un utente amministratore per il lavoro quotidiano.

## I concetti

- L'**autenticazione** risponde a "**chi sei?**": password, MFA, chiavi, token.
- L'**autorizzazione** risponde a "**cosa puoi fare?**": i permessi.
- **IAM** (*Identity and Access Management*) è il servizio che gestisce entrambe: è **globale** e **gratuito** su
  entrambi i cloud.

Il principio guida è il **privilegio minimo** (*least privilege*): ogni persona e ogni programma riceve **solo** i
permessi necessari per il suo compito, e solo per il tempo necessario. Un'applicazione che legge un bucket non deve
poterlo cancellare; uno sviluppatore non deve poter modificare la fatturazione.

Un secondo principio: **le identità umane e quelle delle macchine sono diverse**. Le persone accedono con un account
personale e MFA; i programmi usano **credenziali temporanee** ottenute tramite ruoli (AWS) o service account (Google
Cloud), mai le credenziali di una persona.

## IAM su AWS

=== "AWS"

    | Elemento | Cos'è |
    |---|---|
    | **Utente root** | l'identità con cui è stato creato l'account. Poteri illimitati: da usare quasi mai. |
    | **Utente IAM** | un'identità a lungo termine, con password e/o chiavi di accesso. Oggi sconsigliato per le persone. |
    | **Gruppo IAM** | un insieme di utenti IAM a cui assegnare le stesse policy |
    | **Ruolo IAM** | un'identità **senza credenziali permanenti**, che persone, servizi o applicazioni possono **assumere** ottenendo credenziali temporanee |
    | **Policy** | un documento JSON che elenca permessi |
    | **IAM Identity Center** | il servizio consigliato per gli accessi delle persone: utenti centralizzati, accesso con MFA a uno o più account tramite *permission set*, credenziali sempre temporanee |

    ### Le policy

    Una policy è un documento JSON con una o più **dichiarazioni** (*statement*):

    ```json
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "LeggiIlBucketDelSito",
          "Effect": "Allow",
          "Action": ["s3:GetObject", "s3:ListBucket"],
          "Resource": [
            "arn:aws:s3:::academy-sito",
            "arn:aws:s3:::academy-sito/*"
          ]
        }
      ]
    }
    ```

    - `Effect`: `Allow` o `Deny`.
    - `Action`: le operazioni, nel formato `servizio:Operazione`; `*` è un jolly (`s3:*`).
    - `Resource`: le risorse, identificate da un **ARN** (*Amazon Resource Name*), il nome univoco di ogni risorsa AWS.
    - `Condition` (facoltativa): condizioni aggiuntive, per esempio solo da un certo indirizzo IP o solo con MFA.

    La logica di valutazione: **tutto è negato per impostazione predefinita**; un `Allow` esplicito concede; un `Deny`
    esplicito **vince sempre** su qualsiasi `Allow`.

    Tipi di policy: **gestite da AWS** (pronte, come `AdministratorAccess` o `ReadOnlyAccess`), **gestite dal
    cliente** (le tue, riutilizzabili), **inline** (incorporate in un'identità), **policy basate sulle risorse**
    (attaccate a una risorsa, come la *bucket policy* di S3, e indicano anche *chi* può accedere).

    ### Ruoli per i servizi

    Un'istanza EC2 che deve leggere da S3 **non** riceve chiavi di accesso: le si associa un **ruolo** (tramite un
    *instance profile*), e il software dell'istanza ottiene in automatico credenziali temporanee che si rinnovano da
    sole. Lo stesso vale per Lambda, per i container e, dalla Fase 7, per GitHub Actions.

=== "Google Cloud"

    | Elemento | Cos'è |
    |---|---|
    | **Principal** | chi riceve i permessi: un account Google, un gruppo Google, un **service account**, un dominio, identità federate |
    | **Permesso** | una singola operazione, nel formato `servizio.risorsa.verbo`, per esempio `storage.objects.get` |
    | **Ruolo** | un insieme di permessi. **Non** si assegnano permessi singoli, ma ruoli |
    | **Policy di autorizzazione** (*allow policy*) | l'elenco dei **binding** "ruolo → principal", attaccato a una risorsa (organizzazione, cartella, progetto o singola risorsa) e ereditato verso il basso |
    | **Service account** | un'identità per applicazioni e VM, non per persone |

    ### I tipi di ruolo

    - **Ruoli di base**: Owner, Editor, Viewer. Molto ampi: da evitare in produzione.
    - **Ruoli predefiniti**: creati da Google per ogni servizio, a grana fine, per esempio
      `roles/storage.objectViewer` (lettura degli oggetti) o `roles/compute.instanceAdmin.v1`.
    - **Ruoli personalizzati**: creati da te con un elenco preciso di permessi.

    Una policy di autorizzazione, in YAML:

    ```yaml
    bindings:
    - role: roles/storage.objectViewer
      members:
      - serviceAccount:sito@academy-mario-2026.iam.gserviceaccount.com
    - role: roles/viewer
      members:
      - group:sviluppatori@esempio.it
    ```

    I permessi sono **additivi**: un principal ha l'unione dei ruoli ricevuti a ogni livello della gerarchia. Esistono
    anche le **deny policy**, che negano permessi indipendentemente dai ruoli, e le **policy dell'organizzazione**,
    che limitano le configurazioni possibili (per esempio "vietato creare bucket pubblici").

    ### Service account

    Una VM di Compute Engine o un servizio Cloud Run eseguono con un **service account**: il codice ottiene in
    automatico credenziali temporanee. Le **chiavi** dei service account (file JSON) esistono ma sono un rischio,
    come le access key di AWS: si evitano, e dalla Fase 7 userai la federazione delle identità (*Workload Identity
    Federation*).

## AWS e Google Cloud a confronto

| Concetto | AWS | Google Cloud |
|---|---|---|
| Identità di un'applicazione | ruolo IAM | service account |
| Insieme di permessi | policy (JSON) | ruolo |
| Dove si attaccano i permessi | all'identità (o alla risorsa con le policy basate su risorse) | alla risorsa, con una policy di autorizzazione che elenca i principal |
| Accesso centralizzato delle persone | IAM Identity Center | Cloud Identity / Google Workspace, gruppi Google |
| Limiti massimi per l'organizzazione | Service Control Policies | Organization Policy, deny policy |
| Analisi dei permessi | IAM Access Analyzer | IAM Recommender, Policy Analyzer |

## Le buone pratiche

1. MFA per tutte le persone, a partire dall'utente principale.
2. Mai usare l'utente root (AWS) o un ruolo di base come Owner (Google Cloud) per il lavoro quotidiano.
3. Permessi assegnati a **gruppi**, non a singole persone.
4. Privilegio minimo: parti da zero e aggiungi; usa gli strumenti che suggeriscono di togliere i permessi non usati.
5. **Credenziali temporanee** ovunque possibile; chiavi di accesso a lungo termine solo se inevitabili, ruotate e mai
   in Git.
6. Rivedi periodicamente chi ha accesso a cosa, e disattiva subito gli accessi di chi lascia il team.

## Laboratorio

!!! laboratorio "Laboratorio · Un amministratore per il lavoro quotidiano"
    === "AWS"

        1. Console (come root) → **IAM Identity Center** → **Enable**. Scegli la regione del corso: Identity Center
           vive in una sola regione.
        2. **Users** → **Add user**: il tuo nome utente e la tua email. Riceverai un invito.
        3. **Groups** → crea il gruppo `Amministratori` e aggiungi il tuo utente.
        4. **Permission sets** → **Create** → predefinito **AdministratorAccess**, durata della sessione 4 ore.
        5. **AWS accounts** → seleziona il tuo account → **Assign users or groups** → gruppo `Amministratori`,
           permission set `AdministratorAccess`.
        6. Accetta l'invito, imposta password e MFA. Annota l'**AWS access portal URL**
           (`https://d-xxxxxxxxxx.awsapps.com/start`).
        7. Esci da root e, da ora, accedi solo dal portale.

    === "Google Cloud"

        1. Crea un service account per i laboratori: **IAM e amministrazione** → **Service account** → **Crea**,
           nome `lab-lettore`, nessun ruolo per ora.
        2. **IAM** → **Concedi l'accesso**: principal `lab-lettore@ID-PROGETTO.iam.gserviceaccount.com`, ruolo
           **Visualizzatore oggetti Storage** (`roles/storage.objectViewer`).
        3. Osserva nella pagina **IAM** l'elenco dei principal e dei ruoli: il tuo account è **Proprietario**.
        4. Apri il **Policy Troubleshooter** e verifica se `lab-lettore` ha il permesso `storage.objects.delete` sul
           progetto (risposta attesa: no).

        Per le persone, in un'azienda, si usano i gruppi di Google Workspace o Cloud Identity; con un account
        personale lavori come proprietario del tuo progetto, protetto dalla verifica in due passaggi.

??? success "Verifica"
    - AWS: accedendo dal portale vedi l'account con il ruolo `AdministratorAccess`; in alto a destra compare un nome del
      tipo `AWSReservedSSO_AdministratorAccess_.../tuonome`, cioè un **ruolo assunto**, con credenziali temporanee.
    - Google Cloud: il Policy Troubleshooter risponde che il permesso **non** è concesso, perché il ruolo
      `storage.objectViewer` include solo permessi di lettura.

    **Spiegazione.** Su AWS non hai creato un utente IAM con password permanente, ma un'identità centralizzata che
    ottiene credenziali temporanee: è il metodo che AWS raccomanda per le persone. La sessione scade dopo 4 ore.

## Esercizi

!!! esercizio "Esercizio 1 · Leggi la policy"
    Cosa permette questa policy AWS? Un utente che la riceve insieme alla policy gestita `AmazonS3FullAccess` può
    cancellare oggetti dal bucket `academy-log`?

    ```json
    {
      "Version": "2012-10-17",
      "Statement": [
        {"Effect": "Allow", "Action": "s3:*", "Resource": "*"},
        {"Effect": "Deny", "Action": "s3:DeleteObject", "Resource": "arn:aws:s3:::academy-log/*"}
      ]
    }
    ```

??? success "Soluzione"
    Concede tutte le operazioni S3 su tutte le risorse, **tranne** la cancellazione di oggetti nel bucket
    `academy-log`. No, l'utente **non** può cancellare: un `Deny` esplicito prevale su qualsiasi `Allow`, anche se
    proveniente da un'altra policy.

!!! esercizio "Esercizio 2 · Scegli il ruolo"
    Su Google Cloud, quale tipo di ruolo e quale principal useresti per:

    1. un'applicazione su Cloud Run che deve leggere file da un bucket;
    2. un gruppo di analisti che deve solo eseguire query su BigQuery;
    3. un consulente che deve vedere tutte le risorse del progetto per una settimana, senza modificarle?

??? success "Soluzione"
    1. Un **service account** dedicato all'applicazione con il ruolo predefinito `roles/storage.objectViewer`,
       assegnato **sul bucket** (non sull'intero progetto).
    2. Il **gruppo Google** degli analisti con ruoli predefiniti di BigQuery (per esempio `roles/bigquery.jobUser` sul
       progetto e `roles/bigquery.dataViewer` sui dataset).
    3. Il suo account con il ruolo `roles/viewer` e una **condizione IAM** con data di scadenza; in alternativa
       ricordarsi di rimuoverlo. È uno dei pochi casi in cui un ruolo di base è accettabile, perché è di sola lettura.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa afferma il principio del privilegio minimo?","options":["Ogni utente deve avere pochi file","Ogni identità riceve solo i permessi necessari per il suo compito","Gli amministratori devono essere pochi","Le password devono essere brevi"],"answer":[1],"explanation":"È la base di ogni progettazione IAM."},
 {"q":"Su AWS una policy concede s3:* e un'altra nega s3:DeleteObject sulla stessa risorsa. Il risultato?","options":["La cancellazione è permessa","La cancellazione è negata","Dipende dall'ordine","Errore di configurazione"],"answer":[1],"explanation":"Un Deny esplicito vince sempre."},
 {"q":"Come ottiene i permessi un'istanza EC2 che deve leggere da S3?","options":["Con le access key di un amministratore salvate sul disco","Con un ruolo IAM associato all'istanza","Con la password root","Non è possibile"],"answer":[1],"explanation":"Il ruolo fornisce credenziali temporanee rinnovate automaticamente."},
 {"q":"Su Google Cloud, cosa si assegna ai principal?","options":["Singoli permessi","Ruoli","ARN","Security group"],"answer":[1],"explanation":"I ruoli sono insiemi di permessi."},
 {"q":"Quali sono buone pratiche IAM? (scegli 2)","options":["Assegnare i permessi a gruppi","Usare credenziali temporanee per le applicazioni","Condividere l'utente root nel team","Salvare le chiavi dei service account in Git","Dare il ruolo Owner a tutti gli sviluppatori"],"answer":[0,1],"explanation":"Le altre sono gravi errori di sicurezza."},
 {"q":"Qual è il servizio AWS consigliato per gli accessi delle persone a uno o più account?","options":["Utenti IAM con access key","IAM Identity Center","Amazon Cognito","AWS Organizations da solo"],"answer":[1],"explanation":"Credenziali temporanee, MFA, gestione centralizzata."}
]}
</script>

## Fonti

- [AWS IAM User Guide: Security best practices](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html)
- [AWS IAM: Policy evaluation logic](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_evaluation-logic.html)
- [AWS IAM: Roles](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html)
- [AWS IAM Identity Center: Getting started](https://docs.aws.amazon.com/singlesignon/latest/userguide/getting-started.html)
- [Google Cloud IAM: Overview](https://docs.cloud.google.com/iam/docs/overview)
- [Google Cloud IAM: Roles and permissions](https://docs.cloud.google.com/iam/docs/roles-overview)
- [Google Cloud: Best practices for using service accounts](https://docs.cloud.google.com/iam/docs/best-practices-service-accounts)
- [Google Cloud: Policy Troubleshooter](https://docs.cloud.google.com/policy-intelligence/docs/troubleshoot-access)
