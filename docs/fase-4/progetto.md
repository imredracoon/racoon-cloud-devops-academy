# Progetto 4 · Sito statico nel cloud

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai usare un cloud in modo sicuro e consapevole dei costi: un'identità con privilegi minimi, uno storage a
    oggetti protetto, una CDN con HTTPS, il deploy automatizzato da riga di comando e una pulizia completa. Il risultato è
    anche utile: una pagina personale pubblica da mettere nel CV.

| | |
|---|---|
| **Durata** | 8-12 ore |
| **Consegna** | Repository pubblico `cloud-static-site` e l'indirizzo HTTPS del sito |
| **Costo** | entro i crediti per i nuovi clienti; senza crediti, pochi centesimi al mese |

## Il compito

Pubblica una pagina personale statica (HTML e CSS: chi sei, cosa stai imparando, link a GitHub e LinkedIn) servita in
**HTTPS** attraverso una **CDN**, con lo storage **non** accessibile direttamente.

=== "AWS"

    Architettura:

    ```mermaid
    flowchart LR
      U[Visitatore] -->|HTTPS| CF[CloudFront<br>distribuzione]
      CF -->|Origin Access Control| S3[(Bucket S3<br>privato)]
      D[deploy.sh<br>aws s3 sync] --> S3
      D -->|invalidazione| CF
    ```

    1. Un bucket S3 **privato** (Block Public Access attivo), con versioning.
    2. Una distribuzione **CloudFront** con il bucket come origine tramite **Origin Access Control** (OAC): solo
       CloudFront può leggere il bucket, grazie alla bucket policy che CloudFront propone.
    3. Default root object `index.html`, reindirizzamento da HTTP a HTTPS.
    4. Una pagina `404.html` personalizzata (risposte di errore personalizzate di CloudFront).

=== "Google Cloud"

    Architettura:

    ```mermaid
    flowchart LR
      U[Visitatore] -->|HTTPS| LB[Load balancer HTTP/S esterno<br>con Cloud CDN]
      LB --> B[Backend bucket]
      B --> GCS[(Bucket Cloud Storage)]
      D[deploy.sh<br>gcloud storage rsync] --> GCS
    ```

    1. Un bucket Cloud Storage con accesso uniforme, `index.html` come pagina principale e `404.html` come pagina di
       errore (`gcloud storage buckets update --web-main-page-suffix --web-error-page`).
    2. Un **load balancer** applicativo esterno globale con un **backend bucket** e **Cloud CDN** attivo.
    3. HTTPS: con un tuo dominio, un certificato gestito da Google; senza dominio, è accettabile servire in HTTP
       dall'IP del load balancer e documentare il passo mancante.
    4. Per far leggere il bucket al load balancer, il bucket deve concedere `roles/storage.objectViewer` a `allUsers`
       (gli oggetti diventano pubblici). È la configurazione documentata da Google per i siti statici: metti nel bucket
       **solo** i file del sito.

    Il load balancer ha un costo orario anche senza traffico: controlla il calcolatore e **distruggilo** a fine
    progetto se non vuoi tenerlo. Un'alternativa senza load balancer, più semplice ma senza CDN né dominio, è servire
    gli oggetti pubblici direttamente da `https://storage.googleapis.com/NOME-BUCKET/index.html`.

### Automazione

Nel repository:

- `site/` con i file del sito;
- `scripts/deploy.sh`: sincronizza `site/` con il bucket (cancellando i file rimossi) e, su AWS, crea
  un'**invalidazione** della cache di CloudFront; legge nome del bucket e ID della distribuzione da variabili d'ambiente;
- `scripts/destroy.sh`: elimina **tutte** le risorse create, nell'ordine corretto;
- `README.md` con: diagramma, istruzioni per ricreare tutto (passi in console o comandi), costi stimati con il link al
  calcolatore, misure di sicurezza adottate.

### Sicurezza

- Tutto creato con il tuo utente amministratore di Identity Center (AWS) o il tuo account con MFA (Google Cloud), mai
  con l'utente root.
- Etichetta `progetto=academy` su ogni risorsa che la supporta.
- Un budget attivo con avvisi.
- Nessuna credenziale nel repository.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `curl -sI https://INDIRIZZO-CDN/` | `200`, intestazioni della CDN (`x-cache` su CloudFront, `age` o `via` su Google) |
| `curl -sI http://INDIRIZZO-CDN/` (AWS) | `301` verso HTTPS |
| `curl -sI https://INDIRIZZO-CDN/non-esiste` | la tua pagina 404 |
| AWS: `curl -sI https://NOME-BUCKET.s3.REGIONE.amazonaws.com/index.html` | `403`: il bucket non è accessibile direttamente |
| Modifica `index.html`, `deploy.sh`, ricarica | la modifica è visibile (subito su AWS grazie all'invalidazione) |
| `destroy.sh` e poi la console | nessuna risorsa del progetto rimasta |

## Criteri di valutazione

- [ ] Il sito risponde in HTTPS attraverso la CDN.
- [ ] Lo storage non è raggiungibile direttamente (AWS) o contiene solo i file del sito (Google Cloud).
- [ ] `deploy.sh` e `destroy.sh` funzionano e sono idempotenti.
- [ ] Il README spiega architettura, sicurezza e costi.
- [ ] Il budget è attivo e nessuna risorsa è rimasta orfana dopo la distruzione.

??? success "Una possibile soluzione: deploy.sh per AWS"
    ```bash
    #!/usr/bin/env bash
    # Pubblica site/ sul bucket e invalida la cache di CloudFront.
    # Uso: BUCKET=nome DISTRIBUTION_ID=E123ABC scripts/deploy.sh
    set -euo pipefail
    : "${BUCKET:?imposta BUCKET}"
    : "${DISTRIBUTION_ID:?imposta DISTRIBUTION_ID}"

    cd "$(dirname "$0")/.."
    aws s3 sync site/ "s3://$BUCKET/" --delete
    id=$(aws cloudfront create-invalidation --distribution-id "$DISTRIBUTION_ID" \
           --paths '/*' --query 'Invalidation.Id' --output text)
    echo "Invalidazione $id avviata"
    aws cloudfront wait invalidation-completed --distribution-id "$DISTRIBUTION_ID" --id "$id"
    echo "Pubblicato su https://$(aws cloudfront get-distribution --id "$DISTRIBUTION_ID" \
           --query 'Distribution.DomainName' --output text)"
    ```

    **Spiegazione.** `sync --delete` rimuove dal bucket i file cancellati in locale. L'invalidazione forza CloudFront a
    rileggere i file dall'origine: le prime 1.000 invalidazioni al mese sono gratuite. Nella Fase 7 questo script
    girerà in una pipeline di GitHub Actions a ogni push, e nella Fase 6 l'infrastruttura sarà creata da Terraform.

## Fonti

- [Amazon CloudFront: Restrict access to an Amazon S3 origin (OAC)](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/private-content-restricting-access-to-s3.html)
- [Amazon CloudFront: Invalidating files](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/Invalidation.html)
- [Amazon CloudFront: Custom error responses](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/GeneratingCustomErrorResponses.html)
- [Cloud Storage: Host a static website](https://docs.cloud.google.com/storage/docs/hosting-static-website)
- [Cloud CDN: Set up a backend bucket](https://docs.cloud.google.com/cdn/docs/setting-up-cdn-with-bucket)
