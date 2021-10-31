#!/usr/bin/env bash

# Load the shared environment variables (shared with Makefile).
# Export the vars in .env into the shell.
export $(egrep -v '^#' .env | xargs)

# Source: https://docs.microsoft.com/en-us/azure/azure-functions/functions-create-function-linux-custom-image?tabs=in-process%2Cbash%2Cazure-cli&pivots=programming-language-csharp#create-and-configure-a-function-app-on-azure-with-the-image
# Source: https://docs.microsoft.com/en-US/cli/azure/functionapp?view=azure-cli-latest#az_functionapp_create

# If the function exists, then create it.
echo Detecting existence of Azure Function.
if az functionapp config show --resource-group ${AMB_AZURE_RESOURCE_GROUP} --name ${AMB_AZURE_FUNCTION_NAME} --query 'name' -o tsv >/dev/null 2>&1; then
    echo Azure function found - will update it now.
else
    echo Creating Azure Function. If it returns 'Conflict', then you need a unique Function name. If it returns 'Bad Request', you may need to wait a few minutes before trying again.
    az functionapp create --resource-group ${AMB_AZURE_RESOURCE_GROUP} --name ${AMB_AZURE_FUNCTION_NAME} --storage-account ${AZURE_STORAGE_ACCOUNT} --runtime custom --functions-version 3 --consumption-plan-location ${AMB_AZURE_REGION} --os-type linux
fi

echo Building the Go binary for Linux.
GOOS=linux go build -o deploy/azure/ambient cmd/myapp/main.go

echo Setting environment variables and connection string.
az functionapp config appsettings set --resource-group ${AMB_AZURE_RESOURCE_GROUP} --name ${AMB_AZURE_FUNCTION_NAME} \
    --settings "AzureWebJobsStorage=$(az storage account show-connection-string --name ${AZURE_STORAGE_ACCOUNT} --query 'connectionString' -o tsv)" \
    "AMB_AZURE_CONTAINER=${AMB_AZURE_CONTAINER}" \
    "AMB_SESSION_KEY=${AMB_SESSION_KEY}" \
    "AMB_PASSWORD_HASH=${AMB_PASSWORD_HASH}"

echo Publishing binary to Azure Function.
(cd deploy/azure && func azure functionapp publish ambientapp)