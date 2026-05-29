# Microsoft Teams notifications from Helmsman

Helmsman can send MS Teams notifications to a channel of your choice. To enable the notifications, simply add a `msTeamsWebhook webhook` in the `settings` section of your desired state file. The webhook URL can be passed directly or from an environment variable.

```toml
[settings]
...
msTeamsWebhook = $MY_MS_TEAMS_WEBHOOK
```

```yaml
settings:
  # ...
  msTeamsWebhook : "$MY_MS_TEAMS_WEBHOOK"
  # ...
```

## Getting a MS Teams Webhook URL

Follow the [Microsoft Teams Guide](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook) for generating a webhook URL.
