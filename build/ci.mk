.PHONY : ci/validate

ci/validate:
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis golang:1.9 make validate

.PHONY : ci/test

ci/test:
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis golang:1.9 make test

.PHONY : ci/build

ci/build:
	@docker run --rm -t -v $(CURDIR):/go/src/github.com/newrelic/nri-redis -w /go/src/github.com/newrelic/nri-redis golang:1.9 make release
