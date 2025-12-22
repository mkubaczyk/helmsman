# Using private helm repos in GCS

Helmsman allows you to use private charts from private repos. Currently only repos hosted in S3 or GCS buckets are supported for private repos.

You need to provide one of the following env variables:

- `GOOGLE_APPLICATION_CREDENTIALS` environment variable to contain the absolute path to your Google cloud credentials.json file.
- Or, `GCLOUD_CREDENTIALS` environment variable to contain the content of the credentials.json file.

If running inside GCP helmsman can use metadata server to use Service Account permissions.

Helmsman uses the [helm GCS](https://github.com/nouney/helm-gcs) plugin to work with GCS helm repos.

```toml
[helmRepos]
  gcsRepo = "gs://myrepobucket/charts"
```

```yaml
helmRepos:
  gcsRepo: "gs://myrepobucket/charts"
```

