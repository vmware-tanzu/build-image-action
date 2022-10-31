# build-image-action

[![build-and-test](https://github.com/vmware-tanzu/build-image-action/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/build-and-test.yaml)
[![golangci-lint](https://github.com/vmware-tanzu/build-image-action/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/golangci-lint.yml)
[![Build and Publish](https://github.com/vmware-tanzu/build-image-action/actions/workflows/publish-image.yaml/badge.svg)](https://github.com/vmware-tanzu/build-image-action/actions/workflows/publish-image.yaml)

## Overview

## Try it out

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
