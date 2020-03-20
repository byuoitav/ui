NAME := ui
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}
DOCKER_URL := docker.pkg.github.com

# version:
# use the git tag, if this commit
# doesn't have a tag, use the git hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)
VERSION := $(shell git rev-parse --short HEAD)
ifneq ($(shell git describe --exact-match --tags HEAD 2> /dev/null),)
	VERSION = $(shell git describe --exact-match --tags HEAD)
endif

# go stuff
PKG_LIST := $(shell cd backend && go list ${PKG}/...)

.PHONY: all deps build test test-cov clean

all: clean build

test:
	@cd backend && go test -v ${PKG_LIST}

test-cov:
	@cd backend && go test -coverprofile=coverage.txt -covermode=atomic ${PKG_LIST}

lint:
	@cd backend && golangci-lint run --tests=false

deps:
	@echo Downloading backend dependencies...
	@cd backend && go mod download

	@echo Downloading frontend dependencies for dragonfruit...
	@cd frontend/dragonfruit && npm install

	@echo Downloading frontend dependencies for blueberry...
	@cd frontend/blueberry && npm install

	@echo Downloading frontend dependencies for cherry...
	@cd frontend/cherry && npm install

build: deps
	@mkdir -p dist
	@echo
	@echo Building backend for linux-amd64...
	@cd backend && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../dist/${NAME}-linux-amd64 ${PKG}

	@echo
	@echo Building backend for linux-arm...
	@cd backend && env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -v -o ../dist/${NAME}-linux-arm ${PKG}

	@echo
	@echo Building dragonfruit...
	@cd frontend/dragonfruit && npm run-script build && mv ./dist/dragonfruit ../../dist/ && rmdir ./dist

	@echo
	@echo Building blueberry...
	@cd frontend/blueberry && npm run-script build && mv ./dist/blueberry ../../dist/ && rmdir ./dist

	@echo
	@echo Building cherry...
	@cd frontend/cherry && npm run-script build && mv ./dist/cherry ../../dist/ && rmdir ./dist

	@echo
	@echo Build output is located in ./dist/.

docker: clean build
ifneq (${COMMIT_HASH},${VERSION})
	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${VERSION}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-amd64 -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${VERSION} dist

	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm:${VERSION}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-arm -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm:${VERSION} dist
else
	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-amd64 -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH} dist

	@echo Building container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm-dev:${COMMIT_HASH}
	@docker build -f dockerfile --build-arg NAME=${NAME}-linux-arm -t ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm-dev:${COMMIT_HASH} dist
endif

deploy: docker
	@echo Logging into Github Package Registry
	@docker login ${DOCKER_URL} -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

# if the commit hash and release are different, this is a tagged build and we should build the tagged version
ifneq (${COMMIT_HASH},${VERSION})
	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${VERSION}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}:${VERSION}

	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm:${VERSION}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm:${VERSION}
else
	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-dev:${COMMIT_HASH}

	@echo Pushing container ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm-dev:${COMMIT_HASH}
	@docker push ${DOCKER_URL}/${OWNER}/${NAME}/${NAME}-arm-dev:${COMMIT_HASH}
endif

clean:
	@cd backend && go clean
	@cd frontend/dragonfruit && rm -rf dist node_modules
	@cd frontend/blueberry && rm -rf dist node_modules
	@cd frontend/cherry && rm -rf dist node_modules
	@rm -rf dist/
