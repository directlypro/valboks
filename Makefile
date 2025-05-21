.PHONY: build clean test

BUILD_DIR=bin
BINARY_NAME=valboks-cli

build:
	go build -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/valboks-cli


clean:
	rm -rf ${BUILD_DIR}

test:
	go test ./..

run:
	go run ./cmd/valboks-cli
