# build-image-action

[![build-and-test](https://github.com/vmware-tanzu/build-image-action/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/build-and-test.yaml)
[![golangci-lint](https://github.com/vmware-tanzu/build-image-action/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/golangci-lint.yml)
[![Build and Publish](https://github.com/vmware-tanzu/build-image-action/actions/workflows/publish-image.yaml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/publish-image.yaml)

This GitHub Action creates a TBS Build on the given cluster.

## Overview

## Try it out

### Setup

In order to use this action there are two things that need to be configured:

1. Ensure that the GitHub action runner has access to the kubernetes API Server.
1. Configure a service account that has the required permissions.

#### Access to the kubernetes API server

The GitHub action talks directly to the kubernetes API server, so if you are running this on github.com with the default action runners
you'll need to ensure your API server is accessable from GitHubs [IP ranges](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-githubs-ip-addresses). 
Alternatively it may be possible to runner the action on a custom runner within your firewall (with access to the TAP cluster).

#### Permissions Required

The minimum permissions required on the TBS cluster are documented below:

```
ClusterRole
 └ kpack.io
   └ clusterbuilders verbs=[get]
Role (developer namespace)
 ├ ''
 │ ├ pods verbs=[get watch list] ✔
 │ └ pods/log verbs=[get] ✔
 └ kpack.io
   └ builds verbs=[get watch list create delete] ✔
```

The [example file](https://github.com/vmware-tanzu/build-image-action/blob/main/config/rbac.yaml) contains the minimum
required permissions.

To apply this file to a namespace called `dev`:

```
kubectl apply -f https://raw.githubusercontent.com/vmware-tanzu/build-image-action/main/config/rbac.yaml
```

Then to access the values:

```
SECRET=$(kubectl get sa github-actions -oyaml -n dev | yq '.secrets[0].name')

CA_CERT=$(kubectl get secret $SECRET -oyaml -n dev | yq '.data."ca.crt"')
NAMESPACE=$(kubectl get secret $SECRET -oyaml -n dev | ksd | yq .stringData.namespace)
TOKEN=$(kubectl get secret $SECRET -oyaml -n dev | ksd | yq .stringData.token)
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
- `clusterBuilder`: Name of the cluster builder to use, defaults to `default`

#### Basic Configuration

```yaml
- name: Build Image
  id: build
  uses: vmware-tanzu/build-image-action@v1-alpha
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

The build-image-action project team welcomes contributions from the community. Before you start working with 
this project please read and sign our Contributor License Agreement CLA. If you wish to contribute code and 
you have not signed our contributor license agreement (CLA), our bot will prompt you to do so when you open 
a Pull Request. For any questions about the CLA process, please refer to our FAQ.

## License

The scripts and documentation in this project are released under the [Apache 2](LICENSE).
