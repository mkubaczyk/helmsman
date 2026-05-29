# Running Helmsman inside your k8s cluster

Helmsman can be deployed inside your k8s cluster and can talk to the k8s API using a `bearer token`.

See [connecting to your cluster with bearer token](../settings/creating_kube_context_with_token.md) for more details.

Your desired state will look like:

```toml
[settings]
  kubeContext = "test" # the name of the context to be created
  bearerToken = true
  clusterURI = "https://kubernetes.default"
```

```yaml
settings:
  kubeContext: "test" # the name of the context to be created
  bearerToken: true
  clusterURI: "https://kubernetes.default"
```

To deploy Helmsman into a k8s cluster, few steps are needed:

> The steps below assume default namespace

1. Create a k8s service account

    ```shell
    kubectl create sa helmsman
    ```

2. Create a clusterrolebinding

    ```shell
    kubectl create clusterrolebinding helmsman-cluster-admin --clusterrole=cluster-admin --serviceaccount=default:helmsman
    ```

3. Deploy helmsman

    This command gives an interactive session:

    ```shell
    kubectl run helmsman --restart Never --image ghcr.io/mkubaczyk/helmsman:latest --serviceaccount=helmsman -- helmsman -f -- sleep 3600
    ```

    But you can also create a proper kubernetes deployment and mount a volume to it containing your desired state file(s).
