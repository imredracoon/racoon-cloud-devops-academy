# Certified Kubernetes Administrator (CKA)

L'esame di riferimento per chi amministra cluster Kubernetes. È **interamente pratico**: lavori da riga di comando su
cluster reali e risolvi task a tempo. Per questo è molto apprezzato: non si supera a memoria.

| | |
|---|---|
| **Ente** | Linux Foundation e Cloud Native Computing Foundation (CNCF) |
| **Versione di Kubernetes** | v1.37 |
| **Formato** | task pratici da riga di comando, online con sorvegliante |
| **Durata** | 2 ore |
| **Costo** | 445 USD (esame) |
| **Incluso** | 12 mesi per sostenere l'esame, un nuovo tentativo, due sessioni del simulatore Killer.sh (36 ore ciascuna) |
| **Validità** | 2 anni |
| **Durante l'esame** | puoi consultare la documentazione ufficiale di Kubernetes nel browser dell'ambiente d'esame |
| **Pagina ufficiale** | [training.linuxfoundation.org/certification/certified-kubernetes-administrator-cka](https://training.linuxfoundation.org/certification/certified-kubernetes-administrator-cka/) |

Dati verificati il 28 ottobre 2026. La Linux Foundation propone spesso sconti stagionali: controlla prima di
acquistare. Leggi anche le [istruzioni importanti per CKA e CKAD](https://docs.linuxfoundation.org/tc-docs/certification/tips-cka-and-ckad).

## Quando sostenerla

Dopo la [Fase 8](../../fase-8/index.md), compresa la lezione 8.12 su kubeadm: backup e ripristino di etcd,
aggiornamento del cluster e risoluzione dei guasti dei nodi sono argomenti d'esame che kind da solo non copre.

## Piano di studio (6 settimane)

| Settimana | Attività |
|---|---|
| 1 | Velocità con `kubectl`: alias, completamento, `--dry-run=client -o yaml`, `kubectl explain`. Ripasso 8.2-8.4 |
| 2 | Service, DNS, Gateway API, Ingress, Network Policy (8.5, 8.9, 8.10) |
| 3 | Storage, configurazione, scheduling, autoscaling (8.6-8.8) |
| 4 | RBAC, kubeadm, etcd, aggiornamenti, Helm, Kustomize, CRD (8.11-8.13) |
| 5 | Troubleshooting su cluster kubeadm: nodi, componenti, applicazioni, rete (8.14) |
| 6 | Esami pratici a tempo con questa simulazione, poi le due sessioni Killer.sh incluse nell'acquisto |

!!! tip "Usa le sessioni Killer.sh alla fine"
    Il simulatore incluso è volutamente più difficile dell'esame. Usalo nelle ultime due settimane, quando hai già
    superato le simulazioni del corso: ti serve per abituarti all'ambiente e al ritmo.

## In questa appendice

- [Argomenti specifici](argomenti.md)
- [Simulazione pratica](simulazione.md): 39 task originali con verifica e soluzione
