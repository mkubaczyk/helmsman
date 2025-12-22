# Using pre-configured helm repos

The primary use-case is if you have some helm repositories that require HTTP basic authentication and you don't want to store the password in the desired state file or as an environment variable. In this case you can execute the following sequence to have those repositories configured:

Set up the helmsman configuration:

```toml
preconfiguredHelmRepos = [ "myrepo1", "myrepo2" ]
```

```yaml
preconfiguredHelmRepos:
  - myrepo1
  - myrepo2
```

> In this case you will manually need to execute `helm repo add myrepo1 <URL> --username= --password=`
