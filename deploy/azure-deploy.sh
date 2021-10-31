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

# Source: https://docs.microsoft.com/en-US/cli/azure/functionapp?view=azure-cli-latest#az_functionapp_create
echo Creating Azure Function. If it returns, 'Conflict', then you need a unique Function name.
az functionapp create --resource-group ${AMB_AZURE_RESOURCE_GROUP} --name ${AMB_AZURE_FUNCTION_NAME} --storage-account ${AZURE_STORAGE_ACCOUNT} --runtime custom --functions-version 2 --consumption-plan-location ${AMB_AZURE_REGION} --os-type linux

