---
version: v3.0.0-beta5
---

# Override defined namespaces from command line

If you use different release branches for your releasing/managing your applications in your k8s clusters, then you might want to use the same desired state but with different namespaces on each branch. Instead of duplicating the DSF in multiple branches and adjusting it, you can use the `--ns-override` command line flag when running helmsman.

This flag overrides all namespaces defined in your DSF with the single one you pass from command line.

## Example

`dsf.toml`:

```toml
[metadata]
org = "example.com"
description = "example Desired State File for demo purposes."


[settings]
kubeContext = "minikube"

[namespaces]
  [namespaces.staging]
  protected = false
  [namespaces.production]
  prtoected = true

[helmRepos]
jenkins = https://charts.jenkins.io
center = https://repo.chartcenter.io

[apps]

    [apps.jenkins]
    description = "jenkins"
    namespace = "production" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "jenkins/jenkins" # changing the chart name means delete and recreate this chart
    version = "2.15.1" # chart version
    valuesFile = "" # leaving it empty uses the default chart values

    [apps.artifactory]
    description = "artifactory"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "jfrog/artifactory" # changing the chart name means delete and recreate this chart
    version = "11.4.2" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
```

`dsf.yaml`:

```yaml
metadata:
  org: "example.com"
  description: "example Desired State File for demo purposes."


settings:
  kubeContext: "minikube"

namespaces:
  staging:
    protected: false
  production:
    protected: true

helmRepos:
  jenkins: https://charts.jenkins.io
  jfrog: https://charts.jfrog.io

apps:

  jenkins:
    description: "jenkins"
    namespace: "production" # maps to the namespace as defined in environments above
    enabled: true # change to false if you want to delete this app release [empty: false]
    chart: "jenkins/jenkins" # changing the chart name means delete and recreate this chart
    version: "2.15.1" # chart version
    valuesFile: "" # leaving it empty uses the default chart values

  artifactory:
    description: "artifactory"
    namespace: "staging" # maps to the namespace as defined in environments above
    enabled: true # change to false if you want to delete this app release [empty: false]
    chart: "jfrog/artifactory" # changing the chart name means delete and recreate this chart
    version: "11.4.2" # chart version
    valuesFile: "" # leaving it empty uses the default chart values
```

In command line, we run :

```shell
helmsman -f dsf.toml --debug --ns-override testing
```

This will override the `staging` and `production` namespaces defined in `dsf.toml` :

```console
2018/03/31 17:38:12 INFO: Plan generated at: Sat Mar 31 2018 17:37:57
DECISION: release [ jenkins ] is not present in the current k8s context. Will install it in namespace [[ testing ]] -- priority: 0
DECISION: release [ artifactory ] is not present in the current k8s context. Will install it in namespace [[ testing ]] -- priority: 0
```
