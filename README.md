[![GitHub release](https://img.shields.io/github/v/release/mkubaczyk/helmsman)](https://github.com/mkubaczyk/helmsman/releases)

![helmsman-logo](docs/images/helmsman.png)

> Helmsman v4.x supports Helm 3.x and Helm 4.x. For Helm 2.x, use Helmsman v1.x

# What is Helmsman?

Helmsman is a Helm Charts (k8s applications) as Code tool which allows you to automate the deployment/management of your Helm charts from version controlled code.

# Why has this repository changed the owner?

The previous owner (Praqma company, later Eficode company) of this repository decided to transfer the repository into the hands of current maintainers
to make sure the project can be developed further with no interruptions and unnecessary dependencies.
We'll do our best to get it up and running as soon as possible.
Thank you for your patience and trusting Helmsman with your tasks!

# How does it work?

Helmsman uses a simple declarative [TOML](https://github.com/toml-lang/toml) file to allow you to describe a desired state for your k8s applications as in the [example toml file](https://github.com/mkubaczyk/helmsman/blob/master/examples/example.toml).
Alternatively YAML declaration is also acceptable [example yaml file](https://github.com/mkubaczyk/helmsman/blob/master/examples/example.yaml).

The desired state file (DSF) follows the [desired state specification](https://github.com/mkubaczyk/helmsman/blob/master/docs/desired_state_specification.md).

Helmsman sees what you desire, validates that your desire makes sense (e.g. that the charts you desire are available in the repos you defined), compares it with the current state of Helm and figures out what to do to make your desire come true.

To plan without executing:

```sh
helmsman -f example.toml
```

To plan and execute the plan:

```sh
helmsman --apply -f example.toml
```

To show debugging details:

```sh
helmsman --debug --apply -f example.toml
```

To run a dry-run:

```sh
helmsman --debug --dry-run -f example.toml
```

To limit execution to specific application:

```sh
helmsman --debug --dry-run --target artifactory -f example.toml
```

# Features

- **Built for CD**: Helmsman can be used as a docker image or a binary.
- **Applications as code**: describe your desired applications and manage them from a single version-controlled declarative file.
- **Suitable for Multitenant Clusters**: manage releases across multiple namespaces with fine-grained access control.
- **Easy to use**: deep knowledge of Helm CLI and Kubectl is NOT mandatory to use Helmsman.
- **Plan, View, apply**: you can run Helmsman to generate and view a plan with/without executing it.
- **Portable**: Helmsman can be used to manage charts deployments on any k8s cluster.
- **Protect Namespaces/Releases**: you can define certain namespaces/releases to be protected against accidental human mistakes.
- **Define the order of managing releases**: you can define the priorities at which releases are managed by helmsman (useful for dependencies).
- **Parallelise**: Releases with the same priority can be executed in parallel.
- **Idempotency**: As long your desired state file does not change, you can execute Helmsman several times and get the same result.
- **Continue from failures**: In the case of partial deployment due to a specific chart deployment failure, fix your helm chart and execute Helmsman again without needing to rollback the partial successes first.

# Install

## From binary

Please make sure the following are installed prior to using `helmsman` as a binary (the docker image contains all of them):

- [kubectl](https://github.com/kubernetes/kubectl)
- [helm](https://github.com/helm/helm) (helm v3.x or v4.x for `helmsman` v4.x)
- [helm-diff](https://github.com/databus23/helm-diff) (`helmsman` >= 1.6.0)

If you use private helm repos, you will need either `helm-gcs` or `helm-s3` plugin or you can use basic auth to authenticate to your repos. See the [docs](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/helm_repos) for details.

Check the [releases page](https://github.com/mkubaczyk/helmsman/releases) for the different versions.

```sh
# Set desired version (see latest release link below)
VERSION="4.0.5"

# on Linux
curl -L https://github.com/mkubaczyk/helmsman/releases/download/v${VERSION}/helmsman_${VERSION}_linux_amd64.tar.gz | tar zx
# on MacOS
curl -L https://github.com/mkubaczyk/helmsman/releases/download/v${VERSION}/helmsman_${VERSION}_darwin_amd64.tar.gz | tar zx

mv helmsman /usr/local/bin/helmsman
```

See the [latest release](https://github.com/mkubaczyk/helmsman/releases/latest) for current version.

## As a docker image

Docker images are published to `ghcr.io/mkubaczyk/helmsman` with variants for Helm 3 and Helm 4:

| Tag | Description |
|-----|-------------|
| `latest` | Latest release with Helm 4 |
| `<version>` | Specific release with Helm 4 (default) |
| `<version>-helm3` | Specific release with Helm 3 |
| `<version>-helm4` | Specific release with Helm 4 |
| `<version>-helm<helm-version>` | Specific release with exact Helm version |

```sh
# Latest with Helm 4 (default)
docker pull ghcr.io/mkubaczyk/helmsman:latest

# Specific release with Helm 4
docker pull ghcr.io/mkubaczyk/helmsman:<version>

# Specific release with Helm 3
docker pull ghcr.io/mkubaczyk/helmsman:<version>-helm3
```

Replace `<version>` with the desired Helmsman version. See the [latest release](https://github.com/mkubaczyk/helmsman/releases/latest) for available tags.

## As a package

Helmsman has been packaged in Archlinux under `helmsman-bin` for the latest binary release, and `helmsman-git` for master.

You can also install Helmsman using [Homebrew](https://brew.sh)

```sh
brew install helmsman
```

## As an [asdf-vm](https://asdf-vm.com/) plugin

```sh
asdf plugin-add helmsman
asdf install helmsman latest
```

# Documentation

> Documentation for Helmsman v1.x can be found at: [docs v1.x](https://github.com/mkubaczyk/helmsman/tree/1.x/docs)

- [How-Tos](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/).
- [Desired state specification](https://github.com/mkubaczyk/helmsman/blob/master/docs/desired_state_specification.md).
- [CMD reference](https://github.com/mkubaczyk/helmsman/blob/master/docs/cmd_reference.md)

## Usage

Helmsman can be used in three different settings:

- [As a binary with a hosted cluster](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/settings).
- [As a docker image in a CI system or local machine](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/deployments/ci.md) Always use a tagged docker image from [GHCR](https://github.com/mkubaczyk/helmsman/pkgs/container/helmsman).
- [As a docker image inside a k8s cluster](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/deployments/inside_k8s.md)

# Contributing

Pull requests, feedback/feature requests are welcome. Please check our [contribution guide](CONTRIBUTION.md).
