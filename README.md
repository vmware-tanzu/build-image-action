# build-image-action

[![build-and-test](https://github.com/vmware-tanzu/build-image-action/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/build-and-test.yaml)
[![golangci-lint](https://github.com/vmware-tanzu/build-image-action/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/golangci-lint.yml)
[![Build and Publish](https://github.com/vmware-tanzu/build-image-action/actions/workflows/publish-image.yaml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/publish-image.yaml)

## Overview

## Try it out

### Prerequisites

1. Install and setup [kpack](https://github.com/pivotal/kpack/blob/main/docs/install.md)
1. Install [kpack cli](https://github.com/vmware-tanzu/kpack-cli/releases)
1. Create dev namespace
   ```bash
   kubectl create namespace dev
   ```
1. Create a secret with push credentials for the docker registry that you plan on publishing OCI images to with kpack
    ```bash
   kubectl create secret docker-registry kpack-registry-credentials \                                                130 ↵ emjohnson@emjohnson-a03
    --docker-username=_json_key \
    --docker-password="$(cat key.json)" \
    --docker-server=gcr.io \
    --namespace dev
    ```
1. Create a ClusterBuilder (and ServiceAccount, ClusterStore, ClusterStack)
    ```bash
   ytt -v tag=<my-clusterbuilder-tag> -f hack/kpack.yaml | kubectl apply -f -
   ```

   For example:
   ```bash
   ytt -v tag=gcr.io/kontinue/emj/clusterbuilder -f hack/kpack.yaml | kubectl apply -f -
   ```
1. Get appropriate permissions to run
    ```bash
   kubectl apply -f config/rbac.yaml
   eval "$(./config/get_auth.sh)"
    ```

## Usage

### Run locally

```bash
NAMESPACE=dev GITHUB_SERVER_URL=https://github.com GITHUB_REPOSITORY=<my-repo> GITHUB_SHA=<my-sha> TAG=<my-tag> GITHUB_OUTPUT=<my-output> SERVICE_ACCOUNT_NAME=kpack-service-account go run main.go
```

For example:

```bash
NAMESPACE=dev GITHUB_SERVER_URL=https://github.com/emmjohnson GITHUB_REPOSITORY=github-actions-poc GITHUB_SHA=e84d037eedbbd7fefc8da0e2c7609e05faef5f0e TAG=gcr.io/kontinue/emj/app-action GITHUB_OUTPUT=/Users/emjohnson/sandbox/vmware-tanzu/build-image-action/output.txt SERVICE_ACCOUNT_NAME=kpack-service-account go run main.go
kubectl get builds -n dev
kubectl get pods -n dev
cat /Users/emjohnson/sandbox/vmware-tanzu/build-image-action/output.txt
name=gcr.io/kontinue/emj/app-action@sha256:a37e5abcefaa73417eff08f9771840460334d0543287a777c40d16f15ab0ecca
```

### Run as action


### Setup

In order to use this action a service account will need to exist inside TAP that has permissions to access the required resources. The 
[example file](https://github.com/vmware-tanzu/build-image-action/blob/main/config/rbac.yaml) contains the minimum required permissions.

To apply this file to a namespace called `dev`:

```
kubectl apply -f https://raw.githubusercontent.com/vmware-tanzu/build-image-action/main/config/rbac.yaml
```

Then to access the values:

```
SECRET=$(kubectl get sa github-actions -oyaml | yq '.secrets[0].name')

CA_CERT=$(kubectl get secret $SECRET -oyaml | yq '.data."ca.crt"')
NAMESPACE=$(kubectl get secret $SECRET -oyaml | ksd | yq .stringData.namespace)
TOKEN=$(kubectl get secret $SECRET -oyaml | ksd | yq .stringData.token)
SERVER=$(kubectl config view --minify | yq '.clusters[0].cluster.server')
```

Using the GitHub cli create the required secrets on the repository:

```
gh secret set CA_CERT --app actions --body "$CA_CERT"
gh secret set NAMESPACE --app actions --body "$NAMESPACE"
gh secret set TOKEN --app actions --body "$TOKEN"
gh secret set SERVER --app actions --body "$SERVER"
``` 

### Usage

#### Auth

  - `server`: Host of the API Server.
  - `ca-cert`: CA Certificate of the API Server.
  - `token`: Service Account token to access kubernetes.
  - `namespace`: _(required)_ The namespace to create the build resource in.

#### Image Configuration

  - `destination`: _(required)_
  - `env`:
  - `serviceAccountName`: Name of the service account in the namespace, defaults to `default`

#### Basic Configuration

```yaml
- name: Build Image
  id: build
  uses: vmware-tanzu/build-image-action@v1
  with:
    # auth
    server: ${{ secrets.SERVER }}
    token: ${{ secrets.TOKEN }}
    ca_cert: ${{ secrets.CA_CERT }}
    namespace: ${{ secrets.NAMESPACE }}
    # image config
    destination: gcr.io/project-id/name-for-image
    env: |
      BP_JAVA_VERSION=17
```

##### Outputs

  - `name`: The full name, including sha of the built image.

##### Example

```yaml
- name: Do something with image
  run:
    echo "${{ steps.build.outputs.name }}"
```

## Documentation

TODO

## Contributing

The build-image-action project team welcomes contributions from the community. Before you start working with build-image-action, please
read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

TODO The scripts and documentation in this project are released under the [Apache 2](LICENSE).
