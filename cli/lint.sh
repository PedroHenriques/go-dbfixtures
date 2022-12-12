#!/bin/sh
set -e

docker run --rm -v "${PWD}/":"/usr/src/app/" golangci/golangci-lint:v1.50 /bin/sh -c "go mod tidy && golangci-lint run -v";