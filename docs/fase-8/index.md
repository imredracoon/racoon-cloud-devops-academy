# Fase 8 · Kubernetes

!!! obiettivi "Al termine di questa fase saprai"
    - descrivere l'architettura di un cluster Kubernetes e il ruolo di ogni componente;
    - creare un cluster locale con **kind** e lavorare con `kubectl` in modo rapido e preciso;
    - eseguire applicazioni con Pod, Deployment, StatefulSet, DaemonSet e Job, e aggiornarle senza interruzioni;
    - esporre i servizi con Service, **Gateway API** e Ingress, e isolarli con le **Network Policy**;
    - gestire configurazione, segreti, dischi persistenti, risorse e autoscaling;
    - proteggere il cluster con **RBAC**, service account e Pod Security;
    - installare e aggiornare un cluster con **kubeadm**, fare backup di etcd e risolvere i guasti;
    - impacchettare le applicazioni con **Helm** e **Kustomize** e rilasciarle con **GitOps** (Argo CD);
    - scegliere tra Kubernetes gestito su AWS (**EKS**) e Google Cloud (**GKE**).

| | |
|---|---|
| **Durata indicativa** | 3 mesi, circa 120 ore |
| **Prerequisiti** | Fasi 1, 3 e 7 (Linux, container, CI/CD) |
| **Ambiente** | Dev Container (kind, kubectl, Helm già installati); per la [lezione 8.12](12-kubeadm.md) due o tre macchine virtuali |
| **Progetto** | [Cluster con GitOps](progetto.md) |
| **Certificazioni collegate** | [CKA](../appendici/cka/index.md) · la parte Kubernetes di [Google Associate Cloud Engineer](../appendici/gcp-ace/index.md) |

## Perché Kubernetes

Nella [Fase 3](../fase-3/index.md) hai eseguito la Tasks API con Docker Compose su **una** macchina. Funziona, ma cosa succede quando:

- la macchina si spegne? Nessuno riavvia i container altrove;
- il traffico triplica? Devi aggiungere macchine e distribuire i container a mano;
- rilasci una nuova versione? Devi fermare e ripartire, con un'interruzione;
- dieci squadre condividono le stesse macchine? Nessuno separa risorse e permessi.

**Kubernetes** è un *orchestratore di container*: tu descrivi lo **stato desiderato** ("tre copie della Tasks API,
raggiungibili su questo nome, con al massimo 256 MiB di memoria ciascuna") e il cluster lavora continuamente per
mantenerlo, su decine o migliaia di macchine. Nato in Google sull'esperienza del sistema interno Borg e donato alla
Cloud Native Computing Foundation (CNCF) nel 2015, oggi è lo standard di fatto: lo offrono tutti i cloud (EKS, GKE,
AKS) e lo chiedono gran parte degli annunci DevOps e platform engineering.

!!! warning "Kubernetes non serve sempre"
    Kubernetes è potente ma complesso. Per una sola applicazione piccola, un servizio gestito come **Amazon ECS su
    Fargate** o **Cloud Run** ([Fase 5](../fase-5/index.md)) costa meno fatica. Kubernetes conviene quando ci sono molti servizi, molte
    squadre, esigenze di portabilità tra cloud o un ecosistema di strumenti da standardizzare. Saperlo dire in un
    colloquio è segno di maturità.

## Le versioni usate

| Strumento | Versione |
|---|---|
| Kubernetes | 1.37 (rilasci ogni circa 4 mesi, ognuno supportato per circa 14 mesi) |
| kind | 0.33 |
| kubectl | 1.37 (al massimo una versione minore di distanza dal cluster) |
| Helm | 4.3 |
| Gateway API | 1.5, con Envoy Gateway 1.9 |
| Argo CD | 3.5 |

Tutti i manifest del laboratorio sono nella cartella [`labs/k8s`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/k8s)
e sono stati validati con un API server Kubernetes 1.37.

## Nei tre livelli

Le lezioni 8.1-8.11 e 8.13-8.15 si fanno con kind nel Dev Container in tutti i livelli: è il modo più rapido per
imparare. Nei livelli 2 e 3, dalla [lezione 8.12](12-kubeadm.md#nei-tre-livelli), costruisci un **cluster vero** e ci porti il progetto.

=== "Livello 1 · Un PC"

    kind per tutta la fase. Per la lezione 8.12, se la workstation ha almeno 16 GB di RAM, le tre VM `k8s-cp`,
    `k8s-w1` e `k8s-w2` di `academy-vm`; altrimenti il laboratorio con kind.

=== "Livello 2 · Workstation + server"

    Un cluster kubeadm sulle VM `k8s-cp`, `k8s-w1` e `k8s-w2` di **`srv1`** (servono almeno 12 GB di RAM sul
    server), con lo storage NFS di `srv1` e **MetalLB** per i servizi `LoadBalancer`. Il progetto gira lì.

=== "Livello 3 · Mini data center"

    Un cluster kubeadm **sulle macchine fisiche**: `nodo2` (control plane, Ubuntu) e `nodo3` (worker, AlmaLinux).
    Lo storage è l'NFS di `nodo1`, gli indirizzi dei servizi li assegna MetalLB sulla rete del laboratorio e `nodo1`
    li pubblica verso casa. Se `nodo1` ha 16 GB, puoi aggiungere le VM `k8s-w1` e `k8s-w2` come worker.

## Le lezioni

1. [Architettura di Kubernetes](01-architettura.md)
2. [kind e kubectl](02-kubectl.md)
3. [Pod, namespace e label](03-pod.md)
4. [Deployment e rollout](04-deployment.md)
5. [Service e DNS](05-service.md)
6. [ConfigMap e Secret](06-config.md)
7. [Storage persistente](07-storage.md)
8. [Scheduling e autoscaling](08-scheduling.md)
9. [Ingress e Gateway API](09-ingress.md)
10. [Network Policy](10-networkpolicy.md)
11. [RBAC e service account](11-rbac.md)
12. [Cluster con kubeadm](12-kubeadm.md)
13. [Helm, Kustomize e operatori](13-helm-kustomize.md)
14. [Troubleshooting](14-troubleshooting.md)
15. [GitOps e Kubernetes nel cloud](15-gitops-cloud.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

!!! tip "Come studiare questa fase"
    L'esame CKA è **pratico**: due ore davanti a un terminale. Qui vale più che altrove la regola del corso: scrivi
    ogni comando, rompi le cose di proposito, ripara. Alla fine di ogni lezione chiediti "saprei farlo senza
    guardare, in tre minuti?". Se la risposta è no, ripeti il laboratorio.
