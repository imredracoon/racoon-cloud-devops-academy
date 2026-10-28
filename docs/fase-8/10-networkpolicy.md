# 8.10 Network Policy

!!! obiettivi "Obiettivi"
    - Spiegare il comportamento predefinito della rete dei Pod e perché è un rischio.
    - Scrivere Network Policy per il traffico in ingresso e in uscita.
    - Combinare selettori di Pod, di namespace e blocchi di indirizzi.
    - Applicare il principio "nega tutto, poi apri" alla Tasks API.

## Tutti parlano con tutti

Per impostazione predefinita in Kubernetes **ogni Pod può raggiungere ogni altro Pod**, in qualunque namespace. Se un
attaccante compromette un'applicazione esposta, da lì può contattare il database di un'altra squadra. È l'equivalente
di una rete aziendale senza firewall interni.

Una **NetworkPolicy** è un firewall dichiarativo per i Pod: seleziona alcuni Pod e dice quale traffico possono
**ricevere** (`ingress`) e **inviare** (`egress`).

!!! warning "Serve un CNI che le applichi"
    La risorsa NetworkPolicy è parte dell'API, ma la applica il **plugin di rete**. Calico, Cilium, i CNI di EKS (con
    l'opzione di Network Policy attivata) e GKE Dataplane V2 lo fanno; Flannel da solo no: le policy verrebbero
    accettate e **ignorate** in silenzio. Il CNI di kind, kindnet, le applica nelle versioni recenti.

## Le regole del gioco

1. Un Pod **non selezionato** da alcuna policy accetta tutto.
2. Appena una policy seleziona un Pod per una direzione (`policyTypes: ["Ingress"]`), in quella direzione il Pod
   accetta **solo** ciò che qualche policy permette esplicitamente.
3. Le policy sono **additive**: non esistono regole di negazione; il traffico è permesso se almeno una policy lo
   permette.
4. Per una connessione servono entrambe le parti: l'**egress** del mittente e l'**ingress** del destinatario. Le
   risposte di una connessione permessa passano sempre (le policy sono *stateful*).

## Anatomia

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: postgres-solo-da-api
  namespace: tasks
spec:
  podSelector:                 # a quali Pod si applica (nel namespace della policy)
    matchLabels:
      app.kubernetes.io/name: postgres
  policyTypes: ["Ingress"]
  ingress:
    - from:
        - podSelector:         # quali Pod possono connettersi (stesso namespace)
            matchLabels:
              app.kubernetes.io/name: tasks-api
      ports:
        - protocol: TCP
          port: 5432
```

Le sorgenti (`from`) e le destinazioni (`to`) possono essere:

| Selettore | Seleziona |
|---|---|
| `podSelector` | Pod dello **stesso** namespace della policy |
| `namespaceSelector` | **tutti** i Pod dei namespace con certe label |
| `namespaceSelector` + `podSelector` nello **stesso** elemento | i Pod con certe label **nei** namespace con certe label (AND) |
| `ipBlock` | intervalli CIDR, tipicamente indirizzi fuori dal cluster |

!!! danger "Il trattino che cambia tutto"
    ```yaml
    # UN elemento: Pod "monitor" NEI namespace "osservabilita" (AND)
    - from:
        - namespaceSelector:
            matchLabels: { kubernetes.io/metadata.name: osservabilita }
          podSelector:
            matchLabels: { app: monitor }

    # DUE elementi: TUTTI i Pod di "osservabilita" OPPURE i Pod "monitor" del namespace locale (OR)
    - from:
        - namespaceSelector:
            matchLabels: { kubernetes.io/metadata.name: osservabilita }
        - podSelector:
            matchLabels: { app: monitor }
    ```
    È l'errore più comune con le Network Policy, e un classico all'esame. La label
    `kubernetes.io/metadata.name` è presente automaticamente su ogni namespace e vale il suo nome.

## Le policy della Tasks API

Il laboratorio applica il modello **nega tutto, poi apri**:

```yaml title="labs/k8s/base/networkpolicy.yaml"
# 1. Nessun Pod del namespace accetta traffico...
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: nega-tutto-in-ingresso
spec:
  podSelector: {}              # {} = tutti i Pod del namespace
  policyTypes: ["Ingress"]
---
# 2. ...tranne PostgreSQL dalla API...
# (postgres-solo-da-api, vista sopra)
---
# 3. ...e la API dai proxy del Gateway.
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-da-gateway
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: tasks-api
  policyTypes: ["Ingress"]
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: envoy-gateway-system
      ports:
        - protocol: TCP
          port: 8000           # la porta del CONTAINER, non quella del Service
```

Nota la porta: le policy agiscono sul traffico verso il Pod, quindi si indica la porta del container (8000), non
quella del Service (80).

### E l'uscita?

Limitare anche l'**egress** è più sicuro (un'applicazione compromessa non può scaricare strumenti né inviare dati
fuori), ma ricorda di permettere il **DNS**, altrimenti nessun nome si risolve più:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-uscite
  namespace: tasks
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: tasks-api
  policyTypes: ["Egress"]
  egress:
    - to:                      # il DNS del cluster
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: kube-system
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - { protocol: UDP, port: 53 }
        - { protocol: TCP, port: 53 }
    - to:                      # il database
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: postgres
      ports:
        - { protocol: TCP, port: 5432 }
```

Per esigenze oltre le NetworkPolicy (regole a livello di cluster stabilite dagli amministratori, regole per nome DNS
o per percorso HTTP) esistono le **AdminNetworkPolicy** del SIG Network e le policy proprie di Cilium e Calico.

## Laboratorio

!!! laboratorio "Laboratorio · Firewall per la Tasks API"
    Con la Tasks API e il Gateway delle lezioni precedenti (e le policy dell'overlay `dev` attive):

    1. Da un Pod di prova nel namespace `default` prova a raggiungere la API e il database:
       ```bash
       kubectl run prova --rm -it --image=busybox:1.37 --restart=Never -- sh
       wget -qO- -T 3 http://tasks-api.tasks/tasks     # dentro il Pod
       nc -zv -w 3 postgres.tasks 5432
       ```
    2. Ripeti la richiesta passando dal Gateway (tunnel della [lezione 8.9](09-ingress.md)). Funziona? Perché?
    3. Applica la policy `api-uscite` qui sopra e verifica che la API funzioni ancora. Poi togli temporaneamente
       la regola del DNS: cosa succede alle richieste `/tasks`?

??? success "Soluzione"
    1. Entrambe le connessioni vanno in **timeout**: il Pod `prova` non è né nel namespace del Gateway né una API.
       Le policy non rispondono con un rifiuto, scartano i pacchetti.
    2. Sì: il traffico arriva alla API dal proxy Envoy, che gira in `envoy-gateway-system`, namespace permesso da
       `api-da-gateway`.
    3. Senza la regola del DNS la API non risolve più `postgres`: la readiness fallisce, i Pod escono dagli endpoint
       e il Gateway risponde con un errore 503. Rimetti la regola. Lezione: con l'egress, il DNS va sempre
       permesso per primo.

!!! esercizio "Esercizio · Prometheus deve leggere le metriche"
    Il Prometheus della [Fase 7](../fase-7/index.md) gira ora nel namespace `osservabilita` con la label `app.kubernetes.io/name: prometheus`
    e deve leggere `/metrics` della Tasks API sulla porta 8000. Scrivi la policy minima, senza aprire la API a tutti
    i Pod di quel namespace.

??? success "Soluzione"
    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: api-da-prometheus
      namespace: tasks
    spec:
      podSelector:
        matchLabels:
          app.kubernetes.io/name: tasks-api
      policyTypes: ["Ingress"]
      ingress:
        - from:
            - namespaceSelector:
                matchLabels:
                  kubernetes.io/metadata.name: osservabilita
              podSelector:            # stesso elemento: AND
                matchLabels:
                  app.kubernetes.io/name: prometheus
          ports:
            - protocol: TCP
              port: 8000
    ```
    Le policy sono additive: questa si aggiunge a `api-da-gateway` senza modificarla.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"In un cluster senza Network Policy, quale traffico tra Pod è permesso?","options":["Nessuno","Solo nello stesso namespace","Tutto","Solo verso i Service"],"answer":[2],"explanation":"Per questo si parte da una policy che nega tutto."},
 {"q":"Cosa fa una policy con podSelector: {} e policyTypes: [Ingress] senza regole ingress?","options":["Permette tutto","Blocca tutto il traffico in ingresso verso tutti i Pod del namespace","Blocca solo il traffico in uscita","Non ha effetto"],"answer":[1],"explanation":"È la classica default deny."},
 {"q":"Come si scrive 'i Pod app=monitor nei namespace con label team=ops'?","options":["Un solo elemento from con namespaceSelector e podSelector insieme","Due elementi from separati","Due policy diverse","Con ipBlock"],"answer":[0],"explanation":"Due elementi separati significano OR."},
 {"q":"Una policy di egress blocca tutto tranne il database. Cosa smette di funzionare per primo?","options":["La risoluzione DNS dei nomi","Le sonde del kubelet","I log","Il volume persistente"],"answer":[0],"explanation":"Il DNS va sempre permesso."},
 {"q":"Le Network Policy non hanno effetto nel tuo cluster. Quale può essere la causa?","options":["Il plugin di rete non le supporta","Hai scritto troppe policy","Il namespace è default","Manca metrics-server"],"answer":[0],"explanation":"Per esempio Flannel da solo."},
 {"q":"Quale porta indichi in una policy di ingress per una API con Service sulla porta 80 e container sulla 8000?","options":["80","8000","30080","Nessuna"],"answer":[1],"explanation":"Le policy agiscono sul traffico verso il Pod."}
]}
</script>

## Fonti

- [Kubernetes: Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Kubernetes: Declare Network Policy](https://kubernetes.io/docs/tasks/administer-cluster/declare-network-policy/)
- [Kubernetes: Well-Known Labels (kubernetes.io/metadata.name)](https://kubernetes.io/docs/reference/labels-annotations-taints/)
- [kindnet: Network Policies](https://kindnet.sigs.k8s.io/docs/user/network-policies/)
- [Network Policy API (AdminNetworkPolicy)](https://network-policy-api.sigs.k8s.io/)
