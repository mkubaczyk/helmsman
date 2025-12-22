# Run Helmsman in CI

You can run Helmsman as a job in your CI system using the [helmsman docker image](https://github.com/mkubaczyk/helmsman/pkgs/container/helmsman).
The following example is a `config.yml` file for CircleCI but can be replicated for other CI systems.

```yaml
version: 2
jobs:

    deploy-apps:
      docker:
        - image: ghcr.io/mkubaczyk/helmsman:latest
      steps:
        - checkout
        - run:
            name: Deploy Helm Packages using helmsman
            command: helmsman --apply -f helmsman-deployments.toml


workflows:
  version: 2
  build:
    jobs:
      - deploy-apps
```

> IMPORTANT: If your CI build logs are publicly readable, don't use the `--verbose` together with `--debug` flags as logs any secrets being passed from env vars to the helm charts.

The `helmsman-deployments.toml` is your desired state file which will version controlled in your git repo.
