# Deployment Guide

## Environment Variables Management

It's recommended to install [direnv](https://direnv.net/docs/installation.html) to help manage your environment variables out of a `.envrc` file. The benefit is when you CD out of the folder, the environment variables will be removed so they are just specific to that folder hierarchy. The directions below will assume you have this utility installed.

Once you have `direnv` installed, create a .envrc file in the root of your project:

```bash
# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

export PATH=$PATH:$(pwd)/bin
```

When you open up your terminal, you will be prompted with this message. You should type `direnv allow` to allow the system to load environment variables from the file.

```bash
direnv: error /Users/YOURPATH/.envrc is blocked. Run `direnv allow` to approve its content
```

## GCP Deployment

To deploy an Ambient application to Google Cloud Run:

- Install the [Google Cloud SDK](https://cloud.google.com/sdk/docs/install).
- Generate a [service account key](https://console.cloud.google.com/apis/credentials/serviceaccountkey). Download it on your system and add reference it from your .envrc file: `GOOGLE_APPLICATION_CREDENTIALS=~/gcp-cloud-key.json`.
- Create a Google Cloud project.
- Create/update your .env file with the Google Cloud information - replace the values with your own information:

```bash
# GCP Deployment
## GCP project ID.
AMB_GCP_PROJECT_ID=my-sample-project-191923
## GCP bucket name (this can be one that doesn't exist yet).
AMB_GCP_BUCKET_NAME=sample-bucket
## Name of the docker image that will be created and stored in GCP Repository.
AMB_GCP_IMAGE_NAME=sample-image
## Name of the Cloud Run service to create.
AMB_GCP_CLOUDRUN_NAME=sample-service
## Region (not zone) where the Cloud Run service will be created:
## https://cloud.google.com/compute/docs/regions-zones#available
AMB_GCP_REGION=us-central1
```

- Run these commands:

```bash
# Authenticate with Google Cloud.
gcloud auth login

# Set current project - replace the value with your own information.
gcloud config set project my-sample-project-191923

# Create a bucket in Google Cloud, enable versioning, and create blank site.json
# and session.bin.
make gcp-init

# Run a Google Cloud Build to build the docker image, push to the Container
# Registry, and then deploy a Google Cloud Run service.
make gcp-deploy
```

- You should now be able to access the URL that appeared in your terminal like this: `Service URL: https://ambient-someurl-uc.a.run.app`
- To remove the service and bucket from GCP, run: `make gcp-delete`.