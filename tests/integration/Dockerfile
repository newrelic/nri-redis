FROM golang:1.18 as builder
ARG CGO_ENABLED=0
WORKDIR /go/src/github.com/newrelic/nri-redis
COPY . .
RUN make clean compile

FROM alpine:latest
COPY --from=builder /go/src/github.com/newrelic/nri-redis/bin /
CMD ["sleep", "1h"]
