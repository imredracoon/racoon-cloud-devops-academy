# 8.11 RBAC e service account

!!! obiettivi "Obiettivi"
    - Distinguere autenticazione e autorizzazione nell'API server.
    - Creare Role, ClusterRole e i relativi binding con il minimo privilegio.
    - Usare i service account per le applicazioni e le pipeline.
    - Creare un utente con un certificato client firmato dal cluster.
    - Applicare le Pod Security Standards a un namespace.

## Chi sei e cosa puoi fare

Ogni richiesta all'API server attraversa tre controlli:

1. **Autenticazione**: *chi sei?* Kubernetes non ha un database di utenti. Riconosce:
    - **certificati client** firmati dalla CA del cluster (il nome utente è il `CN`, i gruppi sono gli `O`);
    - **token** dei service account;
    - token **OIDC** di un identity provider esterno (Google, Microsoft Entra ID, Keycloak);
    - nel cloud, le identità IAM: su EKS un utente o ruolo IAM, su GKE un account Google.
2. **Autorizzazione**: *puoi farlo?* Quasi ovunque con **RBAC** (*Role-Based Access Control*).
3. **Ammissione**: *la richiesta rispetta le regole?* Gli admission controller possono modificarla o rifiutarla
   (Pod Security, quote, policy personalizzate con `ValidatingAdmissionPolicy`).

## I quattro oggetti di RBAC

| Oggetto | Ambito | Contiene |
|---|---|---|
| **Role** | un namespace | regole: verbi su risorse |
| **ClusterRole** | tutto il cluster | regole, anche su risorse senza namespace (nodi, PV) |
| **RoleBinding** | un namespace | assegna un Role **o un ClusterRole** a soggetti, **in quel namespace** |
| **ClusterRoleBinding** | tutto il cluster | assegna un ClusterRole a soggetti in **tutti** i namespace |

RBAC è solo **permissivo**: tutto è negato finché una regola non lo permette, e non esistono regole di negazione.

```yaml title="role.yaml"
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: lettore-pod
  namespace: tasks
rules:
  - apiGroups: [""]                  # "" = gruppo core (pods, services, configmaps, secrets...)
    resources: ["pods", "pods/log"]  # pods/log è una sottorisorsa
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: anna-lettore
  namespace: tasks
subjects:
  - kind: User
    name: anna
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: lettore-pod
  apiGroup: rbac.authorization.k8s.io
```

Lo stesso, in modo imperativo, e la verifica:

```bash
kubectl create role lettore-pod -n tasks --verb=get,list,watch --resource=pods,pods/log
kubectl create rolebinding anna-lettore -n tasks --role=lettore-pod --user=anna

kubectl auth can-i list pods -n tasks --as anna        # yes
kubectl auth can-i delete pods -n tasks --as anna      # no
kubectl auth can-i list pods -n default --as anna      # no
kubectl auth can-i --list -n tasks --as anna           # tutto ciò che può fare
```

`--as` (impersonificazione) permette all'amministratore di verificare i permessi di un altro: è il modo più rapido
per controllare il proprio lavoro, anche all'esame.

!!! tip "ClusterRole riusabili con RoleBinding"
    Un **RoleBinding** che punta a un **ClusterRole** concede quelle regole **solo nel namespace del binding**. Così
    si definisce una volta il ruolo "sviluppatore" e lo si assegna squadra per squadra. Kubernetes include ClusterRole
    pronti: `view`, `edit`, `admin` e `cluster-admin` (tutto, ovunque: da dare con estrema parsimonia).

I verbi più usati: `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection`. Attenzione ai
permessi "nascosti": chi può **creare Pod** in un namespace può montare qualunque Secret di quel namespace, e chi può
`list` i Secret ne legge il contenuto.

## Service account

Gli utenti sono persone; i **service account** sono identità per **processi**: Pod, pipeline CI/CD, controller. Ogni
namespace ha un service account `default`, usato dai Pod che non ne indicano un altro.

```bash
kubectl create serviceaccount deployer -n tasks
kubectl create clusterrole gestore-deploy --verb=get,list,update,patch --resource=deployments.apps
kubectl create rolebinding deployer-gestore -n tasks \
  --clusterrole=gestore-deploy --serviceaccount=tasks:deployer
kubectl auth can-i patch deployments -n tasks --as system:serviceaccount:tasks:deployer   # yes
kubectl create token deployer -n tasks --duration=1h     # token a scadenza
```

- Il nome completo di un service account come utente è `system:serviceaccount:NAMESPACE:NOME`.
- Il kubelet monta nel Pod un **token a scadenza**, rinnovato automaticamente, in
  `/var/run/secrets/kubernetes.io/serviceaccount/`. Se l'applicazione non parla con l'API di Kubernetes, come la
  Tasks API, disattiva il montaggio con `automountServiceAccountToken: false`: un token in meno da rubare.
- Nel cloud lo stesso service account può ricevere permessi **del cloud** senza chiavi: **EKS Pod Identity** (o IRSA)
  su AWS, **Workload Identity Federation for GKE** su Google Cloud. È il principio di OIDC della [Fase 7](../fase-7/index.md) applicato ai
  Pod.

## Un utente con certificato

Esercizio tipico dell'esame CKA: creare le credenziali per un nuovo utente.

```bash
# 1. chiave privata e richiesta di firma; CN = utente, O = gruppo
openssl genrsa -out anna.key 2048
openssl req -new -key anna.key -out anna.csr -subj "/CN=anna/O=sviluppo"

# 2. richiesta al cluster
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: anna
spec:
  request: $(base64 -w0 anna.csr)
  signerName: kubernetes.io/kube-apiserver-client
  expirationSeconds: 86400
  usages: ["client auth"]
EOF

# 3. approvazione e certificato
kubectl certificate approve anna
kubectl get csr anna -o jsonpath='{.status.certificate}' | base64 -d > anna.crt

# 4. kubeconfig
kubectl config set-credentials anna --client-key=anna.key --client-certificate=anna.crt --embed-certs=true
kubectl config set-context anna --cluster=kind-academy --user=anna --namespace=tasks
kubectl --context anna get pods
```

Nelle aziende si preferisce l'accesso tramite un identity provider (OIDC) o l'IAM del cloud: i certificati non si
possono revocare singolarmente.

## Pod Security Standards

Le **Pod Security Standards** definiscono tre profili di sicurezza per i Pod:

| Profilo | Cosa permette |
|---|---|
| `privileged` | tutto |
| `baseline` | blocca le escalation note: container privilegiati, `hostNetwork`, `hostPath`, capacità pericolose |
| `restricted` | in più richiede utente non root, `allowPrivilegeEscalation: false`, `capabilities.drop: ["ALL"]`, profilo seccomp |

L'admission controller **Pod Security Admission** li applica per namespace, tramite label, in tre modalità:
`enforce` (rifiuta), `audit` (registra nei log di audit), `warn` (avvisa l'utente).

```yaml title="labs/k8s/base/namespace.yaml"
metadata:
  name: tasks
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/warn: restricted
```

```bash
kubectl -n tasks run insicuro --image=nginx:1.30
# Error from server (Forbidden): pods "insicuro" is forbidden: violates PodSecurity "restricted:latest":
# allowPrivilegeEscalation != false, unrestricted capabilities, runAsNonRoot != true, seccompProfile...
```

Per regole più specifiche (immagini solo dal registro aziendale, label obbligatorie) si usano
`ValidatingAdmissionPolicy` (integrata, con espressioni CEL) o strumenti come Kyverno e OPA Gatekeeper.

## Esercizi

!!! esercizio "Esercizio 1 · La pipeline"
    La pipeline di rilascio deve poter aggiornare l'immagine dei Deployment e leggere i Pod nei namespace `tasks` e
    `shop`, e nient'altro. Crea il service account `ci` nel namespace `ci` con i permessi minimi, e verifica.

??? success "Soluzione"
    ```bash
    kubectl create namespace ci
    kubectl create serviceaccount ci -n ci
    kubectl create clusterrole rilascio --verb=get,list,patch --resource=deployments.apps \
      --dry-run=client -o yaml > rilascio.yaml
    ```
    Aggiungi una seconda regola per i Pod e applica:
    ```yaml
    rules:
      - apiGroups: ["apps"]
        resources: ["deployments"]
        verbs: ["get", "list", "patch"]
      - apiGroups: [""]
        resources: ["pods"]
        verbs: ["get", "list", "watch"]
    ```
    ```bash
    kubectl apply -f rilascio.yaml
    for ns in tasks shop; do
      kubectl create rolebinding ci-rilascio -n $ns --clusterrole=rilascio --serviceaccount=ci:ci
    done
    kubectl auth can-i patch deployments -n shop --as system:serviceaccount:ci:ci     # yes
    kubectl auth can-i delete deployments -n shop --as system:serviceaccount:ci:ci    # no
    kubectl auth can-i patch deployments -n default --as system:serviceaccount:ci:ci  # no
    ```
    Un **ClusterRole** riusato con due **RoleBinding**: niente ClusterRoleBinding, che darebbe accesso a tutti i
    namespace.

!!! esercizio "Esercizio 2 · Trova l'errore"
    Il gruppo `sviluppo` deve poter leggere i nodi. Un collega crea un ClusterRole corretto e un **RoleBinding** nel
    namespace `default` che lo assegna al gruppo. `kubectl auth can-i list nodes --as anna --as-group sviluppo`
    risponde `no`. Perché?

??? success "Soluzione"
    I nodi non appartengono a nessun namespace: un RoleBinding concede permessi solo **dentro** il proprio namespace
    e non può dare accesso a risorse del cluster. Serve un **ClusterRoleBinding**:
    `kubectl create clusterrolebinding sviluppo-nodi --clusterrole=lettore-nodi --group=sviluppo`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale oggetto assegna un ClusterRole a un utente limitatamente a un solo namespace?","options":["ClusterRoleBinding","RoleBinding","Role","ServiceAccount"],"answer":[1],"explanation":"È il modo per riusare ruoli comuni."},
 {"q":"Come verifichi se l'utente anna può eliminare i Deployment nel namespace tasks?","options":["kubectl auth can-i delete deployments -n tasks --as anna","kubectl get rolebindings","kubectl describe user anna","kubectl rbac check anna"],"answer":[0],"explanation":"Kubernetes non ha oggetti User."},
 {"q":"Qual è il nome utente del service account builder nel namespace ci?","options":["ci:builder","system:serviceaccount:ci:builder","builder@ci","serviceaccount/builder"],"answer":[1],"explanation":"Si usa anche con --as."},
 {"q":"Quali permessi danno, di fatto, accesso ai segreti di un namespace?","options":["list sui Secret","get sui Secret","create sui Pod","get sui Service","list sugli eventi"],"answer":[0,1,2],"explanation":"Un Pod può montare qualunque Secret del suo namespace."},
 {"q":"Quale profilo delle Pod Security Standards richiede utente non root e capacità rimosse?","options":["privileged","baseline","restricted","hardened"],"answer":[2],"explanation":"baseline blocca solo le escalation note."},
 {"q":"In un certificato client per Kubernetes, cosa rappresenta il campo O (organizzazione)?","options":["Il namespace","Il gruppo dell'utente","Il cluster","Il ruolo"],"answer":[1],"explanation":"CN è il nome utente."}
]}
</script>

## Fonti

- [Kubernetes: Controlling Access to the Kubernetes API](https://kubernetes.io/docs/concepts/security/controlling-access/)
- [Kubernetes: Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) e [RBAC Good Practices](https://kubernetes.io/docs/concepts/security/rbac-good-practices/)
- [Kubernetes: Service Accounts](https://kubernetes.io/docs/concepts/security/service-accounts/)
- [Kubernetes: Certificates and Certificate Signing Requests](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/)
- [Kubernetes: Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/) e [Pod Security Admission](https://kubernetes.io/docs/concepts/security/pod-security-admission/)
- [Kubernetes: Validating Admission Policy](https://kubernetes.io/docs/reference/access-authn-authz/validating-admission-policy/)
