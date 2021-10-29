# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make
# * make storage
# * make run

# Load the environment variables.
include .env

.PHONY: default
default: run

################################################################################
# Setup application
################################################################################

.PHONY: privatekey
privatekey:
	@echo Generating private key for encrypting sessions.
	@echo You can paste private key this into your .env file:
	@go run plugin/scssession/cmd/privatekey/main.go

# Save the ARGS.
# https://stackoverflow.com/a/14061796
ifeq (mfa,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

.PHONY: mfa
mfa:
	@echo Generating MFA for user.
	@echo You can paste private key this into your .env file:
	@go run plugin/bearblog/cmd/mfa/main.go ${ARGS}

# Save the ARGS.
# https://stackoverflow.com/a/14061796
ifeq (passhash,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

.PHONY: passhash
passhash:
	@echo Generating password hash.
	@echo You can paste private key this into your .env file:
	@go run plugin/bearblog/cmd/passhash/main.go ${ARGS}

.PHONY: storage
storage:
	@echo Creating session and site storage files locally.
	cp storage/initial/session.bin storage/session.bin
	cp storage/initial/site.json storage/site.json

.PHONY: run
run:
	@echo Starting local server.
	LOCALDEV=true go run cmd/myapp/main.go

.PHONY: amb
amb:
	go run cmd/amb/main.go

################################################################################
# Deploy application to Google Cloud
################################################################################

.PHONY: gcp-init
gcp-init:
	@echo Creating the initial files in Google Cloud Storage.
	gsutil mb -p $(AMB_GCP_PROJECT_ID) -l ${AMB_GCP_REGION} -c Standard gs://${AMB_GCP_BUCKET_NAME}
	gsutil versioning set on gs://${AMB_GCP_BUCKET_NAME}
	gsutil cp storage/initial/site.json gs://${AMB_GCP_BUCKET_NAME}/storage/site.json
	gsutil cp storage/initial/session.bin gs://${AMB_GCP_BUCKET_NAME}/storage/session.bin

.PHONY: gcp-deploy
gcp-deploy:
	@echo Deploying to Google Cloud Run.
	gcloud builds submit --tag gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE_NAME}
	gcloud run deploy --image gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE_NAME} \
		--platform managed \
		--allow-unauthenticated \
		--region ${AMB_GCP_REGION} ${AMB_GCP_CLOUDRUN_NAME} \
		--update-env-vars AMB_SESSION_KEY=${AMB_SESSION_KEY} \
		--update-env-vars AMB_PASSWORD_HASH=${AMB_PASSWORD_HASH} \
		--update-env-vars AMB_GCP_PROJECT_ID=${AMB_GCP_PROJECT_ID} \
		--update-env-vars AMB_GCP_BUCKET_NAME=${AMB_GCP_BUCKET_NAME}

.PHONY: gcp-delete
gcp-delete:
	@echo Removing files from Google Cloud.
	-gcloud run services delete --platform managed --region ${AMB_GCP_REGION} ${AMB_GCP_CLOUDRUN_NAME}
	-gsutil -m rm -r -f gs://${AMB_GCP_BUCKET_NAME}

################################################################################
# Deploy application to AWS
################################################################################

.PHONY: aws-init
aws-init:
	@echo Creating the initial files in AWS S3.
ifeq "${AWS_REGION}" "us-east-1")
	aws s3api create-bucket --bucket ${AMB_AWS_BUCKET_NAME}
else
	aws s3api create-bucket --bucket ${AMB_AWS_BUCKET_NAME} --create-bucket-configuration '{"LocationConstraint": "${AWS_REGION}"}'
endif
	aws s3api put-public-access-block --bucket ${AMB_AWS_BUCKET_NAME} --public-access-block-configuration '{"BlockPublicAcls": true,"IgnorePublicAcls": true,"BlockPublicPolicy": true,"RestrictPublicBuckets": true}'
	aws s3 cp storage/initial/site.json s3://${AMB_AWS_BUCKET_NAME}/storage/site.json
	aws s3 cp storage/initial/session.bin s3://${AMB_AWS_BUCKET_NAME}/storage/session.bin

.PHONY: aws-deploy
aws-deploy:
	@echo Deploying to AWS App Runner.
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AMB_AWS_ACCOUNT_NUMBER}.dkr.ecr.${AWS_REGION}.amazonaws.com
	-aws ecr create-repository --repository-name ${AMB_GCP_IMAGE_NAME}
#	docker build -t ${AMB_AWS_ACCOUNT_NUMBER}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AMB_GCP_IMAGE_NAME}:1.0 .
	docker push ${AMB_AWS_ACCOUNT_NUMBER}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AMB_GCP_IMAGE_NAME}:1.0
	-aws cloudformation create-stack --stack-name ${AMB_GCP_CLOUDRUN_NAME} --template-body file://aws-apprunner.json --capabilities CAPABILITY_IAM --parameters ParameterKey=ParameterSessionKey,ParameterValue=${AMB_SESSION_KEY} ParameterKey=ParameterPasswordHash,ParameterValue=${AMB_PASSWORD_HASH} ParameterKey=ParameterAWSS3Bucket,ParameterValue=${AMB_AWS_BUCKET_NAME} ParameterKey=ParameterAWSRegion,ParameterValue=${AWS_REGION}
	-aws cloudformation update-stack --stack-name ${AMB_GCP_CLOUDRUN_NAME} --template-body file://aws-apprunner.json --capabilities CAPABILITY_IAM --parameters ParameterKey=ParameterSessionKey,ParameterValue=${AMB_SESSION_KEY} ParameterKey=ParameterPasswordHash,ParameterValue=${AMB_PASSWORD_HASH} ParameterKey=ParameterAWSS3Bucket,ParameterValue=${AMB_AWS_BUCKET_NAME} ParameterKey=ParameterAWSRegion,ParameterValue=${AWS_REGION}

#	gcloud builds submit --tag gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE_NAME}
#	gcloud run deploy --image gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE_NAME} \
		--platform managed \
		--allow-unauthenticated \
		--region ${AMB_GCP_REGION} ${AMB_GCP_CLOUDRUN_NAME} \
		--update-env-vars AMB_SESSION_KEY=${AMB_SESSION_KEY} \
		--update-env-vars AMB_PASSWORD_HASH=${AMB_PASSWORD_HASH} \
		--update-env-vars AMB_GCP_PROJECT_ID=${AMB_GCP_PROJECT_ID} \
		--update-env-vars AMB_GCP_BUCKET_NAME=${AMB_GCP_BUCKET_NAME}

.PHONY: aws-delete
aws-delete:
	@echo Removing files from AWS.
	aws cloudformation delete-stack --stack-name ${AMB_GCP_CLOUDRUN_NAME}
	-aws s3 rm s3://${AMB_AWS_BUCKET_NAME} --recursive
	-aws s3api delete-bucket --bucket ${AMB_AWS_BUCKET_NAME}