HOSTNAME=github.com
NAMESPACE=jlopez
NAME=unifi
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_arm64

default: build

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test ./... -v

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

lint:
	golangci-lint run

fmt:
	go fmt ./...
	terraform fmt -recursive examples/

pre-commit: fmt lint test

docker-up:
	docker compose up -d

docker-down:
	docker compose down

.PHONY: build install test testacc lint fmt pre-commit docker-up docker-down
