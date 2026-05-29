---
version: v3.8.1
---

# Specification file

Starting from v3.8.0, Helmsman allows you to use Specification file passed with `--spec <file>` flag
in order to define multiple Desired State Files to be merged together.

An example Specification file `spec.yaml`:

```yaml
---
stateFiles:
  - path: examples/example.yaml
  - path: examples/minimal-example.yaml
  - path: examples/minimal-example.toml

```

This file can be then run with:

```shell
helmsman --spec spec.yaml ...
```

It takes the files from `stateFiles` list in the same order they are defined.
Then Helmsman will read each file one by one and merge the previous states with the current file it goes through.

One can take advantage of that and define the state of the environment starting with more general definitions and then reaching more specific cases in the end,
which would overwrite or extend things from previous files.
