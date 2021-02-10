#!make
.DEFAULT_GOAL=start

local.build:
	@docker-compose build

local.start:
	@docker-compose up -d

local.down:
	@docker-compose down

local.test:
	@docker exec gandalf go test ./... -cover

local.coverage.generate_report:
	@docker exec gandalf go test -coverprofile coverage.out ./...

local.coverage.open_report:
	@go tool cover -html=coverage.out

ci.test:
	@docker exec gandalf go test -v -covermode=count -coverprofile=coverage.out ./...

logs:
	@docker logs -f $(shell docker-compose ps -q gandalf)

sh:
	@docker exec -it gandalf /bin/sh

docker_tag_and_push:
	@export TAG=$(date +%d%m%Y-%H%M%S)
	@docker build -f build/docker/dockerfile.prod -t $(REGISTRY):latest -t $(REGISTRY):$(TAG) .
	@echo $(GITHUB_TOKEN) | docker login ghcr.io -u $(GITHUB_USER) --password-stdin
	@docker push $(REGISTRY):$(TAG)
	@docker push $(REGISTRY):latest

start: local.start

stop: local.down

coverage_report: local.coverage.generate_report local.coverage.open_report

ci_check_tests: local.start ci.test

renew: local.down local.build local.start

.PHONY:  start stop sh logs renew coverage_report ci_check_tests docker_tag_and_push
