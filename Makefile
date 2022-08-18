METADATA_PKG=github.com/shipt/tempest-template/internal/metadata

METADATA_BUILD_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
METADATA_BUILD_HASH=$(shell git rev-parse HEAD)
METADATA_BUILD_TAG=tag
METADATA_SERVICE_NAME=tempest-template

LD_FLAGS= -linkmode external \
			-X '$(METADATA_PKG).buildBranch=$(METADATA_BUILD_BRANCH)' \
			-X '$(METADATA_PKG).buildHash=$(METADATA_BUILD_HASH)' \
			-X '$(METADATA_PKG).buildTag=$(METADATA_BUILD_TAG)' \
			-X '$(METADATA_PKG).serviceName=$(METADATA_SERVICE_NAME)'
 
BUILD_DIR=bin
BUILD_ENV=CGO_ENABLED=1 CC_FOR_TARGET=gcc CXX_FOR_TARGET=g++
BUILD_OUT=$(BUILD_DIR)/tempest-template
BUILD_FLAGS=-a -o $(BUILD_OUT) -ldflags="$(LD_FLAGS)"

INSTALL_BIN_PREFIX=/usr/local

TEST_COVERAGE_OUT=$(BUILD_DIR)/coverage.out
TEST_FLAGS=-p 1 -count 1 -coverpkg=./... -coverprofile=$(TEST_COVERAGE_OUT)
TEST_PKGS=$(shell go list ./...)

.PHONY: all clean test install run
all:
	mkdir -p $(BUILD_DIR)
	$(BUILD_ENV) go build $(BUILD_FLAGS) ./cmd/tempest-template

clean:
	go clean
	rm -rf $(BUILD_DIR)

test:
	mkdir -p $(BUILD_DIR)
	go test $(TEST_FLAGS) $(TEST_PKGS)

install: 
	cp $(BUILD_OUT) $(INSTALL_BIN_PREFIX)/bin/

run:
	docker compose up -d && \
	go run cmd/tempest-template/main.go webserver

include migrations.mk
