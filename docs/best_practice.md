# Best Practice

When using Helmsman, we recommend the following best practices:

- Add useful metadata in your desired state files (DSFs) so that others (who have access to them) can understand what your DSF is for. We recommend the following metadata: organization, maintainer (name and email), and description/purpose.

- Define `context` (see [the DSF spec](desired_state_specification.md#context)) for each DSF. This helps prevent different DSFs from operating on each other's releases.

- Store your DSFs in git (or any other VCS) so that you have an audit trail of your deployments. You can also rollback to a previous state by going back to previous commits.
> Rollback can be more complex regarding application data.

- Do not store secrets in your DSFs! Use one of [the supported ways to pass secrets to your releases](how_to/apps/secrets.md).

- To protect against accidental operations, define certain namespaces (e.g, production) as protected namespaces and deploy your production-ready releases there.

- If you use multiple desired state files (DSFs) with the same cluster, make sure your namespace protection definitions are identical across all DSFs.

- When using multiple DSFs, make sure that apps managed in the same namespace are in one DSF. This avoids the need for defining the same namespace (with its settings) across multiple DSFs

- Don't maintain the same release in multiple DSFs.

- While the decision on how many DSFs to use and what each can contain is up to you and depends on your case, we recommend coming up with your own rules for how to split them. For example, you can have one for infra (3rd party tools), one for staging, and one for production apps.

