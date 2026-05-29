# Cluster connection -- creating the kube context with bearer tokens

Helmsman can create the kube context for you (i.e. establish connection to your cluster). This guide describe how its done with bearer tokens. If you want to use certificates, check [this guide](creating_kube_context_with_certs.md).

All you need to do is set `bearerToken` to true and set the `clusterURI` to point to your cluster API endpoint in the `settings` stanza.

> Note: Helmsman and therefore helm will only be able to do what the kubernetes service account (from which the token is taken) allows.

By default, Helmsman will look for a token in `/var/run/secrets/kubernetes.io/serviceaccount/token`. If you have the token else where, you can specify its path with `bearerTokenPath`.

```toml
[settings]
  kubeContext = "test" # the name of the context to be created
  bearerToken = true
  clusterURI = "https://kubernetes.default"
  # bearerTokenPath = "/path/to/custom/bearer/token/file"
```

```yaml
settings:
  kubeContext: "test" # the name of the context to be created
  bearerToken: true
  clusterURI: "https://kubernetes.default"
  # bearerTokenPath: "/path/to/custom/bearer/token/file"
```
