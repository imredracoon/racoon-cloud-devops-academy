# Argomenti specifici · CKA

I cinque domini del programma ufficiale ([curriculum CNCF](https://github.com/cncf/curriculum)) con le competenze
richieste e le lezioni della Fase 8.

## Storage (10%)

| Competenza | Lezioni |
|---|---|
| StorageClass e provisioning dinamico | [8.7](../../fase-8/07-storage.md) |
| Tipi di volume, modalità di accesso e politiche di recupero | [8.7](../../fase-8/07-storage.md) |
| PersistentVolume e PersistentVolumeClaim | [8.7](../../fase-8/07-storage.md) |

## Troubleshooting (30%)

| Competenza | Lezioni |
|---|---|
| Guasti di cluster e nodi (kubelet, runtime, certificati) | [8.12](../../fase-8/12-kubeadm.md), [8.14](../../fase-8/14-troubleshooting.md) |
| Guasti dei componenti del cluster (Pod statici, etcd, scheduler) | [8.1](../../fase-8/01-architettura.md), [8.14](../../fase-8/14-troubleshooting.md) |
| Monitorare l'uso delle risorse dei cluster e delle applicazioni (`kubectl top`) | [8.8](../../fase-8/08-scheduling.md), [8.14](../../fase-8/14-troubleshooting.md) |
| Gestire e valutare l'output dei container (log, eventi) | [8.3](../../fase-8/03-pod.md), [8.14](../../fase-8/14-troubleshooting.md) |
| Guasti di Service e rete (endpoint, DNS, policy) | [8.5](../../fase-8/05-service.md), [8.10](../../fase-8/10-networkpolicy.md), [8.14](../../fase-8/14-troubleshooting.md) |

## Workloads & Scheduling (15%)

| Competenza | Lezioni |
|---|---|
| Deployment, aggiornamenti progressivi e rollback | [8.4](../../fase-8/04-deployment.md) |
| ConfigMap e Secret | [8.6](../../fase-8/06-config.md) |
| Autoscaling dei workload | [8.8](../../fase-8/08-scheduling.md) |
| Primitive di autoriparazione (probe, controller) | [8.3](../../fase-8/03-pod.md), [8.4](../../fase-8/04-deployment.md) |
| Ammissione e scheduling dei Pod: risorse, limiti, affinità, taint e toleration | [8.8](../../fase-8/08-scheduling.md) |

## Cluster Architecture, Installation & Configuration (25%)

| Competenza | Lezioni |
|---|---|
| RBAC | [8.11](../../fase-8/11-rbac.md) |
| Preparare l'infrastruttura e creare cluster con kubeadm | [8.12](../../fase-8/12-kubeadm.md) |
| Gestire il ciclo di vita del cluster (aggiornamenti, backup e ripristino di etcd) | [8.12](../../fase-8/12-kubeadm.md) |
| Installare componenti con Helm e Kustomize | [8.13](../../fase-8/13-helm-kustomize.md) |
| Interfacce di estensione (CNI, CSI, CRI) | [8.1](../../fase-8/01-architettura.md), [8.12](../../fase-8/12-kubeadm.md) |
| CRD e operatori | [8.13](../../fase-8/13-helm-kustomize.md) |

## Services & Networking (20%)

| Competenza | Lezioni |
|---|---|
| Connettività tra i Pod | [8.5](../../fase-8/05-service.md) |
| Network Policy | [8.10](../../fase-8/10-networkpolicy.md) |
| Tipi di Service ed endpoint | [8.5](../../fase-8/05-service.md) |
| Gateway API per il traffico in ingresso | [8.9](../../fase-8/09-ingress.md) |
| Ingress controller e risorse Ingress | [8.9](../../fase-8/09-ingress.md) |
| CoreDNS | [8.5](../../fase-8/05-service.md), [8.14](../../fase-8/14-troubleshooting.md) |

## Dove sbagliano più spesso i candidati

- **Contesto sbagliato**: ogni task indica il cluster da usare. Esegui sempre il comando di cambio contesto
  riportato nel task.
- **Namespace dimenticato**: crea gli oggetti nel namespace richiesto, non in `default`.
- **Modifiche a un Pod esistente**: molti campi non si modificano. Esporta il manifest, correggilo, ricrea il Pod.
- **Etcd**: i percorsi dei certificati si leggono in `/etc/kubernetes/manifests/etcd.yaml`, non si ricordano.
- **Verificare**: un task non verificato è un task che potresti aver sbagliato.
