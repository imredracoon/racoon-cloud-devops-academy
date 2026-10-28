# 5.7 DNS e distribuzione dei contenuti

!!! obiettivi "Obiettivi"
    - Gestire zone DNS pubbliche e private nel cloud.
    - Usare le politiche di instradamento DNS: peso, latenza, geolocalizzazione, failover.
    - Mettere una CDN davanti a un'applicazione, con regole di cache e origini multiple.
    - Accelerare il traffico globale senza cache.

## DNS gestito

Ripassa la [lezione 1.8](../fase-1/08-dns-routing.md): il DNS traduce nomi in indirizzi, con record di tipo A, AAAA, CNAME, MX, TXT, e un **TTL** che
dice per quanto tempo i resolver possono conservare la risposta.

=== "AWS"

    **Amazon Route 53**:

    - **zone ospitate pubbliche** (risolte da Internet) e **private** (risolte solo dalle VPC associate);
    - registrazione dei **domini**;
    - **record alias**: puntano a risorse AWS (ALB, CloudFront, bucket S3) anche **all'apice della zona**
      (`esempio.it`, dove un CNAME non è ammesso) e senza costo per le query;
    - **health check** sugli endpoint, usati dalle politiche di failover;
    - **Route 53 Resolver** per la risoluzione ibrida tra VPC e reti in sede (endpoint in ingresso e in uscita).

=== "Google Cloud"

    **Cloud DNS**:

    - zone **pubbliche** e **private** (visibili alle reti VPC autorizzate), zone di **inoltro** e di **peering** per
      scenari ibridi;
    - **DNSSEC** gestito;
    - **policy di instradamento**: round robin ponderato, geolocalizzazione, failover con controlli di integrità.
    - Per registrare un dominio si usa un registrar esterno, poi si delega la zona a Cloud DNS indicando i suoi
      server dei nomi.

## Le politiche di instradamento

| Politica (Route 53) | Funzionamento | Uso tipico |
|---|---|---|
| **Simple** | un record, uno o più valori | un solo endpoint |
| **Weighted** | distribuisce per peso (per esempio 90/10) | rilasci graduali, test A/B |
| **Latency** | risponde con la regione a latenza minore per l'utente | applicazioni multi-regione |
| **Geolocation** | risponde in base al paese o continente dell'utente | contenuti localizzati, vincoli di legge |
| **Geoproximity** | in base alla distanza, con un fattore di spostamento | regolare quanto traffico attira una regione |
| **Failover** | primario e secondario, con health check | disaster recovery attivo-passivo |
| **Multivalue answer** | fino a otto valori sani scelti a caso | semplice distribuzione con health check |
| **IP-based** | in base all'intervallo IP del client | instradamento per rete di provenienza |

Cloud DNS offre le politiche ponderate, di geolocalizzazione e di failover. Ricorda che il DNS **non** è un load
balancer istantaneo: i client conservano le risposte per il TTL. Per failover rapidi si usano TTL brevi (60 secondi).

## La CDN

Una **CDN** conserva copie dei contenuti nei punti di presenza ([lezione 4.2](../fase-4/02-infrastruttura-globale.md)): l'utente riceve la risposta dal PoP più
vicino e l'**origine** (il bucket, il load balancer) riceve molte meno richieste.

=== "AWS"

    **Amazon CloudFront**:

    - **origini**: bucket S3 (con **Origin Access Control**), ALB, qualsiasi server HTTP; più origini nella stessa
      distribuzione, con **comportamenti** (*cache behavior*) per percorso: `/static/*` dal bucket, `/api/*` dall'ALB;
    - **politiche di cache** (cosa entra nella chiave della cache, TTL) e **di richiesta all'origine**;
    - **gruppi di origini** con failover automatico;
    - **URL e cookie firmati** per contenuti privati, **geo-restrizione**;
    - **CloudFront Functions** e **Lambda@Edge** per eseguire codice ai margini (riscritture, intestazioni);
    - integrazione con **AWS WAF** e **Shield**; certificati da **ACM** in `us-east-1`.

=== "Google Cloud"

    **Cloud CDN** si attiva su un **Application Load Balancer esterno globale** per i backend (servizi o bucket):

    - **modalità di cache**: solo i contenuti statici, tutti, o seguendo le intestazioni dell'origine;
    - **chiavi di cache** personalizzabili, TTL, **invalidazione**;
    - **URL e cookie firmati**;
    - integrazione con **Cloud Armor**.
    - **Media CDN** è dedicata ai grandi volumi di streaming video.

### Cosa mettere in cache

- **Statico** (immagini, CSS, JavaScript, file del sito): sì, con TTL lunghi e nomi dei file che cambiano a ogni
  versione (`app.3f9a2c.js`), così non serve invalidare.
- **API e contenuti personalizzati**: in genere no, oppure con TTL brevissimi e chiavi di cache che includono ciò che
  rende diversa la risposta. La CDN accelera comunque la connessione (TLS terminato vicino all'utente, rete del
  fornitore fino all'origine).

## Accelerare senza cache

**AWS Global Accelerator** fornisce **due IP anycast statici** che fanno entrare il traffico TCP/UDP nella rete AWS
nel punto più vicino all'utente e lo portano all'endpoint regionale più vicino e sano (ALB, NLB, istanze). È utile per
protocolli non HTTP, per IP fissi da comunicare a clienti e per failover tra regioni in pochi secondi. Su Google Cloud
lo stesso risultato si ottiene con i load balancer globali e il livello di rete Premium.

## Laboratorio

!!! laboratorio "Laboratorio · Una CDN davanti alla Tasks API"
    Riprendi l'architettura dei laboratori precedenti (load balancer con le istanze della Tasks API).

    === "AWS"

        1. **CloudFront** → **Create distribution**: origine il DNS dell'ALB, protocollo HTTP verso l'origine (per il
           laboratorio), **Viewer protocol policy: Redirect HTTP to HTTPS**.
        2. **Cache policy: CachingDisabled** per il comportamento predefinito (è un'API), **Origin request policy:
           AllViewerExceptHostHeader**, metodi consentiti **GET, HEAD, OPTIONS, PUT, POST, PATCH, DELETE**.
        3. Attendi il deploy e prova `https://DOMINIO.cloudfront.net/health` e un `POST /tasks`.
        4. Osserva le intestazioni con `curl -sI`: `x-cache` e `x-amz-cf-pop` (il PoP che ti ha servito).

        **Pulizia.** Disattiva la distribuzione, attendi, poi eliminala.

    === "Google Cloud"

        1. Nel servizio di backend del load balancer della [lezione 5.2](02-ha-scaling.md), attiva **Cloud CDN** con modalità di cache
           **Usa le impostazioni dell'origine** (*USE_ORIGIN_HEADERS*): la Tasks API non invia intestazioni di cache,
           quindi le risposte non vengono conservate.
        2. Prova `curl -sI http://IP-DEL-LB/health` più volte.
        3. Per vedere la cache in azione, aggiungi un **backend bucket** con un file statico e una regola della mappa URL
           per `/static/*`; scarica due volte il file e osserva l'intestazione `age`.

        **Pulizia.** Disattiva Cloud CDN ed elimina il backend bucket e il bucket, insieme alle altre risorse.

??? success "Spiegazione"
    Anche senza cache, la CDN termina HTTPS vicino all'utente e porta le richieste all'origine sulla rete del fornitore.
    Con i contenuti statici la differenza è netta: alla seconda richiesta la risposta arriva dal PoP (`Hit from
    cloudfront` in `x-cache`, o un valore di `age` maggiore di zero) senza toccare l'origine.

## Esercizi

!!! esercizio "Esercizio · Progetta l'instradamento"
    1. Vuoi mandare il 5% degli utenti alla nuova versione dell'applicazione, su un ALB diverso.
    2. La tua applicazione gira in `eu-south-1` e `us-east-1`: ogni utente deve raggiungere la più veloce.
    3. Per legge, gli utenti tedeschi devono essere serviti dalla regione di Francoforte.
    4. Un sito di riserva statico su S3 deve entrare in funzione se l'applicazione principale non risponde.
    5. Un partner deve inserire nel suo firewall due IP fissi per raggiungere la tua API su UDP.

??? success "Soluzione"
    1. Route 53 **weighted** (95/5).
    2. Route 53 **latency-based** (o un load balancer globale su Google Cloud).
    3. Route 53 **geolocation** (paese Germania → Francoforte, più un record predefinito per gli altri).
    4. Route 53 **failover**: primario l'ALB con health check, secondario il sito statico (o un gruppo di origini con
       failover in CloudFront).
    5. **AWS Global Accelerator** (due IP anycast statici, supporta UDP).

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Perché in Route 53 si usa un record alias per puntare esempio.it a un load balancer?","options":["È più veloce da digitare","All'apice della zona un CNAME non è permesso, e l'alias verso risorse AWS è ammesso e senza costo per le query","Gli alias sono cifrati","Per usare IPv6"],"answer":[1],"explanation":"Il record alias è un'estensione di Route 53."},
 {"q":"Quale politica di Route 53 usi per un rilascio graduale del 10% del traffico?","options":["Latency","Weighted","Geolocation","Simple"],"answer":[1],"explanation":"I pesi stabiliscono la proporzione."},
 {"q":"Quale politica manda ogni utente alla regione con la risposta più rapida?","options":["Geolocation","Latency","Failover","Multivalue"],"answer":[1],"explanation":"Geolocation decide in base al paese, non alla velocità."},
 {"q":"Come limiti a CloudFront l'accesso a un bucket S3 privato?","options":["Rendendo pubblico il bucket","Con Origin Access Control e la bucket policy","Con un NAT gateway","Con Route 53"],"answer":[1],"explanation":"Lezione 4 e progetto della Fase 4."},
 {"q":"Quale servizio fornisce IP anycast statici e accelera traffico TCP e UDP senza cache?","options":["CloudFront","AWS Global Accelerator","Route 53","Direct Connect"],"answer":[1],"explanation":"CloudFront è una CDN per HTTP."},
 {"q":"Su Google Cloud, dove si attiva Cloud CDN?","options":["Su un servizio di backend di un Application Load Balancer esterno","Su un backend bucket","Su una VM senza load balancer","Su Cloud SQL","Su Pub/Sub"],"answer":[0,1],"explanation":"Cloud CDN lavora insieme al load balancer globale."}
]}
</script>

## Fonti

- [Amazon Route 53: Choosing a routing policy](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy.html)
- [Amazon Route 53: Choosing between alias and non-alias records](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-choosing-alias-non-alias.html)
- [Amazon CloudFront Developer Guide](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/Introduction.html)
- [Amazon CloudFront: Managed cache policies](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/using-managed-cache-policies.html)
- [AWS Global Accelerator](https://docs.aws.amazon.com/global-accelerator/latest/dg/what-is-global-accelerator.html)
- [Cloud DNS overview](https://docs.cloud.google.com/dns/docs/overview) e [routing policies](https://docs.cloud.google.com/dns/docs/zones/manage-routing-policies)
- [Cloud CDN overview](https://docs.cloud.google.com/cdn/docs/overview) e [cache modes](https://docs.cloud.google.com/cdn/docs/caching#cache-modes)
