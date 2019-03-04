FROM golang:1.10 as builder
RUN go get -d github.com/newrelic/nri-redis/... && \
    cd /go/src/github.com/newrelic/nri-redis && \
    make compile && \
    strip ./bin/nr-redis

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-redis/bin/nr-redis /var/db/newrelic-infra/newrelic-integrations/bin/nr-redis
COPY --from=builder /go/src/github.com/newrelic/nri-redis/redis-definition.yml /var/db/newrelic-infra/newrelic-integrations/redis-definition.yml

