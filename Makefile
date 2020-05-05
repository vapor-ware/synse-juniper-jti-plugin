#
# Synse Juniper JTI Plugin
#

PLUGIN_NAME    := juniper-jti
PLUGIN_VERSION := 0.1.0
IMAGE_NAME     := vaporio/juniper-jti-plugin
BIN_NAME       := synse-juniper-jti-plugin

GIT_COMMIT     ?= $(shell git rev-parse --short HEAD 2> /dev/null || true)
GIT_TAG        ?= $(shell git describe --tags 2> /dev/null || true)
BUILD_DATE     := $(shell date -u +%Y-%m-%dT%T 2> /dev/null)
GO_VERSION     := $(shell go version | awk '{ print $$3 }')

PKG_CTX := github.com/vapor-ware/synse-sdk/sdk
LDFLAGS := -w \
	-X ${PKG_CTX}.BuildDate=${BUILD_DATE} \
	-X ${PKG_CTX}.GitCommit=${GIT_COMMIT} \
	-X ${PKG_CTX}.GitTag=${GIT_TAG} \
	-X ${PKG_CTX}.GoVersion=${GO_VERSION} \
	-X ${PKG_CTX}.PluginVersion=${PLUGIN_VERSION}

.PHONY: build build-linux clean cover deploy docker docker-dev fmt
.PHONY: github-tag lint test version help

.DEFAULT_GOAL := help


build:  ## Build the plugin binary
	go build -ldflags "${LDFLAGS}" -o ${BIN_NAME}

build-linux:  ## Build the plugin binary for linux amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} .

clean:  ## Remove temporary files
	go clean -v
	rm -rf dist

cover: test  ## Run tests and open the coverage report
	go tool cover -html=coverage.out

deploy:  ## Run a local deployment of the plugin with Synse Server
	docker-compose up -d

docker:  ## Build the production docker image locally
	docker build -f Dockerfile \
		--label "org.label-schema.build-date=${BUILD_DATE}" \
		--label "org.label-schema.vcs-ref=${GIT_COMMIT}" \
		--label "org.label-schema.version=${PLUGIN_VERSION}" \
		-t ${IMAGE_NAME}:latest .

docker-dev:  ## Build the development docker image locally
	docker build -f Dockerfile.dev -t ${IMAGE_NAME}:dev-${GIT_COMMIT} .

fmt:  ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' -not -wholename '*.pb.go' | while read -r file; do goimports -w "$$file"; done

github-tag:  ## Create and push a tag with the current plugin version
	git tag -a ${PLUGIN_VERSION} -m "${PLUGIN_NAME} plugin version ${PLUGIN_VERSION}"
	git push -u origin ${PLUGIN_VERSION}

lint:  ## Lint project source files
	golint -set_exit_status ./pkg/...

test:  ## Run project tests
	@ # Note: this requires go1.10+ in order to do multi-package coverage reports
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

version:  ## Print the version of the plugin
	@echo "${PLUGIN_VERSION}"

help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
