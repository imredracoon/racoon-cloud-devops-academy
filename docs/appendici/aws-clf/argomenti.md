# Argomenti specifici · AWS Cloud Practitioner

La guida d'esame divide il programma in quattro domini. Per ogni obiettivo trovi cosa devi saper fare e le lezioni
del corso che lo trattano. I nomi degli obiettivi sono quelli ufficiali, in inglese, per ritrovarli nella guida.

## Dominio 1 · Cloud Concepts (24%)

| Obiettivo | Cosa devi sapere | Lezioni |
|---|---|---|
| 1.1 Benefits of the AWS Cloud | economie di scala, elasticità, agilità, portata globale, pagamento a consumo, alta disponibilità | [4.1](../../fase-4/01-cloud-computing.md), [4.2](../../fase-4/02-infrastruttura-globale.md) |
| 1.2 Design principles of the AWS Cloud | i sei pilastri del Well-Architected Framework e cosa chiede ciascuno | [5.10](../../fase-5/10-well-architected.md) |
| 1.3 Benefits of and strategies for migration | AWS Cloud Adoption Framework (CAF), le strategie di migrazione (rehost, replatform, refactor...), servizi di trasferimento dati come la famiglia Snow e DataSync | [4.1](../../fase-4/01-cloud-computing.md), [5.3](../../fase-5/03-storage.md) |
| 1.4 Concepts of cloud economics | costi fissi contro variabili, costo totale di possesso, licenze (BYOL), dimensionamento corretto, automazione | [4.9](../../fase-4/09-costi.md) |

## Dominio 2 · Security and Compliance (30%)

| Obiettivo | Cosa devi sapere | Lezioni |
|---|---|---|
| 2.1 Shared responsibility model | cosa spetta ad AWS e cosa al cliente, e come cambia tra EC2, RDS, Lambda e S3 | [4.4](../../fase-4/04-iam.md), [5.6](../../fase-5/06-sicurezza.md) |
| 2.2 Security, governance, and compliance concepts | AWS Artifact, conformità, cifratura a riposo e in transito, log e audit con CloudTrail e Config | [5.6](../../fase-5/06-sicurezza.md), [5.8](../../fase-5/08-monitoraggio.md) |
| 2.3 Access management capabilities | utente root e sua protezione, IAM (utenti, gruppi, ruoli, policy), MFA, privilegio minimo, IAM Identity Center | [4.4](../../fase-4/04-iam.md) |
| 2.4 Components and resources for security | security group e network ACL, AWS WAF, Shield, GuardDuty, Inspector, Security Hub, KMS, Secrets Manager, Trusted Advisor | [5.1](../../fase-5/01-vpc.md), [5.6](../../fase-5/06-sicurezza.md) |

## Dominio 3 · Cloud Technology and Services (34%)

| Obiettivo | Cosa devi sapere | Lezioni |
|---|---|---|
| 3.1 Methods of deploying and operating | console, CLI, SDK, IaC (CloudFormation), modalità di connessione (VPN, Direct Connect) | [4.10](../../fase-4/10-cli.md), [6.1](../../fase-6/01-iac.md) |
| 3.2 AWS global infrastructure | Region, Availability Zone, edge location, Local Zones, criteri di scelta della Region | [4.2](../../fase-4/02-infrastruttura-globale.md) |
| 3.3 Compute services | EC2 e famiglie di istanze, Auto Scaling, Elastic Load Balancing, container (ECS, EKS, Fargate), Lambda | [4.5](../../fase-4/05-compute.md), [5.2](../../fase-5/02-ha-scaling.md), [5.9](../../fase-5/09-container-serverless.md) |
| 3.4 Database services | RDS, Aurora, DynamoDB, ElastiCache, database gestiti contro installati su EC2, Database Migration Service | [4.8](../../fase-4/08-database-ai.md), [5.4](../../fase-5/04-database.md) |
| 3.5 Network services | VPC, subnet, gateway, Route 53, CloudFront, Global Accelerator | [4.7](../../fase-4/07-rete.md), [5.1](../../fase-5/01-vpc.md), [5.7](../../fase-5/07-dns-cdn.md) |
| 3.6 Storage services | S3 e classi di storage, EBS, EFS, FSx, Storage Gateway, AWS Backup | [4.6](../../fase-4/06-storage.md), [5.3](../../fase-5/03-storage.md) |
| 3.7 AI/ML and data analytics | servizi di AI pronti (Rekognition, Comprehend, Transcribe...), SageMaker AI, Amazon Bedrock, Athena, Kinesis, Glue, QuickSight | [4.8](../../fase-4/08-database-ai.md) |
| 3.8 Other in-scope service categories | SQS, SNS, EventBridge, strumenti per sviluppatori, servizi per utenti finali e IoT ai fini del riconoscimento | [5.5](../../fase-5/05-messaggistica.md) |

## Dominio 4 · Billing, Pricing, and Support (12%)

| Obiettivo | Cosa devi sapere | Lezioni |
|---|---|---|
| 4.1 Pricing models | On-Demand, Savings Plans, Reserved Instances, Spot, Dedicated Hosts, costi di trasferimento dati | [4.9](../../fase-4/09-costi.md) |
| 4.2 Billing, budget, and cost management | AWS Budgets, Cost Explorer, Cost and Usage Report, tag di allocazione dei costi, AWS Organizations e fatturazione consolidata, Pricing Calculator | [4.3](../../fase-4/03-account.md), [4.9](../../fase-4/09-costi.md) |
| 4.3 Technical resources and support options | piani di supporto (vedi l'[avviso sui nuovi piani](index.md#cosa-sapere-sui-cambiamenti-recenti)), documentazione, re:Post, Trusted Advisor, Health Dashboard, AWS Marketplace, partner | [4.9](../../fase-4/09-costi.md) |

## Dove sbagliano più spesso i candidati

- **Responsabilità condivisa**: con RDS l'aggiornamento del sistema operativo è di AWS; con EC2 è tuo. La
  configurazione dei security group e dei dati è sempre tua.
- **Savings Plans e Reserved Instances**: entrambi richiedono un impegno di 1 o 3 anni; i Savings Plans si basano su
  una spesa oraria e sono più flessibili.
- **CloudTrail contro CloudWatch contro Config**: chi ha fatto cosa (CloudTrail), metriche e log (CloudWatch),
  come era configurata una risorsa nel tempo (Config).
- **Shield contro WAF**: Shield protegge dagli attacchi DDoS; WAF filtra le richieste web secondo regole.
- **Region e AZ**: un'AZ è uno o più data center separati all'interno di una Region; per l'alta disponibilità si
  distribuisce su più AZ, per il disaster recovery su più Region.
