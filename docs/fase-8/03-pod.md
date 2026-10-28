# 8.3 Pod, namespace e label

!!! obiettivi "Obiettivi"
    - Descrivere il Pod e quando mettere più container nello stesso Pod.
    - Usare init container e sidecar.
    - Conoscere fasi e stati di un Pod e le politiche di riavvio.
    - Configurare sonde di salute, richieste e limiti di risorse, contesto di sicurezza.
    - Organizzare gli oggetti con namespace, label, selettori e annotazioni.

## Il Pod

Il **Pod** è l'unità più piccola che Kubernetes esegue: uno o più container che condividono

- lo **stesso indirizzo IP** e le stesse porte (si parlano su `localhost`);
- gli stessi **volumi**, se dichiarati;
- lo stesso destino: vengono schedulati insieme sullo stesso nodo.

Nella maggior parte dei casi un Pod contiene **un solo container applicativo**. Non si creano quasi mai Pod "a mano":
li crea un controller (Deployment, StatefulSet, Job) che li ricrea se muoiono. Un Pod da solo, se il nodo cade, non
torna.

```yaml title="pod.yaml"
apiVersion: v1
kind: Pod
metadata:
  name: web
  labels:
    app: web
    ambiente: dev
spec:
  containers:
    - name: nginx
      image: nginx:1.30
      ports:
        - containerPort: 80
      resources:
        requests:
          cpu: 100m
          memory: 64Mi
        limits:
          memory: 128Mi
```

```bash
kubectl apply -f pod.yaml
kubectl get pod web -o wide
kubectl port-forward pod/web 8080:80     # http://localhost:8080 dal tuo terminale
kubectl logs web
kubectl exec -it web -- sh
kubectl delete pod web
```

### Più container nello stesso Pod

| Schema | Esempio |
|---|---|
| **Init container** | parte prima dei container principali e deve terminare con successo: aspettare il database, preparare file, applicare migrazioni |
| **Sidecar** | resta in esecuzione accanto all'applicazione: un agente che invia i log, un proxy di service mesh |

Dalla versione 1.33 i **sidecar nativi** sono stabili: sono init container con `restartPolicy: Always`, che partono
prima dell'applicazione, restano vivi per tutta la sua durata e si chiudono dopo di lei.

```yaml title="Init container e sidecar"
apiVersion: v1
kind: Pod
metadata:
  name: app-con-aiutanti
spec:
  initContainers:
    - name: aspetta-db
      image: busybox:1.37
      command: ["sh", "-c", "until nc -z postgres 5432; do echo attendo; sleep 2; done"]
    - name: spedisci-log          # sidecar nativo
      image: busybox:1.37
      restartPolicy: Always
      command: ["sh", "-c", "tail -F /log/app.log"]
      volumeMounts:
        - name: log
          mountPath: /log
  containers:
    - name: app
      image: busybox:1.37
      command: ["sh", "-c", "while true; do date >> /log/app.log; sleep 5; done"]
      volumeMounts:
        - name: log
          mountPath: /log
  volumes:
    - name: log
      emptyDir: {}
```

## Il ciclo di vita

| Fase (`status.phase`) | Significato |
|---|---|
| `Pending` | accettato, ma i container non sono ancora partiti: in attesa di nodo o di immagine |
| `Running` | assegnato a un nodo, almeno un container in esecuzione |
| `Succeeded` | tutti i container terminati con successo (tipico dei Job) |
| `Failed` | tutti terminati, almeno uno con errore |
| `Unknown` | il nodo non comunica |

Nella colonna `STATUS` di `kubectl get pods` vedi spesso anche il **motivo** di un container:
`ContainerCreating`, `ImagePullBackOff` (immagine inesistente o registro inaccessibile), `CrashLoopBackOff` (il
container parte e muore di continuo; Kubernetes aspetta sempre più a lungo tra un tentativo e l'altro),
`OOMKilled` (superato il limite di memoria), `Error`, `Completed`.

La `restartPolicy` del Pod decide cosa fa il kubelet quando un container termina: `Always` (predefinita, per i
servizi), `OnFailure` o `Never` (per i Job).

## Le sonde di salute

Il kubelet può controllare periodicamente ogni container con tre **sonde** (*probe*), ciascuna di tipo `httpGet`,
`tcpSocket`, `exec` o `grpc`:

| Sonda | Domanda | Se fallisce |
|---|---|---|
| `startupProbe` | "hai finito di avviarti?" | le altre sonde aspettano; dopo troppi fallimenti il container viene riavviato |
| `livenessProbe` | "sei ancora vivo?" | il container viene **riavviato** |
| `readinessProbe` | "sei pronto a ricevere traffico?" | il Pod viene **tolto dagli endpoint** del Service, senza riavvio |

!!! warning "Non far dipendere la liveness dal database"
    Il percorso `/health` della Tasks API controlla anche il database. Se lo usassi come *liveness*, un guasto del
    database farebbe riavviare in continuazione **tutti** i Pod dell'API, che non hanno colpa. Per questo nel
    laboratorio `/health` è la sonda di **readiness**, mentre la liveness controlla solo che la porta risponda.

## Risorse: richieste e limiti

```yaml
resources:
  requests:        # quanto il Pod chiede: lo scheduler lo usa per scegliere il nodo
    cpu: 250m      # 250 millicore = un quarto di CPU
    memory: 128Mi
  limits:          # quanto può usare al massimo
    memory: 256Mi  # oltre: il container viene ucciso (OOMKilled)
```

- La **CPU** oltre il limite viene **rallentata** (*throttling*); la **memoria** oltre il limite fa **uccidere** il
  container. Molti team impostano il limite di memoria e lasciano la CPU senza limite, per non rallentare i picchi.
- In base a richieste e limiti ogni Pod riceve una **classe QoS**: `Guaranteed` (richieste = limiti per tutte le
  risorse), `Burstable` (qualche richiesta), `BestEffort` (nessuna). Quando un nodo resta senza memoria, il kubelet
  sfratta prima i `BestEffort`, poi i `Burstable`.
- Unità: `m` per millicore; `Mi` e `Gi` sono potenze di 2 (mebibyte, gibibyte), `M` e `G` potenze di 10.

## Il contesto di sicurezza

```yaml
spec:
  securityContext:              # vale per tutti i container del Pod
    runAsNonRoot: true
    runAsUser: 10001
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: api
      securityContext:          # vale per questo container
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: ["ALL"]
```

Sono le stesse protezioni della [Fase 3](../fase-3/index.md) (utente senza privilegi, capacità ridotte), applicate dal cluster. Il
namespace `tasks` del laboratorio impone il profilo **restricted** delle *Pod Security Standards*: un Pod senza queste
impostazioni viene **rifiutato** ([lezione 8.11](11-rbac.md)).

## Namespace

Un **namespace** è una divisione logica del cluster: separa i nomi (due Deployment `api` possono esistere in
namespace diversi), i permessi RBAC, le quote di risorse e le Network Policy.

```bash
kubectl get namespaces
kubectl create namespace tasks
kubectl get pods -n kube-system
kubectl get pods -A                # tutti i namespace
```

Namespace presenti di default: `default`, `kube-system` (componenti del cluster), `kube-public`, `kube-node-lease`
(battito cardiaco dei nodi). Alcuni oggetti non appartengono a nessun namespace: nodi, PersistentVolume,
StorageClass, ClusterRole, i namespace stessi (`kubectl api-resources --namespaced=false`).

## Label, selettori e annotazioni

Le **label** sono coppie chiave-valore sugli oggetti. Non hanno significato per Kubernetes finché un **selettore**
non le usa: è così che un Service trova i suoi Pod e un Deployment riconosce i Pod che gestisce.

```bash
kubectl get pods -l app=web
kubectl get pods -l 'ambiente in (dev,test)'
kubectl get pods -l app=web,ambiente!=prod
kubectl label pod web versione=1.0
kubectl get pods --show-labels
```

Kubernetes raccomanda alcune label comuni, usate anche nel laboratorio: `app.kubernetes.io/name`,
`app.kubernetes.io/instance`, `app.kubernetes.io/version`, `app.kubernetes.io/part-of`,
`app.kubernetes.io/managed-by`.

Le **annotazioni** sono anch'esse coppie chiave-valore, ma non si usano per selezionare: servono per informazioni
destinate a persone o strumenti (il commit da cui viene il rilascio, la configurazione di un controller).

## Esercizi

!!! esercizio "Esercizio 1 · Un Pod a regola d'arte"
    Nel namespace `prove` crea un Pod `cache` con immagine `redis:8`, label `app=cache` e `livello=backend`,
    richiesta di 100m CPU e 64Mi di memoria, limite di memoria 128Mi, e una sonda di liveness TCP sulla porta 6379.
    Poi trova la sua classe QoS.

??? success "Soluzione"
    ```bash
    kubectl create namespace prove
    kubectl run cache -n prove --image=redis:8 --labels=app=cache,livello=backend \
      --dry-run=client -o yaml > cache.yaml
    ```
    Aggiungi a mano risorse e sonda:
    ```yaml
    spec:
      containers:
        - name: cache
          image: redis:8
          resources:
            requests: { cpu: 100m, memory: 64Mi }
            limits: { memory: 128Mi }
          livenessProbe:
            tcpSocket: { port: 6379 }
            periodSeconds: 10
    ```
    ```bash
    kubectl apply -f cache.yaml
    kubectl get pod cache -n prove -o jsonpath='{.status.qosClass}'   # Burstable
    ```
    È `Burstable`: ci sono richieste, ma non coincidono con i limiti per tutte le risorse.

!!! esercizio "Esercizio 2 · Leggi lo stato"
    `kubectl get pods` mostra un Pod con `RESTARTS 7` e stato `CrashLoopBackOff`. Quali due comandi usi per capire
    il motivo, e quale opzione ti mostra i log del tentativo precedente?

??? success "Soluzione"
    `kubectl describe pod NOME` (eventi, codice di uscita, motivo come `OOMKilled`) e `kubectl logs NOME`.
    Il container corrente potrebbe essere appena ripartito e non avere ancora log: `kubectl logs NOME --previous`
    mostra quelli del tentativo fallito.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa condividono i container dello stesso Pod?","options":["Indirizzo IP e porte","I volumi dichiarati nel Pod","Il nodo su cui girano","Il filesystem radice dell'immagine","L'utente con cui girano, sempre"],"answer":[0,1,2],"explanation":"Ogni container ha il proprio filesystem dall'immagine; i volumi si condividono se montati da entrambi."},
 {"q":"Cosa succede se la readiness probe fallisce?","options":["Il container viene riavviato","Il Pod viene tolto dagli endpoint del Service","Il Pod viene eliminato","Il nodo viene drenato"],"answer":[1],"explanation":"Il riavvio è l'effetto della liveness."},
 {"q":"Un container supera il limite di memoria. Cosa succede?","options":["Viene rallentato","Viene ucciso con OOMKilled","Riceve più memoria","Nulla"],"answer":[1],"explanation":"Il throttling vale per la CPU."},
 {"q":"Come si chiama la classe QoS di un Pod senza richieste né limiti?","options":["Guaranteed","Burstable","BestEffort","Default"],"answer":[2],"explanation":"È la prima a essere sfrattata."},
 {"q":"Come trova un Service i Pod a cui mandare traffico?","options":["Per nome del Pod","Con un selettore sulle label","Per indirizzo IP fisso","Tramite le annotazioni"],"answer":[1],"explanation":"Le annotazioni non si usano per selezionare."},
 {"q":"Cosa indica ImagePullBackOff?","options":["Il container è in crash","Il nodo non riesce a scaricare l'immagine","Memoria esaurita","La sonda di readiness fallisce"],"answer":[1],"explanation":"Nome o tag sbagliato, registro privato senza credenziali, rete."}
]}
</script>

## Fonti

- [Kubernetes: Pods](https://kubernetes.io/docs/concepts/workloads/pods/) e [Pod Lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
- [Kubernetes: Init Containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) e [Sidecar Containers](https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/)
- [Kubernetes: Configure Liveness, Readiness and Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [Kubernetes: Resource Management for Pods and Containers](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) e [Pod Quality of Service Classes](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/)
- [Kubernetes: Configure a Security Context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
- [Kubernetes: Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/), [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/), [Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/)
