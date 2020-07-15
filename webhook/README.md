# Webhook Null Solver

This Null Solver is necessary because cert-manager tries to trigger a call to some DNS Api that assures that the challenge address has been added, before releasing the Let's Encrypt calls.

Because we don't have a programable DNS API, but instead we've the acme-solver querying Kubernetes objects directly there's no sense to call someone to add those objects, so we need to discard the webhook calls.


