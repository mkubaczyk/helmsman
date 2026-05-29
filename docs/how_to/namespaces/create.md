# Create namespaces

You can define namespaces to be used in your cluster. If they don't exist, Helmsman will create them for you.

```toml
#...

[namespaces]
[namespaces.staging]
[namespaces.production]

#...
```

```yaml

namespaces:
  staging:
  production:

```

The example above will create two namespaces; staging and production.