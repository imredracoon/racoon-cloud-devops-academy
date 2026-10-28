# Argomenti specifici · AWS Solutions Architect Associate

Gli obiettivi sono quelli della guida ufficiale SAA-C03. Ogni riga indica cosa devi saper decidere e dove lo hai
studiato nel corso.

## Dominio 1 · Design Secure Architectures (30%)

| Obiettivo | Cosa devi saper decidere | Lezioni |
|---|---|---|
| 1.1 Secure access to AWS resources | utenti, gruppi e ruoli IAM; policy basate su identità e su risorsa; ruoli tra account; federazione e IAM Identity Center; Organizations e Service Control Policy; protezione dell'utente root | [4.4](../../fase-4/04-iam.md), [5.6](../../fase-5/06-sicurezza.md), [7.4](../../fase-7/04-deploy-oidc.md) |
| 1.2 Secure workloads and applications | subnet pubbliche e private, security group e NACL, endpoint VPC, NAT gateway, WAF, Shield, Cognito, Secrets Manager | [5.1](../../fase-5/01-vpc.md), [5.6](../../fase-5/06-sicurezza.md) |
| 1.3 Data security controls | cifratura con KMS (chiavi gestite da AWS e dal cliente), cifratura in transito con ACM, policy dei bucket, Block Public Access, Object Lock, backup e retention | [4.6](../../fase-4/06-storage.md), [5.3](../../fase-5/03-storage.md), [5.6](../../fase-5/06-sicurezza.md) |

## Dominio 2 · Design Resilient Architectures (26%)

| Obiettivo | Cosa devi saper decidere | Lezioni |
|---|---|---|
| 2.1 Scalable and loosely coupled architectures | SQS, SNS, EventBridge, Step Functions; API Gateway e Lambda; container con ECS ed EKS; cache; microservizi; scalabilità orizzontale | [5.5](../../fase-5/05-messaggistica.md), [5.9](../../fase-5/09-container-serverless.md), [5.2](../../fase-5/02-ha-scaling.md) |
| 2.2 Highly available and/or fault-tolerant architectures | Multi-AZ, Auto Scaling e health check, Route 53 (failover, routing), RDS Multi-AZ e repliche, Aurora globale, strategie di disaster recovery con RPO e RTO | [5.2](../../fase-5/02-ha-scaling.md), [5.4](../../fase-5/04-database.md), [5.7](../../fase-5/07-dns-cdn.md) |

## Dominio 3 · Design High-Performing Architectures (24%)

| Obiettivo | Cosa devi saper decidere | Lezioni |
|---|---|---|
| 3.1 High-performing storage | S3 (prefissi, multipart, Transfer Acceleration), tipi di volume EBS (gp3, io2), EFS e FSx, storage di istanza | [4.6](../../fase-4/06-storage.md), [5.3](../../fase-5/03-storage.md) |
| 3.2 Elastic compute | famiglie di istanze, Auto Scaling, Lambda e i suoi limiti, Fargate, placement group | [4.5](../../fase-4/05-compute.md), [5.2](../../fase-5/02-ha-scaling.md), [5.9](../../fase-5/09-container-serverless.md) |
| 3.3 High-performing database | RDS contro Aurora contro DynamoDB, repliche di lettura, RDS Proxy, ElastiCache, DynamoDB Accelerator | [5.4](../../fase-5/04-database.md) |
| 3.4 High-performing network | CloudFront, Global Accelerator, VPC peering, Transit Gateway, Direct Connect, VPN | [5.1](../../fase-5/01-vpc.md), [5.7](../../fase-5/07-dns-cdn.md) |
| 3.5 Data ingestion and transformation | Kinesis Data Streams e Data Firehose, Glue, Athena, Lake Formation, EMR | [4.8](../../fase-4/08-database-ai.md), [5.5](../../fase-5/05-messaggistica.md) |

## Dominio 4 · Design Cost-Optimized Architectures (20%)

| Obiettivo | Cosa devi saper decidere | Lezioni |
|---|---|---|
| 4.1 Cost-optimized storage | classi S3 e regole di lifecycle, S3 Intelligent-Tiering, tipi EBS, snapshot, backup | [5.3](../../fase-5/03-storage.md), [4.9](../../fase-4/09-costi.md) |
| 4.2 Cost-optimized compute | Spot, Savings Plans, Reserved Instances, dimensionamento con Compute Optimizer, serverless | [4.9](../../fase-4/09-costi.md), [5.10](../../fase-5/10-well-architected.md) |
| 4.3 Cost-optimized database | capacità on-demand o provisioned di DynamoDB, Aurora Serverless, dimensionamento di RDS | [5.4](../../fase-5/04-database.md) |
| 4.4 Cost-optimized network | NAT gateway e costi di trasferimento, endpoint VPC gateway (gratuiti per S3 e DynamoDB), CloudFront per ridurre il traffico in uscita | [5.1](../../fase-5/01-vpc.md), [5.7](../../fase-5/07-dns-cdn.md) |

## Dove sbagliano più spesso i candidati

- **"LEAST operational overhead"** indica quasi sempre un servizio gestito o serverless, anche se costa un po' di più.
- **Multi-AZ contro repliche di lettura**: Multi-AZ serve alla disponibilità (failover automatico), le repliche
  alla scalabilità delle letture.
- **SQS contro SNS contro EventBridge**: coda con un consumatore (SQS), notifica a molti iscritti (SNS), routing di
  eventi con regole (EventBridge).
- **Endpoint gateway contro interface**: S3 e DynamoDB hanno endpoint gateway gratuiti; gli altri servizi usano
  endpoint interface (PrivateLink), a pagamento.
- **Strategie di DR**: più l'RTO è basso, più la strategia è costosa (da backup and restore a multi-site attivo).
