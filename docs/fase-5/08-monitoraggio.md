# 5.8 Monitoraggio e audit

!!! obiettivi "Obiettivi"
    - Raccogliere metriche e log delle risorse e delle applicazioni nel cloud.
    - Creare allarmi e dashboard.
    - Registrare chi ha fatto cosa con gli audit log, e proteggere questi registri.
    - Verificare la conformità delle configurazioni nel tempo.
    - Tracciare le richieste attraverso più servizi.

Nella [Fase 7](../fase-7/index.md) approfondirai l'osservabilità con strumenti open source (Prometheus, Grafana, Loki). Qui usi i servizi
nativi dei due cloud, che ogni architettura dovrebbe avere fin dal primo giorno.

## Metriche, log, allarmi

=== "AWS"

    **Amazon CloudWatch** è il servizio di monitoraggio di AWS:

    - **metriche**: i servizi pubblicano metriche automaticamente (CPU di EC2, richieste e codici di risposta dell'ALB,
      connessioni di RDS...). Le metriche di EC2 hanno risoluzione predefinita di 5 minuti (1 minuto con il
      monitoraggio dettagliato). Memoria e spazio su disco **non** sono visibili dall'esterno della VM: servono
      l'**agente CloudWatch** o metriche personalizzate;
    - **CloudWatch Logs**: gruppi di log e flussi, con conservazione configurabile, ricerca con **Logs Insights**;
    - **allarmi**: su una metrica e una soglia, con azioni (notifica SNS, azione di Auto Scaling, riavvio di
      un'istanza); gli **allarmi compositi** combinano più allarmi;
    - **dashboard**;
    - **Synthetics** (controlli simulati dall'esterno) e **Application Signals**.

=== "Google Cloud"

    **Google Cloud Observability** (ex Stackdriver):

    - **Cloud Monitoring**: metriche dei servizi, **Ops Agent** sulle VM per memoria, disco e log; dashboard;
      **controlli di uptime** dall'esterno; **criteri di avviso** con canali di notifica (email, SMS, Pub/Sub, webhook);
      **SLO** dichiarati sui servizi;
    - **Cloud Logging**: raccolta di tutti i log, **Esplora log** con il linguaggio di query, **bucket di log** con
      conservazione configurabile, **sink** per esportare verso BigQuery, Cloud Storage o Pub/Sub, metriche basate sui
      log;
    - **Error Reporting**, **Cloud Trace**, **Cloud Profiler**.

## Audit: chi ha fatto cosa

Ogni chiamata alle API del cloud (dalla console, dalla CLI, da Terraform, da un'applicazione) può essere registrata:
chi, quando, da dove, cosa, con quale esito. È indispensabile per la sicurezza, le indagini e la conformità.

=== "AWS"

    **AWS CloudTrail**:

    - la **cronologia degli eventi** degli ultimi 90 giorni (eventi di gestione) è attiva e gratuita;
    - un **trail** salva gli eventi in un bucket S3 (e opzionalmente in CloudWatch Logs) senza limiti di tempo;
      un **trail dell'organizzazione** copre tutti gli account;
    - **eventi di gestione** (creare, modificare, cancellare risorse) ed **eventi sui dati** (per esempio ogni lettura
      di un oggetto S3: volumi enormi, si attivano solo dove servono);
    - **convalida dell'integrità** dei file di log, per dimostrare che non sono stati alterati;
    - **CloudTrail Lake** per interrogare gli eventi con SQL.

    Proteggi il bucket dei log: account separato, Object Lock, accesso minimo.

=== "Google Cloud"

    **Cloud Audit Logs**:

    - **Admin Activity**: modifiche alla configurazione. Sempre attivi, gratuiti, non disattivabili, conservati 400
      giorni;
    - **Data Access**: letture e modifiche ai dati. Disattivati per impostazione predefinita (tranne BigQuery), da
      attivare dove servono;
    - **System Event** e **Policy Denied** (accessi negati dalle policy di sicurezza);
    - un **sink a livello di organizzazione** porta i log di audit in un progetto dedicato con conservazione lunga.

## Conformità delle configurazioni

=== "AWS"

    **AWS Config** registra la configurazione di ogni risorsa e le sue modifiche nel tempo ("com'era configurato questo
    security group martedì scorso?") e valuta **regole**: gestite (per esempio "i volumi EBS devono essere cifrati",
    "nessun bucket pubblico") o personalizzate. I **conformance pack** raggruppano regole secondo uno standard. Le non
    conformità possono attivare una **correzione automatica**.

=== "Google Cloud"

    **Cloud Asset Inventory** conserva la storia delle risorse e delle policy; **Security Command Center** rileva le
    configurazioni errate rispetto a standard (CIS e altri); le **Organization Policy** impediscono a monte le
    configurazioni vietate.

## Gestione operativa

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Eseguire comandi e patch su molte VM | **Systems Manager** (Run Command, Patch Manager, Session Manager, Inventory) | **VM Manager** (patch, inventario, OS policy), OS Login |
| Stato dei servizi del fornitore e degli eventi del tuo account | **AWS Health Dashboard** | **Personalized Service Health** |
| Tracciamento distribuito delle richieste | **AWS X-Ray** (e CloudWatch Application Signals) | **Cloud Trace** |

## Laboratorio

!!! laboratorio "Laboratorio · Allarme e indagine"
    Con l'architettura del laboratorio [5.2](02-ha-scaling.md#laboratorio) attiva:

    === "AWS"

        1. **CloudWatch** → **Alarms** → **Create alarm**: metrica `ApplicationELB` → per target group →
           `UnHealthyHostCount` del tuo target group, statistica **Maximum**, periodo 1 minuto, soglia `>= 1`.
           Azione: notifica a un nuovo topic SNS con la tua email (conferma l'iscrizione dalla mail).
        2. Collegati a un'istanza con Session Manager e ferma il container: `sudo docker stop api`.
        3. Entro qualche minuto ricevi l'email; nel frattempo l'ASG sostituisce l'istanza.
        4. **CloudTrail** → **Event history**: filtra per *Event name* `TerminateInstances` e scopri **chi** ha terminato
           l'istanza (l'Auto Scaling, non tu). Cerca anche `StartSession`: la tua sessione di Session Manager.

        **Pulizia.** Elimina l'allarme e il topic SNS con le altre risorse.

    === "Google Cloud"

        1. **Monitoring** → **Controlli di uptime** → crea un controllo HTTP verso `http://IP-DEL-LB/health`, ogni minuto.
        2. **Avvisi** → **Crea criterio**: condizione sul controllo di uptime fallito; canale di notifica email.
        3. Collegati con IAP a una VM e ferma il container: `sudo docker stop api`. Il bilanciatore toglie la VM e la
           riparazione automatica la ricrea; se le fermi tutte, scatta l'avviso.
        4. **Logging** → **Esplora log**: cerca le voci di audit con
           `protoPayload.methodName:"instances.delete"` per vedere chi ha eliminato e ricreato la VM (il gestore del
           MIG, con un service account di Google).

        **Pulizia.** Elimina controllo di uptime, criterio di avviso e canale.

??? success "Spiegazione"
    Hai collegato i tre pilastri: una **metrica** (host non sani, controllo di uptime) che fa scattare un **allarme**, e
    gli **audit log** che spiegano cosa è successo dopo. In un incidente reale la sequenza è proprio questa: l'allarme ti
    avvisa, metriche e log ti dicono cosa è successo, gli audit log chi ha cambiato cosa.

## Esercizi

!!! esercizio "Esercizio · Il servizio giusto"
    1. Vuoi sapere chi ha cancellato un bucket S3 la settimana scorsa.
    2. Vuoi un avviso se la memoria di una VM supera il 90%.
    3. Un revisore chiede la prova che nessun volume EBS sia mai stato creato senza cifratura negli ultimi sei mesi.
    4. Su Google Cloud vuoi conservare per 7 anni i log di audit di tutti i progetti dell'organizzazione.
    5. Una richiesta attraversa API Gateway, tre funzioni Lambda e DynamoDB ed è lenta: dove si perde tempo?

??? success "Soluzione"
    1. **CloudTrail** (cronologia degli eventi, `DeleteBucket`).
    2. **Agente CloudWatch** (o Ops Agent su Google Cloud) per pubblicare la memoria, poi un **allarme**.
    3. **AWS Config** con la regola gestita `encrypted-volumes` e la cronologia delle configurazioni.
    4. Un **sink aggregato a livello di organizzazione** verso un bucket di log (o Cloud Storage con retention policy
       bloccata) in un progetto dedicato.
    5. **AWS X-Ray** (tracciamento distribuito).

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale servizio registra le chiamate alle API di AWS, con utente, ora e origine?","options":["CloudWatch Metrics","CloudTrail","AWS Config","X-Ray"],"answer":[1],"explanation":"Su Google Cloud: Cloud Audit Logs."},
 {"q":"Perché la memoria di un'istanza EC2 non compare tra le metriche predefinite?","options":["Non è importante","Non è visibile dall'hypervisor: serve un agente nella VM","È a pagamento","Compare solo per Windows"],"answer":[1],"explanation":"Si usa l'agente CloudWatch (o l'Ops Agent su Google Cloud)."},
 {"q":"Quale servizio AWS valuta se le risorse rispettano regole di configurazione nel tempo?","options":["AWS Config","CloudTrail","Trusted Advisor","GuardDuty"],"answer":[0],"explanation":"Registra anche la cronologia delle configurazioni."},
 {"q":"Quali log di audit di Google Cloud sono sempre attivi e non disattivabili?","options":["Data Access","Admin Activity","Debug","Nessuno"],"answer":[1],"explanation":"I Data Access vanno attivati esplicitamente (tranne BigQuery)."},
 {"q":"Cosa può fare un allarme CloudWatch?","options":["Inviare una notifica tramite SNS","Attivare una politica di Auto Scaling","Cancellare l'account","Modificare il codice dell'applicazione","Rinnovare un certificato"],"answer":[0,1],"explanation":"Può anche riavviare o recuperare un'istanza."},
 {"q":"Come dimostri che i file di CloudTrail non sono stati alterati?","options":["Con la convalida dell'integrità dei file di log","Con Trusted Advisor","Con un backup","Con Route 53"],"answer":[0],"explanation":"CloudTrail produce file di digest firmati."}
]}
</script>

## Fonti

- [Amazon CloudWatch User Guide](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/WhatIsCloudWatch.html)
- [Amazon CloudWatch: Collect metrics and logs with the CloudWatch agent](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Install-CloudWatch-Agent.html)
- [AWS CloudTrail User Guide](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-user-guide.html)
- [AWS Config Developer Guide](https://docs.aws.amazon.com/config/latest/developerguide/WhatIsConfig.html)
- [AWS Systems Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/what-is-systems-manager.html), [AWS X-Ray](https://docs.aws.amazon.com/xray/latest/devguide/aws-xray.html)
- [Google Cloud Observability](https://docs.cloud.google.com/stackdriver/docs)
- [Cloud Audit Logs overview](https://docs.cloud.google.com/logging/docs/audit)
- [Cloud Monitoring: Uptime checks](https://docs.cloud.google.com/monitoring/uptime-checks) e [Alerting](https://docs.cloud.google.com/monitoring/alerts)
