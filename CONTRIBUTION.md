# Contribution Guide

Pull requests, feeback/feature requests are all welcome. This guide will be updated overtime.

## Build helmsman from source

To build helmsman from source, you need go:1.17+. Follow the steps below:

```sh
git clone https://github.com/mkubaczyk/helmsman.git
make tools # installs few tools for testing, building, releasing
make build
make test
```

## The branches and tags

`master` is where Helmsman latest code lives.
`1.x` this is where Helmsman versions 1.x lives.

> Helmsman v1.x supports helm v2.x only and will no longer be supported except for bug fixes and minor changes.

## Submitting pull requests

- If your PR is for Helmsman v1.x, it should target the `1.x` branch.
- Please make sure you state the purpose of the pull request and that the code you submit is documented. If in doubt, [this guide](https://blog.github.com/2015-01-21-how-to-write-the-perfect-pull-request/) offers some good tips on writing a PR.
- Please make sure you update the documentation with new features or the changes your PR adds. The following places are required.
  - Update existing [how_to](docs/how_to/) guides or create new ones.
  - If necessary, Update the [Desired State File spec](docs/desired_state_specification.md)
  - If adding new flags, Update the [cmd reference](docs/cmd_reference.md)
- Please add tests wherever possible to test your new changes.

## Contribution to documentation

Contribution to the documentation can be done via pull requests or by opening an issue.

## Reporting issues/feature requests

Please provide details of the issue, versions of helmsman, helm and kubernetes and all possible logs.

## Releasing Helmsman

Release is automated from CicrcleCI based on Git tags. [Goreleaser](goreleaser.com) is used to release the binaries and update the release notes on Github while the CircleCI pipeline builds a set of docker images and pushes them to dockerhub.

The following steps are needed to cut a release (They assume that you are on master and the code is up to date):

1. Change the version variable in [main.go](internal/app/main.go) and in [.version](.version)
2. Update the [release-notes.md](release-notes.md) file with new version and changelog.
3. Update the installation section in the [README.md](README.md) file to point to the latest version.
4. (Optional), if new helm versions are required, update the [circleci config](.circleci/config.yml) and add more docker commands.
5. Commit your changes locally.
6. Create a git tag with the following command: `git tag -a <semantic version number> -m "<semantic version number>" <your-last-commit-sha>`
7. Push your commit and tag with `git push --follow-tags`
8. This should trigger the [pipeline on circleci](https://circleci.com/gh/Praqma/workflows/helmsman) which eventually releases to Github and dockerhub.
