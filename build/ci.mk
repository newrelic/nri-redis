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
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis $(BUILDER_TAG) make release

.PHONY : ci/prerelease
ci/prerelease: ci/deps
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis -e PRERELEASE=true $(BUILDER_TAG) make release
