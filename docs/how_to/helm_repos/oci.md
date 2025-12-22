---
version: v4.0.0
---

# Using OCI registries for helm charts

Helmsman allows you to use charts stored in OCI registries. OCI support is built into Helm 3.8+ and Helm 4.

If the registry requires authentication, you must login before running Helmsman:

```sh
helm registry login -u myuser my-registry.local
```

> **Note for Helm 4**: Use only the domain name when logging in (e.g., `my-registry.local`), not a full URL path.

```toml
[apps]
  [apps.my-app]
    chart = "oci://my-registry.local/my-chart"
    version = "1.0.0"
```

```yaml
#...
apps:
  my-app:
    chart: oci://my-registry.local/my-chart
    version: 1.0.0
```

For more information, read the [helm registries documentation](https://helm.sh/docs/topics/registries/).
