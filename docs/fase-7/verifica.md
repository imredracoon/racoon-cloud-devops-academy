# Verifica di fine fase 7

Passa alla Fase 8 quando ottieni **almeno l'80%** (16 risposte corrette su 20). Le spiegazioni indicano la lezione
da ripassare.

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 7 · 20 domande","questions":[
 {"q":"Cosa distingue il continuous deployment dalla continuous delivery?","options":["Il deployment rilascia in produzione automaticamente, senza approvazione manuale","La delivery non usa test","Il deployment non produce artefatti","Nessuna differenza"],"answer":[0],"explanation":"Lezione 7.1."},
 {"q":"Quale strategia di rilascio manda una piccola percentuale del traffico alla nuova versione per verificarla?","options":["Ricreazione","Blue/green","Canary","Rolling"],"answer":[2],"explanation":"Lezione 7.1."},
 {"q":"Quali sono metriche DORA?","options":["Frequenza dei deploy","Tempo di ripristino","Numero di commit al giorno","Copertura dei test","Dimensione delle immagini"],"answer":[0,1],"explanation":"Le altre due sono il tempo di consegna delle modifiche e il tasso di fallimento (lezione 7.1)."},
 {"q":"In GitHub Actions, come fai dipendere un job da un altro?","options":["needs","if","concurrency","matrix"],"answer":[0],"explanation":"Lezione 7.2."},
 {"q":"Perché fissare le azioni di terze parti a uno SHA di commit?","options":["Perché il tag può essere spostato, lo SHA è immutabile","Per velocità","Perché i tag non funzionano","Per risparmiare minuti"],"answer":[0],"explanation":"Lezione 7.2."},
 {"q":"Quale è il modo sicuro di usare il titolo di una pull request in uno script di un workflow?","options":["Inserirlo con ${{ }} direttamente in run","Passarlo come variabile d'ambiente con env e usarlo tra virgolette","Salvarlo come segreto","Non è possibile"],"answer":[1],"explanation":"Evita l'iniezione di comandi (lezione 7.2)."},
 {"q":"Nella pipeline della Tasks API, perché la scansione Trivy precede la pubblicazione?","options":["Per non pubblicare immagini con vulnerabilità gravi","Perché Trivy lo richiede","Per ridurre la cache","Per i tag"],"answer":[0],"explanation":"Lezione 7.3."},
 {"q":"Quale tag d'immagine è adatto a identificare cosa gira in produzione?","options":["latest","main","sha-1a2b3c4 o 1.2.0","dev"],"answer":[2],"explanation":"Tag immutabili (lezione 7.3)."},
 {"q":"Cosa serve perché un job ottenga il token OIDC di GitHub?","options":["permissions: id-token: write","Un segreto AWS_ACCESS_KEY_ID","Un runner self-hosted","Il permesso packages: write"],"answer":[0],"explanation":"Lezione 7.4."},
 {"q":"Nella trust policy AWS per GitHub, cosa limita l'accesso a un solo repository e ambiente?","options":["La condizione sul claim sub","Il nome del ruolo","La regione","Il claim exp"],"answer":[0],"explanation":"Per esempio repo:UTENTE/tasks-api:environment:produzione (lezione 7.4)."},
 {"q":"Su Google Cloud, quale servizio permette a GitHub Actions di autenticarsi senza chiavi?","options":["Workload Identity Federation","Secret Manager","Cloud Build","IAP"],"answer":[0],"explanation":"Lezione 7.4."},
 {"q":"Quale segnale useresti per capire in quale servizio si perde tempo durante una richiesta che ne attraversa cinque?","options":["Metriche","Log","Tracce distribuite","Backup"],"answer":[2],"explanation":"Lezione 7.5."},
 {"q":"Uno SLO del 99,5% su 30 giorni ammette circa...","options":["43 minuti di indisponibilità","3 ore e 36 minuti","7 ore e 12 minuti","1 giorno"],"answer":[1],"explanation":"0,5% di 720 ore (lezione 7.5)."},
 {"q":"Quale alert segue meglio le buone pratiche?","options":["CPU sopra l'80% per 5 secondi","Tasso di errori 5xx sopra il 5% per 5 minuti, con un runbook","Una notifica per ogni riga di log","Memoria sopra il 40%"],"answer":[1],"explanation":"Sintomo visto dagli utenti, azionabile (lezione 7.5)."},
 {"q":"Prometheus raccoglie le metriche...","options":["in pull, interrogando /metrics delle destinazioni","in push, ricevendole dalle applicazioni","dai log","da Grafana"],"answer":[0],"explanation":"Lezione 7.6."},
 {"q":"Quale query restituisce le richieste al secondo per endpoint?","options":["sum by (endpoint) (rate(http_requests_total[5m]))","rate(sum(http_requests_total))","avg(http_requests_total)","http_requests_total / 60"],"answer":[0],"explanation":"rate prima di sum (lezione 7.6)."},
 {"q":"Perché la Tasks API usa la modalità multiprocesso di prometheus_client?","options":["Gunicorn avvia più worker, ognuno con i suoi contatori: vanno sommati","Per cifrare le metriche","Per usare meno CPU","Perché Flask lo richiede"],"answer":[0],"explanation":"Lezione 7.6."},
 {"q":"Come rendi riproducibili sorgenti dati e dashboard di Grafana?","options":["Con il provisioning da file versionati","Con screenshot","Con un backup manuale del volume","Non è possibile"],"answer":[0],"explanation":"Lezione 7.7."},
 {"q":"Cosa indicizza Loki?","options":["Solo le etichette","Tutto il testo","Nulla","Solo gli errori"],"answer":[0],"explanation":"Lezione 7.8."},
 {"q":"Quale query LogQL mostra le righe del servizio api con codice di stato 500 o superiore, dai log di accesso di Gunicorn?","options":["{service=\"api\"} | pattern `<ip> - - [<_>] \"<metodo> <percorso> <_>\" <stato> <_>` | stato >= 500","SELECT * FROM api WHERE stato >= 500","{stato>=500}","rate(api[5m]) > 500"],"answer":[0],"explanation":"Lezione 7.8."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Passa alla [Fase 8 · Kubernetes](../fase-8/index.md): l'ultima fase,
in cui la Tasks API girerà su un orchestratore di container, con tutto ciò che hai imparato finora.
