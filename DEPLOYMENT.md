# Deployment Guide <!-- omit in toc -->

- [Environment Variable Management](#environment-variable-management)
- [Deployments](#deployments)
  - [Google Cloud - Cloud Run](#google-cloud---cloud-run)
  - [AWS - App Runner](#aws---app-runner)
  - [Azure - Functions](#azure---functions)

## Environment Variable Management

It's recommended to install [direnv](https://direnv.net/docs/installation.html) to help manage your environment variables out of a `.envrc` file. The benefit is when you CD out of the folder, the environment variables will be unset so they are available to the root folder and all child folders. The directions below will assume you have this utility installed.

Once you have `direnv` installed, create a .envrc file in the root of your project:

```bash
# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

export PATH=$PATH:$(pwd)/bin:$(pwd)/node_modules/.bin
```

When you open up your terminal, you will be prompted with this message. You should type `direnv allow` to allow the system to load environment variables from the file.

```bash
direnv: error /Users/YOURPATH/.envrc is blocked. Run `direnv allow` to approve its content
```

## Deployments

For all deployments, you need to create a new file called `.env` in the root of the repository with this content:

```bash
# App version.
AMB_APP_VERSION=1.0
# Set this to any value to allow you to do testing locally without cloud access.
AMB_LOCAL=true
# Optional: Enable the Dev Console that amb connects to. Default is: false
AMB_DEVCONSOLE_ENABLE=true
# Optional: Set the URL for the Dev Console that amb connects to. Default is: http://localhost
# AMB_DEVCONSOLE_URL=http://localhost
# Optional: Set the port for the Dev Console that amb connects to. Default is: 8081
# AMB_DEVCONSOLE_PORT=8081
# Session key to encrypt the cookie store. Generate with: make privatekey
AMB_SESSION_KEY=
# Password hash that is base64 encoded. Generate with: make passhash passwordhere
AMB_PASSWORD_HASH=

# Optional: set the web server port.
# PORT=8080
# Optional: set the time zone from here:
# https://golang.org/src/time/zoneinfo_abbrs_windows.go
# AMB_TIMEZONE=America/New_York
# Optional: set the URL prefix if behind a proxy.
# AMB_URL_PREFIX=/api
```

### Google Cloud - Cloud Run

To deploy an Ambient app to Google Cloud Run:

- Install the [Google Cloud SDK](https://cloud.google.com/sdk/docs/install).
- Generate a [service account key](https://console.cloud.google.com/apis/credentials/serviceaccountkey). Download it on your system and add reference it from your .envrc file: `GOOGLE_APPLICATION_CREDENTIALS=~/gcp-cloud-key.json`. This is needed only if you want to test locally.
- Create a Google Cloud project.
- Update your .env file with the Google Cloud information - replace the values with your own information:

```bash
# GCP project ID.
AMB_GCP_PROJECT_ID=myapp-191923
# GCP bucket name (this can be one that doesn't exist yet).
AMB_GCP_BUCKET=myapp-bucket
# Name of the docker image that will be created and stored in GCP Repository.
AMB_GCP_IMAGE=myapp-image
# Name of the Cloud Run service to create.
AMB_GCP_CLOUDRUN_NAME=myapp-service
# Region (not zone) where the Cloud Run service will be created:
# https://cloud.google.com/compute/docs/regions-zones#available
AMB_GCP_REGION=us-central1
```

- Run these commands:

```bash
# Authenticate with Google Cloud.
gcloud auth login

# Set current project - replace the value with your own information.
gcloud config set project myapp-project-191923

# Create a bucket in Google Cloud, enable versioning, and upload a
# blank site.bin and session.bin.
make gcp-init

# Run a Google Cloud Build to build the docker image, push to the Container
# Registry, and then deploy a Google Cloud Run service.
make gcp-deploy
```

- You should now be able to access the URL that appeared in your terminal like this: `Service URL: https://ambient-someurl-uc.a.run.app`
- To remove the service and bucket from GCP, run: `make gcp-delete`.

### AWS - App Runner

To deploy an Ambient app to AWS App Runner:

- Install the [AWS CLI v2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html).
- Generate access keys from your AWS account.
- Update your .envrc file with the AWS credentials - replace the values with your own information:

```bash
# AWS access keys.
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

- Update your .env file with the AWS information - replace the values with your own information:

```bash
# AWS account number where the App Runner will be created.
AMB_AWS_ACCOUNT=121212121212
# AWS region where the App Runner will be created.
AWS_REGION=us-east-1
# AWS S3 bucket name where the site and session files will be stored. New or existing.
AMB_AWS_BUCKET=myapp-storage
```

- Run these commands:

```bash

# You should get an output of your storage access key and connection string.
# Add it to your .envrc file. Then run the command to trust and reload the 
# env variables.
direnv allow

# Create a bucket in AWS, enable versioning, and upload a blank site.bin
# and session.bin.
make aws-init

# Run a docker build to create the docker image, push to the AWS ECR, and then
# deploy to AWS App Runner.
make aws-deploy
```

- You should now be able to access the URL from the App Runner service: like this: `Default domain: https://someurl.us-east-1.awsapprunner.com`
- To remove the service and bucket from AWS, run: `make aws-delete`.

### Azure - Functions

To deploy an Ambient app to an Azure Function, you will need the Azure CLI.

If you don't have the Azure CLI installed, you can either [install it](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) or run it from a container. You can skip the next few steps if you already have the Azure CLI installed.

To run the Azure CLI from a docker container, paste the code below into a new file: `bin/az`. It will allow you to use the the dockerized version of the Azure CLI. The path to the bin folder should be included in your `PATH` env variable already because of the .envrc file above.

```bash
#!/usr/bin/env bash

docker exec azurecli az "$@"
```

You will also need to install the Azure Functions Core Tools. You could [install for your operating system](https://docs.microsoft.com/en-us/azure/azure-functions/functions-run-local) or run this command to install it using npm:

```bash
# Install v3 of the Azure Functions Core Tools locally using npm.
npm install azure-functions-core-tools@3 --unsafe-perm true
```

- Update your .env file with the Azure information - replace the values with your own information:

```bash
# Azure resource group where storage account will be created. New or existing.
AMB_AZURE_RESOURCE_GROUP=myapp-rg
# Azure storage account where the storage container will be created. New or existing. Unique.
AZURE_STORAGE_ACCOUNT=myappstorage191923
# Azure container where the site and session files will be stored. New or existing.
AMB_AZURE_CONTAINER=myapp-container
# Azure function name. Unique.
AMB_AZURE_FUNCTION_NAME=myapp-191923
```

Refresh your terminal session and then run these commands:

```bash
# Start docker in the background so you can use the Azure CLI without installing it.
make azcli-start

# Login to Azure.
az login

# Test acccess by trying to load storage accounts.
az storage account list

# Create a resource group, storage account, storage container, and upload a
# blank site.bin and session.bin.
make az-init

# You should get an output of your storage access key and connection string.
# Add it to your .envrc file. Then run the command to trust and reload the 
# env variables.
direnv allow

# Create or update an Azure Function with a custom runtime, update env variables,
# and then deploy the Go binary.
make az-deploy

# Remove resource group which removes the Azure Function and storage account.
make az-delete

# When you're done, you can stop the Azure CLI docker container from running in the background.
make azcli-stop
```