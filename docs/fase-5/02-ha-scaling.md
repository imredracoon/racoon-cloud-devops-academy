# 5.2 Bilanciamento e scalabilità

!!! obiettivi "Obiettivi"
    - Distinguere scalabilità verticale e orizzontale, alta disponibilità e tolleranza ai guasti.
    - Scegliere il tipo di load balancer giusto e configurare i controlli di integrità.
    - Creare gruppi di istanze che si riparano e si ridimensionano da soli.
    - Scegliere la politica di scalabilità adatta al carico.
    - Progettare applicazioni senza stato, condizione per scalare orizzontalmente.

## I concetti

- **Scalabilità verticale** (*scale up*): una macchina più grande. Semplice, ma ha un limite e spesso richiede un
  riavvio.
- **Scalabilità orizzontale** (*scale out*): più macchine uguali. Nessun limite pratico, ma l'applicazione deve
  essere progettata per girare in più copie.
- **Elasticità**: aggiungere e togliere risorse **in automatico** seguendo il carico.
- **Alta disponibilità**: il servizio resta disponibile nonostante i guasti, con interruzioni minime.
- **Tolleranza ai guasti**: il servizio continua **senza alcuna interruzione** visibile; costa di più.

Due metriche guidano la progettazione del ripristino:

- **RTO** (*Recovery Time Objective*): quanto tempo può restare fermo il servizio;
- **RPO** (*Recovery Point Objective*): quanti dati (misurati in tempo) si possono perdere.

### Applicazioni senza stato

Per scalare orizzontalmente, ogni copia dell'applicazione deve poter rispondere a qualsiasi richiesta: niente dati
conservati **solo** sul disco o nella memoria della singola istanza. Sessioni, file caricati e dati vanno in servizi
condivisi: database, cache (ElastiCache, Memorystore), storage a oggetti. La Tasks API è già **senza stato**: tutti i
dati sono nel database indicato da `DATABASE_URL`.

## I load balancer

Un **load balancer** distribuisce le richieste tra più destinazioni sane, verificandole con **controlli di
integrità** (*health check*): una destinazione che non risponde correttamente smette di ricevere traffico.

=== "AWS"

    **Elastic Load Balancing** offre tre tipi:

    | Tipo | Livello | Quando |
    |---|---|---|
    | **Application Load Balancer** (ALB) | 7, HTTP/HTTPS | applicazioni web e API: instradamento per percorso (`/api/*`) o nome host, terminazione TLS, integrazione con WAF e Cognito |
    | **Network Load Balancer** (NLB) | 4, TCP/UDP/TLS | prestazioni estreme, latenza bassissima, IP statici, protocolli non HTTP |
    | **Gateway Load Balancer** | 3 | inserire appliance di sicurezza di terze parti (firewall) nel percorso del traffico |

    Gli elementi di un ALB: i **listener** (porta e protocollo, per esempio HTTPS:443 con un certificato di **ACM**),
    le **regole**, i **target group** (le destinazioni: istanze, IP, funzioni Lambda) con il loro health check.

=== "Google Cloud"

    **Cloud Load Balancing** è software distribuito, senza istanze da dimensionare:

    | Tipo | Ambito | Quando |
    |---|---|---|
    | **Application Load Balancer** esterno | globale o regionale | HTTP/HTTPS; il globale ha **un unico IP anycast** per tutto il mondo e manda l'utente al backend più vicino |
    | **Application Load Balancer** interno | regionale o tra regioni | HTTP tra servizi interni |
    | **Network Load Balancer** proxy | globale o regionale | TCP con terminazione TLS |
    | **Network Load Balancer** passthrough | regionale | TCP/UDP, preserva l'IP del client |

    Gli elementi: **regola di forwarding** (IP e porta), **proxy di destinazione** (con il certificato), **mappa URL**,
    **servizio di backend** con health check, e i backend (gruppi di istanze, NEG, bucket).

## Gruppi di istanze che si gestiscono da soli

=== "AWS"

    Un **Auto Scaling group** (ASG) mantiene un numero di istanze tra un **minimo** e un **massimo**, con una
    **capacità desiderata**:

    - le istanze nascono da un **launch template** (AMI, tipo, security group, ruolo IAM, user data);
    - il gruppo si estende su **più AZ** e le bilancia;
    - se un'istanza risulta non sana (anche secondo l'health check del load balancer, se attivato), viene **sostituita**;
    - le **politiche di scalabilità** cambiano la capacità desiderata.

    | Politica | Funzionamento |
    |---|---|
    | **Target tracking** | mantiene una metrica a un valore obiettivo, per esempio CPU media al 50% (la più usata) |
    | **Step scaling** | aggiunge o toglie istanze a gradini, in base all'entità dell'allarme |
    | **Scheduled** | cambia la capacità a orari fissi (picco del lunedì mattina) |
    | **Predictive** | prevede il carico dalla storia e scala in anticipo |

=== "Google Cloud"

    Un **gruppo di istanze gestito** (*managed instance group*, MIG) crea VM identiche da un **modello di istanza**:

    - **zonale** o **regionale** (distribuito su più zone: da preferire);
    - **riparazione automatica** (*autohealing*): ricrea le VM che falliscono l'health check;
    - **autoscaler** basato su CPU, capacità di servizio del load balancer, metriche di Cloud Monitoring o
      pianificazioni; anche **predittivo**;
    - **aggiornamenti progressivi** (*rolling update*) verso un nuovo modello.

### Scalabilità senza sorprese

- Imposta un **massimo** ragionevole: protegge da costi fuori controllo (o da un attacco).
- Il **periodo di riscaldamento** (*warm-up* / *cooldown*) evita che il gruppo aggiunga istanze prima che le nuove
  siano pronte.
- Il **database** spesso è il collo di bottiglia: più istanze dell'applicazione significano più connessioni.

## Laboratorio

!!! laboratorio "Laboratorio · La Tasks API dietro un load balancer"
    Usa l'immagine pubblicata nella [Fase 3](../fase-3/index.md) (`ghcr.io/TUO-UTENTE/tasks-api:1.1.0`, pubblica) e la rete del laboratorio
    [5.1](01-vpc.md). Script di avvio delle istanze:

    ```bash title="avvio-tasks.sh"
    #!/bin/bash
    apt-get update && apt-get install -y docker.io
    systemctl enable --now docker
    docker run -d --restart unless-stopped -p 8000:8000 --name api ghcr.io/TUO-UTENTE/tasks-api:1.1.0
    ```

    === "AWS"

        1. **Launch template** `tasks-lt`: Ubuntu 24.04, `t3.micro`, security group `app-sg`, ruolo con
           `AmazonSSMManagedInstanceCore`, user data `avvio-tasks.sh`.
        2. **Target group** `tasks-tg`: tipo *Instances*, HTTP 8000, VPC `tasks`, health check `GET /health`,
           codice atteso `200`.
        3. **Application Load Balancer** `tasks-alb`: *Internet-facing*, subnet pubbliche delle due AZ, security group
           `alb-sg`, listener HTTP:80 → `tasks-tg`.
        4. **Auto Scaling group** `tasks-asg`: template `tasks-lt`, subnet `app-a` e `app-b`, collegato a `tasks-tg`,
           **health check ELB attivo**, capacità desiderata 2, minima 2, massima 4, politica **target tracking** CPU 50%.
        5. Quando i target risultano *healthy*, apri `http://DNS-DEL-ALB/health` e ricarica più volte: il campo `host`
           cambia.
        6. Termina a mano una delle istanze: dopo qualche minuto l'ASG ne crea una nuova.

        **Pulizia.** Elimina l'ASG (termina le istanze), il load balancer, il target group, il launch template.

    === "Google Cloud"

        1. **Modello di istanza** `tasks-tpl`: `e2-micro`, Ubuntu 24.04, rete `tasks`, subnet `app`, nessun IP esterno,
           service account `tasks-app`, script di avvio `avvio-tasks.sh`.
        2. **Controllo di integrità** `tasks-hc`: HTTP porta 8000, percorso `/health`.
        3. **Gruppo di istanze gestito** `tasks-mig`: **regionale**, 2 istanze, autoscaling su CPU al 50% (min 2, max 4),
           riparazione automatica con `tasks-hc`; **porta denominata** `http:8000`.
        4. **Bilanciamento del carico** → Application Load Balancer esterno globale: backend `tasks-mig` sulla porta
           denominata `http`, health check `tasks-hc`; frontend HTTP:80.
        5. Dopo qualche minuto apri `http://IP-DEL-LB/health` e ricarica: il campo `host` cambia.
        6. Elimina una VM del gruppo: il MIG la ricrea.

        **Pulizia.** Elimina il load balancer (regola di forwarding, proxy, mappa URL, servizio di backend), il MIG, il
        modello, l'health check.

??? success "Spiegazione e una sorpresa"
    Il campo `host` di `/health` mostra quale istanza ha risposto: il load balancer distribuisce le richieste. Ora prova
    a creare un'attività con `POST /tasks` e poi a leggere `GET /tasks` più volte: **a volte c'è, a volte no**. Ogni
    istanza usa il proprio database SQLite locale: l'applicazione, con questa configurazione, **ha uno stato**. È
    esattamente il problema che risolverai nella [lezione 5.4](04-database.md) con un database gestito condiviso.

## Esercizi

!!! esercizio "Esercizio · Scegli il bilanciamento"
    1. Un'API REST con percorsi `/v1/*` e `/v2/*` serviti da due gruppi di istanze diversi.
    2. Un server di gioco online su UDP che richiede latenza minima e IP statici da comunicare ai partner.
    3. Un sito con utenti in Europa, America e Asia, che vuoi servire con un solo indirizzo IP e il backend più vicino.
    4. Un negozio online che ha sempre un picco il venerdì alle 20.

??? success "Soluzione"
    1. **ALB** con regole per percorso / **Application Load Balancer** con mappa URL.
    2. **NLB** (supporta UDP e IP statici per zona) / **Network Load Balancer passthrough**.
    3. **Application Load Balancer esterno globale** di Google Cloud (IP anycast). Su AWS: **AWS Global Accelerator**
       (IP anycast statici verso ALB regionali) o **CloudFront** davanti agli ALB.
    4. **Scalabilità pianificata** (scheduled) per anticipare il picco, più target tracking per le variazioni impreviste;
       eventualmente la scalabilità **predittiva**.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale load balancer AWS instrada le richieste HTTP in base al percorso dell'URL?","options":["Network Load Balancer","Application Load Balancer","Gateway Load Balancer","Classic DNS"],"answer":[1],"explanation":"L'ALB lavora al livello 7."},
 {"q":"Cosa fa un Auto Scaling group quando un'istanza non supera l'health check?","options":["La riavvia all'infinito","La sostituisce con una nuova istanza","Avvisa e basta","Ferma il gruppo"],"answer":[1],"explanation":"Su Google Cloud la riparazione automatica del MIG fa lo stesso."},
 {"q":"Quale politica mantiene la CPU media al 50%?","options":["Scheduled","Target tracking","Step scaling","Manual"],"answer":[1],"explanation":"È la più semplice e la più usata."},
 {"q":"Perché un'applicazione deve essere senza stato per scalare orizzontalmente?","options":["Per usare meno CPU","Perché qualsiasi istanza deve poter servire qualsiasi richiesta","Per evitare i log","Per usare IPv6"],"answer":[1],"explanation":"I dati vanno in servizi condivisi."},
 {"q":"Cosa misura l'RPO?","options":["Il tempo massimo di fermo","La quantità massima di dati, espressa in tempo, che si può perdere","Il numero di repliche","La latenza"],"answer":[1],"explanation":"L'RTO misura il tempo di ripristino."},
 {"q":"Quali sono caratteristiche dell'Application Load Balancer esterno globale di Google Cloud?","options":["Un unico IP anycast","Invio dell'utente al backend sano più vicino","Funziona solo con UDP","Richiede istanze dedicate da dimensionare","È solo zonale"],"answer":[0,1],"explanation":"È un servizio software distribuito sulla rete globale di Google."}
]}
</script>

## Fonti

- [Elastic Load Balancing: Product comparisons](https://aws.amazon.com/elasticloadbalancing/features/)
- [Application Load Balancer User Guide](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html)
- [Amazon EC2 Auto Scaling User Guide](https://docs.aws.amazon.com/autoscaling/ec2/userguide/what-is-amazon-ec2-auto-scaling.html)
- [Amazon EC2 Auto Scaling: Target tracking scaling policies](https://docs.aws.amazon.com/autoscaling/ec2/userguide/as-scaling-target-tracking.html)
- [Google Cloud: Load balancer overview e scelta](https://docs.cloud.google.com/load-balancing/docs/choosing-load-balancer)
- [Compute Engine: Managed instance groups](https://docs.cloud.google.com/compute/docs/instance-groups) e [Autoscaling](https://docs.cloud.google.com/compute/docs/autoscaler)
- [AWS Well-Architected: Reliability pillar](https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/welcome.html)
