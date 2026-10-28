# Argomenti specifici · Google Associate Cloud Engineer

Le quattro sezioni della guida ufficiale, divise nei loro sotto-obiettivi.

## Sezione 1 · Setting up a cloud solution environment (~20%)

| Sotto-obiettivo | Cosa devi saper fare | Lezioni |
|---|---|---|
| Progetti e account | gerarchia organizzazione, cartelle, progetti; policy dell'organizzazione; ruoli IAM a livello di progetto; abilitare le API; quote; Cloud Identity e gruppi; Gemini Cloud Assist | [4.3](../../fase-4/03-account.md), [4.4](../../fase-4/04-iam.md), [4.10](../../fase-4/10-cli.md) |
| Fatturazione | account di fatturazione e collegamento ai progetti, budget e avvisi, esportazione della fatturazione in BigQuery, etichette | [4.3](../../fase-4/03-account.md), [4.9](../../fase-4/09-costi.md) |

## Sezione 2 · Planning and implementing a cloud solution (~30%)

| Sotto-obiettivo | Cosa devi saper fare | Lezioni |
|---|---|---|
| Risorse di calcolo | scegliere tra Compute Engine, GKE (Standard e Autopilot), Cloud Run e Cloud Run functions; tipi di macchina, VM Spot, template e gruppi di istanze gestiti | [4.5](../../fase-4/05-compute.md), [5.2](../../fase-5/02-ha-scaling.md), [5.9](../../fase-5/09-container-serverless.md) |
| Dati e storage | Cloud Storage (classi, posizioni, lifecycle), Cloud SQL, AlloyDB, Spanner, Firestore, Bigtable, BigQuery, Memorystore; caricare dati | [4.6](../../fase-4/06-storage.md), [5.3](../../fase-5/03-storage.md), [5.4](../../fase-5/04-database.md) |
| Rete | VPC e subnet in modalità custom, regole firewall, Shared VPC, peering, Cloud NAT, Private Google Access, bilanciatori, Cloud VPN e Interconnect, Cloud DNS | [4.7](../../fase-4/07-rete.md), [5.1](../../fase-5/01-vpc.md), [5.7](../../fase-5/07-dns-cdn.md) |
| Infrastruttura come codice | Terraform con il provider Google, Infrastructure Manager, Config Connector, Cloud Foundation Fabric | [6.1](../../fase-6/01-iac.md)-[6.7](../../fase-6/07-moduli.md) |

## Sezione 3 · Ensuring successful operation of a cloud solution (~30%)

| Sotto-obiettivo | Cosa devi saper fare | Lezioni |
|---|---|---|
| Compute Engine | accesso con OS Login e IAP, snapshot e pianificazioni, immagini, aggiornamenti dei MIG | [4.5](../../fase-4/05-compute.md), [5.2](../../fase-5/02-ha-scaling.md) |
| GKE | `get-credentials`, pool di nodi, autoscaling, distribuire e scalare workload | [8.2](../../fase-8/02-kubectl.md), [8.4](../../fase-8/04-deployment.md), [8.8](../../fase-8/08-scheduling.md), [8.15](../../fase-8/15-gitops-cloud.md) |
| Cloud Run | revisioni, suddivisione del traffico, istanze minime e concorrenza | [5.9](../../fase-5/09-container-serverless.md) |
| Storage e database | versioni degli oggetti, URL firmati, alta disponibilità e backup di Cloud SQL, stima dei costi di BigQuery | [5.3](../../fase-5/03-storage.md), [5.4](../../fase-5/04-database.md) |
| Rete | modificare subnet e regole firewall, diagnosticare connettività | [5.1](../../fase-5/01-vpc.md), [1.9](../../fase-1/09-diagnostica.md) |
| Monitoraggio e log | Ops Agent, avvisi e controlli di uptime, ambiti delle metriche, sink e bucket dei log, log di audit, Trace e Profiler | [5.8](../../fase-5/08-monitoraggio.md), [7.5](../../fase-7/05-osservabilita.md) |

## Sezione 4 · Configuring access and security (~20%)

| Sotto-obiettivo | Cosa devi saper fare | Lezioni |
|---|---|---|
| IAM | ruoli di base, predefiniti e personalizzati; concedere e revocare ruoli; condizioni IAM; privilegio minimo | [4.4](../../fase-4/04-iam.md) |
| Service account | creare e collegare service account, impersonificazione, evitare le chiavi, Workload Identity Federation (anche per GKE) | [4.4](../../fase-4/04-iam.md), [5.6](../../fase-5/06-sicurezza.md), [7.4](../../fase-7/04-deploy-oidc.md) |

## Dove sbagliano più spesso i candidati

- **Priorità del firewall**: numero più basso = priorità più alta. Le regole implicite bloccano tutto l'ingresso e
  permettono tutta l'uscita.
- **VPC globale, subnet regionali**: in AWS la VPC è regionale, in Google Cloud è globale.
- **Ruolo Service Account User**: serve per collegare un service account a una risorsa (permesso `actAs`).
- **Log di audit Data Access**: sono disattivati per impostazione predefinita (tranne BigQuery).
- **`LIMIT` in BigQuery** non riduce i byte letti: servono colonne selezionate e partizioni.
