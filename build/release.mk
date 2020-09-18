BUILD_DIR    := ./bin/
GORELEASER_VERSION := v0.143.0
GORELEASER_BIN ?= bin/goreleaser

bin:
	@mkdir -p $(BUILD_DIR)

$(GORELEASER_BIN): bin
	@echo "=== $(INTEGRATION) === [$(GORELEASER_BIN)] Installing goreleaser $(GORELEASER_VERSION)"
	@(wget -qO /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASER_VERSION)/goreleaser_$(OS_DOWNLOAD)_x86_64.tar.gz)
	@(tar -xf  /tmp/goreleaser.tar.gz -C bin/)
	@(rm -f /tmp/goreleaser.tar.gz)
	@echo "=== $(INTEGRATION) === [$(GORELEASER_BIN)] goreleaser downloaded"

.PHONY : release/deps
release/deps: $(GORELEASER_BIN)

.PHONY : release
release: release/deps
ifeq ($(PRERELEASE), true)
	@echo "pre-release"
	@$(GORELEASER_BIN) release --config $(CURDIR)/build/.goreleaser.yml --rm-dist
else
	@echo "build"
	# this will just compile binaries
	@$(GORELEASER_BIN) build --config $(CURDIR)/build/.goreleaser.yml --snapshot --rm-dist
endif

.PHONY : release/fix-archive
release/fix-archive:
	@echo "=== $(INTEGRATION) === [release/fix-archive] fixing archives internal structure"

.PHONY : release/sign
release/sign:
	@echo "=== $(INTEGRATION) === [release/sign] signing packages"
	@bash $(CURDIR)/build/sign.sh


.PHONY : release/publish
release/publish: release/fix-archive release/sign
	@echo "=== $(INTEGRATION) === [release/publish] publishing artifacts"



OS := $(shell uname -s)
ifeq ($(OS), Darwin)
	OS_DOWNLOAD := "darwin"
	TAR := gtar
else
	OS_DOWNLOAD := "linux"
endif
