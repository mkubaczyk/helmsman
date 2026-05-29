---
version: v3.0.0-beta5
---

# Secrets

## Passing secrets from env variables

Starting from v0.1.3, Helmsman allows you to pass secrets and other user input to helm charts from environment variables as follows:

```toml
# ...
[apps]

   [apps.jira]
    description = "jira"
    namespace = "staging"
    enabled = true
    chart = "myrepo/jira"
    version = "0.1.5"
    valuesFile = "applications/jira-values.yaml"
    test = true
    [apps.jira.set] # the format is [apps.<<release_name (as defined above)>>.set]
    db_username= "$JIRA_DB_USERNAME" # pass any number of key/value pairs where the key is the input expected by the helm charts and the value is an env variable name starting with $
    db_password= "$JIRA_DB_PASSWORD"
# ...
```

```yaml
# ...
apps:

  jira:
    description: "jira"
    namespace: "staging"
    enabled: true
    chart: "myrepo/jira"
    version: "0.1.5"
    valuesFile: "applications/jira-values.yaml"
    test: true
    set:
      db_username: "$JIRA_DB_USERNAME" # pass any number of key/value pairs where the key is the input expected by the helm charts and the value is an env variable name starting with $
      db_password: "$JIRA_DB_PASSWORD"
# ...

```

These input variables will be passed to the chart when it is deployed/upgraded using helm's `--set <<var_name>>=<<var_value_read_from_env_var>>`

## Passing secrets from env files

You can also keep these environment variables in files, by default Helmsman will load variables from a `.env` file but you can also specify files by using the `-e` option:

```shell
helmsman -e myVars
```

Below are some examples of valid env files

```shell
# I am a comment and that is OK
SOME_VAR=someval
FOO=BAR # comments at line end are OK too
export BAR=BAZ
```

Or you can do YAML(ish) style

```yaml
FOO: bar
BAR: baz
```

## Passing secrets using helm secrets plugin

You can also use the [helm secrets plugin](https://github.com/jkroepke/helm-secrets) to pass your secrets.

## Passing secrets using hiera eyaml

An alternative method is to use heira eyaml as described in [this guide](../settings/use-hiera-eyaml-as-secrets-encryption.md).

