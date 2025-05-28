GO_PATH := ./
BASE_DIR := src
BIN_NAME := bootstrap
LAMBDA_DIRS := api/post/create api/post/update api/post/getbyid api/post/getall

all: deps build

deps:
	@for dir in $(GO_PATH); do \
		echo "Installing dependencies in $$dir..."; \
		cd $$dir && go mod tidy; \
	done

build:
	@for dir in $(LAMBDA_DIRS); do \
		echo "Building $(BASE_DIR)/$$dir..."; \
		(cd $(BASE_DIR)/$$dir && env GOARCH=amd64 GOOS=linux go build -o ./build/$(BIN_NAME)); \
	done

clean:
	@for dir in $(LAMBDA_DIRS); do \
		echo "Cleaning $(BASE_DIR)/$$dir..."; \
		(cd $(BASE_DIR)/$$dir && rm -rf build); \
	done

.PHONY: all deps build clean
