# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

GIT_HASH=`git rev-parse HEAD`
GIT_TAG=`git tag --points-at HEAD`
OUTPUT=build
BUILD_FLAGS=-ldflags="-s -w -X main.sha1ver=$(GIT_HASH) -X main.gitTag=$(GIT_TAG)"

SONGE_BINARY_NAME=songe-converter
SONGE_BINARY_WIN=$(SONGE_BINARY_NAME).exe
SONGE_BINARY_MAC=$(SONGE_BINARY_NAME)-mac
SIMPLE_BINARY_NAME=simple-converter
SIMPLE_BINARY_WIN=$(SIMPLE_BINARY_NAME).exe
SIMPLE_BINARY_MAC=$(SIMPLE_BINARY_NAME)-mac

all: build-songe build-simple
build-songe: build-songe-win build-songe-linux build-songe-mac
build-simple: build-simple-win build-simple-linux build-simple-mac

# Build songe-converter
build-songe-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(SONGE_BINARY_WIN) -v ./cmd/songe-converter

build-songe-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(SONGE_BINARY_NAME) -v ./cmd/songe-converter

build-songe-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(SONGE_BINARY_MAC) -v ./cmd/songe-converter

# Build simple-converter
build-simple-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(SIMPLE_BINARY_WIN) -v ./cmd/simple-converter

build-simple-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(SIMPLE_BINARY_NAME) -v ./cmd/simple-converter

build-simple-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(SIMPLE_BINARY_MAC) -v ./cmd/simple-converter

clean:
	$(GOCLEAN)
	rm -rf $(OUTPUT)

release:
	@read -p "Enter version: " version; \
	git tag "v$$version" && \
	git push && \
	git push --tags && \
	make clean && \
	make

deps:
	$(GOGET) github.com/bmatcuk/doublestar
	$(GOGET) github.com/TomOnTime/utfutil
	$(GOGET) github.com/ttacon/chalk
	$(GOGET) github.com/otiai10/copy
