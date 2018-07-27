INTEGRATION     := redis
BINARY_NAME      = nr-$(INTEGRATION)
SRC_DIR          = ./src/
GO_PKGS         := $(shell go list ./... | grep -v "/vendor/")
GO_FILES        := $(shell find src -type f -name "*.go")
VALIDATE_DEPS    = github.com/golang/lint/golint
TEST_DEPS        = github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml
INTEGRATIONS_DIR = /var/db/newrelic-infra/newrelic-integrations/
CONFIG_DIR       = /etc/newrelic-infra/integrations.d

all: build

build: clean validate compile test

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

validate-deps:
	@echo "=== $(INTEGRATION) === [ validate-deps ]: installing validation dependencies..."
	@go get -v $(VALIDATE_DEPS)

validate-only:
	@printf "=== $(INTEGRATION) === [ validate ]: running gofmt... "
# `gofmt` expects files instead of packages. `go fmt` works with
# packages, but forces -l -w flags.
	@OUTPUT="$(shell gofmt -l $(GO_FILES))" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Incorrect syntax in the following files:" ;\
		echo "$$OUTPUT" ;\
		exit 1 ;\
	fi
	@printf "=== $(INTEGRATION) === [ validate ]: running golint... "
	@OUTPUT="$(shell golint $(SRC_DIR)...)" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Issues found:" ;\
		echo "$$OUTPUT" ;\
		exit 1 ;\
	fi
	@printf "=== $(INTEGRATION) === [ validate ]: running go vet... "
	@OUTPUT="$(shell go vet $(SRC_DIR)...)" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Issues found:" ;\
		echo "$$OUTPUT" ;\
		exit 1;\
	fi

validate: validate-deps validate-only

compile-deps:
	@echo "=== $(INTEGRATION) === [ compile-deps ]: installing build dependencies..."
	@go get -v -d -t ./...

compile-only:
	@echo "=== $(INTEGRATION) === [ compile ]: building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

compile: compile-deps compile-only

test-deps: compile-deps
	@echo "=== $(INTEGRATION) === [ test-deps ]: installing testing dependencies..."
	@go get -v $(TEST_DEPS)

test-only:
	@echo "=== $(INTEGRATION) === [ test ]: running unit tests..."
	@gocov test $(SRC_DIR)/... | gocov-xml > coverage.xml

test: test-deps test-only

install: bin/$(BINARY_NAME)
	@echo "=== $(INTEGRATION) === [ install ]: installing bin/$(BINARY_NAME)..."
	@sudo install -D --mode=755 --owner=root --strip $(ROOT)bin/$(BINARY_NAME) $(INTEGRATIONS_DIR)/bin/$(BINARY_NAME)
	@sudo install -D --mode=644 --owner=root $(ROOT)$(INTEGRATION)-definition.yml $(INTEGRATIONS_DIR)/$(INTEGRATION)-definition.yml
	@sudo install -D --mode=644 --owner=root $(ROOT)$(INTEGRATION)-config.yml.sample $(CONFIG_DIR)/$(INTEGRATION)-config.yml.sample

.PHONY: all build clean validate-deps validate-only validate compile-deps compile-only compile test-deps test-only test
