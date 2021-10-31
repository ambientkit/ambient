#!/usr/bin/env bash

# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

# Set the Azure alias if using docker.
if ! command -v az &> /dev/null; then
    az() {
        docker exec azurecli az "$@"
    }
fi

echo Creating resource group.
az group create --location eastus --resource-group ${AZURE_RESOURCE_GROUP}

echo Creating storeage account.
az storage account create --name ${AZURE_STORAGE_ACCOUNT} --resource-group ${AZURE_RESOURCE_GROUP}

# This method may take a few minutes to propogate so we will skip it
# az role assignment create \
#     --role "Storage Blob Data Contributor" \
#     --assignee $(az ad signed-in-user show --query objectId -o tsv) \
#     --scope "/subscriptions/SUBSCRIPTIONID/resourceGroups/${AZURE_RESOURCE_GROUP}/providers/Microsoft.Storage/storageAccounts/${AZURE_STORAGE_ACCOUNT}"

echo Creating storage container.
az storage container create --name ${AZURE_CONTAINER_NAME} --account-name ${AZURE_STORAGE_ACCOUNT} --account-name ${AZURE_STORAGE_ACCOUNT} --account-key $(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)

echo Uploading files.
az storage blob upload --container-name ${AZURE_CONTAINER_NAME} --file /root/storage/initial/site.json --name storage/site.json --account-name ${AZURE_STORAGE_ACCOUNT} --account-key $(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)
az storage blob upload --container-name ${AZURE_CONTAINER_NAME} --file /root/storage/initial/session.bin --name storage/session.bin --account-name ${AZURE_STORAGE_ACCOUNT} --account-key $(az storage account keys list --account-name ${AZURE_STORAGE_ACCOUNT} --query '[0].value' -o tsv)