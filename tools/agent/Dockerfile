FROM golang:1.17-alpine AS dlv-builder

RUN apk add gcc musl-dev && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# ########################################################
FROM newrelic/infrastructure:1.22.0

COPY --from=dlv-builder /go/bin/dlv /usr/local/bin/
COPY ./bin/ /var/db/newrelic-infra/newrelic-integrations/bin/

# ENTRYPOINT overwritten by tilt
