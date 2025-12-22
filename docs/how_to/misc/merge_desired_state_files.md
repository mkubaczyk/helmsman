# Supply multiple desired state files

Helmsman allows you to pass the `-f` flag multiple times to specify multiple desired state files
that should be merged. This allows us to do things like specify our non-environment-specific config in a `common.toml` file
and environment specific info in a `nonprod.toml` or `prod.toml` file. This process uses [this library](https://github.com/imdario/mergo)
to do the merging, and is subject to the limitations described there.

For example:

`common.toml`:

```toml
[metadata]
org = "Organization Name"
maintainer = "project-owners@example.com"
description = "Project charts"

[settings]
storageBackend = "secret"
...
```

`nonprod.toml`:

```toml
[settings]
kubeContext = "cluster-nonprod"

[apps]
  [apps.external-dns]
  valuesFiles = ["./external-dns/values.yaml", "./external-dns/nonprod.yaml"]

  [apps.cert-issuer]
  valuesFile = "./cert-issuer/nonprod.yaml"
...
```

One can then run the following to use the merged config of the above files, with later files override values of earlier ones:

```shell
helmsman -f common.toml -f nonprod.toml ...
```

## Distinguishing releases deployed from different Desired State Files

When using multiple DSFs -and since Helmsman doesn't maintain any external state-, it has been possible for operations from one DSF to cause problems to releases deployed by other DSFs. A typical example is that releases deployed by other DSFs are considered `untracked` and get scheduled for deleting. Workarounds existed (e.g. using the `--keep-untracked-releases`, `--target` and `--group` flags).

Starting from Helmsman v3.0.0-beta5, `context` is introduced to define the context in which a DSF is used. This context is used as the ID of that specific DSF and must be unique across the used DSFs. The context is then used to label the different releases to link them to the DSF they were first deployed from. These labels are then checked by Helmsman on each run to make sure operations are limited to releases from a specific context.

Here is how it is used:

`infra.yaml`:

```yaml
context: infra-apps
settings:
  kubeContext: "cluster"
  storageBackend: "secret"

namespaces:
  infra:
    protected: true

apps:
  external-dns:
    namespace: infra
    valuesFile: "./external-dns/values.yaml"
    ...

  cert-issuer:
    namespace: infra
    valuesFile: "./cert-issuer/nonprod.yaml"
    ...
...
```

`prod.yaml`:

```yaml
context: prod-apps
settings:
  kubeContext: "cluster"
  storageBackend: "secret"

namespaces:
  prod:
    protected: true

apps:
  my-prod-app:
    namespace: prod
    valuesFile: "./my-prod-app/values.yaml"
    ...
...
```

> If you need to migrate releases from one Helmsman's context to another, check this [guide](../apps/migrate_contexts.md).

### Limitations

* If no context is provided in DSF (or merged DSFs), `default` is applied as a default context. This means any set of DSFs that don't define custom contexts can still operate on each other's releases.

* When merging multiple DSFs, context from the firs DSF in the list gets overridden by the context in the last DSF.

* If multiple DSFs use the same context name, they will mess up each other's releases. You can use `--keep-untracked-releases` to avoid that. However, it is recommended to avoid having multiple DSFs using the same context name.
