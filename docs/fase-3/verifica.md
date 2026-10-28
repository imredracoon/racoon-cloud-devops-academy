# Verifica di fine fase 3

Passa alla [Fase 4](../fase-4/index.md) quando ottieni **almeno l'80%** (16 risposte corrette su 20). Le spiegazioni indicano la lezione
da ripassare.

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 3 · 20 domande","questions":[
 {"q":"Perché un container si avvia molto più velocemente di una VM?","options":["Usa meno CPU","Non avvia un sistema operativo: è un gruppo di processi isolati sul kernel dell'host","È compresso","Non ha un filesystem"],"answer":[1],"explanation":"Lezione 3.1."},
 {"q":"Quale meccanismo del kernel fa sì che un container veda solo i propri processi?","options":["cgroup","namespace PID","OverlayFS","seccomp"],"answer":[1],"explanation":"Lezione 3.1."},
 {"q":"Cosa definiscono le specifiche OCI?","options":["Il linguaggio YAML","I formati e il comportamento standard di immagini, runtime e registry","Le API di Kubernetes","Il formato dei log"],"answer":[1],"explanation":"Per questo un'immagine costruita con Docker gira con Podman e Kubernetes (lezione 3.1)."},
 {"q":"`docker run -d -p 127.0.0.1:8080:80 nginx`: da dove è raggiungibile Nginx?","options":["Da Internet sulla porta 80","Solo dall'host stesso, sulla porta 8080","Da tutta la rete locale sulla 8080","Da nessuna parte"],"answer":[1],"explanation":"Lezione 3.2."},
 {"q":"Un container è sparito da `docker ps`. Quali comandi usi per capire cosa è successo?","options":["docker ps -a","docker logs NOME","docker system prune","docker rmi","docker push"],"answer":[0,1],"explanation":"Lezione 3.2."},
 {"q":"Un container termina con codice 137. Qual è la causa più frequente?","options":["Errore di sintassi nel Dockerfile","Il processo è stato terminato con SIGKILL, spesso per aver superato il limite di memoria","Porta già occupata","Immagine non trovata"],"answer":[1],"explanation":"137 = 128 + 9 (lezione 3.2)."},
 {"q":"Ordina il Dockerfile per sfruttare la cache: quale sequenza è corretta?","options":["COPY . . → RUN pip install → CMD","COPY requirements.txt . → RUN pip install → COPY . . → CMD","RUN pip install → COPY requirements.txt . → CMD","CMD → COPY . . → RUN pip install"],"answer":[1],"explanation":"Ciò che cambia raramente in alto (lezione 3.3)."},
 {"q":"Hai scritto `ENV DB_PASSWORD=segreta` nel Dockerfile. Qual è il problema?","options":["Nessuno","La password resta nell'immagine ed è leggibile da chiunque la possieda","ENV non accetta password","Il container non parte"],"answer":[1],"explanation":"I segreti si passano all'avvio (lezione 3.3)."},
 {"q":"Qual è la differenza tra CMD ed ENTRYPOINT?","options":["Nessuna","CMD è il comando predefinito sostituibile; ENTRYPOINT è l'eseguibile fisso a cui si aggiungono argomenti","ENTRYPOINT viene eseguito durante la build","CMD funziona solo con Podman"],"answer":[1],"explanation":"Lezione 3.3."},
 {"q":"Perché nel container Gunicorn ascolta su 0.0.0.0 e non su 127.0.0.1?","options":["Per sicurezza","Perché 127.0.0.1 nel container è raggiungibile solo dall'interno del container, e la porta pubblicata non funzionerebbe","Perché Gunicorn non supporta 127.0.0.1","Per usare IPv6"],"answer":[1],"explanation":"Lezione 3.3."},
 {"q":"In quale caso una build multi-stage riduce di più la dimensione dell'immagine?","options":["Un'applicazione compilata, come un programma Go","Un sito HTML statico","Un'immagine già basata su scratch","Mai"],"answer":[0],"explanation":"Il compilatore resta nella prima fase, l'immagine finale contiene solo l'eseguibile (lezione 3.3)."},
 {"q":"Cosa succede a un volume con nome quando elimini il container che lo usava?","options":["Viene eliminato","Resta, con i suoi dati","Viene svuotato","Diventa in sola lettura"],"answer":[1],"explanation":"Lezione 3.4."},
 {"q":"Due container devono parlarsi per nome. Cosa serve?","options":["Pubblicare le porte con -p","Collegarli alla stessa rete definita dall'utente","Usare la rete bridge predefinita","Modificare /etc/hosts a mano"],"answer":[1],"explanation":"Le reti definite dall'utente offrono il DNS interno (lezione 3.4)."},
 {"q":"Su RHEL con SELinux, come monti la cartella ~/sito in un container Podman?","options":["-v ~/sito:/dati:Z","-v ~/sito:/dati --privileged","setenforce 0 e poi -v","Non è possibile"],"answer":[0],"explanation":"Z rietichetta la cartella per l'uso da parte del container (lezione 3.4)."},
 {"q":"In Compose, l'API parte prima che PostgreSQL accetti connessioni e fallisce. Come risolvi?","options":["depends_on con condition: service_healthy e un healthcheck sul database","sleep 30 nel Dockerfile","restart: always sul database","Pubblicando la porta 5432"],"answer":[0],"explanation":"Lezione 3.5."},
 {"q":"Quale comando ferma lo stack Compose mantenendo i dati del database?","options":["docker compose down -v","docker compose down","docker volume prune","docker system prune -a --volumes"],"answer":[1],"explanation":"-v cancella i volumi (lezione 3.5)."},
 {"q":"Cosa fa `${POSTGRES_PASSWORD:?messaggio}` in un compose.yaml?","options":["Usa una password casuale","Fa fallire Compose con il messaggio se la variabile non è impostata","Cripta la password","Chiede la password a terminale"],"answer":[1],"explanation":"Meglio un errore chiaro che una configurazione sbagliata (lezione 3.5)."},
 {"q":"In produzione, quale riferimento all'immagine è il più affidabile?","options":["tasks-api:latest","tasks-api:1.1.0","tasks-api@sha256:...","tasks-api"],"answer":[2],"explanation":"Il digest non può cambiare (lezione 3.6)."},
 {"q":"Quale opzione NON va usata con container di cui non ti fidi completamente?","options":["--read-only","--cap-drop ALL","--privileged","--memory 256m"],"answer":[2],"explanation":"--privileged dà al container quasi tutti i poteri dell'host (lezione 3.6)."},
 {"q":"Trivy trova molte vulnerabilità in un'immagine basata su python:3.12. Qual è la prima azione efficace?","options":["Passare a python:3.12-slim e ricostruire","Ignorare tutto","Aggiungere --privileged","Disattivare il firewall"],"answer":[0],"explanation":"Meno pacchetti, meno vulnerabilità (lezione 3.6)."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Passa alla [Fase 4 · Cloud](../fase-4/index.md).
