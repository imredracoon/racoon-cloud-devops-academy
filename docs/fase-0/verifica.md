# Verifica di fine fase 0

Rispondi a tutte le domande senza consultare le lezioni. Passa alla [Fase 1](../fase-1/index.md) quando ottieni **almeno l'80%**
(12 risposte corrette su 15). Se non ci arrivi, ripassa le lezioni indicate nelle spiegazioni e riprova
dopo un giorno.

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 0 · 15 domande","questions":[
 {"q":"Quale memoria è volatile, cioè perde il contenuto allo spegnimento?","options":["SSD","HDD","RAM","Chiavetta USB"],"answer":[2],"explanation":"La RAM è la memoria di lavoro e si svuota allo spegnimento (lezione 0.1)."},
 {"q":"Un disco venduto come 500 GB, in GiB vale circa:","options":["500 GiB","537 GiB","466 GiB","512 GiB"],"answer":[2],"explanation":"500 × 10^9 / 2^30 ≈ 465,7 GiB. I produttori usano le unità decimali (lezione 0.1)."},
 {"q":"Il kernel è:","options":["Un'applicazione grafica","Il nucleo del sistema operativo che gestisce l'hardware","Un tipo di shell","Un gestore di pacchetti"],"answer":[1],"explanation":"Lezione 0.2."},
 {"q":"Quale distribuzione appartiene alla famiglia Red Hat?","options":["Ubuntu","Debian","Rocky Linux","openSUSE"],"answer":[2],"explanation":"Rocky Linux è compatibile con RHEL e usa dnf (lezione 0.2)."},
 {"q":"Quale protocollo traduce un nome di dominio in indirizzo IP?","options":["HTTP","DNS","TCP","SSH"],"answer":[1],"explanation":"Lezione 0.3."},
 {"q":"Quali indirizzi sono privati secondo la RFC 1918?","options":["10.255.0.1","192.168.100.1","172.15.0.1","172.16.0.1","100.64.0.1"],"answer":[0,1,3],"explanation":"Gli intervalli privati sono 10.0.0.0/8, 172.16.0.0/12 (172.16-172.31) e 192.168.0.0/16. 172.15.0.1 è fuori intervallo; 100.64.0.0/10 è uno spazio condiviso per i provider (RFC 6598), non RFC 1918 (lezione 0.3)."},
 {"q":"Su quale porta ascolta di solito SSH?","options":["21","22","80","443"],"answer":[1],"explanation":"Lezione 0.3."},
 {"q":"Ti trovi in /var/log. Quale comando ti porta in /var?","options":["cd /","cd ..","cd ~","cd -"],"answer":[1],"explanation":".. è la cartella superiore (lezione 0.4)."},
 {"q":"Quale comando cancella una cartella non vuota con tutto il contenuto?","options":["rmdir cartella","rm cartella","rm -r cartella","del cartella"],"answer":[2],"explanation":"rmdir funziona solo con cartelle vuote; rm senza -r rifiuta le cartelle (lezione 0.4)."},
 {"q":"Dopo `echo uno > f.txt` e `echo due > f.txt`, cosa contiene f.txt?","options":["uno","due","uno e due su due righe","Niente"],"answer":[1],"explanation":"> sostituisce il contenuto; per aggiungere serve >> (lezione 0.4)."},
 {"q":"Dove trovi la documentazione completa di un comando direttamente nel terminale?","options":["help.txt","man comando","info.html","Nel file /etc/os-release"],"answer":[1],"explanation":"Lezione 0.4."},
 {"q":"Quale sintassi Markdown crea un collegamento?","options":["<link>testo</link>","[testo](https://indirizzo)","(testo)[https://indirizzo]","{testo:https://indirizzo}"],"answer":[1],"explanation":"Lezione 0.5."},
 {"q":"Secondo il NIST, la possibilità di ottenere risorse da soli, senza interazione umana con il fornitore, si chiama:","options":["Elasticità rapida","Servizio misurato","Self-service on demand","Resource pooling"],"answer":[2],"explanation":"Lezione 0.6."},
 {"q":"Un'azienda carica il codice della sua app e il fornitore gestisce server e sistema operativo. È un esempio di:","options":["IaaS","PaaS","SaaS","On-premises"],"answer":[1],"explanation":"Lezione 0.6."},
 {"q":"Qual è un obiettivo centrale di DevOps?","options":["Rilasciare meno spesso per evitare problemi","Rilasciare più spesso e in modo affidabile grazie ad automazione e collaborazione","Separare sviluppatori e sistemisti","Eliminare i test"],"answer":[1],"explanation":"Lezione 0.6."}
]}
</script>

## E adesso?

Hai completato la Fase 0 e il [progetto](progetto.md)? Passa alla [Fase 1 · Linux e reti](../fase-1/index.md).
