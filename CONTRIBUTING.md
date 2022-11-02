# Contributing to build-image-action

We welcome contributions from the community and first want to thank you for taking the time to contribute!

Please familiarize yourself with the [Code of Conduct](https://github.com/vmware/.github/blob/main/CODE_OF_CONDUCT.md)
before contributing.

Before you start working with build-image-action, please read and sign our Contributor License
Agreement [CLA](https://cla.vmware.com/cla/1/preview). If you wish to contribute code and you have not signed our
contributor license agreement (CLA), our bot will prompt you to do so when you open a Pull Request. For any questions
about the CLA process, please refer to our [FAQ]([https://cla.vmware.com/faq](https://cla.vmware.com/faq)).

## Ways to contribute

We welcome many types of contributions and not all of them need a Pull request. Contributions may include:

* New features and proposals
* Documentation
* Bug fixes
* Issue Triage
* Answering questions and giving feedback
* Helping to onboard new contributors
* Other related activities

## Getting started

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

### Usage

#### Run locally

```bash
NAMESPACE=dev GITHUB_SERVER_URL=https://github.com GITHUB_REPOSITORY=<my-repo> GITHUB_SHA=<my-sha> TAG=<my-tag> GITHUB_OUTPUT=<my-output> SERVICE_ACCOUNT_NAME=kpack-service-account go run main.go
```

For example:

```bash
NAMESPACE=dev GITHUB_SERVER_URL=https://github.com GITHUB_REPOSITORY=emmjohnson/github-actions-poc GITHUB_SHA=e84d037eedbbd7fefc8da0e2c7609e05faef5f0e TAG=gcr.io/kontinue/emj/app-action GITHUB_OUTPUT=/Users/emjohnson/sandbox/vmware-tanzu/build-image-action/output.txt SERVICE_ACCOUNT_NAME=kpack-service-account go run main.go
kubectl get builds -n dev
kubectl get pods -n dev
cat /Users/emjohnson/sandbox/vmware-tanzu/build-image-action/output.txt
name=gcr.io/kontinue/emj/app-action@sha256:a37e5abcefaa73417eff08f9771840460334d0543287a777c40d16f15ab0ecca
```

### Run as action

[//]: # (TODO: See POC)

## Contribution Flow

This is a rough outline of what a contributor's workflow looks like:

* Make a fork of the repository within your GitHub account
* Create a topic branch in your fork from where you want to base your work
* Make commits of logical units
* Make sure your commit messages are with the proper format, quality and descriptiveness (see below)
* Push your changes to the topic branch in your fork
* Create a pull request containing that commit

We follow the GitHub workflow and you can find more details on
the [GitHub flow documentation](https://docs.github.com/en/get-started/quickstart/github-flow).

### Pull Request Checklist

Before submitting your pull request, we advise you to use the following:

1. Check if your code changes will pass both code linting checks and unit tests.
2. Ensure your commit messages are descriptive. We follow the conventions
   on [How to Write a Git Commit Message](http://chris.beams.io/posts/git-commit/). Be sure to include any related
   GitHub issue references in the commit message.
   See [GFM syntax](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax#GitHub-flavored-markdown)
   for referencing issues and commits.
3. Check the commits and commits messages and ensure they are free from typos.

## Reporting Bugs and Creating Issues

For specifics on what to include in your report, please follow the guidelines in the issue and pull request templates
when available.

## Ask for Help

The best way to reach us with a question when contributing is to ask on:

* The original GitHub issue
