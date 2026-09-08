# Workshop "Leeft je app nog?" - probes en health checks in Kubernetes.
# Draai alles vanuit de root van de repo. `make help` toont het overzicht.

CLUSTER ?= probe-workshop
IMAGE   ?= probe-demo:dev
APP     ?= probe-demo
DOCKER  ?= docker
URL     ?= http://127.0.0.1:30080

GOED  := demo/k8s/deployment.yaml
KAPOT := demo/k8s/kapot

.DEFAULT_GOAL := help

.PHONY: help alles cluster build load deploy kapot-1 kapot-2 kapot-3 kapot-4 \
        herstel logs watch endpoints belasting unready dependency \
        slides format format-check schoon

##@ Opzetten

alles: cluster build load deploy ## Cluster, image en het goede manifest in één keer

cluster: ## Maak het kind-cluster (1 control plane, 2 workers)
	kind create cluster --name $(CLUSTER) --config kind-cluster.yaml

build: ## Bouw het image als probe-demo:dev
	$(DOCKER) build -t $(IMAGE) demo/

load: ## Laad het image in het kind-cluster
	kind load docker-image $(IMAGE) --name $(CLUSTER)

deploy: ## Rol het goede manifest uit
	kubectl apply -f $(GOED)
	kubectl rollout status deployment/$(APP)

##@ Scenario's uitrollen (tonen eerst de diff)

kapot-1: ## Scenario 1: readinessProbe weggelaten
	@diff -u $(GOED) $(KAPOT)/1-zonder-readiness.yaml || true
	kubectl apply -f $(KAPOT)/1-zonder-readiness.yaml

kapot-2: ## Scenario 3: liveness checkt de dependency
	@diff -u $(GOED) $(KAPOT)/2-liveness-checkt-dependency.yaml || true
	kubectl apply -f $(KAPOT)/2-liveness-checkt-dependency.yaml

kapot-3: ## Scenario 4: trage start zonder startupProbe
	@diff -u $(GOED) $(KAPOT)/3-trage-start.yaml || true
	kubectl apply -f $(KAPOT)/3-trage-start.yaml

kapot-4: ## Scenario 5: geen preStop, grace period van 1 seconde
	@diff -u $(GOED) $(KAPOT)/4-zonder-prestop.yaml || true
	kubectl apply -f $(KAPOT)/4-zonder-prestop.yaml

herstel: ## Terug naar het goede manifest
	kubectl apply -f $(GOED)
	kubectl rollout status deployment/$(APP)

##@ Meekijken (elk in een eigen terminal)

watch: ## Terminal 1: kijk mee met de pods
	kubectl get pods -l app=$(APP) -o wide -w

endpoints: ## Terminal 2: kijk mee met de EndpointSlice
	kubectl get endpointslice -l kubernetes.io/service-name=$(APP) -w

belasting: ## Terminal 3: genereer verkeer en meld elke foutmelding
	@echo "Verkeer naar $(URL)/work - stoppen met ctrl-c"
	@goed=0; fout=0; \
	while true; do \
	  code=$$(curl -s -o /dev/null -w '%{http_code}' --max-time 2 $(URL)/work || true); \
	  if [ "$$code" = "200" ]; then \
	    goed=$$((goed+1)); \
	    if [ $$((goed % 25)) -eq 0 ]; then echo "$$(date +%T)  ok    $$goed goed, $$fout fout"; fi; \
	  else \
	    fout=$$((fout+1)); \
	    echo "$$(date +%T)  FOUT  $${code:-000}   ($$goed goed, $$fout fout)"; \
	  fi; \
	  sleep 0.1; \
	done

logs: ## Volg de logs van alle replica's
	kubectl logs -l app=$(APP) -f --prefix --max-log-requests=10

##@ Live stukmaken zonder opnieuw uit te rollen

unready: ## Zet readiness om op één pod (POST /debug/unready)
	@pod=$$(kubectl get pods -l app=$(APP) -o jsonpath='{.items[0].metadata.name}'); \
	echo "pod: $$pod"; \
	kubectl port-forward "pod/$$pod" 18080:8080 >/dev/null 2>&1 & \
	pf=$$!; sleep 1; \
	curl -sS -X POST http://127.0.0.1:18080/debug/unready; \
	kill $$pf 2>/dev/null || true

dependency: ## Zet de dependency om op ALLE pods (POST /debug/dependency-down)
	@for pod in $$(kubectl get pods -l app=$(APP) -o name); do \
	  kubectl port-forward "$$pod" 18080:8080 >/dev/null 2>&1 & \
	  pf=$$!; sleep 1; \
	  curl -sS -X POST http://127.0.0.1:18080/debug/dependency-down; \
	  kill $$pf 2>/dev/null || true; \
	  sleep 1; \
	done

##@ Overig

slides: ## Serveer de slides lokaal met marp-cli
	npx --yes @marp-team/marp-cli@latest -s presentatie/

format: ## Formatteer alles: gofmt voor Go, Prettier voor Markdown en YAML
	gofmt -w demo
	npx --yes prettier@3 --write . --ignore-unknown

format-check: ## Controleer de opmaak zonder iets te wijzigen
	@scheef=$$(gofmt -l demo); \
	if [ -n "$$scheef" ]; then echo "niet gofmt-schoon:"; echo "$$scheef"; exit 1; fi
	npx --yes prettier@3 --check . --ignore-unknown

schoon: ## Verwijder het kind-cluster
	kind delete cluster --name $(CLUSTER)

help: ## Toon dit overzicht
	@echo "Workshop: leeft je app nog?"
	@awk 'BEGIN {FS = ":.*?## "} \
	     /^##@ / {printf "\n\033[1m%s\033[0m\n", substr($$0, 5); next} \
	     /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo
	@echo "Voor de sessie: make alles"
