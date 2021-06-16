# Current Operator version
VERSION ?= 0.0.1

# Image URL to use all building/pushing image targets
IMG ?= kamwiel:latest
KAMWIEL_NAMESPACE ?= kamwiel

#Use bash as shell
SHELL = /bin/bash

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

# Build Kamwiel binary
build: generate fmt vet
	go build -o bin/kamwiel main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Download vendor dependencies
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build: vendor
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

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

# Start a local Kubernetes cluster using Kind
.PHONY: local-cluster-up
local-cluster-up: kind local-cluster-down
	kind create cluster --name $(KIND_CLUSTER_NAME) --config ./utils/kind-cluster.yaml


# Pushes a local container image of Kamwiel to the registry of the Kind-started local Kubernetes cluster
.PHONY: local-push
local-push: kind
	kind load docker-image $(IMG) --name $(KIND_CLUSTER_NAME)

# Builds the image, pushes to the local cluster and deploys Kamwiel.
# Sets the imagePullPolicy to 'IfNotPresent' so it doesn't try to pull the image again (just pushed into the server registry)
.PHONY: deploy
local-deploy: docker-build local-push deploy
	kubectl -n $(KAMWIEL_NAMESPACE) patch deployment kamwiel -p '{"spec": {"template": {"spec":{"containers":[{"name": "kamwiel", "imagePullPolicy":"IfNotPresent"}]}}}}'

# Rebuild and push the docker image and redeploy Kamwiel to kind
.PHONY: local-rollout
local-rollout: docker-build local-push
	kubectl -n $(KAMWIEL_NAMESPACE) rollout restart deployment.apps/kamwiel

# Set up a test/dev local Kubernetes server loaded up with a freshly built Kamwiel image protected with Authorino
.PHONY: local-setup
local-setup: vendor kind local-cluster-up docker-build local-push
	utils/local-setup.sh

# Deletes the local Kubernetes cluster started using Kind
.PHONY: local-cleanup
local-cleanup: kind
	kind delete cluster --name $(KIND_CLUSTER_NAME)
