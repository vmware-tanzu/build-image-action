#!/bin/bash

set -euo pipefail

SECRET=$(kubectl get sa github-actions -oyaml -n dev | yq '.secrets[0].name')

CA_CERT=$(kubectl get secret $SECRET -oyaml -n dev| yq '.data."ca.crt"')
NAMESPACE=$(kubectl get secret $SECRET -oyaml -n dev| ksd | yq .stringData.namespace)
TOKEN=$(kubectl get secret $SECRET -oyaml -n dev| ksd | yq .stringData.token)
SERVER=$(kubectl config view --minify | yq '.clusters[0].cluster.server')

echo "export CA_CERT=$CA_CERT"
echo "export NAMESPACE=$NAMESPACE"
echo "export TOKEN=$TOKEN"
echo "export SERVER=$SERVER"
