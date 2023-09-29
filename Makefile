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
# Configure local environment
#
###############################################################################

	
###############################################################################
#
# Docker compose commands
#
###############################################################################

.PHONY: clean-all
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
	go run cmd/micro/*.go  \
		--input examples/bank/microservice.yaml \
		--output examples/bank/bank.gen.go  \
		--package bank 
		
###############################################################################
#
# Build commands 
#
###############################################################################


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
