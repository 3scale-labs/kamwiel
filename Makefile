# Current Operator version
VERSION ?= 0.0.1

# Image URL to use all building/pushing image targets
IMG ?= kamwiel:latest

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

# Deletes the local Kubernetes cluster started using Kind
.PHONY: local-cluster-down
local-cluster-down: kind
	kind delete cluster --name $(KIND_CLUSTER_NAME)

# Pushes a local container image of Kamwiel to the registry of the Kind-started local Kubernetes cluster
.PHONY: local-push
local-push: kind
	kind load docker-image $(IMG) --name $(KIND_CLUSTER_NAME)

# Set up a test/dev local Kubernetes server loaded up with a freshly built Kamwiel image protected with Authorino
.PHONY: local-setup
local-setup: vendor kind local-cluster-up docker-build local-push
	utils/local-setup.sh

