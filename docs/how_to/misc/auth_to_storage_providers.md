# Authenticating to cloud storage providers

Helmsman can read files like certificates for connecting to the cluster from cloud storage providers: GCS, S3 and Azure blob storage. Below is the authentication requirement for each provider:

## AWS S3

You need to provide ALL the following AWS env variables:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_DEFAULT_REGION`

## Google GCS

You need to provide ONE of the following env variables:

- `GOOGLE_APPLICATION_CREDENTIALS` the absolute path to your Google cloud credentials.json file.
- Or, `GCLOUD_CREDENTIALS` the content of the credentials.json file.

If running inside GCP helmsman can use metadata server to use Service Account permissions.

## Microsoft Azure

You need to provide ALL of the following env variables:

- `AZURE_STORAGE_ACCOUNT`
- `AZURE_STORAGE_ACCESS_KEY`