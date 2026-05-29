---
version: 3.3.0
---

# Define resource quotas for namespaces

You can define namespaces to be used in your cluster. If they don't exist, Helmsman will create them for you. You can also define how much resource limits to set for each namespace.

You can read more about the `Quotas` specification [here](https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/quota-memory-cpu-namespace/#create-a-resourcequota).

```toml
#...
[namespaces]

  [namespaces.helmsman1]

    [namespaces.helmsman1.quotas]
    "limits.cpu" = "10"
    "limits.memory" = "30Gi"
    pods = "25"
    "requests.cpu" = "10"
    "requests.memory" = "30Gi"

      [[namespaces.helmsman1.quotas.customQuotas]]
      name = "requests.nvidia.com/gpu"
      value = "2"
#...
```

```yaml
namespaces:
  helmsman1:
    quotas:
      limits.cpu: '10'
      limits.memory: '30Gi'
      pods: '25'
      requests.cpu: '10'
      requests.memory: '30Gi'
      customQuotas:
        - name: 'requests.nvidia.com/gpu'
          value: '2'
```

The example above will create one namespace - helmsman1 - with resource quotas defined for the helmsman1 namespace.