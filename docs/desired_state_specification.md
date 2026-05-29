# Helmsman desired state specification

This document describes the specification for how to write your Helm charts' desired state file. This can be either a [Toml](https://github.com/toml-lang/toml) or [Yaml](http://yaml.org/) formatted file. The desired state file consists of:

- [Metadata](#metadata) [Optional] -- metadata for any human reader of the desired state file.
- [Certificates](#certificates) [Optional] -- only needed when you want Helmsman to connect kubectl to your cluster for you.
- [Context](#context) [optional] -- define the context in which a DSF is used.
- [Settings](#settings) [Optional] -- data about your k8s cluster and how to deploy Helm on it if needed.
- [Namespaces](#namespaces) -- defines the namespaces where you want your Helm charts to be deployed.
- [Helm Repos](#helm-repos) [Optional] -- defines the repos where you want to get Helm charts from.
- [Apps](#apps) -- defines the applications/charts you want to manage in your cluster.

> You can use environment variables in the desired state files. The environment variable name should start with "$", or encapsulated in "${", "}". "$" characters can be escaped like "$$".

> You can also use environment variables in your helm values/secrets files.

## Metadata

Optional : Yes.

Synopsis: Metadata is used for the human reader of the desired state file. While it is optional, we recommend having a maintainer and scope/cluster metadata.

Options:

- you can define any key/value pairs.

Example:

```toml
[metadata]
scope = "cluster foo"
maintainer = "k8s-admin"
```

```yaml
metadata:
  scope: "cluster foo"
  maintainer: "k8s-admin"
```

## Certificates

Optional : Yes, only needed if you want Helmsman to connect kubectl to your cluster for you.

Synopsis: defines where to find the certificates needed for connecting kubectl to a k8s cluster. If connection settings (username/password/clusterAPI) are provided in the Settings section below, then you need **AT LEAST** to provide caCrt and caKey. You can optionally provide a client certificate (caClient) depending on your cluster connection setup.

Options:

- **caCrt** : a valid URL/S3/GCS/Azure bucket or local relative file path to a certificate file.
- **caKey** : a valid URL/S3/GCS/Azure bucket or local relative file path to a client key file.
- **caClient**: a valid URL/S3/GCS/Azure bucket or local relative file path to a client certificate file.

> bucket format is: [s3 or gs or az]://bucket-name/dir1/dir2/.../file.extension

Example:

```toml
[certificates]
caCrt = "s3://myS3bucket/mydir/ca.crt"
caKey = "gs://myGCSbucket/ca.key"
#caKey = "az://myAzureContainer/ca.key
caClient ="../path/to/my/local/client-certificate.crt"
#caClient = "$CA_CLIENT"
```

```yaml
certificates:
  caCrt: "s3://myS3bucket/mydir/ca.crt"
  caKey: "gs://myGCSbucket/ca.key"
  #caKey:  "az://myAzureContainer/ca.key
  caClient: "../path/to/my/local/client-certificate.crt"
  #caClient: "$CA_CLIENT"
```

## Context

Optional : Yes.

Synopsis: defines the context in which a DSF is used. This context is used as the ID of that specific DSF and must be unique across the used DSFs. If not defined, `default` is used. Check [here](how_to/misc/merge_desired_state_files.md) for more details on the limitations.

> Renaming the Helmsman context can be done from v3.2.0 using the `--migrate-context` flag. Check [this guide](how_to/apps/migrate_contexts.md) for details.

```yaml
context: prod-apps
# ...
```

## Settings

Optional : Yes.

Synopsis: provides settings for connecting to your k8s cluster.

> If you don't provide the `settings` stanza, helmsman would use your current kube context.

Options:

- **kubeContext** : the kube context you want Helmsman to use or create. Helmsman will try connect to this context first, if it does not exist, it will try to create it (i.e. connect to a k8s cluster) using the options below.

The following options can be skipped if your kubectl context is already created and you don't want Helmsman to connect kubectl to your cluster for you.

- **username**   : the username to be used for kubectl credentials.
- **password**   : an environment variable name (starting with `$`) where your password is stored. Get the password from your k8s admin or consult k8s docs on how to get/set it.
- **clusterURI** : the URI for your cluster API or the name of an environment variable (starting with `$`) containing the URI.
- **bearerToken**: whether you want helmsman to connect to the cluster using a bearer token. Default is `false`
- **bearerTokenPath**: optional. If bearer token is used, you can specify a custom location (URL, cloud bucket, local file path) for the token file.
- **storageBackend** : by default Helm v3 stores release information in secrets, using secrets for storage is recommended for security.
- **slackWebhook** : a [Slack](http://slack.com) Webhook URL to receive Helmsman notifications. This can be passed directly or in an environment variable.
- **msTeamsWebhook** : a [Microsoft Teams](https://www.microsoft.com/pl-pl/microsoft-teams/group-chat-software) Webhook URL to receive Helmsman notifications. This can be passed directly or in an environment variable.
- **reverseDelete** : if set to `true` it will reverse the priority order whilst deleting.
- **namespaceLabelsAuthoritative** : if set to `true` it will remove all the namespace's labels that are not defined in DSL for particular namespace
- **vaultEnabled**: if set to `true` it will use [helm-vault](https://github.com/Just-Insane/helm-vault) to decrypt secret files instead of using default helm-secrets
- **vaultDeliminator**: secret deliminator used when parsing value files. See [helm-vault](https://github.com/Just-Insane/helm-vault#available-flags) docs
- **vaultPath**: vault path (secret mount location in Vault). See [helm-vault](https://github.com/Just-Insane/helm-vault#available-flags) docs
- **vaultMountPoint**: vault secret engine mount point. See [helm-vault](https://github.com/Just-Insane/helm-vault#available-flags) docs
- **vaultTemplate**: substring with path to vault key instead of deliminator. See [helm-vault](https://github.com/Just-Insane/helm-vault#available-flags) docs
- **vaultKvVersion**: version of the KV secrets engine in Vault. See [helm-vault](https://github.com/Just-Insane/helm-vault#available-flags) docs
- **vaultEnvironment**: environment that secrets should be stored under. See [helm-vault](https://github.com/Just-Insane/helm-vault#available-flags) docs
- **eyamlEnabled** : if set to `true` it will use [hiera-eyaml](https://github.com/voxpupuli/hiera-eyaml) to decrypt secret files instead of using default helm-secrets based on sops
- **eyamlPrivateKeyPath** : if set with path to the eyaml private key file, it will use it instead of looking for default one in ./keys directory relative to where Helmsman were run. It needs to be defined in conjunction with eyamlPublicKeyPath.
- **eyamlPublicKeyPath** : if set with path to the eyaml public key file, it will use it instead of looking for default one in ./keys directory relative to where Helmsman were run. It needs to be defined in conjunction with eyamlPrivateKeyPath.
- **globalHooks** : defines global lifecycle hooks to apply yaml manifest before and/or after different helmsman operations. Check [here](how_to/apps/lifecycle_hooks.md) for more details.
- **globalMaxHistory** : defines the **global** maximum number of helm revisions state (secrets/configmap) to keep. Releases can override this global value by setting `maxHistory`. If both are not set or are set to `0`, it is defaulted to 10.
- **skipIgnoredApps** : if set to true apps, that would normally be listed in the plan as `ignored`, will be skipped. They won't show up on the plan output and won't be considered in decisions. This is especially useful when using `-target` or `-group` flags with significant amount of apps where most of them show up as `ignored` in the plan output making it hard to read.
- **skipPendingApps** : if set to true apps that are in a pending (install/upgrade/rollback) state or being deleted, will be ignored, when set to false Helmsman will stop if apps are found in these states.

Example:

```toml
[settings]
kubeContext = "minikube"
# username = "admin"
# password = "$K8S_PASSWORD"
# clusterURI = "https://192.168.99.100:8443"
## clusterURI= "$K8S_URI"
# storageBackend = "secret"
# slackWebhook = $MY_SLACK_WEBHOOK
# msTeamsWebhook = $MY_MS_TEAMS_WEBHOOK
# reverseDelete = false
# eyamlEnabled = true
# eyamlPrivateKeyPath = "../keys/custom-key.pem"
# eyamlPublicKeyPath = "../keys/custom-key.pub"
# vaultEnabled = false
# [settings.globalHooks]
#   successCondition= "Complete"
#   deleteOnSuccess= true
#   postInstall= "job.yaml"
globalMaxHistory= 10
```

```yaml
settings:
  kubeContext: "minikube"
  #username: "admin"
  #password: "$K8S_PASSWORD"
  #clusterURI: "https://192.168.99.100:8443"
  ##clusterURI: "$K8S_URI"
  #storageBackend: "secret"
  #slackWebhook: "$MY_SLACK_WEBHOOK"
  #msTeamsWebhook: "$MY_MS_TEAMS_WEBHOOK"
  #reverseDelete: false
  # eyamlEnabled: true
  # eyamlPrivateKeyPath: ../keys/custom-key.pem
  # eyamlPublicKeyPath: ../keys/custom-key.pub
  # vaultEnabled: false
  # globalHooks:
  #   successCondition: "Complete"
  #   deleteOnSuccess: true
  #   preInstall: "job.yaml"
  globalMaxHistory: 10
```

## Namespaces

Optional : No.

Synopsis: defines the namespaces to be used/created in your k8s cluster and whether they are protected or not.  You can add as many namespaces as you need.
If a namespace does not already exist, Helmsman will create it.

Options:

- **protected** : defines if a namespace is protected (true or false). Default false.

> For the definition of what a protected namespace means, check the [protection guide](how_to/misc/protect_namespaces_and_releases.md)

- **labels** : defines labels to be added to the namespace, doesn't remove existing labels but updates them if the label key exists with any other different value. You can define any key/value pairs. Default is empty.

- **annotations** : defines annotations to be added to the namespace. It behaves the same way as the labels option.

- **limits** : defines a [LimitRange](https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/memory-default-namespace/) to be configured on the namespace

Example:

```toml
[namespaces]
[namespaces.staging]
[namespaces.dev]
protected = false
[namespaces.production]
protected = true
[namespaces.production.labels]
env = "prod"
[namespaces.production.annotations]
iam.amazonaws.com/role = "dynamodb-reader"
[[namespaces.production.limits]]
type = "Container"
[namespaces.production.limits.default]
cpu = "300m"
memory = "200Mi"
[namespaces.production.limits.defaultRequest]
cpu = "200m"
memory = "100Mi"
[[namespaces.production.limits]]
type = "Pod"
[namespaces.production.limits.max]
memory = "300Mi"
```

```yaml
namespaces:
  staging:
  dev:
    protected: false
  production:
    protected: true
    limits:
      - type: Container
        default:
          cpu: "300m"
          memory: "200Mi"
        defaultRequest:
          cpu: "200m"
          memory: "100Mi"
      - type: Pod
        max:
          memory: "300Mi"
    labels:
      env: "prod"
    annotations:
      iam.amazonaws.com/role: "dynamodb-reader"
```

## Helm Repos

Optional : Yes.

Synopsis: defines the Helm repos where your charts can be found. You can add as many repos as you need. Public repos can be added without any additional setup. Private repos require authentication.

> AWS S3 and Google GCS buckets can be used for private repos (using the [Helm S3](https://github.com/hypnoglow/helm-s3) and [Helm GCS](https://github.com/nouney/helm-gcs) plugins). You can also use private repos with basic auth and pre-configured helm repos.

Authenticating to private cloud helm repos:

- **For S3 repos**: you need to have valid AWS access keys in your environment variables. See [here](https://github.com/hypnoglow/helm-s3#note-on-aws-authentication) for more details.
- **For GCS repos**: check [here](https://www.terraform.io/docs/providers/google/index.html#authentication-json-file) for getting the required authentication file. Once you have the file, you have two options, either:
  - set `GOOGLE\_APPLICATION\_CREDENTIALS` environment variable to contain the absolute path to your Google cloud credentials.json file.
  - Or, set `GCLOUD\_CREDENTIALS` environment variable to contain the content of the credentials.json file.

> You can also provide basic auth to access private repos that support basic auth. See the example below.

Options:

- you can define any key/value pair where the key is the repo name and value is a valid URI for the repo. Basic auth info can be added in the repo URL as in the example below.

Example:

```toml
[helmRepos]
autoscaler = https://kubernetes.github.io/autoscaler
grafana = https://grafana.github.io/helm-charts
myS3repo = "s3://my-S3-private-repo/charts"
myGCSrepo = "gs://my-GCS-private-repo/charts"
myPrivateRepo = "https://user:$TOP_SECRET_PASSWORD@mycustomprivaterepo.org"
```

```yaml
helmRepos:
  autoscaler: https://kubernetes.github.io/autoscaler
  grafana: https://grafana.github.io/helm-charts
  myS3repo: "s3://my-S3-private-repo/charts"
  myGCSrepo: "gs://my-GCS-private-repo/charts"
  myPrivateRepo: "https://user:$TOP_SECRET_PASSWORD@mycustomprivaterepo.org"
```

## Preconfigured Helm Repos

Optional : Yes.

Synopsis: defines the list of helm repositories that the helmsman will consider already preconfigured and thus will not try to overwrite it's configuration.

The primary use-case is if you have some helm repositories that require HTTP basic authentication and you don't want to store the password in the desired state file or as an environment variable. In this case you can execute the following sequence to have those repositories configured:

> In this case you will need to execute `helm repo add myrepo1 <URL> --username= --password=` manually first.

Set up the helmsman configuration:

```toml
preconfiguredHelmRepos = [ "myrepo1", "myrepo2" ]
```

```yaml
preconfiguredHelmRepos:
- myrepo1
- myrepo2
```

## AppsTemplates

> This feature is only for YAML.

Optional : Yes.

Synopsis: allows for YAML (TOML has no variable reference support) object creation, that is ignored by state file importer, but can be used as a reference with YAML anchors to not repeat yourself. Read [this](https://blog.daemonl.com/2016/02/yaml.html) example about YAML anchors.

Examples:

```yaml
helmRepos:
  jenkins: https://charts.jenkins.io

appsTemplates:

  default: &template
    valuesFile: ""
    test: true
    protected: false
    wait: true
    enabled: true

  custom: &template_custom
    valuesFile: ""
    test: true
    protected: false
    wait: false
    enabled: true

apps:
  jenkins:
    <<: *template
    name: "jenkins-stage"
    namespace: "staging"
    chart: "jenkins/jenkins"
    version: "0.9.2"
    priority: -3

  jenkins2:
    <<: *template_custom
    name: "jenkins-prod"
    namespace: "production"
    chart: "jenkins/jenkins"
    version: "0.9.0"
    priority: -2

```

## Apps

Optional : Yes.

Synopsis: defines the releases (instances of Helm charts) you would like to manage in your k8s cluster.

Releases must have unique names which are defined under `apps`. Example: in `[apps.jenkins]`, the release name will be `jenkins` and it should be unique within the DSF.

Options:

**Required**

- **namespace**     : the namespace where the release should be deployed. The namespace should map to one of the ones defined in [namespaces](#namespaces).
- **enabled**       : describes the required state of the release (true for enabled, false for disabled). Once a release is deployed, you can change it to false if you want to delete this release [default is false].
- **chart**         : the chart name. It should contain the repo name as well. Example: repoName/chartName. Changing the chart name means delete and reinstall this release using the new Chart.
- **version**       : the chart version.

**Optional**

- **group**         : group name this apps belongs to. It has no effect until Helmsman's flag `-group` is passed. Check this [doc](how_to/misc/limit-deployment-to-specific-group-of-apps.md) for more details.
- **description**   : a release metadata for human readers.
- **valuesFile**    : a valid path (URL, cloud bucket, local absolute/relative file path) to custom Helm values.yaml file. File extension must be `yaml`. Cannot be used with valuesFiles together. Leaving it empty uses the default chart values.
- **valuesFiles**   : array of valid paths (URL, cloud bucket, local absolute/relative file path) to custom Helm values.yaml file. File extension must be `yaml`. Cannot be used with valuesFile together. Leaving it empty uses the default chart values.
  > The values file(s) path is resolved when the DSF yaml/toml file is loaded, relative to the path that the dsf was loaded from.
- **secretsFile**   : a valid path (URL, cloud bucket, local absolute/relative file path) to custom Helm secrets.yaml file. File extension must be `yaml`. Cannot be used with secretsFiles together. Leaving it empty uses the default chart secrets.
- **secretsFiles**  : array of valid paths (URL, cloud bucket, local absolute/relative file path) to custom Helm secrets.yaml file. File extension must be `yaml`. Cannot be used with secretsFile together. Leaving it empty uses the default chart secrets.
  > The secrets file(s) path is resolved when the DSF yaml/toml file is loaded, relative to the path that the dsf was loaded from.
  > To use the secrets files you must have the helm-secrets plugin
- **test**          : defines whether to run the chart tests whenever the release is installed. Default is false.
- **protected**     : defines if the release should be protected against changes. Namespace-level protection has higher priority than this flag. Check the [protection guide](how_to/misc/protect_namespaces_and_releases.md) for more details. Default is false.
- **wait**          : defines whether Helmsman should block execution until all k8s resources are in a ready state. Default is false.
- **timeout**       : helm timeout in seconds. Default 300 seconds.
- **noHooks**       : helm noHooks option. If true, it will disable pre/post upgrade hooks. Default is false.
- **priority**      : defines the priority of applying operations on this release. Only negative values allowed and the lower the value, the higher the priority. Default priority is 0. Apps with equal priorities will be applied in the order they were added in your state file (DSF).
- **set**           : is used to override certain values from values.yaml with values from environment variables or directly provided in the Desired State File. This is particularly useful for passing secrets to charts. If an environment variable with the same name as the provided value exists, the environment variable value will be used, otherwise, the provided value will be used as-is. The TOML stanza for this is `[apps.<app_name>.set]`
- **setString**     : is used to override String values from values.yaml or chart's defaults. This uses the `--set-string` flag in helm. This option is useful for image tags and the like. The TOML stanza for this is `[apps.<app_name>.setString]`
- **setFile**       : is used to override values from values.yaml or chart's defaults from provided file. This uses the `--set-file` flag in helm. This option is useful for embedding file contents in the values. The TOML stanza for this is `[apps.<app_name>.setFile]`
  > set, setString and setFile can't take nested elements. If you need to provide nested values, you can combine them in one line with dots e.g. `TOML: "image.tag"=some\_value` `YAML: "image.tag": some\_value`
- **helmFlags**     : array of `helm upgrade` flags, is used to pass flags to helm install/upgrade commands. **These flags are not passed to helm diff**. For setting values, use **set**, **setString** or **setFile** instead.
- **helmDiffFlags** : array of `helm diff upgrade` flags, is used to pass flags to helm diff upgrade commands. **These flags are not passed to helm during upgrade**. For setting values, use **set**, **setString** or **setFile** instead.
    > helmDiffFlags can be useful for example if you need to use the `--disable-openapi-validation` flag, in that case you would need to set it both in helmFlags and helmDiffFlags
- **hooks** : defines global lifecycle hooks to apply yaml manifest before and/or after different helmsman operations. Check [here](how_to/apps/lifecycle_hooks.md) for more details. Unset hooks for a release are inherited from `globalHooks` in the [settings](#Settings) stanza.
- **maxHistory**    : defines the maximum number of helm revisions state (secrets/configmap) to keep. If unset, it will inherit the value of `settings.globalMaxHistory`, if that's also unset, it defaults to 10.
- **postRenderer**  : for Helm 3, the path to an executable to be used for post rendering. For Helm 4, this must be a plugin name (post-renderers are now plugins in Helm 4). Requires helm-diff v3.1.2+.

Example:

> Whitespace does not matter in TOML files. You could use whatever indentation style you prefer for readability.

```toml
[apps]

  [apps.jenkins]
    name = "jenkins"
    description = "jenkins"
    namespace = "staging"
    enabled = true
    group = "critical"
    chart = "jenkins/jenkins"
    version = "0.9.0"
    valuesFile = ""
    test = true
    maxHistory = 4
    protected = false
    wait = true
    priority = -3
    helmFlags = [
      "--timeout=600s",
    ]
  [apps.jenkins.set]
    secret1="$SECRET_ENV_VAR1"
    secret2="SECRET_ENV_VAR2" # works with/without $ at the beginning
  [apps.jenkins.setString]
    longInt = "1234567890"
    "image.tag" = "1.0.0"
  [apps.jenkins.hooks]
    successCondition= "Complete"
    successTimeout= "90s"
    deleteOnSuccess= true
    postInstall="job.yaml"
    preInstall="https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.crds.yaml"
```

```yaml
apps:
  jenkins:
    name: "jenkins"
    description: "jenkins"
    namespace: "staging"
    enabled: true
    group: "critical"
    chart: "jenkins/jenkins"
    version: "0.9.0"
    valuesFile: ""
    test: true
    maxHistory: 4
    protected: false
    wait: true
    priority: -3
    helmFlags: [
      "--timeout=600s",
    ]
    set:
      secret1: "$SECRET_ENV_VAR1"
      secret2: "$SECRET_ENV_VAR2"
    setString:
      longInt: "1234567890"
      "image.tag": "1.0.0"
    hooks:
      successCondition: "Complete"
      successTimeout: "90s"
      deleteOnSuccess: true
      postInstall: "job.yaml"
      preInstall: "https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.crds.yaml"
```
