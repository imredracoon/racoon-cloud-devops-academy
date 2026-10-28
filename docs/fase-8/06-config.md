# 8.6 ConfigMap e Secret

!!! obiettivi "Obiettivi"
    - Separare la configurazione dall'immagine con ConfigMap e Secret.
    - Passare i valori ai container come variabili d'ambiente o come file.
    - Capire cosa protegge (e cosa non protegge) un Secret.
    - Conoscere le soluzioni per non tenere segreti in Git: gestori esterni e cifratura.

## Perché separare la configurazione

Un principio delle applicazioni *twelve-factor*, già incontrato nella Fase 3: **la stessa immagine** deve girare in
sviluppo, collaudo e produzione; cambia solo la configurazione. In Kubernetes la configurazione vive in due oggetti:

| Oggetto | Per cosa |
|---|---|
| **ConfigMap** | valori non sensibili: livello di log, URL di servizi, file di configurazione interi |
| **Secret** | valori sensibili: password, token, chiavi, certificati TLS |

Entrambi sono limitati a 1 MiB e appartengono a un namespace: un Pod può usare solo quelli del proprio namespace.

## ConfigMap

```bash
kubectl create configmap tasks-config \
  --from-literal=LOG_LEVEL=info \
  --from-literal=PAGINA_MAX=50 \
  --from-file=nginx.conf                     # chiave = nome del file, valore = contenuto
kubectl get configmap tasks-config -o yaml
```

### Come variabili d'ambiente

```yaml
containers:
  - name: api
    image: tasks-api:dev
    envFrom:                       # tutte le chiavi come variabili
      - configMapRef:
          name: tasks-config
    env:
      - name: LIVELLO              # una sola chiave, con un nome diverso
        valueFrom:
          configMapKeyRef:
            name: tasks-config
            key: LOG_LEVEL
```

### Come file

```yaml
containers:
  - name: proxy
    image: nginx:1.30
    volumeMounts:
      - name: conf
        mountPath: /etc/nginx/conf.d     # ogni chiave diventa un file in questa cartella
        readOnly: true
volumes:
  - name: conf
    configMap:
      name: tasks-config
      items:                             # facoltativo: solo alcune chiavi
        - key: nginx.conf
          path: default.conf
```

!!! warning "Gli aggiornamenti"
    Se modifichi una ConfigMap, i file montati come volume si aggiornano da soli dopo qualche decina di secondi (ma
    l'applicazione deve rileggerli), mentre le **variabili d'ambiente no**: si leggono solo all'avvio del container.
    Montaggi con `subPath` non si aggiornano mai. Il modo più robusto è far cambiare il nome della ConfigMap a ogni
    modifica, così il Deployment fa un rollout: è quello che fa il `configMapGenerator` di Kustomize aggiungendo un
    suffisso, come `tasks-api-g8g6m986bk` nel laboratorio.

## Secret

```bash
kubectl create secret generic tasks-db --from-literal=password='S3gr3t0-lungo'
kubectl create secret tls tasks-tls --cert=tls.crt --key=tls.key
kubectl create secret docker-registry ghcr --docker-server=ghcr.io \
  --docker-username=TUO-UTENTE --docker-password="$TOKEN"
kubectl get secret tasks-db -o yaml
```

```yaml
data:
  password: UzNncjN0MC1sdW5nbw==     # base64
```

!!! danger "base64 non è cifratura"
    `echo UzNncjN0MC1sdW5nbw== | base64 -d` restituisce la password in chiaro. Chiunque possa leggere il Secret con
    `kubectl get` la conosce. Un Secret è più protetto di una ConfigMap solo perché:

    - si può dare il permesso di leggere le ConfigMap ma **non** i Secret (RBAC, lezione 8.11);
    - si può configurare la **cifratura a riposo** in etcd (`EncryptionConfiguration`, meglio con un provider KMS);
      su EKS e GKE si attiva con una chiave KMS del cloud;
    - il kubelet li tiene in memoria (`tmpfs`) sui nodi, non su disco.

Si usano esattamente come le ConfigMap, con `secretRef`, `secretKeyRef` e volumi di tipo `secret`. Nel laboratorio la
password arriva al container da `secretKeyRef`, e la variabile `DATABASE_URL` la riusa con la sintassi
`$(POSTGRES_PASSWORD)`:

```yaml title="labs/k8s/base/api.yaml (estratto)"
env:
  - name: POSTGRES_PASSWORD
    valueFrom:
      secretKeyRef:
        name: tasks-db
        key: password
  - name: DATABASE_URL
    value: postgresql+psycopg://tasks:$(POSTGRES_PASSWORD)@postgres:5432/tasks
```

Montare i segreti come **file** è leggermente più sicuro delle variabili d'ambiente, che possono finire nei log di
errore o essere lette da processi figli.

Le ConfigMap e i Secret che non cambiano si possono dichiarare `immutable: true`: proteggono da modifiche accidentali
e alleggeriscono il carico sull'API server.

## I segreti e Git

In un flusso GitOps tutto sta in Git, ma un Secret in Git è una password in chiaro (base64) nella cronologia per
sempre. Le soluzioni diffuse:

| Soluzione | Come funziona |
|---|---|
| **External Secrets Operator** | in Git metti un oggetto `ExternalSecret` che dice *quale* segreto leggere; l'operatore lo prende da AWS Secrets Manager, Google Secret Manager, Vault e crea il Secret nel cluster |
| **Secrets Store CSI Driver** | monta i segreti del gestore esterno direttamente come file nel Pod |
| **Sealed Secrets** | cifri il Secret con la chiave pubblica del controller del cluster; in Git va solo la versione cifrata, che solo il cluster può decifrare |
| **SOPS** | cifra i valori dei file YAML con KMS o age; Argo CD e Flux lo supportano con plugin |

Nel laboratorio, per semplicità, l'overlay `dev` genera il Secret da un file locale escluso da Git, e l'overlay
`gitops` si aspetta un Secret creato a mano. Nel cloud la scelta naturale è External Secrets con il gestore del
provider, autorizzato con **EKS Pod Identity** o **Workload Identity Federation for GKE** (lezione 8.15): nessuna
chiave da custodire.

## Esercizi

!!! esercizio "Esercizio 1 · Configurazione da file"
    Crea una ConfigMap `saluto` con la chiave `index.html` che contiene `<h1>Ciao da Kubernetes</h1>` e un Pod `nginx`
    che la monta in `/usr/share/nginx/html`. Verifica con `kubectl port-forward`. Poi modifica la ConfigMap: la
    pagina cambia senza riavviare il Pod?

??? success "Soluzione"
    ```bash
    echo '<h1>Ciao da Kubernetes</h1>' > index.html
    kubectl create configmap saluto --from-file=index.html
    kubectl run nginx --image=nginx:1.30 --dry-run=client -o yaml > nginx.yaml
    ```
    ```yaml
    spec:
      containers:
        - name: nginx
          image: nginx:1.30
          volumeMounts:
            - name: html
              mountPath: /usr/share/nginx/html
      volumes:
        - name: html
          configMap:
            name: saluto
    ```
    ```bash
    kubectl apply -f nginx.yaml
    kubectl port-forward pod/nginx 8080:80 &
    curl localhost:8080
    kubectl create configmap saluto --from-literal=index.html='<h1>Aggiornata</h1>' \
      --dry-run=client -o yaml | kubectl apply -f -
    sleep 60; curl localhost:8080
    ```
    Sì: il volume si aggiorna entro circa un minuto e Nginx serve i file dal disco a ogni richiesta. Nota il trucco
    `--dry-run=client -o yaml | kubectl apply -f -` per aggiornare un oggetto creato in modo imperativo.

!!! esercizio "Esercizio 2 · Leggi un segreto"
    Con un solo comando, stampa in chiaro la chiave `password` del Secret `tasks-db` nel namespace `tasks`.

??? success "Soluzione"
    ```bash
    kubectl get secret tasks-db -n tasks -o jsonpath='{.data.password}' | base64 -d
    ```
    Con l'overlay `dev` il nome reale ha un suffisso generato da Kustomize: trovalo con
    `kubectl get secrets -n tasks`. Proprio perché è così facile, chi può leggere i Secret va limitato con RBAC.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Come sono memorizzati i valori nel campo data di un Secret?","options":["Cifrati con AES","Codificati in base64","In chiaro","Con hash SHA-256"],"answer":[1],"explanation":"La cifratura a riposo in etcd va configurata a parte."},
 {"q":"Modifichi una ConfigMap usata come variabili d'ambiente. Cosa succede ai Pod in esecuzione?","options":["Ricevono subito i nuovi valori","Mantengono i vecchi valori fino al riavvio","Vengono riavviati automaticamente","Vanno in errore"],"answer":[1],"explanation":"I volumi invece si aggiornano (tranne con subPath)."},
 {"q":"Quali soluzioni evitano di salvare segreti in chiaro in Git?","options":["External Secrets Operator","Sealed Secrets","SOPS","Codificare in base64","Rinominare il file"],"answer":[0,1,2],"explanation":"base64 non protegge nulla."},
 {"q":"Perché il configMapGenerator di Kustomize aggiunge un suffisso al nome?","options":["Per evitare conflitti con altri namespace","Perché ogni modifica cambi il nome e provochi un rollout dei Pod","Per cifrare i valori","È un requisito di Kubernetes"],"answer":[1],"explanation":"Il Deployment si aggiorna perché cambia il riferimento nel modello del Pod."},
 {"q":"Quale opzione rende una ConfigMap non modificabile?","options":["readOnly: true","immutable: true","locked: true","final: true"],"answer":[1],"explanation":"Per cambiarla bisogna eliminarla e ricrearla."}
]}
</script>

## Fonti

- [Kubernetes: ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) e [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
- [Kubernetes: Secrets](https://kubernetes.io/docs/concepts/configuration/secret/) e [Good practices for Kubernetes Secrets](https://kubernetes.io/docs/concepts/security/secrets-good-practices/)
- [Kubernetes: Encrypting Confidential Data at Rest](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/)
- [Kubernetes: Define Dependent Environment Variables](https://kubernetes.io/docs/tasks/inject-data-application/define-interdependent-environment-variables/)
- [External Secrets Operator](https://external-secrets.io/latest/) · [Secrets Store CSI Driver](https://secrets-store-csi-driver.sigs.k8s.io/)
