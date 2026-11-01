# Comandi di comodo dell'operator (lezione G.11). Kubebuilder genera un Makefile simile, più completo.
CONTROLLER_GEN = go run sigs.k8s.io/controller-tools/cmd/controller-gen@v0.23.0

.PHONY: generate manifests test run docker-build kind-deploy

generate:   ## rigenera zz_generated.deepcopy.go dopo aver cambiato i tipi
	$(CONTROLLER_GEN) object paths=./...

manifests:  ## rigenera CRD e ClusterRole dai marcatori +kubebuilder
	$(CONTROLLER_GEN) crd rbac:roleName=probe-operator paths=./... \
		output:crd:dir=config/crd output:rbac:dir=config/rbac

test: generate
	go test ./...

run: manifests generate  ## esegue l'operator fuori dal cluster, con il tuo kubeconfig
	kubectl apply -f config/crd/
	go run ./cmd

docker-build:
	docker build -t probe-operator:dev .

kind-deploy: manifests docker-build
	kind load docker-image probe-operator:dev --name academy
	kubectl apply -f config/crd/ -f config/rbac/ -f config/manager.yaml
	kubectl -n probe-operator rollout restart deployment/probe-operator
