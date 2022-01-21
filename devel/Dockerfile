FROM golang:1.17-alpine AS dlv-builder

RUN apk add gcc musl-dev && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# ########################################################
FROM newrelic/infrastructure:1.22.0

COPY --from=dlv-builder /go/bin/dlv /usr/local/bin/
COPY ./bin/ /var/db/newrelic-infra/newrelic-integrations/bin/

ENTRYPOINT [ "dlv", "--listen=0.0.0.0:50100", "--headless=true", "--api-version=2", "--check-go-version=false", "--only-same-user=false", "exec", "/var/db/newrelic-infra/newrelic-integrations/bin/nri-redis", "--" ]
