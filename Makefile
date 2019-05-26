# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

GIT_HASH=`git rev-parse HEAD`
GIT_TAG=`git tag --points-at HEAD`
OUTPUT=build
BINARY_NAME=songe-converter
BINARY_WIN=$(BINARY_NAME).exe
BINARY_MAC=$(BINARY_NAME)-mac
BUILD_FLAGS=-ldflags="-s -w -X main.sha1ver=$(GIT_HASH) -X main.gitTag=$(GIT_TAG)"

all: build-win build-linux build-mac

build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(BINARY_WIN) -v ./src

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(BINARY_NAME) -v ./src

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o ./$(OUTPUT)/$(BINARY_MAC) -v ./src

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
