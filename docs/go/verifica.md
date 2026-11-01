# Verifica finale del percorso Go

Il percorso è completo quando ottieni **almeno l'80%** (16 risposte corrette su 20). Puoi fare la verifica anche a
metà strada: le domande indicano la lezione da ripassare, quindi salta pure quelle delle lezioni che non hai ancora
studiato.

<script type="application/json" class="quiz-data">
{"title":"Verifica finale del percorso Go · 20 domande","questions":[
 {"q":"Perché un binario Go compilato con CGO_ENABLED=0 si copia su un altro server Linux e funziona senza installare nulla?","options":["Perché Go è interpretato","Perché il binario contiene il programma, le librerie usate e il runtime, senza dipendere dalla libc","Perché i server Linux hanno Go preinstallato","Perché usa Docker"],"answer":[1],"explanation":"Lezioni G.1 e G.6."},
 {"q":"Come compili `probe` per un'istanza AWS Graviton dal Dev Container x86-64?","options":["Non si può","GOOS=linux GOARCH=arm64 go build ./cmd/probe","go build --target graviton","Serve un'istanza Graviton"],"answer":[1],"explanation":"Lezione G.1, esercizio 3."},
 {"q":"Quale di queste righe NON compila dentro una funzione?","options":["n := 5","var n int","var n = 5","n := 5 seguito da nessun uso di n"],"answer":[3],"explanation":"Le variabili inutilizzate sono errori di compilazione (lezione G.1)."},
 {"q":"Hai `ram := map[string]int{\"db\": 16}`. Cosa vale `v, ok := ram[\"web\"]`?","options":["v=0, ok=false","Un panic","v=nil, ok=false","Un errore di compilazione"],"answer":[0],"explanation":"Valore zero e ok=false (lezione G.2)."},
 {"q":"Una funzione restituisce `fmt.Errorf(\"lettura: %w\", err)` con err = os.ErrNotExist. Cosa restituisce `errors.Is(risultato, os.ErrNotExist)`?","options":["false","true","Dipende dal sistema operativo","Va in panic"],"answer":[1],"explanation":"%w conserva l'errore originale (lezione G.2)."},
 {"q":"Un metodo deve modificare un campo della struct su cui è chiamato. Come deve essere il ricevitore?","options":["(v VM)","(v *VM)","(VM)","Non si può modificare"],"answer":[1],"explanation":"Con un ricevitore valore il metodo modifica una copia (lezione G.2)."},
 {"q":"Il tipo ControlloTCP ha i metodi Nome() string ed Esegui() error. Cosa serve perché soddisfi l'interfaccia Controllo che richiede proprio quei metodi?","options":["Scrivere implements Controllo","Niente: la soddisfa già","Ereditare da Controllo","Registrarlo con una funzione init"],"answer":[1],"explanation":"Le interfacce sono implicite (lezione G.3)."},
 {"q":"net.DialTimeout verso una porta restituisce `i/o timeout`. Qual è la causa più probabile?","options":["Nessun processo ascolta sulla porta","Un firewall scarta i pacchetti","Il nome DNS non esiste","Il programma non ha i permessi"],"answer":[1],"explanation":"Una porta chiusa senza firewall risponde subito con connection refused (lezione G.3)."},
 {"q":"Nel file probe.yaml una chiave è scritta `expect_staus`. Cosa succede alla lettura con yaml.Unmarshal?","options":["Errore di sintassi","La chiave viene ignorata in silenzio: per questo la configurazione si valida dopo averla letta","Il programma va in panic","Il campo prende il valore della chiave"],"answer":[1],"explanation":"Lezione G.4."},
 {"q":"Quale struttura di test è la più diffusa nei progetti Go?","options":["Un test per ogni riga di codice","Il test a tabella: un elenco di casi e un ciclo che li prova tutti","Test solo end-to-end","Test scritti in un altro linguaggio"],"answer":[1],"explanation":"Lezione G.4."},
 {"q":"Nel metodo All di probe, 100 goroutine scrivono ciascuna in results[i]. Serve un mutex?","options":["Sì, sempre","No: ogni goroutine scrive in una posizione diversa e si legge solo dopo wg.Wait()","Sì, ma solo con più di 8 target","No, perché Go protegge le slice"],"answer":[1],"explanation":"Lezione G.5. Il race detector lo confermerebbe."},
 {"q":"Come fa probe serve a finire le richieste in corso quando Kubernetes lo ferma?","options":["Ignora i segnali","Usa signal.NotifyContext per SIGTERM e poi http.Server.Shutdown","Chiama os.Exit(0) subito","Aspetta SIGKILL"],"answer":[1],"explanation":"Lezione G.5."},
 {"q":"Quale opzione dell'http.Server protegge da client che inviano gli header lentissimamente?","options":["MaxHeaderBytes","ReadHeaderTimeout","IdleConnections","KeepAlive"],"answer":[1],"explanation":"Lezione G.5."},
 {"q":"Perché lo stadio di build del Dockerfile usa `FROM --platform=$BUILDPLATFORM`?","options":["Per usare l'immagine più piccola","Perché il compilatore giri nativo e produca il binario per TARGETARCH senza emulazione","Perché Docker lo richiede","Per eseguire i test"],"answer":[1],"explanation":"Lezione G.6."},
 {"q":"In un programma che usa l'AWS SDK for Go v2 su un'istanza EC2, dove devono stare le credenziali?","options":["In una costante","In un file .env nel repository","Nel ruolo IAM dell'istanza, letto dalla catena predefinita","In una variabile d'ambiente scritta nel Dockerfile"],"answer":[2],"explanation":"Lezione G.7."},
 {"q":"Perché in un test Terratest `defer terraform.DestroyContext(...)` si scrive prima di InitAndApplyContext?","options":["Per velocità","Perché il destroy avvenga anche se l'apply o le verifiche falliscono","Per ordine alfabetico","Perché il destroy va eseguito prima"],"answer":[1],"explanation":"Lezione G.8."},
 {"q":"Quale label rende pericolosa una metrica Prometheus?","options":["target, con pochi valori decisi da te","Un ID diverso per ogni richiesta","type, con i valori http e tcp","result, con success e failure"],"answer":[1],"explanation":"Valori illimitati: problema di cardinalità (lezione G.9)."},
 {"q":"Cosa segnala govulncheck?","options":["Ogni dipendenza vecchia","Le vulnerabilità note nel codice che il programma chiama davvero","Gli errori di formattazione","Le race condition"],"answer":[1],"explanation":"Lezione G.9."},
 {"q":"Un tool Go gira in un Pod e deve solo elencare i Pod di tutti i namespace. Quale RBAC gli dai?","options":["cluster-admin","Un ClusterRole con get e list su pods, legato al suo service account","Un Role con * su tutte le risorse","Nessuno: i Pod vedono tutto"],"answer":[1],"explanation":"Minimo privilegio (lezione G.10)."},
 {"q":"Nell'operator Probe, come si ripete il controllo ogni interval?","options":["Con time.Sleep dentro Reconcile","Restituendo ctrl.Result{RequeueAfter: interval}","Con un CronJob","Con un goroutine per ogni Probe"],"answer":[1],"explanation":"Lezione G.11."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Aggiungi `racoon-probe` e `probe-operator` ai repository fissati
in evidenza sul tuo profilo GitHub, e indica Go tra le competenze nel CV accanto a Python. Per continuare ad
allenarti, il modo migliore è leggere e contribuire al codice degli strumenti che usi: molti progetti della CNCF
segnalano le *issue* adatte ai nuovi arrivati con l'etichetta `good first issue`.
