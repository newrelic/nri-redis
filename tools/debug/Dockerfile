FROM golang:1.17-alpine AS dlv-builder

RUN apk add gcc musl-dev && \
    go install github.com/go-delve/delve/cmd/dlv@latest

COPY ./bin/ /go/bin/

# ENTRYPOINT overwritten by tilt