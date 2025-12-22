# Define resource limits for namespaces

You can define namespaces to be used in your cluster. If they don't exist, Helmsman will create them for you. You can also define how much resource limits to set for each namespace.

You can read more about the `LimitRange` specification [here](https://docs.openshift.com/container-platform/3.11/dev_guide/compute_resources.html#dev-limit-ranges).

```toml
#...

[namespaces]
[namespaces.staging]
  [[namespaces.staging.limits]]
    type = "Container"
    [namespaces.staging.limits.default]
      cpu = "300m"
      memory = "200Mi"
    [namespaces.staging.limits.defaultRequest]
      cpu = "200m"
      memory = "100Mi"
  [[namespaces.staging.limits]]
    type = "Pod"
    [namespaces.staging.limits.max]
      memory = "300Mi"
[namespaces.production]

#...
```

```yaml

namespaces:
  staging:
    limits:
      - type: Container
        default:
          cpu: "300m"
          memory: "200Mi"
        defaultRequest:
          cpu: "200m"
          memory: "100Mi"
      - type: Pod
        max:
          memory: "300Mi"
  production:

```

The example above will create two namespaces - staging and production - with resource limits defined for the staging namespace.
