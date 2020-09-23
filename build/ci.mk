BUILDER_TAG ?= nri-$(INTEGRATION)-builder

.PHONY : ci/deps
ci/deps:
	@docker build -t $(BUILDER_TAG) -f $(CURDIR)/build/Dockerfile $(CURDIR)

.PHONY : ci/validate
ci/validate: ci/deps
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis $(BUILDER_TAG) make validate

.PHONY : ci/test
ci/test: ci/deps
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis $(BUILDER_TAG) make test

.PHONY : ci/build
ci/build: ci/deps
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis $(BUILDER_TAG) make release/compile

.PHONY : ci/prerelease
ci/prerelease: ci/deps
ifdef TAG
	@docker run --rm -t \
			-v $(CURDIR):/go/src/github.com/newrelic/nri-redis \
			-w /go/src/github.com/newrelic/nri-redis \
			-e PRERELEASE=true \
			-e GITHUB_TOKEN=$(GH_TOKEN) \
			-e TAG \
			-e GPG_MAIL \
			-e GPG_PASSPHRASE \
			-e GPG_PRIVATE_KEY \
			$(BUILDER_TAG) make release
else
	@echo "TAG env variable expected to be set"
endif