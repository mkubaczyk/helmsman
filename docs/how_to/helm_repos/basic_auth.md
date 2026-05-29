# Using private helm repos with basic auth

Helmsman allows you to use any private helm repo hosting which supports basic auth (e.g. Artifactory).

For such repos, you need to add the basic auth information in the repo URL as in the example below:

> Be aware that some special characters in the username or password can make the URL invalid.

```toml
[helmRepos]
# PASS is an env var containing the password
myPrivateRepo = "https://user:$PASS@myprivaterepo.org"
```

```yaml
helmRepos:
  # PASS is an env var containing the password
  myPrivateRepo: "https://user:$PASS@myprivaterepo.org"
```

