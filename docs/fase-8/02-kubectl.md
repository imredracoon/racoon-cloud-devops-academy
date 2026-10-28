# 8.2 kind e kubectl

!!! obiettivi "Obiettivi"
    - Creare, usare ed eliminare un cluster locale con kind.
    - Capire il file kubeconfig: cluster, utenti e contesti.
    - Usare i comandi fondamentali di `kubectl` e i formati di output.
    - Generare YAML corretto in pochi secondi, come all'esame CKA.

## kind: Kubernetes in Docker

**kind** (*Kubernetes IN Docker*) è un progetto ufficiale di Kubernetes (SIG Testing) che crea un cluster usando
**container Docker come nodi**. Ogni "nodo" è un container con dentro kubelet, containerd e i componenti del control
plane. È lo strumento con cui il progetto Kubernetes stesso collauda il proprio codice: leggero, gratuito, si crea e
si distrugge in un minuto.

Nel Dev Container del corso kind è già installato. Il file del laboratorio descrive un cluster con un control plane e
due worker:

```yaml title="labs/k8s/kind/cluster.yaml"
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: academy
nodes:
  - role: control-plane
  - role: worker
    labels:
      zona: a
  - role: worker
    labels:
      zona: b
```

```bash
cd labs/k8s
kind create cluster --config kind/cluster.yaml   # circa un minuto
kind get clusters                                 # academy
docker ps --format '{{.Names}}'                   # i tre "nodi" sono container
kubectl get nodes -o wide
```

```text
NAME                    STATUS   ROLES           AGE   VERSION   INTERNAL-IP   ...
academy-control-plane   Ready    control-plane   60s   v1.37.0   172.18.0.3
academy-worker          Ready    <none>          40s   v1.37.0   172.18.0.2
academy-worker2         Ready    <none>          40s   v1.37.0   172.18.0.4
```

Per eliminarlo: `kind delete cluster --name academy`. Le immagini costruite in locale non sono visibili ai nodi:
si caricano con `kind load docker-image NOME:TAG --name academy`.

!!! note "Alternative locali"
    **minikube** (anche lui progetto ufficiale), **k3d** (k3s in Docker) e Docker Desktop offrono cluster locali
    simili. Usiamo kind perché è leggero, multi-nodo e funziona identico nel Dev Container e in Codespaces.

## Il kubeconfig

`kubectl` legge le informazioni di connessione dal file `~/.kube/config` (o dal file indicato nella variabile
`KUBECONFIG`). kind lo aggiorna da solo. Il file ha tre sezioni:

```yaml
clusters:   # DOVE: indirizzo dell'API server e certificato della CA
  - name: kind-academy
    cluster: { server: https://127.0.0.1:41234, certificate-authority-data: ... }
users:      # CHI: credenziali (certificato client, token, plugin come aws o gke-gcloud-auth-plugin)
  - name: kind-academy
    user: { client-certificate-data: ..., client-key-data: ... }
contexts:   # COMBINAZIONE di cluster + utente + namespace predefinito
  - name: kind-academy
    context: { cluster: kind-academy, user: kind-academy }
current-context: kind-academy
```

```bash
kubectl config get-contexts                     # elenco, con * su quello attivo
kubectl config use-context kind-academy         # cambia cluster
kubectl config set-context --current --namespace=tasks   # namespace predefinito
kubectl config view --minify                    # solo il contesto attivo
```

!!! danger "Il contesto sbagliato è un classico incidente"
    Con più cluster configurati (sviluppo, produzione) è facile lanciare un comando su quello sbagliato. Controlla
    sempre il contesto prima di comandi distruttivi. All'esame CKA **ogni domanda** inizia con un comando
    `kubectl config use-context ...` (o `ssh` sul nodo indicato): eseguilo sempre.

## Il comando kubectl

La forma è sempre `kubectl VERBO TIPO [NOME] [opzioni]`.

| Verbo | Esempio | Cosa fa |
|---|---|---|
| `get` | `kubectl get pods -A` | elenca oggetti (`-A` tutti i namespace) |
| `describe` | `kubectl describe pod web` | dettagli ed **eventi**: il primo posto dove guardare quando qualcosa non va |
| `apply` | `kubectl apply -f file.yaml` | crea o aggiorna in modo dichiarativo |
| `create` | `kubectl create deployment web --image=nginx` | crea in modo imperativo |
| `delete` | `kubectl delete -f file.yaml` | elimina |
| `edit` | `kubectl edit deployment web` | apre l'oggetto nell'editor e salva le modifiche |
| `logs` | `kubectl logs web -f` | log del container |
| `exec` | `kubectl exec -it web -- sh` | comando dentro il container |
| `explain` | `kubectl explain pod.spec.containers.resources` | documentazione dei campi, **disponibile all'esame** |
| `api-resources` | `kubectl api-resources` | tutti i tipi, con nome breve e gruppo |

Nomi brevi utili: `po` (pods), `deploy`, `svc`, `ns`, `cm`, `pv`, `pvc`, `sa`, `no` (nodes), `netpol`.

### I formati di output

```bash
kubectl get pods -o wide                     # più colonne: IP e nodo
kubectl get deploy web -o yaml               # l'oggetto completo
kubectl get pods -l app=web                  # filtro per etichetta
kubectl get pods --sort-by=.metadata.creationTimestamp
kubectl get pods -o jsonpath='{.items[*].metadata.name}'
kubectl get nodes -o custom-columns=NOME:.metadata.name,CPU:.status.capacity.cpu
kubectl get events -A --sort-by=.lastTimestamp
```

**JSONPath** permette di estrarre un campo preciso: all'esame capita di dover scrivere in un file, per esempio, i
nomi di tutti i nodi o l'immagine di un Pod.

### Generare YAML invece di scriverlo

Scrivere YAML a mano è lento e soggetto a errori di indentazione. Il metodo dei professionisti (e l'unico sensato
all'esame) è farlo generare a `kubectl` e poi modificarlo:

```bash
kubectl run web --image=nginx:1.30 --port=80 --dry-run=client -o yaml > pod.yaml
kubectl create deployment api --image=nginx:1.30 --replicas=3 --dry-run=client -o yaml > deploy.yaml
kubectl create service clusterip api --tcp=80:8080 --dry-run=client -o yaml
kubectl create configmap conf --from-literal=LIVELLO=debug --dry-run=client -o yaml
kubectl create job backup --image=busybox:1.37 --dry-run=client -o yaml -- echo ciao
```

- `--dry-run=client` non invia nulla al cluster; `-o yaml` stampa l'oggetto.
- `--dry-run=server` invece lo invia all'API server, che lo **valida** senza salvarlo: utile per controllare un
  manifest prima di applicarlo.
- `kubectl diff -f file.yaml` mostra cosa cambierebbe nel cluster.

### Prepara il terminale

```bash
alias k=kubectl
source <(kubectl completion bash)
complete -o default -F __start_kubectl k
export do="--dry-run=client -o yaml"      # k run web --image=nginx $do
```

L'alias `k` e il completamento sono già configurati nell'ambiente d'esame CKA. Nel Dev Container li trovi già
attivi.

## Laboratorio

!!! laboratorio "Laboratorio · Primo contatto"
    1. Crea il cluster `academy` e verifica che i tre nodi siano `Ready`.
    2. Elenca i Pod di sistema: `kubectl get pods -n kube-system -o wide`. Riconosci etcd, API server, scheduler,
       controller manager, kube-proxy e CoreDNS? Su quale nodo gira ciascuno?
    3. Leggi i manifest dei componenti statici del control plane:
       `docker exec academy-control-plane ls /etc/kubernetes/manifests`.
    4. Crea un Deployment `prova` con immagine `nginx:1.30` e due repliche, **generando prima il YAML**.
    5. Scrivi in `/tmp/nodi.txt` i nomi dei nodi, uno per riga, con un solo comando.
    6. Elimina il Deployment.

??? success "Soluzione"
    ```bash
    kind create cluster --config kind/cluster.yaml
    kubectl get nodes
    kubectl get pods -n kube-system -o wide
    docker exec academy-control-plane ls /etc/kubernetes/manifests
    # etcd.yaml  kube-apiserver.yaml  kube-controller-manager.yaml  kube-scheduler.yaml

    kubectl create deployment prova --image=nginx:1.30 --replicas=2 --dry-run=client -o yaml > prova.yaml
    kubectl apply -f prova.yaml
    kubectl get pods -l app=prova -o wide

    kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' > /tmp/nodi.txt
    kubectl delete -f prova.yaml
    ```

    I componenti del control plane girano sul nodo `academy-control-plane` come **static pod**: il kubelet li avvia
    direttamente dai file in `/etc/kubernetes/manifests`, senza passare dallo scheduler. kube-proxy gira su tutti i
    nodi (è un DaemonSet), CoreDNS su due Pod qualsiasi.

!!! esercizio "Esercizio · Trova il campo"
    Senza aprire il browser, scopri come si chiama il campo che imposta il tempo concesso a un Pod per chiudersi in
    modo ordinato prima di essere terminato, e qual è il valore predefinito.

??? success "Soluzione"
    ```bash
    kubectl explain pod.spec | grep -i -A3 grace
    kubectl explain pod.spec.terminationGracePeriodSeconds
    ```
    Il campo è `terminationGracePeriodSeconds`, predefinito **30 secondi**: alla chiusura il container riceve
    `SIGTERM`, e se dopo 30 secondi è ancora vivo riceve `SIGKILL`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa sono i nodi di un cluster kind?","options":["Macchine virtuali","Container Docker","Processi sull'host","Istanze cloud"],"answer":[1],"explanation":"kind = Kubernetes IN Docker."},
 {"q":"Quale sezione del kubeconfig combina cluster, utente e namespace?","options":["clusters","users","contexts","preferences"],"answer":[2],"explanation":"Si cambia con kubectl config use-context."},
 {"q":"Quale opzione genera il YAML di un oggetto senza crearlo nel cluster?","options":["--dry-run=client -o yaml","--output=none","--validate=false","--force"],"answer":[0],"explanation":"--dry-run=server invece lo fa validare all'API server."},
 {"q":"Quale comando mostra la documentazione del campo resources di un container?","options":["kubectl help resources","kubectl explain pod.spec.containers.resources","kubectl describe resources","kubectl docs pod"],"answer":[1],"explanation":"kubectl explain funziona anche all'esame, senza browser."},
 {"q":"Dove guardi per primo quando un Pod non parte?","options":["kubectl describe pod, sezione Events","Il codice sorgente","Il file kubeconfig","kubectl version"],"answer":[0],"explanation":"Gli eventi dicono spesso subito il motivo: immagine non trovata, risorse insufficienti, volume mancante."},
 {"q":"Come rendi disponibile ai nodi di kind un'immagine costruita in locale?","options":["docker push su localhost","kind load docker-image","kubectl cp","Non è possibile"],"answer":[1],"explanation":"Copia l'immagine nel containerd di ogni nodo."}
]}
</script>

## Fonti

- [kind: Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/) e [Configuration](https://kind.sigs.k8s.io/docs/user/configuration/)
- [Kubernetes: Organizing Cluster Access Using kubeconfig Files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)
- [Kubernetes: kubectl Quick Reference](https://kubernetes.io/docs/reference/kubectl/quick-reference/)
- [Kubernetes: JSONPath Support](https://kubernetes.io/docs/reference/kubectl/jsonpath/)
- [Kubernetes: Version Skew Policy](https://kubernetes.io/releases/version-skew-policy/)
