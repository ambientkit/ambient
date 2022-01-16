// Package envdetect detects in which cloud an app is running based on environment variables.
// Source: https://www.josephspurrier.com/cloud-environment-variables
package envdetect

import (
	"os"
	"strconv"
)

// LoadDotEnv returns true if the AMB_DOTENV environment variable is set.
func LoadDotEnv() bool {
	result, _ := strconv.ParseBool(os.Getenv("AMB_DOTENV"))
	return result
}

// RunningLocalDev returns true if the AMB_LOCAL environment variable is set.
func RunningLocalDev() bool {
	result, _ := strconv.ParseBool(os.Getenv("AMB_LOCAL"))
	return result
}

// DevConsoleEnabled returns true if the AMB_DEVCONSOLE_ENABLE environment variable is set.
func DevConsoleEnabled() bool {
	result, _ := strconv.ParseBool(os.Getenv("AMB_DEVCONSOLE_ENABLE"))
	return result
}

// DevConsoleURL returns the URL used for the Dev Console that amb connects to.
func DevConsoleURL() string {
	URL := os.Getenv("AMB_DEVCONSOLE_URL")
	if len(URL) == 0 {
		URL = "http://localhost"
	}
	return URL
}

// DevConsolePort returns the port used for the Dev Console that amb connects to.
func DevConsolePort() string {
	port := os.Getenv("AMB_DEVCONSOLE_PORT")
	if len(port) == 0 {
		port = "8081"
	}
	return port
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
