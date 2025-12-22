---
version: v4.0.0
---

# Helm repos

Helm does not add any repos by default. You must explicitly define all repos you want to use.

This example defines a custom repo:

```toml
[helmRepos]
  custom = "https://mycustomrepo.org"
```

```yaml
helmRepos:
  custom: "https://mycustomrepo.org"
```

You can name your repos however you like. The name is used to reference charts from that repo:

```toml
#...
[helmRepos]
  bitnami = "https://charts.bitnami.com/bitnami"
  grafana = "https://grafana.github.io/helm-charts"
#...
```

```yaml
# ...
helmRepos:
  bitnami: "https://charts.bitnami.com/bitnami"
  grafana: "https://grafana.github.io/helm-charts"
# ...
```

Then reference charts using `repoName/chartName` format:

```yaml
apps:
  my-app:
    chart: "bitnami/nginx"
    version: "15.0.0"
```
