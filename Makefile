docker_image  := helm-hub-sync
build_dir     := $(CURDIR)/bin
dist_dir      := $(CURDIR)/dist
github_repo   := jfrog/helm-hub-sync
version       := $(shell git describe --tags --always --dirty="-dev")
date          := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
version_flags := -ldflags='-X "main.version=$(version)" -X "main.buildTime=$(date)"'
os_archs      := "darwin/amd64 linux/amd64 windows/amd64"

# GOPROXY defines which URL to use to retrieve Go Modules from
GOPROXY=https://gocenter.io

# Sets PWD to pwd_unknown if it doesn't have a value. Normally this should not
# happen. If you do see pwd_unknown showing up, you'll need to make sure your
# system understand the PWD command.
PWD ?= pwd_unknown

# List all PHONY targets
.PHONY: help build build-docker build-dist install tags list clean clean-build clean-dist deps score setup test

# Help
help: ## Displays the help for each target (this message).
	@echo 
	@echo Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo

# Build
build: ## builds the executable
	@echo "Building..."
	$Q rm -f $(build_dir)
	$Q go build $(if $V,-v) $(version_flags)

build-docker: ## builds a docker container
	@echo "Building..."
	$Q docker build -t $(docker_image):$(version) .

build-dist: ## builds executables for all environments
	@echo "Building..."
	$Q rm -rf $(dist_dir)
	$Q gox -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch=$(os_archs) $(version_flags)

test: ## Run all testcases.
	env TESTDIR=${PWD}/test go test ./...

# Install
install: ## installs the executable in your GOPATH
	@echo "Building..."
	$Q rm -f $(GOPATH)/bin/$(docker_image)
	$Q go install $(if $V,-v) $(version_flags)

# Clean
clean: clean-dist clean-build ## cleans all generated files

clean-build: ## deletes all builds
	@echo "Removing cross-compilation files..."
	$Q rm -rf $(build_dir)

clean-dist: ## deletes all distributions
	@echo "Removing distribution files..."
	$Q rm -rf $(dist_dir)

# Utilities
tags: ## list the git tags
	@echo "Listing tags..."
	$Q @git tag

list: ## display a list of all modules in this app and all required modules
	@echo "Listing modules..."
	$Q go list ./...
	$Q go mod graph

deps: ## gets all dependencies
	@echo "Getting all modules..."
	$Q go get ./...

score: ## gets a score based on GoReportcard.
	goreportcard-cli -v

setup: clean ## creates folders and downloads tools
	@echo "Setup..."
	$Q if ! grep "/bin" .gitignore > /dev/null 2>&1; then \
		echo "/bin" >> .gitignore; \
	fi
	$Q if ! grep "/dist" .gitignore > /dev/null 2>&1; then \
		echo "/dist" >> .gitignore; \
	fi
	$Q if ! grep "/cover" .gitignore > /dev/null 2>&1; then \
		echo "/cover" >> .gitignore; \
	fi
	$Q if ! grep "/test" .gitignore > /dev/null 2>&1; then \
		echo "/test" >> .gitignore; \
	fi
	$Q mkdir -p cover
	$Q mkdir -p bin
	$Q mkdir -p test
	$Q mkdir -p dist
	$Q go get github.com/mitchellh/gox
	$Q go get github.com/gojp/goreportcard/cmd/goreportcard-cli

# comment this line out for quieter things
V := 1
Q := $(if $V,,@)
