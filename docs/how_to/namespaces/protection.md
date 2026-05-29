# Protecting namespaces

You can define namespaces to be used in your cluster. If they don't exist, Helmsman will create them for you.

You can also define certain namespaces to be protected using the `protected` field. Please check [this doc](../misc/protect_namespaces_and_releases.md) for details about what protection means and the difference between namespace-level and release-level protection.


```toml
#...

[namespaces]
[namespaces.staging]
[namespaces.production]
  protected = true

#...
```

```yaml

namespaces:
  staging:
  production:
    protected: true

```

The example above will create two namespaces; staging and production. Where Helmsman sees the production namespace as a protected namespace.
