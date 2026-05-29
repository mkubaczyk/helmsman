# Slack notifications from Helmsman

Helmsman can send slack notifications to a channel of your choice. To enable the notifications, simply add a `slack webhook` in the `settings` section of your desired state file. The webhook URL can be passed directly or from an environment variable.

```toml
[settings]
...
slackWebhook = $MY_SLACK_WEBHOOK
```

```yaml
settings:
  # ...
  slackWebhook : "$MY_SLACK_WEBHOOK"
  # ...
```

## Getting a Slack Webhook URL

Follow the [slack guide](https://api.slack.com/incoming-webhooks) for generating a webhook URL.
