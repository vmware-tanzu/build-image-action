#!/usr/bin/env bash

set -euo pipefail

/usr/bin/builder kpack --github-server-url="${GITHUB_SERVER_URL}" \
  --github-repository="${GITHUB_REPOSITORY}" \
  --github-sha="${GITHUB_SHA}" \
  --github-action-output="${GITHUB_OUTPUT}" \
  --ca-cert="${INPUT_CA_CERT}" \
  --server="${INPUT_SERVER}" \
  --token="${INPUT_TOKEN}" \
  --namespace="${INPUT_NAMESPACE}" \
  --tag="${INPUT_DESTINATION}" \
  --env-vars="${INPUT_ENV}" \
  --service-account-name="${INPUT_SERVICEACCOUNTNAME}" \
  --cluster-builder="${INPUT_CLUSTERBUILDER}" \
  --timeout="${INPUT_TIMEOUT}"


