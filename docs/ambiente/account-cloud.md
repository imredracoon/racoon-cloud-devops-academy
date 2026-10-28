# Account cloud in sicurezza

!!! warning "Non creare ancora gli account"
    Gli account cloud servono dalla **[Fase 4](../fase-4/index.md)**. Le offerte per i nuovi clienti durano pochi mesi: se li apri adesso,
    le perderai prima di poterle usare. Leggi questa pagina quando arrivi alla [lezione 4.3](../fase-4/03-account.md).

!!! obiettivi "Obiettivi"
    - Sapere cosa offrono AWS e Google Cloud ai nuovi clienti.
    - Creare gli account con le protezioni minime: autenticazione a più fattori, budget, allarmi.
    - Conoscere le regole che evitano addebiti imprevisti.

## Perché servono regole di sicurezza

Un account cloud è collegato a una **carta di pagamento**. Se qualcuno ruba le tue credenziali, o se dimentichi
accesa una risorsa costosa, il conto può crescere velocemente. Con le regole di questa pagina il rischio diventa
molto basso. Sono le stesse pratiche che un'azienda seria applica ai propri account.

## Le offerte per i nuovi clienti

=== "AWS"

    Dal 15 luglio 2025 i nuovi account AWS scelgono tra due piani:

    - **Piano gratuito (Free plan)**: fino a **200 USD di crediti** (100 alla registrazione, altri 100 completando
      attività guidate). Dura **6 mesi o fino all'esaurimento dei crediti**. Non vengono addebitati costi: alla
      scadenza l'account va aggiornato al piano a pagamento, altrimenti viene chiuso.
    - **Piano a pagamento (Paid plan)**: stessi crediti, ma accesso a tutti i servizi. Quando i crediti finiscono
      si paga l'utilizzo.

    In entrambi i casi restano disponibili oltre 30 servizi con una quota **sempre gratuita**.
    Per il corso consigliamo il **piano gratuito**: è impossibile ricevere addebiti.

=== "Google Cloud"

    I nuovi clienti ricevono la **prova gratuita**: **300 USD di crediti** validi **90 giorni**. Durante la prova
    non vengono addebitati costi e l'account **non** passa automaticamente a pagamento: alla scadenza le risorse
    vengono fermate e, dopo 30 giorni, cancellate. Si paga solo dopo aver attivato manualmente un account di
    fatturazione a pagamento.

    In più esiste il **Free Tier** sempre gratuito, con limiti mensili: per esempio una VM `e2-micro` in alcune
    regioni USA e 5 GB di Cloud Storage regionale in USA.

## Le dieci regole d'oro

1. **Autenticazione a più fattori (MFA)** su ogni utente, a partire dall'utente principale. Usa un'app di
   autenticazione sul telefono o una chiave di sicurezza fisica.
2. **Non usare l'utente principale** (*root user* su AWS) per il lavoro quotidiano: crea un utente amministratore
   separato. Lo vedrai nella [lezione 4.4](../fase-4/04-iam.md).
3. **Budget e allarmi** prima di qualsiasi altra cosa: un avviso via email a 1, 5 e 10 USD.
4. **Mai pubblicare chiavi di accesso**. Non vanno mai in un repository Git, nemmeno privato. Nel corso userai
   credenziali temporanee e, dalla [Fase 7](../fase-7/index.md), l'accesso senza chiavi tramite OIDC.
5. **Distruggi le risorse alla fine di ogni laboratorio**. Ogni laboratorio cloud termina con una sezione
   *Pulizia*: non saltarla mai.
6. **Una sola regione** per tutto il corso, così trovi sempre le tue risorse. Consigliamo `eu-south-1` (Milano)
   o `eu-central-1` (Francoforte) su AWS e `europe-west8` (Milano) su Google Cloud, salvo diversa indicazione
   del laboratorio per sfruttare le quote gratuite.
7. **Etichetta le risorse** con un'etichetta `progetto=academy`: ti aiuta a trovarle e a cancellarle.
8. **Controlla la fatturazione ogni settimana** (AWS Billing and Cost Management, Google Cloud Billing).
9. **Diffida delle risorse "grandi"**: nel corso non servono mai macchine potenti, GPU o database di grandi dimensioni.
10. **Se ricevi un addebito inatteso**, contatta subito il supporto di fatturazione: sia AWS sia Google Cloud
    hanno procedure per gli errori in buona fede, soprattutto per i nuovi clienti.

## La procedura passo per passo

La creazione degli account è spiegata in dettaglio, con budget e MFA, nella
[lezione 4.3 · Account, progetti e fatturazione](../fase-4/03-account.md).

## Fonti

- [AWS Free Tier](https://aws.amazon.com/free/)
- [AWS What's New: crediti e piano gratuito di 6 mesi (luglio 2025)](https://aws.amazon.com/about-aws/whats-new/2025/07/aws-free-tier-credits-month-free-plan/)
- [Google Cloud: funzionalità gratuite e prova gratuita](https://docs.cloud.google.com/free)
- [Google Cloud: domande frequenti sulla prova gratuita](https://cloud.google.com/signup-faqs)
