---
version: v3.2.0
---

# Migrating releases from Helmsman context to another

The `context` stanza has been introduced in v3.0.0 to allow you to [distinguish releases managed by different Helmsman's files](misc/merge_desired_state_files.md#distinguishing-releases-deployed-from-different-desired-state-files). However, once a context is defined, it couldn't be modified. 

From v.3.2.0, you can migrate releases in a DSF to another context (or rename the context) using the `--migrate-context`. This option can be combined with any other Helmsman flags. Behind the scenes, it will just update the Helmsman labels that contain the context name on the release secrets/configmaps before proceeding with the regular execution.

# Remember
- It is safe to run the `--migrate-context` flag multiple times.
- It can be used in conjunction with other cmd flags.
- It will respect `--target` & `--group` flags if specified (i.e. context migration will only be applied to the selected releases).
- The flag introduces an extra operation done before any other operations are done. So to reduce execution time, don't use it when it's not needed.