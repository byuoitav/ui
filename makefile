NAME := ui
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}

# go stuff
PKG_LIST := $(shell cd backend && go list ${PKG}/...)

.PHONY: all deps build test test-cov clean

all: build

deps:
	@cd backend && go mod download

build: deps
	@mkdir -p dist
	@echo Building backend...
	@cd backend && env GOOS=linux GOARCH=amd64 go build -v -i -o ../dist/${NAME}-linux-amd64 ${PKG}
	@echo Done.

	@echo Building frontend...
	@echo Done.
	@echo Build output is located in ./dist/.

test:
	@cd backend && go test -v ${PKG_LIST} && pwd

test-cov:
	@cd backend && go test -coverprofile=coverage.txt -covermode=atomic ${PKG_LIST}

clean:
	@cd backend && go clean
	@rm -rf dist/

lint:
	@cd backend && golangci-lint run --tests=false
