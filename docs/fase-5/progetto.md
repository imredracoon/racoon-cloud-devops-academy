# Progetto 5 · Applicazione a tre livelli

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai progettare, costruire e documentare un'architettura cloud reale: rete isolata su più zone, applicazione
    scalabile dietro un load balancer, database gestito privato, segreti fuori dal codice, monitoraggio, costi sotto
    controllo e una revisione con il Well-Architected Framework. È il progetto che un selezionatore si aspetta da chi
    si candida come Cloud Engineer junior.

| | |
|---|---|
| **Durata** | 20-30 ore |
| **Consegna** | Repository pubblico `tasks-api-cloud` con documentazione, script e diagramma |
| **Costo** | nei crediti per i nuovi clienti, **se distruggi tutto** a fine sessione; il README deve riportare la stima |

## Il compito

Porta in produzione la Tasks API della [Fase 3](../fase-3/index.md) con un'architettura a tre livelli:

1. **presentazione**: load balancer pubblico (con CDN facoltativa);
2. **applicazione**: la Tasks API in container, su almeno **due zone**, in subnet **private**;
3. **dati**: PostgreSQL **gestito**, privato, con backup automatici.

Puoi scegliere il livello applicativo tra due varianti. La variante A è più vicina all'esame SAA e ACE (VM e gruppi
di autoscaling), la variante B a come si lavora oggi (container gestiti). Se hai tempo, realizzale entrambe.

=== "AWS"

    ```mermaid
    flowchart LR
      U[Utenti] -->|HTTP/S| ALB[Application Load Balancer<br>subnet pubbliche A e B]
      subgraph VPC
        ALB --> APP1[Tasks API<br>subnet privata A]
        ALB --> APP2[Tasks API<br>subnet privata B]
        APP1 --> RDS[(RDS PostgreSQL<br>Multi-AZ, subnet dati)]
        APP2 --> RDS
      end
      APP1 -.-> SM[Secrets Manager]
      APP1 -.-> CW[CloudWatch]
    ```

    - **Rete** ([lezione 5.1](01-vpc.md)): VPC con subnet pubbliche, private e dati su **2 AZ**; un NAT gateway (uno solo, per
      contenere i costi: documenta il compromesso) o, per la variante B, endpoint VPC al posto del NAT.
    - **Variante A** ([lezione 5.2](02-ha-scaling.md)): launch template con Amazon Linux, user data che installa Docker, legge il segreto
      e avvia il container; **Auto Scaling group** min 2, max 4, su entrambe le subnet private, con *target tracking*
      sulla CPU al 50%; **ALB** con health check su `/health`.
    - **Variante B** ([lezione 5.9](09-container-serverless.md)): servizio **ECS su Fargate** con 2 task nelle subnet private, dietro l'ALB.
    - **Dati** ([lezione 5.4](04-database.md)): **RDS for PostgreSQL** `db.t4g.micro`, **Multi-AZ** (attivalo per la prova di failover
      e spegnilo subito dopo se vuoi risparmiare), cifratura attiva, backup automatici 7 giorni, non accessibile
      pubblicamente.
    - **Sicurezza** ([lezione 5.6](06-sicurezza.md)): `DATABASE_URL` in **Secrets Manager**; ruolo IAM delle istanze (o ruolo di
      esecuzione dei task) che può leggere **solo** quel segreto; security group concatenati: ALB ← Internet,
      app ← solo dal SG dell'ALB, database ← solo dal SG dell'app; accesso alle istanze solo con Session Manager.
    - **Monitoraggio** ([lezione 5.8](08-monitoraggio.md)): allarme su `UnHealthyHostCount` e su `HTTPCode_Target_5XX_Count`, notifica SNS
      via email; log del container su CloudWatch Logs (driver `awslogs` o log di ECS).

=== "Google Cloud"

    ```mermaid
    flowchart LR
      U[Utenti] -->|HTTP/S| LB[Application Load Balancer<br>esterno globale]
      subgraph VPC
        LB --> MIG[MIG regionale<br>2 zone, senza IP pubblici]
        MIG --> SQL[(Cloud SQL PostgreSQL<br>IP privato, alta disponibilità)]
      end
      MIG -.-> SEC[Secret Manager]
      MIG -.-> MON[Cloud Monitoring<br>e Logging]
    ```

    - **Rete** ([lezione 5.1](01-vpc.md)): VPC in modalità personalizzata con una subnet nella tua regione, **Cloud NAT** per le
      uscite, regole firewall che ammettono il traffico verso le VM solo dagli intervalli dei controlli di integrità e
      del load balancer (`35.191.0.0/16`, `130.211.0.0/22`) e da IAP per SSH (`35.235.240.0/20`).
    - **Variante A** ([lezione 5.2](02-ha-scaling.md)): template di istanza senza IP esterno, script di avvio che installa Docker, legge il
      segreto e avvia il container; **MIG regionale** su almeno 2 zone, autoscaling min 2 max 4 sulla CPU al 50%,
      riparazione automatica con controllo di integrità su `/health`; **load balancer** HTTP esterno globale.
    - **Variante B** ([lezione 5.9](09-container-serverless.md)): servizio **Cloud Run** con accesso alla VPC (*Direct VPC egress*), segreto montato
      come variabile d'ambiente, istanze minime 0 o 1.
    - **Dati** ([lezione 5.4](04-database.md)): **Cloud SQL per PostgreSQL** con **IP privato** (accesso privato ai servizi),
      **alta disponibilità regionale** (per la prova di failover), backup automatici e recupero point-in-time.
    - **Sicurezza** ([lezione 5.6](06-sicurezza.md)): `DATABASE_URL` in **Secret Manager**; un **service account dedicato** per le VM o
      per Cloud Run con `roles/secretmanager.secretAccessor` **solo** su quel segreto, più i ruoli per scrivere log e
      metriche.
    - **Monitoraggio** ([lezione 5.8](08-monitoraggio.md)): controllo di uptime su `/health`, criterio di avviso con canale email, Ops Agent
      sulle VM (variante A).

### Requisiti comuni

- Ogni risorsa ha le etichette `progetto=tasks-api`, `ambiente=dev`, `responsabile=TUO-NOME`.
- Un **budget** con avvisi al 50%, 80% e 100%.
- Nessuna credenziale nel repository: nemmeno la password del database (sta solo nel gestore dei segreti).
- L'immagine del container è quella pubblicata su GHCR nella [Fase 3](../fase-3/index.md), con un **tag di versione**, mai `latest`.

### La documentazione (è metà del voto)

Nel repository `tasks-api-cloud`:

- `README.md`: obiettivo, **diagramma** (Mermaid), servizi usati e perché, istruzioni passo passo per ricreare tutto,
  **costo stimato** al mese con il link al calcolatore salvato, istruzioni di distruzione;
- `docs/decisioni.md`: le scelte e i compromessi (un solo NAT gateway, Multi-AZ sì o no, variante A o B), nel formato
  *contesto → decisione → conseguenze*;
- `docs/revisione.md`: la revisione con i pilastri **Reliability**, **Security** e **Cost Optimization** ([lezione
  5.10](10-well-architected.md)), i rischi accettati e quelli da risolvere;
- `docs/dr.md`: **RTO e RPO** dichiarati, la strategia di DR scelta e il risultato della prova di failover;
- `scripts/startup.sh` (variante A) e `scripts/destroy.sh`;
- screenshot essenziali (load balancer con destinazioni sane, allarme, grafico del failover) in `docs/img/`.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `curl -s http://INDIRIZZO-LB/health` | `{"status":"ok"}` |
| `POST /tasks` più volte, poi `GET /tasks` ripetuto | tutte le attività, qualunque istanza risponda: il database è condiviso |
| Un'istanza o un task fermato a mano | il load balancer lo esclude, il gruppo lo sostituisce, il servizio non si interrompe |
| Failover del database (*Reboot with failover* su RDS, *Failover* su Cloud SQL) | dopo un'interruzione di circa un minuto (misurala) l'API torna a rispondere; i dati ci sono tutti |
| Carico con `hey` o `ab` dal Dev Container per qualche minuto | l'autoscaling aggiunge istanze; poi le rimuove |
| Connessione diretta al database dal tuo PC | impossibile: non ha indirizzo pubblico |
| Il ruolo o service account dell'applicazione prova a leggere un altro segreto | accesso negato |
| Allarme | ricevi l'email quando fermi tutte le istanze |
| `destroy.sh` | nessuna risorsa rimasta; controlla anche snapshot finali, IP, volumi, log group |

## Criteri di valutazione

- [ ] Rete con subnet pubbliche e private su almeno due zone; applicazione e database senza IP pubblici.
- [ ] L'applicazione sopravvive alla perdita di un'istanza e al failover del database.
- [ ] Il segreto è letto a runtime, con permessi limitati a quel solo segreto.
- [ ] Security group o regole firewall ammettono solo il traffico necessario.
- [ ] Allarmi e log funzionano.
- [ ] Etichette e budget presenti; costo stimato documentato.
- [ ] README, decisioni, revisione e DR scritti in modo che un collega possa ricreare tutto senza chiederti nulla.
- [ ] Tutto distrutto alla fine.

??? success "Una possibile soluzione: startup.sh per la variante A su AWS"
    ```bash
    #!/bin/bash
    # User data del launch template: installa Docker, legge DATABASE_URL da Secrets Manager e avvia la Tasks API.
    # Il ruolo dell'istanza deve poter eseguire secretsmanager:GetSecretValue sul solo segreto tasks/database-url.
    set -euo pipefail
    exec > >(tee /var/log/startup.log) 2>&1

    REGION=$(TOKEN=$(curl -sX PUT http://169.254.169.254/latest/api/token \
               -H 'X-aws-ec2-metadata-token-ttl-seconds: 60') &&
             curl -s -H "X-aws-ec2-metadata-token: $TOKEN" \
               http://169.254.169.254/latest/meta-data/placement/region)
    IMAGE=ghcr.io/TUO-UTENTE/tasks-api:1.1.0

    dnf install -y docker jq
    systemctl enable --now docker

    DATABASE_URL=$(aws secretsmanager get-secret-value --region "$REGION" \
      --secret-id tasks/database-url --query SecretString --output text | jq -r .DATABASE_URL)

    docker run -d --name api --restart unless-stopped -p 8000:8000 \
      --log-driver awslogs --log-opt awslogs-region="$REGION" \
      --log-opt awslogs-group=/tasks-api --log-opt awslogs-create-group=true \
      -e DATABASE_URL="$DATABASE_URL" "$IMAGE"
    ```

    **Spiegazione.**

    - La regione viene letta dal **servizio dei metadati** con **IMDSv2** (prima un token con `PUT`, poi la lettura):
      lo stesso script funziona in qualunque regione.
    - Il segreto non compare né nel launch template né nel codice: lo legge l'istanza all'avvio grazie al suo ruolo.
      Nei log di avvio non viene stampato.
    - Il driver `awslogs` manda i log del container a CloudWatch Logs; per creare il gruppo il ruolo deve avere
      anche `logs:CreateLogGroup`, `logs:CreateLogStream` e `logs:PutLogEvents` (oppure crea il gruppo in anticipo).
    - Se l'istanza non diventa sana, collegati con Session Manager e leggi `/var/log/startup.log` e `docker logs api`.

    Su Google Cloud lo script è analogo: la lettura del segreto diventa
    `gcloud secrets versions access latest --secret=tasks-database-url` (con il service account della VM) e i log
    del container si raccolgono con l'Ops Agent o con `--log-driver=gcplogs`.

??? success "Una possibile soluzione: costi della variante A su AWS"
    Per una stima realistica, usa il [calcolatore](https://calculator.aws/) con la tua regione: i prezzi cambiano nel
    tempo e da regione a regione, quindi non riportiamo cifre fisse. Le voci da inserire sono:

    | Voce | Note |
    |---|---|
    | 2 istanze `t4g.small` sempre attive | il grosso del costo di calcolo; valuta `t4g.micro` |
    | Application Load Balancer | costo orario fisso più le LCU |
    | NAT gateway | costo orario più costo per GB elaborato; spesso la voce più sorprendente |
    | RDS `db.t4g.micro` Multi-AZ | Multi-AZ raddoppia circa il costo dell'istanza |
    | Storage RDS e backup | i backup fino alla dimensione del database sono inclusi |
    | Secrets Manager, CloudWatch | pochi centesimi in questo progetto |

    Nel README spiega quali voci elimineresti in un ambiente di sviluppo (Multi-AZ, NAT gateway sostituito da
    endpoint VPC) e quali terresti in produzione.

## Fonti

- [AWS: Amazon EC2 Auto Scaling with Elastic Load Balancing](https://docs.aws.amazon.com/autoscaling/ec2/userguide/autoscaling-load-balancer.html)
- [Amazon EC2: Use IMDSv2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html)
- [Amazon RDS: Multi-AZ DB instance deployments](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.MultiAZSingleStandby.html)
- [Docker: Amazon CloudWatch Logs logging driver](https://docs.docker.com/engine/logging/drivers/awslogs/)
- [Google Cloud: Regional managed instance groups](https://docs.cloud.google.com/compute/docs/instance-groups/regional-migs)
- [Cloud SQL: High availability](https://docs.cloud.google.com/sql/docs/postgres/high-availability)
- [Google Cloud: Firewall rules for health checks](https://docs.cloud.google.com/load-balancing/docs/health-check-concepts#ip-ranges)
- [Cloud Run: Direct VPC egress](https://docs.cloud.google.com/run/docs/configuring/vpc-direct-vpc)
