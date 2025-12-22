# Namespace and Release Protection

Since helmsman is used with version controlled code and is often configured to be triggered as part of a CI pipeline, accidental mistakes could happen by the user (e.g, disabling a production application and taking out of service as a result of a mistaken change in the desired state file).

Helmsman provides a fine-grained mechanism to protect releases/namespaces from accidental desired state file changes.

## Protection definition

- When a release (application) is protected, it CANNOT:
    - deleted
    - upgraded
    - moved to another namespace

- A release CAN be moved into protection from a non-protected state.
- If a protected release need to be updated/changed or even deleted, this is possible, but the protection has to be removed first (i.e. remove the namespace/release from the protected state). This explained further below.

> A release is an instance (installation) of an application which has been packaged as a helm chart.

## Protection mechanism
Protection is supported in two forms:

- **Namespace-level Protection**: is defined at the namespace level. A namespace can be declaratively defined to be protected in the desired state file as in the example below:

```toml
[namespaces]
  [namespaces.staging]
  protected = false
  [namespaces.production]
  protected = true

```

- **Release-level Protection** is defined at the release level as in the example below:

```toml
[apps]

    [apps.jenkins]
    namespace = "staging"
    enabled = true
    chart = "jenkins/jenkins"
    version = "2.15.1"
    protected = true # defining this release to be protected.
```

```yaml
apps:

  jenkins:
    namespace: "staging"
    enabled: true
    chart: "jenkins/jenkins"
    version: "2.15.1"
    protected: true # defining this release to be protected.
```

> All releases in a protected namespace are automatically protected. Namespace protection has higher priority than the release-level protection.

## Important Notes

- You can combine both types of protection in your desired state file. The namespace-level protection always has a higher priority.
- Removing the protection from a namespace means all releases deployed in that namespace are no longer protected.
- We recommend using namespace-level protection for production namespace(s) and release-level protection for releases deployed in other namespaces.
- Release/namespace protection is only applied on single desired state files. It is your responsibility to make sure that multiple desired state files (if used) do not conflict with each other (e.g, one defines a particular namespace as protected and another defines it unprotected.) If you use multiple desired state files with the same cluster, please refer to [deployment strategies](../../deployment_strategies.md) and [best practice](../../best_practice.md) documentation.
