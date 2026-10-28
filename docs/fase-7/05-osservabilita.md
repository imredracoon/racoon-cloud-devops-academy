# 7.5 Osservabilità

!!! obiettivi "Obiettivi"
    - Distinguere monitoraggio e osservabilità.
    - Descrivere i tre segnali: metriche, log, tracce.
    - Usare i metodi RED e USE per scegliere cosa misurare.
    - Definire SLI, SLO e budget di errore.
    - Scrivere alert utili, che segnalano i sintomi e non svegliano nessuno senza motivo.

## Monitoraggio e osservabilità

- Il **monitoraggio** risponde a domande **note in anticipo**: "la CPU supera l'80%?", "il sito risponde?". Funziona
  bene per i guasti che hai già visto.
- L'**osservabilità** è la capacità di capire lo stato interno di un sistema dai dati che produce, anche per problemi
  **mai visti prima**: "perché le richieste degli utenti di un certo paese sono lente da ieri sera?".

Un sistema è osservabile se produce dati ricchi e collegabili tra loro. Il monitoraggio è una parte
dell'osservabilità.

## I tre segnali

| Segnale | Cos'è | Domanda tipica | Strumenti del corso |
|---|---|---|---|
| **Metriche** | numeri nel tempo, aggregati ed economici da conservare | quante richieste al secondo? quanti errori? | Prometheus, CloudWatch, Cloud Monitoring |
| **Log** | eventi testuali con data e ora, dettagliati | cosa è successo esattamente alle 10:42? | Loki, CloudWatch Logs, Cloud Logging |
| **Tracce** | il percorso di una richiesta attraverso più servizi, con i tempi di ogni passaggio | dove si perdono i 2 secondi di questa richiesta? | X-Ray, Cloud Trace, Grafana Tempo, Jaeger |

Si usano insieme: un **alert** su una metrica ti avvisa, il **grafico** ti dice quando è iniziato, i **log** ti dicono
cosa è successo, le **tracce** dove.

!!! info "OpenTelemetry"
    **OpenTelemetry** (OTel) è lo standard aperto della CNCF per produrre e trasportare metriche, log e tracce in modo
    indipendente dal fornitore: strumenti l'applicazione una volta e scegli dove inviare i dati. È supportato da
    Prometheus, Grafana, AWS, Google Cloud e da quasi tutti i fornitori. In questo corso usiamo direttamente il
    formato di Prometheus, il più semplice per cominciare.

## Cosa misurare: RED e USE

Per i **servizi** (un'API, un sito) il metodo **RED**:

- **Rate**: richieste al secondo;
- **Errors**: richieste fallite al secondo (o percentuale);
- **Duration**: quanto durano le richieste (distribuzione, non solo media).

Per le **risorse** (CPU, memoria, disco, rete) il metodo **USE**:

- **Utilization**: quanto è occupata la risorsa;
- **Saturation**: quanto lavoro è in coda;
- **Errors**: errori della risorsa.

I **quattro segnali d'oro** del libro *Site Reliability Engineering* di Google (latenza, traffico, errori,
saturazione) uniscono le due prospettive.

### Medie e percentili

La **media** della latenza nasconde i problemi: se 99 richieste durano 10 ms e una dura 5 secondi, la media è circa
60 ms e sembra tutto bene, ma un utente su cento aspetta 5 secondi. Si usano i **percentili**:

- **p50** (mediana): metà delle richieste è più veloce;
- **p95**: il 95% delle richieste è più veloce di questo valore;
- **p99**: rivela le code lente che colpiscono gli utenti più attivi.

## SLI, SLO, SLA

| Termine | Significato | Esempio |
|---|---|---|
| **SLI** (*Service Level Indicator*) | una misura della qualità del servizio vista dall'utente | percentuale di richieste `/tasks` con risposta non 5xx in meno di 300 ms |
| **SLO** (*Service Level Objective*) | l'obiettivo interno per quello SLI, su una finestra di tempo | 99,5% su 30 giorni |
| **SLA** (*Service Level Agreement*) | un impegno contrattuale con il cliente, con penali | 99% al mese, altrimenti rimborso del 10% |

L'SLO è più severo dell'SLA, così te ne accorgi prima di dover pagare penali.

### Il budget di errore

Uno SLO del 99,5% significa che lo 0,5% delle richieste **può** fallire: è il **budget di errore**. In 30 giorni,
lo 0,5% del tempo sono circa **3 ore e 36 minuti**.

| SLO | Indisponibilità ammessa in 30 giorni |
|---|---|
| 99% | 7 ore e 12 minuti |
| 99,5% | 3 ore e 36 minuti |
| 99,9% | 43 minuti e 12 secondi |
| 99,99% | 4 minuti e 19 secondi |

Il budget guida le decisioni: se è ampio, si può rilasciare più spesso e sperimentare; se è esaurito, la priorità
torna all'affidabilità. Nessun servizio ha bisogno del 100%: ogni "nove" in più costa molto di più.

## Alert che funzionano

Un alert sveglia una persona: deve valerne la pena.

- **Allerta sui sintomi**, non sulle cause: "il 5% delle richieste fallisce" (lo vede l'utente) è meglio di "CPU al
  90%" (può essere normale).
- Ogni alert deve essere **azionabile**: chi lo riceve deve sapere cosa fare. Collega un **runbook** (una pagina con
  i passi di diagnosi) nelle annotazioni.
- Usa una **durata** (`for: 5m`) per evitare falsi allarmi su picchi di pochi secondi.
- **Gravità** diverse: `critical` sveglia qualcuno, `warning` apre un ticket per l'orario di lavoro.
- Pochi alert buoni battono tanti alert ignorati: la **stanchezza da alert** è uno dei problemi più comuni.

## Esercizi

!!! esercizio "Esercizio · Metrica, log o traccia?"
    1. Sapere quante richieste al secondo riceve la Tasks API nelle ultime 24 ore.
    2. Capire perché la richiesta di un cliente alle 15:03 ha restituito 500.
    3. Scoprire quale dei cinque microservizi chiamati da una pagina rallenta il caricamento.
    4. Ricevere un avviso se lo spazio su disco del database scende sotto il 10%.

??? success "Soluzione"
    1. **Metrica** (un contatore di richieste e la sua velocità).
    2. **Log** (il messaggio d'errore e lo stack trace di quella richiesta), magari trovato partendo da una traccia.
    3. **Traccia** distribuita.
    4. **Metrica** con un **alert** (metodo USE: utilizzo del disco).

!!! esercizio "Esercizio · Definisci uno SLO"
    La Tasks API riceve circa 1 milione di richieste al mese. Proponi uno SLI e uno SLO di disponibilità, e calcola
    quante richieste possono fallire in un mese.

??? success "Soluzione"
    SLI: percentuale di richieste che **non** ricevono un errore 5xx, misurata dal load balancer (o da Prometheus con
    `http_requests_total`). SLO: **99,5%** su 30 giorni. Budget: 0,5% di 1.000.000 = **5.000 richieste** fallite al
    mese. Se a metà mese ne sono già fallite 4.000, è il momento di rallentare i rilasci e investire in affidabilità.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quali sono i tre segnali dell'osservabilità?","options":["Metriche","Log","Tracce","Backup","Snapshot"],"answer":[0,1,2],"explanation":"Si usano insieme per passare dall'allarme alla causa."},
 {"q":"Cosa misura il metodo RED per un servizio?","options":["Rate, Errors, Duration","RAM, Errori, Disco","Richieste, Eventi, Dati","Rete, Energia, Disponibilità"],"answer":[0],"explanation":"Il metodo USE è invece per le risorse."},
 {"q":"Perché si preferisce il p95 o il p99 della latenza alla media?","options":["Sono più facili da calcolare","La media nasconde le richieste lente che colpiscono una parte degli utenti","La media non esiste in Prometheus","Per risparmiare spazio"],"answer":[1],"explanation":"Un utente su cento può aspettare secondi senza che la media se ne accorga."},
 {"q":"Qual è la differenza tra SLO e SLA?","options":["Nessuna","L'SLO è un obiettivo interno, l'SLA un impegno contrattuale con penali","L'SLA è interno, l'SLO è contrattuale","L'SLO riguarda solo i costi"],"answer":[1],"explanation":"L'SLO di solito è più severo dell'SLA."},
 {"q":"Con uno SLO del 99,9% su 30 giorni, quanta indisponibilità è ammessa?","options":["Circa 4 minuti","Circa 43 minuti","Circa 7 ore","Circa 3 giorni"],"answer":[1],"explanation":"0,1% di 43.200 minuti = 43,2 minuti."},
 {"q":"Quale alert è migliore secondo le buone pratiche?","options":["CPU sopra il 70% per 10 secondi","Più del 5% delle richieste fallisce da 5 minuti","Un alert per ogni riga di log con WARNING","Memoria sopra il 50%"],"answer":[1],"explanation":"Segnala un sintomo visto dagli utenti, con una durata che evita falsi allarmi."}
]}
</script>

## Fonti

- [Google SRE Book: Monitoring Distributed Systems (i quattro segnali d'oro)](https://sre.google/sre-book/monitoring-distributed-systems/)
- [Google SRE Book: Service Level Objectives](https://sre.google/sre-book/service-level-objectives/)
- [Google SRE Workbook: Implementing SLOs](https://sre.google/workbook/implementing-slos/) e [Alerting on SLOs](https://sre.google/workbook/alerting-on-slos/)
- [Brendan Gregg: The USE Method](https://www.brendangregg.com/usemethod.html)
- [Grafana Labs: The RED Method](https://grafana.com/blog/2018/08/02/the-red-method-how-to-instrument-your-services/)
- [OpenTelemetry: What is OpenTelemetry?](https://opentelemetry.io/docs/what-is-opentelemetry/)
- [AWS Well-Architected: Operational Excellence, observability](https://docs.aws.amazon.com/wellarchitected/latest/operational-excellence-pillar/implement-observability.html)
