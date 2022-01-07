#!/usr/bin/env bash

# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

# echo Creating resource group.
az group create --location ${AMB_AZURE_REGION} --resource-group ${AMB_AZURE_RESOURCE_GROUP}

# echo Creating storage account.
az storage account create --name ${AZURE_STORAGE_ACCOUNT} --resource-group ${AMB_AZURE_RESOURCE_GROUP}

# echo Creating storage container.
az storage container create --name ${AMB_AZURE_CONTAINER} --account-name ${AZURE_STORAGE_ACCOUNT} --account-key $(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)

# echo Uploading files.
az storage blob upload --container-name ${AMB_AZURE_CONTAINER} --file /root/storage/initial/site.bin --name storage/site.bin --account-name ${AZURE_STORAGE_ACCOUNT} --account-key $(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)
az storage blob upload --container-name ${AMB_AZURE_CONTAINER} --file /root/storage/initial/session.bin --name storage/session.bin --account-name ${AZURE_STORAGE_ACCOUNT} --account-key $(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)

echo Getting access key. You should add these to your .envrc file:
echo export AZURE_STORAGE_ACCESS_KEY=$(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)
echo export AzureWebJobsStorage=\'$(az storage account show-connection-string --name ${AZURE_STORAGE_ACCOUNT} --query 'connectionString' -o tsv)\'