---
version: v3.4.0
---

# Using Environment Variables in Helmsman DSF and Helm values files

You can use environment variables in any Helmsman desired state file or helm values files or [lifecycle hooks](lifecycle_hooks.md) files (K8S manifests). Both formats `${MY_VAR}` and `$MY_VAR` are accepted.

> To expand environment variables in helm values files and lifecycle hooks files, you have to enable the `--subst-env-values`.

## How does it work?

Helmsman will expand those variables at run time. For helm values files and Helmsman lifecycle hooks files, the variables are expanded into temporary files which are used during runtime and removed at the end of execution.

## Validating against unset env variables

By default, Helmsman will validate that your environment variables are set before using them. If they are unset, an error will be produced.
The validation will parse Helmsman DSF files and other files (values files, lifecycle hooks files) line-by-line. This maybe become slow if you have very large files.

## Skipping env variables validation

Validation of environment variables being set is skipped in the following cases:
- If `--skip-validation` flag is used, no env variables validation is performed on any file.
- If `--no-env-subst` flag is used, no env variables validation is performed on Helmsman desired state files.
- If `--subst-env-values` flag is NOT used, no env variables validation is performed on helm values files and lifecycle hooks files.

## Escaping the `$` sign

### In Helmsman desired state files

If you want to pass the `$` as is, you can escape it like so: `$$` 

### In Helm values files and lifecycle hooks files

If you don't enable `--subst-env-values`, the `$` is passed as is without the need to escape it. However, if you enable `--subst-env-values` and want to pass the `$` as is, you have to escape it like so `$$`