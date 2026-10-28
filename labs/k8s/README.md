# Laboratorio Kubernetes (Fase 8)

La Tasks API delle fasi precedenti, pronta per Kubernetes. Le lezioni della
[Fase 8](../../docs/fase-8/index.md) spiegano ogni file.

| Cartella | Contenuto | Lezione |
|---|---|---|
| `kind/` | cluster locale con un control plane e due worker | 8.2 |
| `base/` | namespace, PostgreSQL (StatefulSet), API (Deployment, Service, PDB), HPA, NetworkPolicy | 8.3-8.10 |
| `rotta/` | HTTPRoute che pubblica la API sul Gateway | 8.9 |
| `overlays/dev`, `overlays/prod`, `overlays/gitops` | varianti Kustomize per ambiente | 8.13 |
| `gateway/` | GatewayClass e Gateway di Envoy Gateway | 8.9 |
| `chart/tasks-api` | la stessa applicazione come chart Helm | 8.13 |
| `argocd/` | Application di Argo CD per il flusso GitOps | 8.15 |

## Avvio rapido

```bash
# 1. Cluster
kind create cluster --config kind/cluster.yaml

# 2. Immagine dell'API (costruita nella Fase 3) caricata nei nodi di kind
docker build -t tasks-api:dev ../tasks-api
kind load docker-image tasks-api:dev --name academy

# 3. Applicazione, ambiente dev
echo -n "una-password-lunga-e-casuale" > overlays/dev/db-password.txt
kubectl apply -k overlays/dev
kubectl -n tasks get pods -w

# 4. Envoy Gateway, che installa anche le CRD della Gateway API (lezione 8.9)
helm install eg oci://docker.io/envoyproxy/gateway-helm --version v1.9.3 \
  -n envoy-gateway-system --create-namespace
kubectl wait --timeout=5m -n envoy-gateway-system deployment/envoy-gateway --for=condition=Available
kubectl apply -f gateway/gateway.yaml

kubectl apply -k rotta     # pubblica la API: curl -H "Host: tasks.localhost" ...
```

Il file `db-password.txt` è escluso da Git (`.gitignore`): i segreti non si versionano.

## Versioni di riferimento

Kubernetes 1.37 · kind 0.33 · Gateway API 1.5 · Envoy Gateway 1.9 · Helm 4.3 · Argo CD 3.5.
