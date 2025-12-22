# Label & annotate namespaces

You can define namespaces to be used in your cluster. If they don't exist, Helmsman will create them for you. You can also set some labels to apply for those namespaces.

```toml
#...

[namespaces]
[namespaces.staging]
  [namespaces.staging.labels]
    env = "staging"
[namespaces.production]
  [namespaces.production.annotations]
    "iam.amazonaws.com/role" = "dynamodb-reader"


#...
```

```yaml

namespaces:
  staging:
    labels:
      env: "staging"
  production:
    annotations:
      iam.amazonaws.com/role: "dynamodb-reader"

```

The above examples create two namespaces; staging and production. The staging namespace has one label `env`= `staging` while the production namespace has one annotation `iam.amazonaws.com/role`=`dynamodb-reader`.
