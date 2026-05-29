---
version: v3.3.0-beta1
---

# Multiple value files

You can include multiple yaml value files to separate configuration for different environments.

> file paths can be a URL (e.g. to a public git repo) , cloud bucket, local absolute/relative file path.

```toml
...
[apps]

    [apps.jenkins-prod]
    description = "production jenkins"
    namespace = "production"
    enabled = true
    chart = "jenkins/jenkins"
    version = "2.15.1" # chart version
    valuesFiles = [
        "../my-jenkins-common-values.yaml",
        "../my-jenkins-production-values.yaml"
    ]

    # the jenkins release below is being tested in the staging namespace
    [apps.jenkins-test]
    description = "test release of jenkins, testing xyz feature"
    namespace = "staging"
    enabled = true
    chart = "jenkins/jenkins"
    version = "2.15.1" # chart version
    valuesFiles = [
        "../my-jenkins-common-values.yaml",
        "../my-jenkins-testing-values.yaml"
    ]

#...
```

```yaml
# ...
apps:

  jenkins-prod:
    description: "production jenkins"
    namespace: "production"
    enabled: true
    chart: "jenkins/jenkins"
    version: "2.15.1" # chart version
    valuesFiles:
      - "../my-jenkins-common-values.yaml"
      - "../my-jenkins-production-values.yaml"

  # the jenkins release below is being tested in the staging namespace
  jenkins-test:
    name: "jenkins-test" # should be unique across all apps
    description: "test release of jenkins, testing xyz feature"
    namespace: "staging"
    enabled: true
    chart: "jenkins/jenkins"
    version: "2.15.1" # chart version
    valuesFiles:
      - "../my-jenkins-common-values.yaml"
      - "../my-jenkins-testing-values.yaml"
# ...
```
