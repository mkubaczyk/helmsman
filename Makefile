.DEFAULT_GOAL := help

PKGS    := $(shell go list ./... | grep -v /vendor/)
TAG     := $(shell git describe --always --tags --abbrev=0 HEAD)
LAST    := $(shell git describe --always --tags --abbrev=0 HEAD^)
BODY    := "`git log ${LAST}..HEAD --oneline --decorate` `printf '\n\#\#\# [Build Info](${BUILD_URL})'`"
DATE    := $(shell date +'%d%m%y')
PRJNAME := $(shell basename "$(PWD)")

# Ensure we have an unambiguous GOPATH.
GOPATH := $(shell go env GOPATH)

ifneq ($(strip $(CIRCLE_WORKING_DIRECTORY)),)
  GOPATH := $(subst /src/$(PRJNAME),,$(CIRCLE_WORKING_DIRECTORY))
  $(info "Using CIRCLE_WORKING_DIRECTORY for GOPATH")
endif

ifneq ($(OS),Windows_NT)
  # Before we start test that we have the mandatory executables available
  EXECUTABLES = go
  OK := $(foreach exec,$(EXECUTABLES),\
    $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH, please install $(exec)")))
endif

export CGO_ENABLED=0
export GO111MODULE=on
export GOFLAGS=-mod=vendor

help:
	@echo "Available options:"
	@grep -E '^[/1-9a-zA-Z._%-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	   | sort \
	   | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-45s\033[0m %s\n", $$1, $$2}'
.PHONY: help

clean: ## Remove build artifacts
	@git clean -fdX
.PHONY: clean

fmt: ## Reformat package sources
	@go fmt ./...
.PHONY: fmt

vet: fmt
	@go vet ./...
.PHONY: vet

imports: ## Ensure imports are present and formatted
	@goimports -w $(shell find . -type f -name '*.go' -not -path './vendor/*')
.PHONY: goimports

deps: ## Install depdendencies. Runs `go get` internally.
	@GOFLAGS="" go get -t -v ./...
	@GOFLAGS="" go mod tidy
	@GOFLAGS="" go mod vendor


update-deps: ## Update depdendencies. Runs `go get -u` internally.
	@GOFLAGS="" go get -t -u ./...
	@GOFLAGS="" go mod tidy
	@GOFLAGS="" go mod vendor

build: deps vet schema ## Build the package
	@go build -o helmsman -ldflags '-X github.com/mkubaczyk/helmsman/internal/app.appVersion="${TAG}-${DATE}" -extldflags "-static"' cmd/helmsman/main.go

generate:
	@go generate #${PKGS}
.PHONY: generate

schema: ## Generate the schema.json file
	@go run schema.go

repo:
	@helm repo list | grep -q "^prometheus-community " || helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	@helm repo update
.PHONY: repo

test: deps vet schema repo ## Run unit tests
	@go test -v -cover -p=1 ./... -args -f ../../examples/example.toml
.PHONY: test

test-%: deps vet repo ## Run a specific unit test
	@go test -v -cover -p=1 -run $(*) ./... -args -f ../../examples/example.toml

cross: deps ## Create binaries for all OSs
	@gox -os 'windows linux darwin' -arch 'arm64 amd64' -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags '-X github.com/mkubaczyk/helmsman/internal/app.appVersion=${TAG}-${DATE}' ./...
.PHONY: cross

release: ## Generate a new release
	@goreleaser --release-notes release-notes.md --rm-dist

tools: ## Get extra tools used by this makefile
	@go get -d -u github.com/golang/dep/cmd/dep
	@go get -d -u github.com/mitchellh/gox
	@go get -d -u github.com/goreleaser/goreleaser
	@gem install hiera-eyaml
.PHONY: tools

helmPlugins: ## Install helm plugins used by Helmsman
	@mkdir -p ~/.helm/plugins
	@helm plugin install --verify=false https://github.com/hypnoglow/helm-s3.git
	@helm plugin install --verify=false https://github.com/nouney/helm-gcs
	@helm plugin install --verify=false https://github.com/databus23/helm-diff
	@helm plugin install --verify=false https://github.com/jkroepke/helm-secrets
.PHONY: helmPlugins
