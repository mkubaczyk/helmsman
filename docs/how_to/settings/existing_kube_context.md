# Cluster connection -- Using an existing kube context

Helmsman can use any predefined kube context in the environment. All you need to do is set the context name in the `settings` stanza.

```toml
[settings]
  kubeContext = "minikube"
```

```yaml
settings:
  kubeContext: "minikube"
```

In the examples above, Helmsman tries to set the kube context to `minikube`. If that fails, it will attempt to create that kube context. Creating kube context requires more information provided. See [this guide](creating_kube_context_with_certs.md) for more details on creating a context with certs or [here](creating_kube_context_with_token.md) for details on creating context with bearer token.
