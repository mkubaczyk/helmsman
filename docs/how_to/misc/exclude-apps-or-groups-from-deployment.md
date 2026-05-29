---
version: v3.10.0
---

# Exclude specific apps or groups from execution

Starting from v3.10.0, Helmsman allows you to pass the `--exclude-target` or `--exclude-group` flag multiple times 
to specify which apps or groups should be excluded from execution.
Thanks to this one can exclude specific applications among all defined for an environment.

## Example

Having environment defined with such apps:

example.yaml:

```yaml
# ...
apps:
    jenkins:
      namespace: "staging" # maps to the namespace as defined in namespaces above
      enabled: true # change to false if you want to delete this app release empty: false:
      chart: "jenkins/jenkins" # changing the chart name means delete and recreate this chart
      version: "2.15.1" # chart version

    artifactory:
      namespace: "production" # maps to the namespace as defined in namespaces above
      enabled: true # change to false if you want to delete this app release empty: false:
      chart: "jfrog/artifactory" # changing the chart name means delete and recreate this chart
      version: "11.4.2" # chart version
# ...
```

running Helmsman with `-f example.yaml` would result in checking state and invoking deployment for both jenkins and artifactory application.

With `--exclude-target` flag in command like

```shell
helmsman -f example.yaml --exclude-target artifactory ...
```

one can execute Helmsman's environment defined with example.yaml limited to only one `jenkins` app by excluding second one - `artifactory` from the execution.

Multiple applications can be excluded with `--exclude-target`, like

```shell
helmsman -f example.yaml --exclude-target artifactory --exclude-target jenkins ...
```

Same rules apply for `--exclude-groups`.
