-include .env

GCP_PROJECT=finchat-api-development
REGION=europe-central2
SHA=`git rev-parse HEAD`

build:
	go build -o api cmd/api/*

run-dev:
	go run cmd/api/*

clean-test-db:
	docker-compose -f docker-compose.test.yml rm -v --stop --force mysql

test-e2e: clean-test-db test-e2e-ci

test-e2e-ci:
	echo 'Test temporary disabled'
	
test-e2e-ci-OLD:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit e2e_tests

configure-gcloud:
	gcloud auth activate-service-account --key-file=gcp-credentials.json

submit-build: configure-gcloud
	gcloud builds submit -t eu.gcr.io/$(GCP_PROJECT)/finchat-api:$(SHA) --project $(GCP_PROJECT)

VARS=MYSQL_CONN_STR=$(REMOTE_MYSQL_CONN_STR),
VARS:=$(VARS)TWILIO_SID=$(TWILIO_SID),TWILIO_TOKEN=$(TWILIO_TOKEN),
VARS:=$(VARS)TWILIO_VERIFY=$(TWILIO_VERIFY),STRIPE_KEY=$(STRIPE_KEY),
VARS:=$(VARS)PUB_KEY=$(PUB_KEY),SUB_KEY=$(SUB_KEY),SEC_KEY=$(SEC_KEY),
VARS:=$(VARS)SERVER_UUID=$(SERVER_UUID),WEBHOOK_TOKEN=$(WEBHOOK_TOKEN)

deploy-cloud-run:
	gcloud run deploy finchat-api \
		--image=eu.gcr.io/$(GCP_PROJECT)/finchat-api:$(SHA) \
		--platform=managed --region=$(REGION) \
		--add-cloudsql-instances=finchat-api-development:europe-central2:db \
		--set-env-vars=$(VARS) \
		--allow-unauthenticated \
		--project $(GCP_PROJECT)

deploy: submit-build deploy-cloud-run

deploy-webhooks:
	gcloud functions deploy create-user-webhook \
		--entry-point=Webhook --runtime=go113 --region=$(REGION) \
		--trigger-event providers/firebase.auth/eventTypes/user.create \
  	--trigger-resource $(GCP_PROJECT) \
		--memory=128MB \
		--set-env-vars=WEBHOOK_TOKEN=secret,WEBHOOK_METHOD=POST,WEBHOOK_ENDPOINT=https://finchat-api-mp5dctunea-lm.a.run.app/auth/v1/users \
		--source=gcf/webhooks \
		--project $(GCP_PROJECT)

	gcloud functions deploy delete-user-webhook \
		--entry-point=Webhook --runtime=go113 --region=$(REGION) \
		--trigger-event providers/firebase.auth/eventTypes/user.delete \
  	--trigger-resource $(GCP_PROJECT) \
		--memory=128MB \
		--set-env-vars=WEBHOOK_TOKEN=secret,WEBHOOK_METHOD=DELETE,WEBHOOK_ENDPOINT=https://finchat-api-mp5dctunea-lm.a.run.app/auth/v1/users \
		--source=gcf/webhooks \
		--project $(GCP_PROJECT)

.PHONY: build