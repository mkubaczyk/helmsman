# Cluster connection -- creating the kube context with certificates

Helmsman can create the kube context for you (i.e. establish connection to your cluster). This guide describe how its done with certificates. If you want to use bearer tokens, check [this guide](creating_kube_context_with_token.md).

Creating the context with certs, requires both the `settings` and `certificates` stanzas.

> If you use GCS, S3, or Azure blob storage for your certificates, you will need to provide means to authenticate to the respective cloud provider in the environment. See [authenticating to cloud storage providers](../misc/auth_to_storage_providers.md) for details.

```toml
[settings]
  kubeContext = "mycontext" # the name of the context to be created
  username = "admin" # the cluster user name
  password = "$K8S_PASSWORD" # the name of an environment variable containing the k8s password
  clusterURI = "${CLUSTER_URI}" # the name of an environment variable containing the cluster API endpoint
  #clusterURI = "https://192.168.99.100:8443" # equivalent to the above

[certificates]
  caClient = "gs://mybucket/client.crt" # GCS bucket path
  caCrt = "s3://mybucket/ca.crt" # S3 bucket path
  # caCrt = "az://myblobcontainer/ca.crt" # Azure blob object
  caKey = "../ca.key" # valid local file relative path to the DSF file
```

```yaml
settings:
  kubeContext: "mycontext" # the name of the context to be created
  username: "admin" # the cluster user name
  password: "$K8S_PASSWORD" # the name of an environment variable containing the k8s password
  clusterURI: "${CLUSTER_URI}" # the name of an environment variable containing the cluster API endpoint
  #clusterURI: "https://192.168.99.100:8443" # equivalent to the above

certificates:
  caClient: "gs://mybucket/client.crt" # GCS bucket path
  caCrt: "s3://mybucket/ca.crt" # S3 bucket path
  #caCrt: "az://myblobcontainer/ca.crt" # Azure blob object
  caKey: "../ca.key" # valid local file relative path to the DSF file

```
