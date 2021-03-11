build:
	go build -o api cmd/api/*

run-dev:
	go run cmd/api/*

clean-test-db:
	docker-compose -f docker-compose.test.yml rm -v --stop --force mysql

test-e2e: clean-test-db test-e2e-ci

test-e2e-ci:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit e2e_tests

.PHONY: build