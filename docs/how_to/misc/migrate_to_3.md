---
version: v3.2.0
---

# Migrate from Helm2 (Helmsman v1.x) to Helm3 (Helmsman v3.x)

This guide describes the process of migrating your Helmsman managed releases from Helm v2 to v3.
Helmsman v3.x is Helm v3-compatible, while Helmsman v1.x is Helm v2-compatible.

The migration process can go as follows: 

## Migrate Helm v2 release state to Helm v3
- Go through the [Helm's v2 to v3 migration guide](https://helm.sh/docs/topics/v2_v3_migration/)
- Manually migrate your releases state/history with the [helm3 2to3 plugin](https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/) (e.g. usage helm3 2to3 convert <helm2-release-name>).

> At this stage, Helm v3 can see and operate on your releases, but Helmsman can't. This is because Helmsman defined labels haven't been migrated to the Helm v3 releases state.

## Migrate to Helmsman v3.x
- Download the latest Helmsman v3.x release from [Github releases](https://github.com/mkubaczyk/helmsman/releases) 
- Modify your Helmsman's TOML/YAML desired state files (DSFs) to be Helmsman v3.x compatible. You can check [v3.0.0 release notes](https://github.com/mkubaczyk/helmsman/blob/v3.0.0/release-notes.md) for what's changed and verify from the [Desired State Spec](https://github.com/mkubaczyk/helmsman/blob/master/docs/desired_state_specification.md) that your DSF files are compatible. 

> Everything related to Tiller will be removed from your DSFs at this stage.

- Helmsman v3.x introduces the [`context` stanza](../../desired_state_specification.md#context) to logically group different groups of applications managed by Helmsman. It is highly recommended that you define a unique `context` for each of your DSFs at this stage.

- In order for Helmsman to recognize the Helm v3 releases state, you need to use the `--migrate-context` flag on your first Helmsman v3.x run. This flag will recreate the Helmsman labels needed to recognize the Helm v3 releases state. 
   - Make sure that `helm` binary points to Helm v3 in your environment before you run this command. 
   - The `--migrate-context` flag is only available in Helmsman v.3.2.0 and above. If you are using Helmsman v3.0.x or v3.1.x, you can recreate the labels manually or with a script, but we highly recommend using Helmsman v3.2.0 or above.
   - You only need to use `--migrate-context` once. However, the flag is safe to use multiple times.

   ```bash
     $ helmsman --debug --migrate-context -f <path-to-your-dsf>.yaml
   ``` 
 
At this stage, you should have your release migrate to Helm v3 and Helmsman v3.x is able to see and manage those releases as usual. The next step would be to clean up your Helm v2 state and remove Tiller deployment which you can do with the [Helm v3 2to3 plugin](https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/). 