-include .env
export

# disables built-in rules
.SUFFIXES: 

###############################################################################
#
# Dependencies
#
###############################################################################

GOLANGCI-LINT_VERSION := v1.52.2

GOPATH := $(shell go env GOPATH)

golangci-lint := $(GOPATH)/bin/golangci-lint

$(golangci-lint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI-LINT_VERSION)

.PHONY: install	
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI-LINT_VERSION)

###############################################################################
#
# Docker compose commands
#
###############################################################################

.PHONY: clean-nats
clean-all: down-natsbox
	docker compose down --volumes

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down: down-natsbox 
	docker compose down

.PHONY: up-natsbox
up-natsbox:
	docker compose -f docker-natsbox.yaml up -d 

.PHONY: down-natsbox
down-natsbox:
	docker compose -f docker-natsbox.yaml down 

.PHONY: natsbox
natsbox: up-natsbox
	docker exec -it nats-gen-nats-box-1 /bin/sh

###############################################################################
#
# Run commands 
#
###############################################################################

.PHONY: bank
bank:
	@mkdir -p gen/bank
	go run cmd/nats-gen/*.go --config examples/bank/spec.yaml 

.PHONY: clean
clean:
	@rm -rf dist
	@rm -rf gen/*

###############################################################################
#
# Release commands 
#
###############################################################################

TAGS := $(shell git show-ref --tags | wc -l)
RELEASE_NUMBER := $(shell expr $(TAGS) + 1)
VERSION := v0.0.$(RELEASE_NUMBER)-alpha

.PHONY: version
version:
	@echo version tag: $(VERSION)

.PHONY: release
release:
	@echo version tag: $(VERSION)
	@git tag ${VERSION}
	@git push origin ${VERSION}
	goreleaser release --clean
	
###############################################################################
#
# Linting & testing
#
###############################################################################

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: test
test:
	go test ./... -timeout 30s -failfast
