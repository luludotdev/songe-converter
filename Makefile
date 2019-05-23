# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

OUTPUT=build
BINARY_NAME=songe-converter
BINARY_WIN=$(BINARY_NAME).exe
BINARY_MAC=$(BINARY_NAME)-mac

all: build-win build-linux build-mac

build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ./$(OUTPUT)/$(BINARY_WIN) -v ./src

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ./$(OUTPUT)/$(BINARY_NAME) -v ./src

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ./$(OUTPUT)/$(BINARY_MAC) -v ./src

clean:
	$(GOCLEAN)
	rm -rf $(OUTPUT)

deps:
	$(GOGET) github.com/bmatcuk/doublestar
	$(GOGET) github.com/TomOnTime/utfutil
