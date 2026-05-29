---
version: v3.2.0
---

# Override Helmsman context name from CMD flags using `--context-override`

There are two main use cases for this flag:

1. To speed up Helmsman's execution when you have too many release (see [issue #418](https://github.com/mkubaczyk/helmsman/issues/418))
   This flag works by skipping the search for the context information which Helmsman adds in the form of labels to the helm release state (secrets/configmaps). 
   
   > Use this option with caution. You must be sure that this won't cause conflicts.

2. [Not recommended] If ,for whatever reason, you want to temporarily override the context defined on the release state (in labels on secrets/configmaps) with something. **Use [`--migrate-context`](migrate_contexts.md) instead to permanently rename your context**