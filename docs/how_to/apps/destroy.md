---
version: v3.0.0-beta5
---

# Delete all deployed releases

Helmsman allows you to delete all the helm releases that were deployed by Helmsman from a given desired state.

The `--destroy` flag will remove all deployed releases from a given desired state file (DSF). Note that this does not currently delete the namespaces nor the Kubernetes contexts created.

This was originally requested in issue [#88](https://github.com/mkubaczyk/helmsman/issues/88).

