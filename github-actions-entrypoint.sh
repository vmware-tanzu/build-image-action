#!/usr/bin/env bash

set -euo pipefail

/usr/bin/builder kpack --ca-cert="${INPUT_CA_CERT}" \
  --server="${INPUT_SERVER}" \
  --token="${INPUT_TOKEN}" \
  --namespace="${INPUT_NAMESPACE}" \
  --tag="${INPUT_DESTINATION}" \
  --env-vars="${INPUT_ENV}" \
  --service-account-name="${INPUT_SERVICEACCOUNTNAME}"


