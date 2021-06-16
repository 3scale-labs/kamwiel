#!/bin/bash

set -euo pipefail

export NAMESPACE="kamwiel"
export AUTHORINO_IMG="quay.io/3scale/authorino:latest"

echo
echo
echo "***************************************************************************"
echo "Creating namespace"
kubectl create namespace "${NAMESPACE}"
echo "***************************************************************************"
echo


echo "***************************************************************************"
echo "Deploying Envoy"
kubectl -n "${NAMESPACE}" apply -f examples/envoy.yaml
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Deploying Kamwiel"
kubectl -n "${NAMESPACE}" apply -f examples/kamwiel.yaml
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Deploying Authorino"
kubectl config set-context --current --namespace="${NAMESPACE}"
kubectl apply -f examples/authorino
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Wait for all deployments to be up"
kubectl -n "${NAMESPACE}" wait --timeout=500s --for=condition=Available deployments --all
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Applying Kuadrant's CRD and sample data"
kubectl apply -f examples/kuadrant/networking.kuadrant.io_apis.yaml
kubectl -n "${NAMESPACE}" apply -f examples/kuadrant/api_samples.yaml
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Applying authorino config"
kubectl -n "${NAMESPACE}" apply -f examples/authorino.yaml
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Creating API KEY"
export API_KEY=$(openssl rand -hex 32)
kubectl -n "${NAMESPACE}" apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: kamwiel-friend-api-1
  labels:
    authorino.3scale.net/managed-by: authorino
    custom-label: friends
stringData:
  api_key: $API_KEY
type: Opaque
EOF
echo "***************************************************************************"
echo

echo "***************************************************************************"
echo "Now you can export the envoy service by doing:"
echo "kubectl port-forward --namespace ${NAMESPACE} deployment/envoy 8000:8000 &"
echo "After that, you can curl kamwiel with the created API KEY like:"
echo
echo "curl -H 'X-API-KEY: ${API_KEY}' http://kamwiel-authorino.127.0.0.1.nip.io:8000/ping -v"
echo
echo "***************************************************************************"
echo "************************** Voilà, profit!!! *******************************"
echo "***************************************************************************"
