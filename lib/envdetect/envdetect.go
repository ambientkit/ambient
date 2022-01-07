// Package envdetect detects in which cloud an app is running based on environment variables.
// Source: https://www.josephspurrier.com/cloud-environment-variables
package envdetect

import (
	"os"
	"strconv"
)

// RunningLocalDev returns true if the AMB_LOCAL environment variable is set.
func RunningLocalDev() bool {
	result, _ := strconv.ParseBool(os.Getenv("AMB_LOCAL"))
	return result
}

// RunningInAWS returns true if running in AWS services. When running in
// App Runner, it will be set: AWS_EXECUTION_ENV=AWS_ECS_FARGATE.
func RunningInAWS() bool {
	_, exists := os.LookupEnv("AWS_EXECUTION_ENV")
	return exists
}

// RunningInAWSLambda returns true if running in AWS Lambda.
func RunningInAWSLambda() bool {
	_, exists := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	return exists
}

// RunningInGoogle returns true if running in Google. When running in
// Google Cloud Run, will be set: K_SERVICE=NAME.
func RunningInGoogle() bool {
	_, exists := os.LookupEnv("K_SERVICE")
	return exists
}

// RunningInAzureFunction returns true if running in Azure Functions.
func RunningInAzureFunction() bool {
	_, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	return exists
}
