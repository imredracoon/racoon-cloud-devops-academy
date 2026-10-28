# 4.9 Costi, supporto e buone pratiche

!!! obiettivi "Obiettivi"
    - Capire da cosa dipende il costo del cloud: calcolo, storage, trasferimento dati, richieste.
    - Stimare i costi prima di costruire e controllarli dopo.
    - Usare etichette, budget e strumenti di analisi della spesa.
    - Conoscere i piani di supporto e le risorse di aiuto.
    - Conoscere i framework di buona progettazione: AWS Well-Architected e Google Cloud Architecture Framework.

## Come si paga il cloud

I principi comuni:

- **paghi ciò che usi**: secondi di calcolo, GB al mese, milioni di richieste;
- il **traffico in entrata** è in genere **gratuito**; il traffico in **uscita** verso Internet (*egress*) si paga, e
  anche quello **tra zone** o **tra regioni**;
- più usi, meno paghi per unità (sconti a scaglioni sullo storage, sconti per impegni sul calcolo);
- le risorse **ferme ma esistenti** costano: un disco non collegato, un IP pubblico, uno snapshot, un NAT gateway.

!!! warning "Le fonti di costo che sorprendono"
    - **IP pubblici IPv4**: su AWS ogni IPv4 pubblico si paga a ore, anche se collegato a un'istanza accesa.
    - **NAT gateway**: costo orario più costo per GB elaborato.
    - **Trasferimento dati in uscita** e tra zone.
    - **Snapshot e backup** dimenticati.
    - **Log** conservati per sempre.

## Stimare prima di costruire

- **AWS Pricing Calculator** ([calculator.aws](https://calculator.aws/)) e **Google Cloud Pricing Calculator**: stimano
  il costo mensile di un'architettura. Sono il primo passo di ogni progetto.
- Le pagine dei prezzi di ogni servizio, per regione.

## Controllare dopo

=== "AWS"

    | Strumento | Funzione |
    |---|---|
    | **Billing and Cost Management** | la home della fatturazione: fatture, pagamenti, crediti |
    | **Cost Explorer** | grafici e analisi della spesa, per servizio, account, etichetta; previsioni |
    | **Budgets** | soglie con avvisi e azioni |
    | **Cost Anomaly Detection** | avvisi automatici sulle spese anomale, con il machine learning |
    | **Data Exports** (Cost and Usage Report) | il dettaglio più completo della spesa, esportato su S3 per l'analisi |
    | **Cost Allocation Tags** | etichette attivate per ripartire i costi |
    | **Compute Optimizer** | suggerisce dimensioni più adatte per istanze e altre risorse |
    | **Trusted Advisor** | controlli su costi, prestazioni, sicurezza, tolleranza ai guasti, limiti di servizio |
    | **Organizations**: fatturazione consolidata | un'unica fattura e sconti per volume su tutti gli account |

=== "Google Cloud"

    | Strumento | Funzione |
    |---|---|
    | **Fatturazione** (Cloud Billing) | account di fatturazione, report, crediti |
    | **Report** di fatturazione | grafici della spesa per progetto, servizio, etichetta; previsioni |
    | **Budget e avvisi** | soglie con email e notifiche Pub/Sub (per automatizzare azioni) |
    | **Esportazione della fatturazione in BigQuery** | il dettaglio completo, interrogabile con SQL |
    | **Recommender** e **Active Assist** | suggerimenti: VM sottoutilizzate, dischi inutilizzati, impegni convenienti |
    | **Quote** | limiti di utilizzo per progetto, che proteggono anche dagli errori |

### Le etichette

Le **etichette** (*tag* su AWS, *label* su Google Cloud) sono coppie chiave-valore sulle risorse:
`progetto=academy`, `ambiente=prod`, `team=web`. Servono a ripartire i costi, trovare e automatizzare le risorse.
Definisci una **politica di etichettatura** fin dal primo giorno: aggiungerle dopo è faticoso.

### FinOps

**FinOps** è la disciplina che unisce finanza, tecnica e business per gestire i costi del cloud: visibilità (chi spende
cosa), ottimizzazione (dimensioni giuste, impegni, spegnere ciò che non serve), responsabilità dei team. Le leve
principali:

1. **spegnere** ciò che non serve (ambienti di sviluppo la notte e nel fine settimana);
2. **dimensionare correttamente** (*rightsizing*);
3. usare **impegni** per la parte stabile e **spot** per quella interrompibile;
4. scegliere la **classe di storage** giusta con i cicli di vita;
5. preferire i servizi **gestiti** e **serverless**, che non costano quando non lavorano.

## Il supporto

=== "AWS"

    A dicembre 2025 AWS ha rinnovato i piani di supporto:

    | Piano | Caratteristiche principali |
    |---|---|
    | **Basic** | gratuito: documentazione, forum, AWS re:Post, controlli base di Trusted Advisor, AWS Health Dashboard, supporto per fatturazione e account |
    | **Business Support+** | da 29 USD al mese: supporto tecnico 24/7 con assistenza basata sull'AI e tecnici AWS |
    | **Enterprise Support** | da 5.000 USD al mese: tempi di risposta più rapidi, un **Technical Account Manager** (TAM) dedicato, guida proattiva |
    | **Unified Operations Support** | da 50.000 USD al mese: per i carichi più critici, con un team dedicato |

    I piani precedenti (**Developer**, **Business** classico, **Enterprise On-Ramp**) restano attivi per chi li ha
    sottoscritti fino al 1° gennaio 2027. Le domande d'esame scritte prima del cambiamento possono ancora citarli:
    ricorda che il TAM dedicato è tipico dei piani Enterprise e che il piano Basic è gratuito.

    Altre risorse: **AWS re:Post** (domande e risposte della comunità), **AWS Knowledge Center**, **AWS Health
    Dashboard** (lo stato dei servizi e gli eventi che riguardano il tuo account), **AWS Marketplace** (software di
    terze parti), **AWS Partner Network**, **AWS Professional Services**.

=== "Google Cloud"

    Google Cloud offre un supporto **Basic** gratuito (fatturazione, documentazione, comunità) e piani a pagamento,
    oggi chiamati **Standard**, **Enhanced** e **Premium**, con tempi di risposta crescenti; il piano Premium include un
    **Technical Account Manager**. Controlla piani e prezzi aggiornati nella pagina ufficiale.

    Altre risorse: **Google Cloud Community**, **Service Health** (stato dei servizi), **Google Cloud Marketplace**,
    partner certificati.

## Progettare bene: i framework

Entrambi i fornitori pubblicano un insieme di principi e domande per valutare un'architettura. Li studierai a fondo
nella Fase 5; per ora ricorda i nomi dei pilastri.

| AWS Well-Architected Framework | Google Cloud Well-Architected Framework |
|---|---|
| Eccellenza operativa | Eccellenza operativa |
| Sicurezza | Sicurezza, privacy e conformità |
| Affidabilità | Affidabilità |
| Efficienza delle prestazioni | Ottimizzazione delle prestazioni |
| Ottimizzazione dei costi | Ottimizzazione dei costi |
| Sostenibilità | Sostenibilità |

Il framework di Google aggiunge prospettive trasversali, per esempio su AI e ML. L'**AWS Well-Architected Tool** nella console guida la revisione di un carico di lavoro con le domande dei pilastri.

## Conformità

I fornitori ottengono certificazioni e attestazioni (ISO 27001, SOC 1/2/3, PCI DSS, e per l'Italia la qualificazione
dei servizi cloud per la Pubblica Amministrazione da parte dell'**ACN**, l'Agenzia per la Cybersicurezza Nazionale).
I documenti sono disponibili ai clienti tramite **AWS Artifact** e il **Compliance Reports Manager** di Google Cloud.
Ricorda la responsabilità condivisa: la certificazione del fornitore non rende automaticamente conforme la tua
applicazione.

## Esercizi

!!! esercizio "Esercizio 1 · Stima un'architettura"
    Con il calcolatore del tuo cloud stima il costo mensile, nella regione del corso, di: due VM piccole accese sempre
    (`t3.small` / `e2-small`), un database PostgreSQL gestito piccolo con alta disponibilità, 50 GB di storage a oggetti
    Standard, 100 GB al mese di traffico in uscita. Qual è la voce più costosa? Come la ridurresti?

??? success "Soluzione"
    I valori esatti dipendono da regione e listino del momento: l'importante è il metodo. Di solito la voce più
    pesante è il **database con alta disponibilità** (paghi due istanze, la principale e lo standby), seguita dalle VM.
    Leve: dimensioni più piccole dopo aver misurato l'uso reale, impegni a 1 anno per le risorse stabili, alta
    disponibilità solo in produzione (non negli ambienti di test), per le VM eventualmente un servizio serverless.
    Salva la stima (il calcolatore genera un link condivisibile): è un ottimo documento per il portfolio.

!!! esercizio "Esercizio 2 · Caccia allo spreco"
    Il conto di un team è raddoppiato. Elenca cinque controlli che faresti, e con quali strumenti.

??? success "Soluzione"
    1. Quali servizi e progetti/account sono cresciuti: **Cost Explorer** / **report di fatturazione**, raggruppando
       per servizio ed etichetta.
    2. Risorse inutilizzate: dischi non collegati, snapshot vecchi, IP pubblici liberi, load balancer senza traffico
       (**Trusted Advisor**, **Recommender**).
    3. VM sovradimensionate: **Compute Optimizer** / suggerimenti di dimensionamento di **Recommender**.
    4. Traffico in uscita e tra zone anomalo: dettaglio dei costi di trasferimento dati.
    5. Ambienti di sviluppo accesi di notte: etichette e pianificazione dello spegnimento.

    E per il futuro: **budget** per team e **rilevamento delle anomalie**.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale traffico è in genere gratuito nel cloud?","options":["In uscita verso Internet","In entrata dall'Internet","Tra regioni","Tra zone"],"answer":[1],"explanation":"Si pagano l'uscita e i trasferimenti tra zone e regioni."},
 {"q":"Quale strumento usi per stimare i costi PRIMA di costruire un'architettura?","options":["Cost Explorer","Pricing Calculator","Trusted Advisor","Budgets"],"answer":[1],"explanation":"Cost Explorer analizza la spesa già sostenuta."},
 {"q":"A cosa servono le etichette (tag/label) sulle risorse? (scegli 2)","options":["Ripartire i costi per progetto, team o ambiente","Trovare e automatizzare le risorse","Aumentare le prestazioni","Cifrare i dati","Ridurre la latenza"],"answer":[0,1],"explanation":"Vanno pianificate dal primo giorno."},
 {"q":"Quale servizio AWS fa controlli automatici su costi, sicurezza, prestazioni e limiti di servizio?","options":["AWS Artifact","Trusted Advisor","Amazon Inspector","AWS Config"],"answer":[1],"explanation":"Con il piano Basic sono disponibili solo i controlli di base."},
 {"q":"Dove trovi i report di conformità di AWS, come ISO 27001 e SOC?","options":["AWS Artifact","AWS Health Dashboard","AWS Marketplace","Amazon Macie"],"answer":[0],"explanation":"Su Google Cloud: Compliance Reports Manager."},
 {"q":"Quale piano di supporto AWS è gratuito?","options":["Basic","Business Support+","Enterprise Support","Unified Operations Support"],"answer":[0],"explanation":"Gli altri sono a pagamento."}
]}
</script>

## Fonti

- [AWS: How AWS Pricing Works (whitepaper)](https://docs.aws.amazon.com/whitepapers/latest/how-aws-pricing-works/how-aws-pricing-works.html)
- [AWS Cost Management](https://aws.amazon.com/aws-cost-management/)
- [AWS News Blog: New and enhanced AWS Support plans (dicembre 2025)](https://aws.amazon.com/blogs/aws/new-and-enhanced-aws-support-plans-add-ai-capabilities-to-expert-guidance)
- [AWS Well-Architected Framework](https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html)
- [AWS Artifact](https://aws.amazon.com/artifact/)
- [Google Cloud: Cloud Billing documentation](https://docs.cloud.google.com/billing/docs)
- [Google Cloud: Customer Care (piani di supporto)](https://cloud.google.com/support)
- [Google Cloud Well-Architected Framework](https://docs.cloud.google.com/architecture/framework)
- [FinOps Foundation: What is FinOps?](https://www.finops.org/introduction/what-is-finops/)
