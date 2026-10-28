# Verifica di fine fase 8

Hai finito il percorso quando ottieni **almeno l'80%** (20 risposte corrette su 25) **e** sai svolgere a tempo i
laboratori della fase. Le spiegazioni indicano la lezione da ripassare. Per l'esame CKA, che è interamente pratico,
allenati con la [simulazione pratica CKA](../appendici/cka/simulazione.md).

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 8 · 25 domande","questions":[
 {"q":"Quale componente sceglie il nodo su cui eseguire un nuovo Pod?","options":["kubelet","kube-scheduler","kube-controller-manager","kube-proxy"],"answer":[1],"explanation":"Lezione 8.1. Il kubelet poi lo avvia."},
 {"q":"Due dei tre membri di etcd sono irraggiungibili. Cosa succede?","options":["Il cluster continua a funzionare normalmente","L'API server non può salvare modifiche; i Pod già in esecuzione continuano a girare","Tutti i Pod si fermano","I nodi worker si spengono"],"answer":[1],"explanation":"Senza quorum etcd non accetta scritture (lezione 8.1)."},
 {"q":"Quale comando genera il YAML di un Deployment senza crearlo?","options":["kubectl create deployment web --image=nginx --dry-run=client -o yaml","kubectl get deployment web -o yaml","kubectl explain deployment","kubectl apply --dry-run"],"answer":[0],"explanation":"Lezione 8.2."},
 {"q":"Una sonda fallisce e il Pod viene tolto dagli endpoint del Service, senza riavvio. Quale sonda?","options":["livenessProbe","readinessProbe","startupProbe","healthProbe"],"answer":[1],"explanation":"Lezione 8.3."},
 {"q":"Perché nel laboratorio la liveness della Tasks API non usa /health?","options":["Perché /health è lento","Perché /health controlla il database: un guasto del database farebbe riavviare tutti i Pod della API","Perché la liveness non supporta HTTP","Per risparmiare memoria"],"answer":[1],"explanation":"Lezione 8.3."},
 {"q":"Un Deployment con replicas: 3, maxSurge: 1, maxUnavailable: 0 viene aggiornato a un'immagine inesistente. Cosa succede?","options":["L'applicazione si ferma","I tre Pod vecchi restano in servizio e il nuovo resta in ImagePullBackOff","Kubernetes torna indietro da solo subito","Il Deployment viene eliminato"],"answer":[1],"explanation":"Lezione 8.4."},
 {"q":"Quale controller usi per un agente di log su ogni nodo?","options":["Deployment","DaemonSet","StatefulSet","CronJob"],"answer":[1],"explanation":"Lezione 8.4."},
 {"q":"Un Service ha zero endpoint. Quali sono le cause probabili?","options":["Il selettore non corrisponde alle label dei Pod","I Pod non sono pronti","Il namespace ha troppi Service","Manca un Ingress","Il Service è di tipo ClusterIP"],"answer":[0,1],"explanation":"Lezione 8.5."},
 {"q":"Da un Pod nel namespace web, come raggiungi il Service db del namespace dati?","options":["db","db.dati","dati.db","db.web.svc"],"answer":[1],"explanation":"O il nome completo db.dati.svc.cluster.local (lezione 8.5)."},
 {"q":"Cosa protegge davvero i Secret di Kubernetes?","options":["La codifica base64","RBAC che limita chi può leggerli e la cifratura a riposo di etcd","Il nome del Secret","Il fatto che non compaiono in kubectl get"],"answer":[1],"explanation":"base64 non è cifratura (lezione 8.6)."},
 {"q":"Modifichi una ConfigMap montata come variabili d'ambiente. Come fai arrivare i nuovi valori ai Pod?","options":["Non serve nulla","Con un rollout: kubectl rollout restart o un nome nuovo della ConfigMap (generatore di Kustomize)","Con kubectl apply sulla ConfigMap","Riavviando l'API server"],"answer":[1],"explanation":"Lezione 8.6."},
 {"q":"Quale modalità di accesso permette la scrittura da Pod su più nodi?","options":["ReadWriteOnce","ReadWriteMany","ReadOnlyMany","ReadWriteOncePod"],"answer":[1],"explanation":"Serve un file system di rete (lezione 8.7)."},
 {"q":"Quali affermazioni sugli StatefulSet sono vere?","options":["I Pod hanno nomi stabili come postgres-0","Ogni Pod ha il suo PVC dai volumeClaimTemplates","Eliminando lo StatefulSet, i PVC restano per impostazione predefinita","I Pod condividono sempre lo stesso volume","Non hanno bisogno di un Service"],"answer":[0,1,2],"explanation":"Lezione 8.7."},
 {"q":"Un Pod resta Pending con '3 Insufficient memory'. Su cosa si basa lo scheduler?","options":["Sul consumo reale di memoria dei nodi","Sulle richieste (requests) dei Pod già assegnati","Sui limiti","Sullo swap"],"answer":[1],"explanation":"Lezione 8.8."},
 {"q":"Vuoi riservare un nodo ai soli Pod con GPU. Cosa usi?","options":["Solo una toleration sui Pod GPU","Un taint sul nodo, più toleration e affinità sui Pod GPU","Un ResourceQuota","Una NetworkPolicy"],"answer":[1],"explanation":"La toleration permette, l'affinità impone (lezione 8.8)."},
 {"q":"Cosa serve perché un HPA basato sulla CPU funzioni?","options":["metrics-server e requests.cpu nei container","Prometheus obbligatoriamente","Un Ingress","Un PodDisruptionBudget"],"answer":[0],"explanation":"L'utilizzo si calcola rispetto alla richiesta (lezione 8.8)."},
 {"q":"Nella Gateway API, quale risorsa definisce porte, protocolli e da quali namespace accettare le route?","options":["GatewayClass","Gateway","HTTPRoute","ReferenceGrant"],"answer":[1],"explanation":"Lezione 8.9."},
 {"q":"Come si realizza un rilascio canary al 10% con la Gateway API?","options":["Con il campo weight nei backendRefs di un'HTTPRoute","Con un'annotazione","Con un NodePort","Non è possibile"],"answer":[0],"explanation":"Lezione 8.9."},
 {"q":"Dopo una policy che nega tutto il traffico in uscita, quale regola aggiungi per prima?","options":["Verso Internet","Verso il DNS del cluster su UDP e TCP 53","Verso l'API server","Verso il Gateway"],"answer":[1],"explanation":"Senza DNS nessun nome si risolve (lezione 8.10)."},
 {"q":"In una NetworkPolicy, namespaceSelector e podSelector nello stesso elemento di from significano:","options":["OR","AND: i Pod con quelle label in quei namespace","Che la policy è non valida","Tutti i Pod del cluster"],"answer":[1],"explanation":"Due elementi separati significano OR (lezione 8.10)."},
 {"q":"Vuoi dare il ClusterRole edit a un gruppo solo nel namespace shop. Cosa crei?","options":["Un ClusterRoleBinding","Un RoleBinding nel namespace shop che punta al ClusterRole edit","Un Role nuovo","Un service account"],"answer":[1],"explanation":"Lezione 8.11."},
 {"q":"Come verifichi se il service account ci del namespace ci può eliminare Pod nel namespace tasks?","options":["kubectl auth can-i delete pods -n tasks --as system:serviceaccount:ci:ci","kubectl get sa ci -o yaml","kubectl describe role ci","kubectl auth whoami"],"answer":[0],"explanation":"Lezione 8.11."},
 {"q":"Quali passi fanno parte dell'aggiornamento del control plane con kubeadm?","options":["Aggiornare il pacchetto kubeadm","kubeadm upgrade plan e kubeadm upgrade apply","Drain, aggiornamento di kubelet e kubectl, riavvio del kubelet, uncordon","Saltare direttamente due versioni minori","Reinstallare etcd da zero"],"answer":[0,1,2],"explanation":"Una minore alla volta (lezione 8.12)."},
 {"q":"I nuovi Pod restano Pending senza eventi FailedScheduling. Il manifest di quale componente controlli?","options":["/etc/kubernetes/manifests/kube-scheduler.yaml","/etc/kubernetes/manifests/etcd.yaml","/var/lib/kubelet/config.yaml del worker","Il Corefile di CoreDNS"],"answer":[0],"explanation":"Lezione 8.14."},
 {"q":"In Argo CD, cosa fa l'opzione prune: true?","options":["Annulla le modifiche manuali","Elimina dal cluster le risorse tolte da Git","Comprime le immagini","Cancella le vecchie revisioni Git"],"answer":[1],"explanation":"selfHeal annulla le modifiche manuali (lezione 8.15)."}
]}
</script>

## E adesso?

Hai completato il percorso. Ora:

1. Rifinisci i progetti delle fasi 1-8 nel tuo GitHub: README chiari, diagrammi, istruzioni per riprodurli.
2. Prepara le certificazioni con le [appendici](../appendici/index.md): simulazioni a tempo finché non superi la
   soglia con margine, più volte di fila.
3. Riprendi la [roadmap](../percorso.md): curriculum, profilo LinkedIn e ricerca del lavoro sono i passi successivi
   del piano.
