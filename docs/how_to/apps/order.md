---
version: v3.0.0-beta5
---

# Using the priority key for Apps

The `priority` flag in Apps definition allows you to define the order at which apps operations will be applied. This is useful if you have dependencies between your apps/services.

Priority is an optional flag and has a default value of 0 (zero). If set, it can only use a negative value. The lower the value, the higher the priority.

If you want your apps to be deleted in the reverse order as they where created, you can also use the optional `Settings` flag `reverseDelete`, to achieve this, set it to `true`

## Example

```toml
[metadata]
  org = "example.com"
  description = "example Desired State File for demo purposes."

[settings]
  kubeContext = "minikube"
  reverseDelete = false # Optional flag to reverse the priorities when deleting

[namespaces]
  [namespaces.staging]
    protected = false
  [namespaces.production]
    protected = true

[helmRepos]
jenkins = https://charts.jenkins.io
center = https://repo.chartcenter.io

[apps]
  [apps.jenkins]
    description = "jenkins"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "jenkins/jenkins" # changing the chart name means delete and recreate this chart
    version = "2.15.1" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    priority= -2

  [apps.jenkins1]
    description = "jenkins"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "jenkins/jenkins" # changing the chart name means delete and recreate this chart
    version = "2.15.1" # chart version
    valuesFile = "" # leaving it empty uses the default chart values

  [apps.jenkins2]
    description = "jenkins"
    namespace = "production" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "jenkins/jenkins" # changing the chart name means delete and recreate this chart
    version = "2.15.1" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    priority= -3

  [apps.artifactory]
    description = "artifactory"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "jfrog/artifactory" # changing the chart name means delete and recreate this chart
    version = "11.4.2" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    priority= -2
```

The above example will generate the following plan:

```console
DECISION: release [ jenkins2 ] is not present in the current k8s context. Will install it in namespace [[ production ]] -- priority: -3
DECISION: release [ jenkins ] is not present in the current k8s context. Will install it in namespace [[ staging ]] -- priority: -2
DECISION: release [ artifactory ] is not present in the current k8s context. Will install it in namespace [[ staging ]] -- priority: -2
DECISION: release [ jenkins1 ] is not present in the current k8s context. Will install it in namespace [[ staging ]] -- priority: 0
```
