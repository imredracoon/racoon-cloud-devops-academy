# 5.6 Sicurezza e crittografia

!!! obiettivi "Obiettivi"
    - Gestire le chiavi di crittografia con un servizio KMS e capire la crittografia a busta.
    - Conservare e ruotare i segreti con un gestore dedicato.
    - Ottenere e rinnovare certificati TLS gestiti.
    - Proteggere le applicazioni web da attacchi applicativi e DDoS.
    - Conoscere i servizi di rilevamento delle minacce, delle vulnerabilità e dei dati sensibili.
    - Gestire le identità degli utenti finali di un'applicazione.

## Gestione delle chiavi

Cifrare è facile; il problema difficile è **proteggere le chiavi**. Un servizio **KMS** (*Key Management Service*)
conserva le chiavi in hardware protetto (HSM), non le fa mai uscire in chiaro e registra ogni utilizzo.

### La crittografia a busta

Per cifrare grandi quantità di dati non si invia tutto al KMS. Si usa la **crittografia a busta** (*envelope
encryption*):

1. il servizio genera una **chiave dati** (*data key*) e la usa per cifrare i dati localmente;
2. la chiave dati viene a sua volta cifrata con la **chiave principale** nel KMS (la *KEK*, *key encryption key*);
3. si conservano insieme i dati cifrati e la chiave dati cifrata;
4. per leggere, si chiede al KMS di decifrare la chiave dati (operazione autorizzata e registrata), poi si decifrano i
   dati.

```mermaid
flowchart LR
  D[Dati in chiaro] -->|cifrati con| DK[Chiave dati]
  DK -->|cifrata con| KEK[Chiave principale<br>nel KMS, non esce mai]
  D2[(Dati cifrati +<br>chiave dati cifrata)]
  DK -.-> D2
```

È ciò che S3, EBS, RDS, Cloud Storage e gli altri servizi fanno dietro le quinte quando scegli una chiave gestita dal
cliente.

=== "AWS"

    - **AWS KMS**: chiavi **gestite da AWS** (create dai servizi, rotazione automatica), chiavi **gestite dal cliente**
      (le tue: controlli la **key policy**, la rotazione annuale o personalizzata, la disattivazione e la cancellazione
      programmata, con attesa da 7 a 30 giorni). Chiavi simmetriche e asimmetriche, chiavi **multi-regione**.
    - Ogni chiave ha una **key policy**: senza un permesso esplicito, nemmeno un amministratore IAM può usarla.
    - **AWS CloudHSM**: HSM dedicati a un singolo cliente, per requisiti che impongono il controllo esclusivo
      dell'hardware.

=== "Google Cloud"

    - **Cloud KMS**: **portachiavi** (*key ring*) regionali o globali, **chiavi** con **versioni**, rotazione automatica
      programmabile. Livelli di protezione: software, **HSM** (Cloud HSM), **esterno** (Cloud External Key Manager:
      la chiave resta in un sistema esterno a Google).
    - Permessi IAM sulle chiavi, per esempio `roles/cloudkms.cryptoKeyEncrypterDecrypter`.
    - **Autokey** crea e assegna automaticamente le chiavi CMEK alle nuove risorse.

## I segreti

Password di database, chiavi API, token: mai nel codice, mai nelle immagini, mai negli script di avvio in chiaro.
Si conservano in un **gestore dei segreti** e l'applicazione li legge all'avvio con la sua identità (ruolo o service
account).

=== "AWS"

    | Servizio | Quando |
    |---|---|
    | **Secrets Manager** | segreti con **rotazione automatica** (integrata per RDS, Aurora, Redshift, DocumentDB), replica tra regioni; a pagamento per segreto |
    | **Systems Manager Parameter Store** | configurazione e segreti semplici (`SecureString` cifrati con KMS), gerarchie di parametri; il livello standard è gratuito |

    RDS può anche **gestire da solo** la password dell'utente principale in Secrets Manager, con rotazione automatica.

=== "Google Cloud"

    **Secret Manager**: segreti con **versioni**, replica automatica o in regioni scelte, CMEK opzionale, notifiche di
    rotazione tramite Pub/Sub, permessi IAM per segreto (`roles/secretmanager.secretAccessor`). Cloud Run e GKE
    possono montare i segreti come variabili d'ambiente o file.

## Certificati TLS

- **AWS Certificate Manager** (ACM): certificati pubblici **gratuiti** per ALB, CloudFront, API Gateway, con
  **rinnovo automatico**. Per CloudFront il certificato deve trovarsi in `us-east-1`. Le chiavi private non si
  esportano (salvo i certificati esportabili, a pagamento).
- **Certificate Manager** e i **certificati gestiti da Google** per i load balancer, con rinnovo automatico.

La validazione del dominio avviene di solito con un record **DNS**: per questo DNS e certificati vanno spesso insieme
(lezione 5.7).

## Proteggere il perimetro

| Minaccia | AWS | Google Cloud |
|---|---|---|
| Attacchi applicativi: SQL injection, XSS, bot, limitazione delle richieste | **AWS WAF** su CloudFront, ALB, API Gateway | **Cloud Armor** sui load balancer |
| DDoS di rete e trasporto | **AWS Shield Standard** (gratuito, sempre attivo); **Shield Advanced** con team di risposta e protezione dei costi | protezione DDoS integrata nei load balancer; **Cloud Armor Enterprise** |
| Gestione centralizzata delle regole tra account | **Firewall Manager** | policy firewall gerarchiche |

## Rilevare minacce e vulnerabilità

=== "AWS"

    | Servizio | Cosa fa |
    |---|---|
    | **GuardDuty** | rilevamento delle minacce analizzando CloudTrail, flow log, DNS e altro (per esempio un'istanza che comunica con server di mining noti) |
    | **Inspector** | scansione delle vulnerabilità di istanze EC2, immagini in ECR e funzioni Lambda |
    | **Macie** | individua dati sensibili (dati personali, numeri di carte) nei bucket S3 |
    | **Security Hub** | raccoglie i risultati di tutti i servizi e controlla la conformità a standard (CIS, AWS Foundational Security Best Practices) |
    | **Detective** | indagini: ricostruisce relazioni ed eventi attorno a un risultato sospetto |
    | **AWS Config** | registra la configurazione delle risorse nel tempo e verifica **regole** di conformità |

=== "Google Cloud"

    | Servizio | Cosa fa |
    |---|---|
    | **Security Command Center** | inventario degli asset, configurazioni errate, vulnerabilità, minacce, conformità: il punto centrale |
    | **Sensitive Data Protection** | individua, classifica e maschera i dati sensibili |
    | **Artifact Analysis** | scansione delle vulnerabilità delle immagini in Artifact Registry |
    | **Cloud Asset Inventory** | storico delle risorse e delle loro configurazioni |

## Le identità degli utenti finali

Gli utenti di un'applicazione (i clienti) non sono utenti IAM. Per registrazione, accesso, MFA e accesso con Google o
altri provider si usano:

- **Amazon Cognito**: *user pool* (directory degli utenti, login, token JWT) e *identity pool* (credenziali AWS
  temporanee per gli utenti dell'app);
- **Identity Platform** (e Firebase Authentication) su Google Cloud;
- **Identity-Aware Proxy** per proteggere applicazioni interne con le identità aziendali, senza VPN.

## Laboratorio

!!! laboratorio "Laboratorio · La password esce dallo script di avvio"
    Riprendi il laboratorio 5.4.

    === "AWS"

        1. **Secrets Manager** → **Store a new secret** → *Other type of secret*, chiave `DATABASE_URL`, valore
           `postgresql+psycopg://tasks:PASSWORD@ENDPOINT:5432/tasks`; nome `tasks/database-url`.
        2. Al ruolo IAM delle istanze aggiungi una policy inline che permette **solo** la lettura di quel segreto:

            ```json
            {
              "Version": "2012-10-17",
              "Statement": [{
                "Effect": "Allow",
                "Action": "secretsmanager:GetSecretValue",
                "Resource": "arn:aws:secretsmanager:REGIONE:ACCOUNT:secret:tasks/database-url-*"
              }]
            }
            ```

        3. Nuova versione del launch template:

            ```bash
            #!/bin/bash
            apt-get update && apt-get install -y docker.io jq
            snap install aws-cli --classic
            systemctl enable --now docker
            DATABASE_URL=$(aws secretsmanager get-secret-value --secret-id tasks/database-url \
              --region REGIONE --query SecretString --output text | jq -r .DATABASE_URL)
            docker run -d --restart unless-stopped -p 8000:8000 --name api \
              -e DATABASE_URL="$DATABASE_URL" ghcr.io/TUO-UTENTE/tasks-api:1.1.0
            ```

        4. **Instance refresh** dell'ASG e verifica che l'API funzioni.

    === "Google Cloud"

        1. ```bash
           gcloud services enable secretmanager.googleapis.com
           printf '%s' 'postgresql+psycopg://tasks:PASSWORD@IP-PRIVATO:5432/tasks' | \
             gcloud secrets create tasks-database-url --data-file=-
           gcloud secrets add-iam-policy-binding tasks-database-url \
             --member="serviceAccount:tasks-app@ID-PROGETTO.iam.gserviceaccount.com" \
             --role="roles/secretmanager.secretAccessor"
           ```
        2. Nel modello di istanza, aggiungi l'ambito di accesso *cloud-platform* (o "Consenti l'accesso completo a
           tutte le API Cloud": i permessi effettivi restano quelli IAM del service account) e usa lo script:

            ```bash
            #!/bin/bash
            apt-get update && apt-get install -y docker.io
            systemctl enable --now docker
            DATABASE_URL=$(gcloud secrets versions access latest --secret=tasks-database-url)
            docker run -d --restart unless-stopped -p 8000:8000 --name api \
              -e DATABASE_URL="$DATABASE_URL" ghcr.io/TUO-UTENTE/tasks-api:1.1.0
            ```

            Le immagini Ubuntu di Compute Engine includono già `gcloud`.
        3. Aggiorna il MIG al nuovo modello e verifica.

    **Pulizia.** Elimina il segreto insieme alle altre risorse del laboratorio. (Su AWS, la cancellazione di un segreto
    ha un periodo di attesa da 7 a 30 giorni: puoi forzarla con `--force-delete-without-recovery` per un laboratorio.)

??? success "Spiegazione"
    La password non compare più nel launch template o nel modello: chi li legge vede solo il **nome** del segreto. Il
    valore lo ottiene solo un'identità autorizzata (il ruolo delle istanze, il service account), e ogni lettura è
    registrata negli audit log. Per cambiare la password basta aggiornare il segreto e ricreare le istanze, senza
    modificare il codice. Un passo ulteriore è far leggere il segreto all'applicazione stessa, con l'SDK, invece che
    allo script di avvio.

## Esercizi

!!! esercizio "Esercizio · Il servizio giusto"
    1. Un sito subisce tentativi di SQL injection e un bot che fa migliaia di richieste al minuto da pochi indirizzi.
    2. Il revisore chiede di sapere chi ha usato la chiave che cifra i backup, e quando.
    3. La password del database deve cambiare automaticamente ogni 30 giorni.
    4. Bisogna scoprire se in qualche bucket ci sono documenti d'identità caricati per errore.
    5. Un'istanza ha iniziato a comunicare con indirizzi noti per il controllo di botnet.
    6. Un'app mobile ha bisogno di registrazione utenti con email e accesso con Google.

??? success "Soluzione"
    1. **AWS WAF** con regole gestite contro SQL injection e una **regola di limitazione** basata sulla frequenza;
       **Cloud Armor** con regole preconfigurate e *rate limiting*.
    2. **AWS KMS** con una chiave gestita dal cliente: gli usi sono registrati in **CloudTrail**; **Cloud KMS** con gli
       **audit log** di accesso ai dati.
    3. **Secrets Manager** con rotazione automatica; su Google Cloud **Secret Manager** con notifiche di rotazione e una
       funzione che esegue il cambio.
    4. **Amazon Macie** / **Sensitive Data Protection**.
    5. **Amazon GuardDuty** / **Security Command Center** (rilevamento delle minacce).
    6. **Amazon Cognito** / **Identity Platform** (o Firebase Authentication).

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Nella crittografia a busta, cosa viene cifrato con la chiave principale del KMS?","options":["Tutti i dati","La chiave dati","Le password degli utenti","I log"],"answer":[1],"explanation":"I dati sono cifrati localmente con la chiave dati."},
 {"q":"Quale servizio AWS offre la rotazione automatica integrata della password di un database RDS?","options":["Parameter Store standard","Secrets Manager","KMS","ACM"],"answer":[1],"explanation":"Parameter Store non ha rotazione integrata."},
 {"q":"Quale servizio protegge un'applicazione web da SQL injection e cross-site scripting?","options":["AWS Shield Standard","AWS WAF","Amazon Inspector","AWS Config"],"answer":[1],"explanation":"Su Google Cloud: Cloud Armor."},
 {"q":"Quale servizio AWS individua dati personali sensibili nei bucket S3?","options":["GuardDuty","Macie","Detective","Inspector"],"answer":[1],"explanation":"Inspector cerca vulnerabilità del software."},
 {"q":"Quali affermazioni su AWS Certificate Manager sono corrette?","options":["I certificati pubblici sono gratuiti per i servizi integrati","Il rinnovo è automatico","Puoi esportare la chiave privata di qualsiasi certificato gratuito","Funziona solo con EC2","Per CloudFront il certificato deve essere in qualunque regione"],"answer":[0,1],"explanation":"Per CloudFront il certificato deve essere in us-east-1."},
 {"q":"Quale servizio gestisce registrazione e accesso degli utenti finali di un'app su AWS?","options":["IAM Identity Center","Amazon Cognito","AWS Organizations","IAM users"],"answer":[1],"explanation":"IAM e Identity Center servono per chi amministra il cloud."}
]}
</script>

## Fonti

- [AWS KMS: Concepts e envelope encryption](https://docs.aws.amazon.com/kms/latest/developerguide/kms-cryptography.html#enveloping)
- [AWS Secrets Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html) e [Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html)
- [AWS Certificate Manager](https://docs.aws.amazon.com/acm/latest/userguide/acm-overview.html)
- [AWS WAF, AWS Shield e Firewall Manager](https://docs.aws.amazon.com/waf/latest/developerguide/what-is-aws-waf.html)
- [AWS Security Hub](https://docs.aws.amazon.com/securityhub/latest/userguide/what-is-securityhub.html), [GuardDuty](https://docs.aws.amazon.com/guardduty/latest/ug/what-is-guardduty.html), [Macie](https://docs.aws.amazon.com/macie/latest/user/what-is-macie.html)
- [Amazon Cognito](https://docs.aws.amazon.com/cognito/latest/developerguide/what-is-amazon-cognito.html)
- [Cloud KMS: Envelope encryption](https://docs.cloud.google.com/kms/docs/envelope-encryption)
- [Secret Manager overview](https://docs.cloud.google.com/secret-manager/docs/overview)
- [Cloud Armor overview](https://docs.cloud.google.com/armor/docs/cloud-armor-overview), [Security Command Center](https://docs.cloud.google.com/security-command-center/docs/security-command-center-overview)
