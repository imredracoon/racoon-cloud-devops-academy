# Google Cloud Associate Cloud Engineer

La certificazione tecnica di base di Google Cloud. Verifica che sai **creare, configurare e gestire** progetti,
risorse di calcolo, rete, storage, Kubernetes, monitoraggio e permessi, dalla console e soprattutto con `gcloud`.

| | |
|---|---|
| **Ente** | Google Cloud |
| **Formato** | 50-60 domande a scelta singola o multipla |
| **Durata** | 2 ore |
| **Punteggio** | superato o non superato; la soglia non è pubblicata |
| **Costo** | 125 USD (più tasse locali) |
| **Validità** | 3 anni |
| **Modalità** | online con sorvegliante oppure in un centro d'esame |
| **Pagina ufficiale** | [cloud.google.com/learn/certification/cloud-engineer](https://cloud.google.com/learn/certification/cloud-engineer) |

Formato e prezzo verificati il 28 ottobre 2026. Google consiglia almeno sei mesi di esperienza pratica: il corso te
la dà con i laboratori delle Fasi 4-8 nella scheda Google Cloud.

## Quando sostenerla

Dopo la [Fase 8](../../fase-8/index.md): l'esame include GKE e Cloud Run. Se lavori soprattutto con Google Cloud,
puoi anticiparla dopo la [Fase 6](../../fase-6/index.md) e ripassare Kubernetes con le lezioni [8.1](../../fase-8/01-architettura.md)-[8.5](../../fase-8/05-service.md) e [8.15](../../fase-8/15-gitops-cloud.md).

## Piano di studio (6 settimane)

| Settimana | Attività |
|---|---|
| 1 | Sezione 1: progetti, gerarchia, IAM di base, API, quote, fatturazione e budget ([4.3](../../fase-4/03-account.md), [4.4](../../fase-4/04-iam.md), [4.9](../../fase-4/09-costi.md), [4.10](../../fase-4/10-cli.md)) |
| 2-3 | Sezione 2: Compute Engine, MIG, GKE, Cloud Run, storage e database, VPC, bilanciatori, Terraform ([4.5](../../fase-4/05-compute.md)-[4.7](../../fase-4/07-rete.md), [5.1](../../fase-5/01-vpc.md)-[5.4](../../fase-5/04-database.md), [5.9](../../fase-5/09-container-serverless.md), [Fase 6](../../fase-6/index.md)) |
| 4 | Sezione 3: operazioni su VM, GKE, Cloud Run, Cloud SQL, Cloud Storage; Monitoring e Logging ([5.8](../../fase-5/08-monitoraggio.md), [8.2](../../fase-8/02-kubectl.md)-[8.8](../../fase-8/08-scheduling.md), [8.15](../../fase-8/15-gitops-cloud.md)) |
| 5 | Sezione 4: ruoli, ruoli personalizzati, service account, impersonificazione, Workload Identity Federation ([4.4](../../fase-4/04-iam.md), [5.6](../../fase-5/06-sicurezza.md), [7.4](../../fase-7/04-deploy-oidc.md)) |
| 6 | Esami completi, ripasso degli errori, prove pratiche in un progetto di prova |

!!! tip "Usa davvero gcloud"
    Molte domande chiedono il comando o l'opzione esatta. Ripeti i laboratori della scheda Google Cloud con
    `gcloud` invece che con la console, e leggi l'output: è il modo più veloce per ricordare i nomi giusti.

## Risorse ufficiali gratuite

- [Guida d'esame](https://cloud.google.com/learn/certification/guides/cloud-engineer).
- Percorso *Associate Cloud Engineer* su [Google Cloud Skills Boost](https://www.cloudskillsboost.google/).
- [Riferimento di gcloud](https://cloud.google.com/sdk/gcloud/reference).
- [Livello gratuito di Google Cloud](https://cloud.google.com/free/docs/free-cloud-features) per esercitarti con
  costi minimi.

## In questa appendice

- [Argomenti specifici](argomenti.md)
- [Simulazione d'esame](simulazione.md): oltre 110 domande originali
