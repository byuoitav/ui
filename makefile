NAME := ui
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}
DOCKER_URL := docker.pkg.github.com

# version:
# use the git tag, if this commit
# doesn't have a tag, use the git hash
VERSION := $(shell git rev-parse HEAD)
ifneq ($(shell git describe --exact-match --tags HEAD 2> /dev/null),)
	VERSION = $(shell git describe --exact-match --tags HEAD)
endif

# go stuff
PKG_LIST := $(shell cd backend && go list ${PKG}/...)

.PHONY: all deps build test test-cov clean

all: clean build

deps:
	@echo Downloading backend dependencies...
	@cd backend && go mod download

	@echo Downloading frontend dependencies...
	@cd frontend/dragonfruit && npm install

build: deps
	@mkdir -p dist
	@echo Building backend...
	@cd backend && env GOOS=linux GOARCH=amd64 go build -v -i -o ../dist/${NAME}-linux-amd64 ${PKG}

	@echo Building dragonfruit...
	@cd frontend/dragonfruit && npm run-script build && mv ./dist/dragonfruit ../../dist/ && rmdir ./dist

	@echo Build output is located in ./dist/.

docker: clean build
	@echo Building docker container ${OWNER}/${NAME}:${VERSION}
	docker build -f dockerfile -t ${DOCKER_URL}/${OWNER}/${NAME}/amd64:${VERSION} dist

	@echo Logging into Dockerhub
	docker login ${DOCKER_URL} -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

	@echo Pushing container to Dockerhub
	docker push ${DOCKER_URL}/${OWNER}/${NAME}/amd64:${VERSION}

test:
	@cd backend && go test -v ${PKG_LIST} && pwd

test-cov:
	@cd backend && go test -coverprofile=coverage.txt -covermode=atomic ${PKG_LIST}

lint:
	@cd backend && golangci-lint run --tests=false

clean:
	@cd backend && go clean
	@cd frontend/dragonfruit && rm -rf dist node_modules
	@rm -rf dist/
