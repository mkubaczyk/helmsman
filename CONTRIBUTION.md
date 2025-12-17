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

Release is automated via GitHub Actions based on Git tags. [Goreleaser](https://goreleaser.com) builds and publishes binaries to GitHub Releases, while the Docker workflow builds and pushes images to GHCR.

To cut a release:

1. Create a PR updating [release-notes.md](release-notes.md) with the new version and changelog.
2. Get approval and merge the PR.
3. Create and push the tag on the merged commit:
   ```bash
   git checkout master && git pull
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push --tags
   ```

The tag triggers the build pipeline which runs tests, creates the GitHub release with binaries, and pushes the Docker image.
