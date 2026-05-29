# DRY-ed Example

## Execution Plan

To deploy the DRY-ed example, the app templates placed in [config](config) and [apps](apps) need to be used. To do so - please run:

```bash
helmsman -apply -f config/helmsman.yaml
```

> **Tip**: That kind of DRY-ed code can be achieved only using YAML desired state files.
