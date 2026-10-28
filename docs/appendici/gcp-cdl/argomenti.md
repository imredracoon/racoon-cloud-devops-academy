# Argomenti specifici · Google Cloud Digital Leader

Le sei sezioni della guida ufficiale, con i temi principali e le lezioni del corso dove li trovi.

| Sezione (peso) | Temi principali | Lezioni |
|---|---|---|
| 1 · Digital Transformation with Google Cloud (~18%) | perché il cloud cambia il modo di lavorare, modelli IaaS, PaaS, SaaS, cloud ibrido e multicloud, infrastruttura globale di Google (regioni, zone, rete), costo totale di possesso, spesa operativa contro investimenti, sostenibilità | [0.6](../../fase-0/06-cloud-devops.md), [4.1](../../fase-4/01-cloud-computing.md), [4.2](../../fase-4/02-infrastruttura-globale.md) |
| 2 · Exploring Data Transformation with Google Cloud (~18%) | valore dei dati, dati strutturati e non strutturati, database (Cloud SQL, AlloyDB, Spanner, Firestore, Bigtable), data warehouse e data lake (BigQuery, Cloud Storage), elaborazione in streaming (Pub/Sub, Dataflow), business intelligence con Looker | [4.6](../../fase-4/06-storage.md), [4.8](../../fase-4/08-database-ai.md), [5.4](../../fase-5/04-database.md), [5.5](../../fase-5/05-messaggistica.md) |
| 3 · Innovating with Google Cloud Artificial Intelligence (~18%) | differenza tra AI, machine learning e AI generativa, qualità dei dati, AI responsabile, API pronte (Vision, Speech, Translation, Natural Language, Document AI), Vertex AI, Model Garden, Gemini e gli agenti, BigQuery ML | [4.8](../../fase-4/08-database-ai.md) |
| 4 · Modernize Infrastructure and Applications with Google Cloud (~18%) | Compute Engine, container e GKE, serverless (Cloud Run, Cloud Run functions), API (Apigee), migrazione (Migration Center, Migrate to Virtual Machines, VMware Engine, Bare Metal Solution), ibrido con GKE Enterprise e Google Distributed Cloud | [4.5](../../fase-4/05-compute.md), [5.9](../../fase-5/09-container-serverless.md), [8.15](../../fase-8/15-gitops-cloud.md) |
| 5 · Trust and Security with Google Cloud (~18%) | responsabilità condivisa e "destino condiviso", difesa in profondità, IAM e privilegio minimo, cifratura e Cloud KMS, zero trust e Identity-Aware Proxy, Security Command Center, Google SecOps, conformità e sovranità dei dati (Assured Workloads, Access Transparency) | [4.4](../../fase-4/04-iam.md), [5.6](../../fase-5/06-sicurezza.md) |
| 6 · Scaling with Google Cloud Operations (~10%) | gestione dei costi (budget, etichette, sconti per uso continuativo e impegni), affidabilità e SRE (SLI, SLO, error budget, cultura blameless), Google Cloud Observability, piani di assistenza, sostenibilità | [4.9](../../fase-4/09-costi.md), [5.8](../../fase-5/08-monitoraggio.md), [7.5](../../fase-7/05-osservabilita.md) |

I pesi sono quelli indicati dalla guida ufficiale al momento della verifica; controlla la guida aggiornata.

## Dove sbagliano più spesso i candidati

- **Spanner, Bigtable, Firestore, BigQuery**: database relazionale globale, NoSQL ad altissimo volume, documenti per
  app mobili e web, analisi su grandi volumi. Associa ognuno al suo caso d'uso tipico.
- **Compute Engine, GKE, Cloud Run**: passando dal primo all'ultimo hai meno controllo sull'infrastruttura ma
  anche meno lavoro di gestione. Lo scenario indica quale compromesso serve.
- **Budget**: avvisano, ma **non bloccano** la spesa da soli.
- **SLO ed error budget**: l'error budget è la quantità di inaffidabilità che ci si può permettere; quando finisce,
  si privilegia la stabilità rispetto alle nuove funzioni.
