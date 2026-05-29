# Limit execution to explicitly defined group of apps

Helmsman allows you to pass the `--group` flag to specify group of apps
the execution of Helmsman deployment will be limited to.
Thanks to this one can deploy specific applications among all defined for an environment.

## Example

Having environment defined with such apps:

example.yaml:

```yaml
# ...
apps:
    jenkins:
      namespace: "staging" # maps to the namespace as defined in namespaces above
      group: "critical" # group name
      enabled: true # change to false if you want to delete this app release empty: false:
      chart: "jenkins/jenkins" # changing the chart name means delete and recreate this chart
      version: "2.15.1" # chart version

    artifactory:
      namespace: "production" # maps to the namespace as defined in namespaces above
      group: "sidecar" # group name
      enabled: true # change to false if you want to delete this app release empty: false:
      chart: "jfrog/artifactory" # changing the chart name means delete and recreate this chart
      version: "11.4.2" # chart version
# ...
```

running Helmsman with `-f example.yaml` would result in checking state and invoking deployment for both jenkins and artifactory application.

With `--group` flag in command like

```shell
helmsman -f example.yaml --group critical ...
```

one can execute Helmsman's environment defined with example.yaml limited to only one `jenkins` app, since its group is `critical`.
Others are ignored until the flag is defined.

Multiple applications can be set with `--group`, like

```shell
helmsman -f example.yaml --group critical --group sidecar ...
```
