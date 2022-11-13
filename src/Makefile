BINARY_NAME=ecl310-rest
BIN_DIR=../bin
OPENAPI_JAR=../openapi-generator-cli.jar
OPENAPI_VER=6.2.1
OPENAPI=java -jar ${OPENAPI_JAR}
OPENAPI_SPEC=ecl310-rest-openapi.yaml
OPENAPI_DIR=openapi

.DEFAULT_GOAL=all

build-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o ${BIN_DIR}/${BINARY_NAME}-darwin main.go

build-darwin-arm:
	GOARCH=arm64 GOOS=darwin go build -o ${BIN_DIR}/${BINARY_NAME}-silicon main.go

build-windows:
	GOARCH=amd64 GOOS=windows go build -o ${BIN_DIR}/${BINARY_NAME}-windows main.go

build-linux:
	GOARCH=amd64 GOOS=linux go build -o ${BIN_DIR}/${BINARY_NAME}-linux main.go

all: build-linux build-darwin-amd64 build-darwin-arm build-windows

build: build-linux

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

openapi-validate:
	$(OPENAPI) validate --input-spec $(OPENAPI_SPEC)

openapi-generate: openapi-clean
	$(OPENAPI) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name go-server \
		--output openapi \
		--additional-properties=isGoSubmodule=true \
		--additional-properties=outputAsLibrary=true

openapi-clean:
	rm -fr $(OPENAPI_DIR)/.openapi-generator
	rm -fr $(OPENAPI_DIR)/api
	rm -fr $(OPENAPI_DIR)/docs
	rm -fr $(OPENAPI_DIR)/go
	rm -fr $(OPENAPI_DIR)/Dockerfile
	rm -fr $(OPENAPI_DIR)/go.mod
	rm -fr $(OPENAPI_DIR)/main.go
	rm -fr $(OPENAPI_DIR)/README.md
