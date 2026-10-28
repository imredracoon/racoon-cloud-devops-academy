# Verifica di fine fase 4

Passa alla Fase 5 quando ottieni **almeno l'80%** (20 risposte corrette su 25). Le spiegazioni indicano la lezione
da ripassare. Se vuoi sostenere l'esame AWS Cloud Practitioner o Google Cloud Digital Leader, prosegui con le
appendici: [AWS CLF-C02](../appendici/aws-clf/index.md) e [Google Cloud Digital Leader](../appendici/gcp-cdl/index.md).

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 4 · 25 domande","questions":[
 {"q":"Quale vantaggio del cloud descrive il passaggio da grandi investimenti iniziali a costi variabili?","options":["Economie di scala","Da CapEx a OpEx","Elasticità","Globalità"],"answer":[1],"explanation":"Lezione 4.1."},
 {"q":"Chi è responsabile delle patch del sistema operativo di un'istanza EC2?","options":["AWS","Il cliente","Il produttore del sistema operativo","Nessuno"],"answer":[1],"explanation":"EC2 è IaaS (lezione 4.1)."},
 {"q":"Chi è responsabile delle patch del motore di un database su Amazon RDS?","options":["Il cliente","AWS","Entrambi","Il produttore del database"],"answer":[1],"explanation":"È un servizio gestito: il cliente sceglie la finestra di manutenzione (lezione 4.1)."},
 {"q":"Un'azienda usa Gmail e Google Docs. Quale modello di servizio è?","options":["IaaS","PaaS","SaaS","On premises"],"answer":[2],"explanation":"Lezione 4.1."},
 {"q":"Come rendi un'applicazione resistente al guasto di un data center?","options":["Usando una VM più grande","Distribuendola su almeno due zone di disponibilità","Facendo backup settimanali","Usando una CDN"],"answer":[1],"explanation":"Lezione 4.2."},
 {"q":"Qual è il primo criterio per scegliere una regione?","options":["Il prezzo più basso","I requisiti di conformità e di legge","Il numero di servizi","La vicinanza al fornitore"],"answer":[1],"explanation":"Lezione 4.2."},
 {"q":"Quale servizio porta in cache i contenuti nei punti di presenza vicino agli utenti?","options":["Route 53","CloudFront","Direct Connect","EBS"],"answer":[1],"explanation":"Su Google Cloud: Cloud CDN (lezioni 4.2 e 4.7)."},
 {"q":"Cosa permette di imporre limiti massimi ai permessi di tutti gli account di un'organizzazione AWS?","options":["IAM groups","Service Control Policies","Security group","AWS Budgets"],"answer":[1],"explanation":"Lezione 4.3."},
 {"q":"Su Google Cloud, quale entità paga i costi di un progetto?","options":["La cartella","L'account di fatturazione collegato","Il service account","L'organizzazione, sempre"],"answer":[1],"explanation":"Lezione 4.3."},
 {"q":"Un budget ha superato la soglia del 100%. Cosa succede, con le impostazioni predefinite?","options":["Le risorse vengono fermate","Ricevi un avviso, ma le risorse continuano a funzionare","L'account viene chiuso","Il budget raddoppia"],"answer":[1],"explanation":"Lezione 4.3."},
 {"q":"Una funzione Lambda deve scrivere su una tabella DynamoDB. Come le dai il permesso?","options":["Con le access key di un amministratore nel codice","Con un ruolo IAM di esecuzione con una policy che permette solo le operazioni necessarie su quella tabella","Con l'utente root","Rendendo pubblica la tabella"],"answer":[1],"explanation":"Ruoli e privilegio minimo (lezione 4.4)."},
 {"q":"Su Google Cloud, a cosa si assegnano i ruoli?","options":["Ai bucket soltanto","Ai principal, tramite le policy di autorizzazione su organizzazione, cartelle, progetti o risorse","Alle zone","Ai tipi di macchina"],"answer":[1],"explanation":"Lezione 4.4."},
 {"q":"Quali sono buone pratiche per l'utente root di AWS?","options":["Attivare l'MFA","Non creare access key per root","Usarlo per il lavoro quotidiano","Condividerlo con il team","Usarlo nella CLI"],"answer":[0,1],"explanation":"Lezioni 4.3 e 4.4."},
 {"q":"Un lavoro di elaborazione video può essere interrotto e ripreso. Quale opzione riduce di più i costi?","options":["On demand","Spot","Dedicated Hosts","Savings Plans a 3 anni"],"answer":[1],"explanation":"Lezione 4.5."},
 {"q":"Quale servizio esegue codice in risposta a eventi, senza server da gestire, pagando solo l'esecuzione?","options":["EC2","AWS Lambda","EBS","Redshift"],"answer":[1],"explanation":"Su Google Cloud: Cloud Run functions (lezione 4.5)."},
 {"q":"Quale tipo di storage è un disco collegato a una sola VM?","options":["A oggetti","A blocchi","A file","Archivio"],"answer":[1],"explanation":"EBS, Persistent Disk (lezione 4.6)."},
 {"q":"Quale classe S3 è la più economica per archivi consultati quasi mai, con recupero in ore?","options":["S3 Standard","S3 Standard-IA","S3 Glacier Deep Archive","S3 Express One Zone"],"answer":[2],"explanation":"Lezione 4.6."},
 {"q":"Dieci server Linux devono condividere la stessa cartella via NFS. Quale servizio AWS usi?","options":["S3","EBS","EFS","Glacier"],"answer":[2],"explanation":"Su Google Cloud: Filestore (lezione 4.6)."},
 {"q":"Una VM in una subnet privata deve scaricare aggiornamenti da Internet senza essere raggiungibile dall'esterno. Cosa serve?","options":["Un Internet gateway nella sua route table","Un NAT gateway o Cloud NAT","Un security group aperto","Un IP pubblico"],"answer":[1],"explanation":"Lezione 4.7."},
 {"q":"Quale firewall di AWS è stateless e si applica alle subnet?","options":["Security group","Network ACL","AWS WAF","AWS Shield"],"answer":[1],"explanation":"Lezione 4.7."},
 {"q":"Quale servizio collega un data center aziendale ad AWS con una connessione fisica privata?","options":["Site-to-Site VPN","Direct Connect","CloudFront","Transit Gateway"],"answer":[1],"explanation":"Su Google Cloud: Cloud Interconnect (lezione 4.7)."},
 {"q":"Un'azienda vuole analizzare con SQL petabyte di dati senza gestire infrastruttura. Quale servizio Google Cloud sceglie?","options":["Cloud SQL","BigQuery","Firestore","Memorystore"],"answer":[1],"explanation":"Lezione 4.8."},
 {"q":"Quale servizio AWS offre accesso a modelli di base di più fornitori per costruire applicazioni di AI generativa?","options":["Amazon SageMaker AI","Amazon Bedrock","Amazon Rekognition","Amazon Comprehend"],"answer":[1],"explanation":"Lezione 4.8."},
 {"q":"Quale strumento analizza la spesa AWS già sostenuta, per servizio ed etichetta?","options":["Pricing Calculator","Cost Explorer","AWS Artifact","Trusted Advisor"],"answer":[1],"explanation":"Il calcolatore stima prima di costruire (lezione 4.9)."},
 {"q":"Con quale comando la AWS CLI ottiene credenziali temporanee da IAM Identity Center?","options":["aws configure","aws configure sso","aws iam create-access-key","aws sts assume-root"],"answer":[1],"explanation":"Lezione 4.10."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Passa alla [Fase 5 · Architetture](../fase-5/index.md).
