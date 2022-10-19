#!/usr/bin/env bash

set -euo pipefail

export CA_CERT=${INPUT_CA_CERT}
export SERVER=${INPUT_SERVER}
export NAMESPACE=${INPUT_NAMESPACE}
export TOKEN=${INPUT_TOKEN}
export TAG=${INPUT_DESTINATION}
export ENV_VARS=${INPUT_ENV}
export SERVICE_ACCOUNT_NAME=${INPUT_SERVICEACCOUNTNAME}

/usr/bin/builder

