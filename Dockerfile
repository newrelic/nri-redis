FROM golang:1.18 as builder
COPY . /go/src/github.com/newrelic/nri-redis/
RUN cd /go/src/github.com/newrelic/nri-redis && \
    make compile && \
    strip ./bin/nri-redis

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-redis/bin/nri-redis /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-redis

USER 1000
