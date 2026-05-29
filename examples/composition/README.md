# Composition

Desired state configuration can be split into multiplle files and applied with:

```sh
helmsman --apply -f main.yaml -f argo.yaml -f artifactory.yaml
```

or using a spec file:

```sh
helmsman --apply --spec spec.yaml
```
