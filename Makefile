build:
	go build -o api cmd/api/*

test-e2e:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit e2e_tests

.PHONY: build