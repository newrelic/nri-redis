FROM golang:1.24.6-bookworm as builder
ARG CGO_ENABLED=0
WORKDIR /go/src/github.com/newrelic/nri-redis
COPY . .
RUN make clean compile

FROM newrelic/infrastructure-bundle:latest

COPY --from=builder /go/src/github.com/newrelic/nri-redis/bin /var/db/newrelic-infra/newrelic-integrations/bin/
COPY tests/integration/newrelic-infra.yml /etc/newrelic-infra.yml
COPY tests/integration/redis-config.yml /etc/newrelic-infra/integrations.d/redis-config.yml