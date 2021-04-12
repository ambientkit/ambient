# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make
# * make gcp-init
# * make gcp-push
# * make privatekey
# * make mfa
# * make passhash passwordhere
# * make local-init
# * make local-run

# Load the environment variables.
include .env

.PHONY: default
default: gcp-push

################################################################################
# Deploy application
################################################################################

.PHONY: gcp-init
gcp-init:
	@echo Pushing the initial files to Google Cloud Storage.
	gsutil mb -p $(AMB_GCP_PROJECT_ID) -l ${AMB_GCP_REGION} -c Standard gs://${AMB_GCP_BUCKET_NAME}
	gsutil versioning set on gs://${AMB_GCP_BUCKET_NAME}
	gsutil cp testdata/empty.json gs://${AMB_GCP_BUCKET_NAME}/storage/site.json
	gsutil cp testdata/empty.json gs://${AMB_GCP_BUCKET_NAME}/storage/session.json

.PHONY: gcp-push
gcp-push:
	@echo Pushing to Google Cloud Run.
	gcloud builds submit --tag gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE_NAME}
	gcloud run deploy --image gcr.io/$(AMB_GCP_PROJECT_ID)/${AMB_GCP_IMAGE_NAME} \
		--platform managed \
		--allow-unauthenticated \
		--region ${AMB_GCP_REGION} ${AMB_GCP_CLOUDRUN_NAME} \
		--update-env-vars AMB_USERNAME=${AMB_USERNAME} \
		--update-env-vars AMB_SESSION_KEY=${AMB_SESSION_KEY} \
		--update-env-vars AMB_PASSWORD_HASH=${AMB_PASSWORD_HASH} \
		--update-env-vars AMB_MFA_KEY="${AMB_MFA_KEY}" \
		--update-env-vars AMB_GCP_PROJECT_ID=${AMB_GCP_PROJECT_ID} \
		--update-env-vars AMB_GCP_BUCKET_NAME=${AMB_GCP_BUCKET_NAME} \
		--update-env-vars AMB_ALLOW_HTML=${AMB_ALLOW_HTML}

.PHONY: privatekey
privatekey:
	@echo Generating private key for encrypting sessions.
	@echo You can paste private key this into your .env file:
	@go run cmd/privatekey/main.go

.PHONY: mfa
mfa:
	@echo Generating MFA for user.
	@echo You can paste private key this into your .env file:
	@go run cmd/mfa/main.go

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
	@go run cmd/passhash/main.go ${ARGS}

.PHONY: local-init
local-init:
	@echo Creating session and site storage files locally.
	cp storage/initial/session.bin storage/session.bin
	cp storage/initial/site.json storage/site.json

.PHONY: local-run
local-run:
	@echo Starting local server.
	LOCALDEV=true go run main.go

.PHONY: amb
amb:
	go run cmd/amb/main.go

.PHONY: ambient
ambient:
	go run cmd/ambient/main.go