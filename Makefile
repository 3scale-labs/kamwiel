SHELL = /bin/bash

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))

CLUSTER_NAMESPACE ?= kamwiel
KAMWIEL_IMG ?= kamwiel:latest

AUTHORINO_IMAGE ?= quay.io/3scale/authorino:latest
AUTHORINO_DEPLOYMENT ?= namespaced
AUTHORINO_REPLICAS ?= 1

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: help build

.PHONY : help
## help:				Prints this guide
help: Makefile
	@{ \
	echo "These Makefile will help you to build, test, deploy Kamwiel in a Kind cluster for development"; \
	echo ;\
	echo "Usage: make [target]"; \
	echo ;\
	echo "Available targets:";\
	echo ;\
	sed -n 's/^##//p' $<;\
  	}

## build:				Build Kamwiel binary
build: fmt vet
	go build -o bin/kamwiel main.go

## vendor:			Download vendor dependencies
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

## fmt:				Run go fmt against code
fmt:
	go fmt ./...

## vet:				Run go vet against code
vet:
	go vet ./...

## docker-build:			Build kamwiel docker image
docker-build: vendor
	docker build . -t ${KAMWIEL_IMG}

## docker-push:			Push kamwiel image to docker repo
docker-push:
	docker push ${KAMWIEL_IMG}

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

kind:
ifeq (, $(shell which kind))
	@{ \
	set -e ;\
	KIND_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KIND_GEN_TMP_DIR ;\
	go mod init tmp ;\
	GO111MODULE="on" go get sigs.k8s.io/kind@v0.10.0 ;\
	rm -rf $$KIND_GEN_TMP_DIR ;\
	}
KIND=$(GOBIN)/kind
else
KIND=$(shell which kind)
endif

KIND_CLUSTER_NAME ?= kamwiel-cluster

## namespace:			Creates a namespace where to deploy Kamwiel
.PHONY: namespace
namespace:
	kubectl create namespace $(CLUSTER_NAMESPACE)


## cert-manager:		 	Install CertManager to the Kubernetes cluster
.PHONY: cert-manager
cert-manager:
	kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
	kubectl -n cert-manager wait --timeout=300s --for=condition=Available deployments --all

## generate-kuadrant-manifest:	Generates Kuadrant CRDs.
KUADRANT_VERSION=v0.0.1-pre2
.PHONY: generate-kuadrant-manifest
generate-kuadrant-manifest: kustomize
	$(eval TMP := $(shell mktemp -d))
	cd $(TMP); git clone --depth 1 --branch $(KUADRANT_VERSION) https://github.com/kuadrant/kuadrant-controller.git
	cd $(TMP)/kuadrant-controller; make kustomize; $(KUSTOMIZE) build config/crd -o $(PROJECT_PATH)/examples/kuadrant/autogenerated/kuadrant-manifest.yaml
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/kuadrant/autogenerated/kuadrant-manifest.yaml
	-rm -rf $(TMP)

## generate-authorino-manifest:	Generates Authorino CRDs, RBAC, etc.
AUTHORINO_VERSION=v0.2.1-pre
.PHONY: generate-authorino-manifest
generate-authorino-manifest: kustomize
	$(eval TMP := $(shell mktemp -d))
	cd $(TMP); git clone --depth 1 --branch $(AUTHORINO_VERSION) https://github.com/kuadrant/authorino.git
	cd $(TMP)/authorino; make kustomize; $(KUSTOMIZE) build install -o $(PROJECT_PATH)/examples/authorino/autogenerated/authorino-manifest.yaml
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/authorino/autogenerated/authorino-manifest.yaml
	-rm -rf $(TMP)


# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
.PHONY: deploy-authorino
deploy-authorino: kustomize
	$(eval TMP := $(shell mktemp -d))
	cd $(TMP); git clone --depth 1 --branch $(AUTHORINO_VERSION) https://github.com/kuadrant/authorino.git
	cd $(TMP)/authorino; make kustomize; $(KUSTOMIZE) build install -o $(PROJECT_PATH)/examples/authorino/autogenerated/authorino-manifest.yaml
	cd $(TMP)/authorino/deploy/base && $(KUSTOMIZE) edit set image authorino=$(AUTHORINO_IMAGE) && $(KUSTOMIZE) edit set namespace $(CLUSTER_NAMESPACE) && $(KUSTOMIZE) edit set replicas authorino-controller-manager=$(AUTHORINO_REPLICAS)
	cd $(TMP)/authorino/deploy/overlays/$(AUTHORINO_DEPLOYMENT) && $(KUSTOMIZE) edit set namespace $(CLUSTER_NAMESPACE)
	cd $(TMP)/authorino/ && $(KUSTOMIZE) build deploy/overlays/$(AUTHORINO_DEPLOYMENT) | kubectl -n $(CLUSTER_NAMESPACE) apply -f -
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/authorino/autogenerated/authorino-manifest.yaml
	-rm -rf $(TMP)

# Deploy Envoy
.PHONY: deploy-envoy
deploy-envoy:
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/envoy-deployment.yaml

# Deploy Kamwiel
.PHONY: deploy-kamwiel
deploy-kamwiel:
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/kamwiel-deployment.yaml

## create-apikey:			Creates a new api key in the cluster, API_KEY_NAME env specifies the identifier of the key
API_KEY_NAME ?= kamwiel-apikey-1
API_KEY ?= $(eval API_KEY := $(shell openssl rand -hex 32))$(API_KEY)
.PHONY: create-apikey
create-apikey:
	@{ \
	echo "***************************************************************************"; \
	echo '{"apiVersion": "v1", "kind": "Secret", "metadata": {"name": "$(API_KEY_NAME)", "labels": { "authorino.3scale.net/managed-by": "authorino", "custom-label": "friends" } }, "stringData": { "api_key": "$(API_KEY)" }, "type": "Opaque"}' | kubectl -n $(CLUSTER_NAMESPACE) apply -f - ; \
	echo "API KEY successfully created with name: ${API_KEY_NAME} and value: ${API_KEY}"; \
	echo "***************************************************************************"; \
	echo; \
	}

## example-config:		Applies authorino protection and kuadrant API Product samples
.PHONY: example-config
example-config:
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/authorino-protection.yaml
	kubectl -n $(CLUSTER_NAMESPACE) apply -f examples/kuadrant/samples/api_samples.yaml

## local-cluster-up:		Start a local Kubernetes cluster using Kind
.PHONY: local-cluster-up
local-cluster-up: kind local-cleanup
	kind create cluster --name $(KIND_CLUSTER_NAME) --config ./utils/kind-cluster.yaml


# Pushes a local container image of Kamwiel to the registry of the Kind-started local Kubernetes cluster
.PHONY: local-push
local-push: kind
	kind load docker-image $(KAMWIEL_IMG) --name $(KIND_CLUSTER_NAME)

## local-deploy:			Builds the image, pushes to the local cluster and deploys Kamwiel.
# Sets the imagePullPolicy to 'IfNotPresent' so it doesn't try to pull the image again (just pushed into the server registry)
.PHONY: local-deploy
local-deploy: docker-build local-push deploy-kamwiel
	kubectl -n $(CLUSTER_NAMESPACE) patch deployment kamwiel -p '{"spec": {"template": {"spec":{"containers":[{"name": "kamwiel", "imagePullPolicy":"IfNotPresent"}]}}}}'

## local-rollout:			Rebuild and push the docker image and redeploy Kamwiel to kind
.PHONY: local-rollout
local-rollout: docker-build local-push
	kubectl -n $(CLUSTER_NAMESPACE) rollout restart deployment.apps/kamwiel

## local-setup:			Set up a test/dev local Kubernetes server loaded up with a freshly built Kamwiel image plus dependencies
.PHONY: local-setup
local-setup: local-cluster-up namespace local-deploy cert-manager deploy-envoy deploy-authorino example-config create-apikey
	kubectl -n $(CLUSTER_NAMESPACE) wait --timeout=500s --for=condition=Available deployments --all
	@{ \
	echo "Now you can export the envoy service by doing:"; \
	echo "kubectl port-forward --namespace $(CLUSTER_NAMESPACE) deployment/envoy 8000:8000 &"; \
	echo "After that, you can curl kamwiel with the created API KEY like:"; \
	echo "curl -H 'X-API-KEY: $(API_KEY)' http://kamwiel-authorino.127.0.0.1.nip.io:8000/ping -v"; \
	echo ;\
	echo "***************************************************************************"; \
	echo "************************** Voilà, profit!!! *******************************"; \
	echo "***************************************************************************"; \
	}

## local-cleanup:			Deletes the local Kubernetes cluster started using Kind
.PHONY: local-cleanup
local-cleanup: kind
	kind delete cluster --name $(KIND_CLUSTER_NAME)
