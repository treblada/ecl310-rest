SHELL=/bin/bash
BINARY_NAME=ecl310-rest
BIN_DIR=../bin
OPENAPI_JAR=../openapi-generator-cli.jar
OPENAPI_VER=6.2.1
OPENAPI=java -jar ${OPENAPI_JAR}
OPENAPI_SPEC=ecl310-rest-openapi.yaml
GENERATED_DIR=generated

.DEFAULT_GOAL=all

.PHONY: openapi-clean openapi-generate xx-openapi-generate openapi-validate openapi-download all build test run install list-arch

build-darwin-amd64: openapi-generate
	GOARCH=amd64 GOOS=darwin go build -o ${BIN_DIR}/${BINARY_NAME}-darwin .

build-darwin-arm: openapi-generate
	GOARCH=arm64 GOOS=darwin go build -o ${BIN_DIR}/${BINARY_NAME}-silicon .

build-windows: openapi-generate
	GOARCH=amd64 GOOS=windows go build -o ${BIN_DIR}/${BINARY_NAME}-windows .

build-linux: openapi-generate
	GOARCH=amd64 GOOS=linux go build -o ${BIN_DIR}/${BINARY_NAME}-linux .

all: build-linux build-darwin-amd64 build-darwin-arm build-windows

build: build-linux

test:
	go test ./...

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BIN_DIR}/${BINARY_NAME}-*

install:
	go install

list-arch:
	go tool dist list

openapi-download:
	wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/$(OPENAPI_VER)/openapi-generator-cli-$(OPENAPI_VER).jar -O $(OPENAPI_JAR)

openapi-config-help:
	$(OPENAPI) config-help --generator-name go

openapi-generate: generated/api/openapi.yaml

generated/api/openapi.yaml: ecl310-rest-openapi.yaml
	$(MAKE) openapi-validate
	$(MAKE) openapi-clean
	$(OPENAPI) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name go-server \
		--output generated \
		--additional-properties=sourceFolder=openapi \
		--additional-properties=isGoSubmodule=true \
		--additional-properties=outputAsLibrary=true
	goimports -w generated/openapi/*.go

openapi-validate:
	$(OPENAPI) validate --input-spec $(OPENAPI_SPEC) --recommend

openapi-clean:
	rm -frv $(GENERATED_DIR)/.openapi-generator
	rm -frv $(GENERATED_DIR)/api
	rm -frv $(GENERATED_DIR)/docs
	rm -frv $(GENERATED_DIR)/Dockerfile
	rm -frv $(GENERATED_DIR)/go.mod
	rm -frv $(GENERATED_DIR)/main.go
	rm -frv $(GENERATED_DIR)/README.md
	find $(GENERATED_DIR)/openapi -maxdepth 1 -type f \! -name "controllers.go" -exec rm -fv {} \;
