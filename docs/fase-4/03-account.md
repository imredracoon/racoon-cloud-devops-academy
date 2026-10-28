# 4.3 Account, progetti e fatturazione

!!! obiettivi "Obiettivi"
    - Creare un account AWS e un account Google Cloud con le offerte per i nuovi clienti.
    - Proteggere l'utente principale con l'autenticazione a più fattori.
    - Impostare budget e avvisi prima di creare qualsiasi risorsa.
    - Capire come sono organizzate le risorse: account e organizzazioni su AWS, organizzazione, cartelle e progetti su
      Google Cloud.

!!! warning "Prima di iniziare"
    Rileggi [Account cloud in sicurezza](../ambiente/account-cloud.md): le dieci regole valgono per tutto il corso.
    Servono un indirizzo email dedicato (consigliato), un numero di telefono e una carta di pagamento (anche
    prepagata, se accettata: verifica al momento dell'iscrizione).

## Come sono organizzate le risorse

=== "AWS"

    Su AWS l'unità fondamentale è l'**account**: un contenitore di risorse, con la propria fatturazione e un **confine
    di sicurezza** netto. Le risorse di un account non sono visibili da un altro, salvo permessi espliciti.

    Le aziende usano **molti account** (uno per ambiente, per team, per progetto) e li raggruppano con **AWS
    Organizations**:

    - una **fatturazione consolidata**: un'unica fattura per tutti gli account, con gli sconti per volume calcolati
      sul totale;
    - **unità organizzative** (*OU*) per raggruppare gli account;
    - **Service Control Policies** (SCP): limiti massimi ai permessi, che valgono anche per gli amministratori degli
      account (per esempio "nessuno può usare regioni fuori dall'Europa").

    Per il corso basta un solo account.

=== "Google Cloud"

    Su Google Cloud le risorse sono organizzate in una **gerarchia**:

    ```mermaid
    flowchart TB
      O[Organizzazione<br>esempio.it] --> F1[Cartella<br>Produzione]
      O --> F2[Cartella<br>Sviluppo]
      F1 --> P1[Progetto<br>tasks-prod]
      F2 --> P2[Progetto<br>tasks-dev]
      P1 --> R1[VM, bucket, database...]
      P2 --> R2[VM, bucket, database...]
    ```

    - Il **progetto** è l'unità fondamentale: ogni risorsa appartiene a un progetto, che ha un **nome**, un **ID**
      univoco a livello mondiale (scelto da te, non modificabile) e un **numero** assegnato da Google.
    - Le **cartelle** raggruppano progetti; l'**organizzazione** è la radice, legata a un dominio aziendale (Google
      Workspace o Cloud Identity). Un account personale lavora senza organizzazione.
    - I permessi e le **policy dell'organizzazione** assegnati a un livello si **ereditano** verso il basso.
    - La **fatturazione** è separata: un **account di fatturazione** paga per uno o più progetti.

    Per il corso basta un progetto, per esempio `academy-NOME-2026`.

## Creare l'account

=== "AWS"

    1. Vai su [aws.amazon.com/free](https://aws.amazon.com/free/) e scegli **Create a free account**.
    2. Inserisci l'email (diventerà l'**utente root**) e il nome dell'account.
    3. Scegli il **Free plan** (piano gratuito): nessun addebito, crediti fino a 200 USD, durata 6 mesi.
    4. Inserisci i dati di contatto e di pagamento, verifica il telefono.
    5. Accedi alla console come utente root.

=== "Google Cloud"

    1. Vai su [cloud.google.com/free](https://cloud.google.com/free) e scegli **Inizia gratuitamente** con un account
       Google (meglio uno dedicato).
    2. Accetta i termini e inserisci i dati di pagamento: la prova gratuita **non** passa automaticamente a pagamento.
    3. Google crea un primo progetto (`My First Project`). Crea invece un progetto tuo: in alto, selettore dei
       progetti → **Nuovo progetto**, e scegli un ID leggibile.

## Proteggere l'utente principale

=== "AWS"

    L'**utente root** ha poteri illimitati sull'account e non può essere limitato da alcuna policy IAM. Si usa solo
    per le poche attività che lo richiedono (per esempio cambiare il piano di supporto o chiudere l'account).

    1. Menu in alto a destra (nome dell'account) → **Security credentials**.
    2. **Multi-factor authentication (MFA)** → **Assign MFA device**: scegli un'app di autenticazione (o una chiave di
       sicurezza FIDO) e segui la procedura.
    3. Verifica che **non** esistano *access keys* per l'utente root: non devono mai esistere.

    Nella prossima lezione creerai un utente amministratore separato per il lavoro quotidiano.

=== "Google Cloud"

    L'account Google con cui ti sei iscritto è **proprietario** (*Owner*) del progetto.

    1. Vai su [myaccount.google.com/security](https://myaccount.google.com/security) e attiva la **verifica in due
       passaggi**, preferibilmente con una passkey o una chiave di sicurezza.
    2. Google Cloud impone già la verifica in due passaggi per gli accessi alla console: attivandola subito eviti
       blocchi e proteggi l'account.

## Budget e avvisi: subito

Prima di creare qualsiasi risorsa, imposta un avviso di spesa.

=== "AWS"

    1. Console → **Billing and Cost Management** → **Budgets** → **Create budget**.
    2. Scegli il modello **Zero spend budget**: ti avvisa appena viene speso anche un centesimo oltre i crediti
       gratuiti. Inserisci la tua email.
    3. Crea anche un **Monthly cost budget** di 5 USD, con avvisi all'80% e al 100%.
    4. In **Billing preferences** attiva gli avvisi del Free Tier via email.

=== "Google Cloud"

    1. Console → **Fatturazione** → **Budget e avvisi** → **Crea budget**.
    2. Ambito: il tuo progetto. Importo: 5 USD (o EUR, secondo la valuta del conto). Soglie: 50%, 90%, 100%.
    3. Lascia selezionato l'invio delle email agli amministratori della fatturazione.

!!! danger "Un budget avvisa, non blocca"
    Sia su AWS sia su Google Cloud un budget **invia avvisi** ma, per impostazione predefinita, **non ferma** le
    risorse. Se ricevi un avviso, intervieni subito. (È possibile automatizzare azioni al superamento di un budget,
    ma è una configurazione avanzata.) Il piano gratuito AWS e la prova gratuita Google Cloud, invece, non addebitano
    nulla finché non passi esplicitamente a pagamento: sono la tua vera rete di sicurezza.

## Esercizi

!!! laboratorio "Laboratorio · Account pronto per il corso"
    Svolgi il laboratorio sul cloud che hai scelto (o su entrambi):

    1. crea l'account (AWS Free plan oppure prova gratuita di Google Cloud);
    2. attiva l'MFA sull'utente principale;
    3. su Google Cloud, crea il progetto del corso;
    4. crea i budget e gli avvisi descritti sopra;
    5. annota in un file `account.md` del tuo repository di appunti: ID dell'account o del progetto, regione scelta,
       data di scadenza dei crediti. **Nessuna password né chiave**.

??? success "Verifica"
    - AWS: in **IAM → Dashboard** compare "Root user has MFA" senza avvisi; in **Budgets** ci sono due budget attivi.
    - Google Cloud: in **Fatturazione → Budget e avvisi** il budget è collegato al progetto; nel selettore dei progetti
      vedi il tuo progetto con il suo ID.
    - Segna in calendario la scadenza dei crediti: 6 mesi (AWS) o 90 giorni (Google Cloud) dalla registrazione.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è il primo passo dopo aver creato un account AWS?","options":["Creare una VM","Attivare l'MFA sull'utente root","Creare access key per root","Passare al piano a pagamento"],"answer":[1],"explanation":"L'utente root ha poteri illimitati e va protetto subito."},
 {"q":"Cosa offre AWS Organizations?","options":["Fatturazione consolidata di più account","Service Control Policies per limitare i permessi degli account","Un database gestito","VM più economiche","Backup automatici"],"answer":[0,1],"explanation":"Organizations gestisce in modo centralizzato molti account."},
 {"q":"Su Google Cloud, a cosa appartiene sempre una risorsa come una VM?","options":["A una cartella","A un progetto","A un account di fatturazione","A una regione"],"answer":[1],"explanation":"Il progetto è l'unità fondamentale."},
 {"q":"Un budget AWS o Google Cloud, con le impostazioni predefinite:","options":["Ferma le risorse al superamento","Invia avvisi ma non ferma le risorse","Cancella l'account","Limita la CPU"],"answer":[1],"explanation":"Bisogna intervenire quando arriva l'avviso."},
 {"q":"Le policy assegnate a una cartella di Google Cloud:","options":["Valgono solo per la cartella","Sono ereditate dai progetti e dalle risorse sottostanti","Valgono solo per l'organizzazione","Non esistono"],"answer":[1],"explanation":"La gerarchia propaga permessi e policy verso il basso."}
]}
</script>

## Fonti

- [AWS: Create an AWS account](https://docs.aws.amazon.com/accounts/latest/reference/manage-acct-creating.html)
- [AWS IAM: Root user best practices](https://docs.aws.amazon.com/IAM/latest/UserGuide/root-user-best-practices.html)
- [AWS Budgets: Creating a budget](https://docs.aws.amazon.com/cost-management/latest/userguide/budgets-create.html)
- [AWS Organizations: What is AWS Organizations?](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_introduction.html)
- [Google Cloud: Resource hierarchy](https://docs.cloud.google.com/resource-manager/docs/cloud-platform-resource-hierarchy)
- [Google Cloud: Creating and managing projects](https://docs.cloud.google.com/resource-manager/docs/creating-managing-projects)
- [Google Cloud: Create, edit, or delete budgets and budget alerts](https://docs.cloud.google.com/billing/docs/how-to/budgets)
