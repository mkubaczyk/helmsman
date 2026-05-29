---
version: v3.5.2
---

# Helmsman Lifecycle hooks

With lifecycle hooks, you can declaratively define certain operations to perform before and/or after helmsman operations.
These operations can be running installing dependencies (e.g. CRDs), executing certain tests, sending custom notifications, etc.
Another useful use-case is if you are using a 3rd party chart which does not define native helm lifecycle hooks that you wish to have.

## Prerequisites

- Hook operations can be defined in a Kubernetes manifest. They can be any kubernetes resource(s) (jobs, cron jobs, deployments, pods, etc).
  - You can only define one manifest file for each lifecycle hook. So make sure all your needed resources are in this manifest.
- Hook operations can also be a script or a command.
- Script or manifest paths must be either absolute or relative to the DSF.
- Hook k8s manifests can also be defined as an URL.

## Supported lifecycle stages

> hook types are case sensitive. Also, note the camleCase.

- `preInstall` : before installing a release.
- `postInstall`: after installing a release.
- `preUpgrade`: before upgrading a release.
- `postUpgrade`: after upgrading a release.
- `preDelete`: before uninstalling a release.
- `postDelete`: after uninstalling a release.

## Hooks stanza details

The following items can be defined in the hooks stanza:

**pre/postInstall, pre/postUpgrade, pre/postDelete**:

A valid path (URL, cloud bucket, local file path) to your hook's k8s manifest or a valid path to a script or a shell command.

The following options only apply to kubernetes manifest type of hooks.

**successCondition**:

The Kubernetes status condition that indicates that your resources have finished their job successfully. You can find out what the status conditions are for different k8s resources with a kubectl command similar to: `kubectl get job -o=jsonpath='{range .items[*]}{.status.conditions[0].type}{"\n"}{end}'`

- For jobs, it is `Complete`
- For pods, it is `Initialized`
- For deployments, it is `Available`

**successTimeout**: (default 30s)

How much time to wait for the `successCondition`

**deleteOnSuccess**: (true/false)

Indicates if you wish to delete the hook's manifest after the hook succeeds. This is only used if you define `successCondition`

> Note: successCondition, deleteOnSuccess and successTimeout are ignored when the `--dry-run` flag is used.

## Global vs App-specific hooks

You can define two types of hooks in your desired state file:

**Global** hooks:

Are defined in the `settings` stanza and are inherited by all releases in the DSF if they haven't defined their own.

These are defined as follows:

```toml
[settings]
 #...
 [settings.globalHooks]
    successCondition= "Initialized"
    deleteOnSuccess= true
    postInstall= "job.yaml"
```

```yaml
settings:
  #...
  globalHooks:
    successCondition: "Initialized"
    deleteOnSuccess: true
    postInstall: "job.yaml"
    #...
```

**App-specific** hooks:

Each app (release) can define its own hooks which **override any global ones**.

These are defined as follows:

```toml
[apps]
  [apps.argo]
    namespace = "production" # maps to the namespace as defined in namespaces above
    enabled = true # change to false if you want to delete this app release [default = false]
    chart = "argo/argo" # changing the chart name means delete and recreate this release
    version = "0.6.4" # chart version
    [apps.argo.hooks]
    successCondition= "Complete"
    successTimeout= "90s"
    deleteOnSuccess= true
    preInstall="job.yaml"
    preInstall="https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.crds.yaml"
    postInstall="https://raw.githubusercontent.com/jetstack/cert-manager/release-0.14/deploy/manifests/00-crds.yaml"
    preUpgrade="job.yaml"
    postUpgrade="job.yaml"
    preDelete="job.yaml"
    postDelete="job.yaml"
```

```yaml
apps:
    argo:
      namespace: "staging" # maps to the namespace as defined in namespaces above
      enabled: true # change to false if you want to delete this app release empty: false:
      chart: "argo/argo" # changing the chart name means delete and recreate this chart
      version: "0.6.5" # chart version
      hooks:
        successCondition: "Complete"
        successTimeout: "90s"
        deleteOnSuccess: true
        preInstall: "job.yaml"
        preInstall: "https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.crds.yaml"
        postInstall: "https://raw.githubusercontent.com/jetstack/cert-manager/release-0.14/deploy/manifests/00-crds.yaml"
        postInstall: "job.yaml"
        preUpgrade: "job.yaml"
        postUpgrade: "job.yaml"
        preDelete: "job.yaml"
        postDelete: "job.yaml"
```

## Enforcing hook manifests deletion on all apps

You can do that by setting `deleteOnSuccess` to true in the `globalHooks` stanza under `settings`. If you need to make an exception for some app, you can set it to `false` in the `hooks` stanza of this app. This overrides the global hooks.

## Expanding variables in hook manifests

You can expand variables/parameters in the hook manifests at run time in one of the following ways:

- use env variables (defined as `$MY_VAR` in your manifests) and run helmsman with `--subst-env-values`. Environment variables can be read from the environment or you can [load them from an env file](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/apps/secrets.md#passing-secrets-from-env-files)

- use [AWS SSM parameters](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html) (defined as `{{ssm: MY_PARAM }}` in your manifests) and run helmsman with `--subst-ssm-values`.

- Pass encrypted values with [hiera-eyaml](https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/settings/use-hiera-eyaml-as-secrets-encryption.md)

## Limitations

- You can only have one manifest file per lifecycle.
- If you have multiple k8s resources in your hook manifest file, `successCondition` may not work.
- pre/postDelete hooks are not respected before/after deleting untracked releases (releases which are no longer defined in your desired state file).
