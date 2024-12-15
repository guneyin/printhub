BINARY_NAME=printhub

.PHONY: build

MAIN_FILE=main.go
PACKAGE=github.com/guneyin/printhub
VERSION=$(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
COMMIT_HASH=$(shell git rev-list -1 HEAD)
BUILD_TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S')

LDFLAG_VERSION='${PACKAGE}/utils.Version=${VERSION}'
LDFLAG_COMMIT_HASH='${PACKAGE}/utils.CommitHash=${COMMIT_HASH}'
LDFLAG_BUILD_TIMESTAMP='${PACKAGE}/utils.BuildTime=${BUILD_TIMESTAMP}'

tidy:
	go mod tidy

vet:
	go vet ./...

doc:
	swag init

run:
	go run ${MAIN_FILE}

build:
	go build -o ${BINARY_NAME} -ldflags "-X ${LDFLAG_VERSION} -X ${LDFLAG_COMMIT_HASH} -X ${LDFLAG_BUILD_TIMESTAMP}" ${MAIN_FILE}

clean:
	go clean
	rm -f ${BINARY_NAME}

