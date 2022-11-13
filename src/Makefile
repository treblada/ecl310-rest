BINARY_NAME=ecl310-rest
BIN_DIR=../bin

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