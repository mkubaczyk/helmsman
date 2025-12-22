# Cluster connection -- Using the current kube context

Helmsman can use the current configured kube context. In this case, the `kubeContext` field in the `settings` stanza needs to be left empty. If no other `settings` fields are needed, you can delete the whole `settings` stanza.


If you want Helmsman to create the kube context for you, see [this guide](creating_kube_context_with_certs.md) for more details on creating a context with certs or [here](creating_kube_context_with_token.md) for details on creating context with bearer token.