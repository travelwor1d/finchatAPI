include .env

SHA=`git rev-parse HEAD`
LATEST_TAG=`gcloud container images list-tags --format=json --limit=1 gcr.io/${GCP_PROJECT}/finchat-api | jq -r '.[0].tags[0]'`

build: gen
	go build -o api cmd/api/*

configure-gcloud:
	gcloud config set project $(GCP_PROJECT)

submit-build: configure-gcloud
	gcloud builds submit -t gcr.io/$(GCP_PROJECT)/finchat-api:$(SHA)

init-terraform:
	terraform -chdir=terraform/staging init

plan-terraform:
	terraform -chdir=terraform/staging plan -var image_tag=$(LATEST_TAG)

apply-terraform:
	terraform -chdir=terraform/staging apply -var image_tag=$(LATEST_TAG) -auto-approve

gen: init-swag

run-dev: gen
	go run cmd/api/*

clean-test-db:
	docker-compose -f docker-compose.test.yml rm -v --stop --force mysql

test-e2e: clean-test-db test-e2e-ci

test-e2e-ci:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit e2e_tests

init-swag:
	swag init -g internal/app/app.go

install-swag:
	go get -u github.com/swaggo/swag/cmd/swag

install: install-swag

.PHONY: build