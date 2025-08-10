FROM golang:1.24.5-alpine3.22 AS builder

ENV SERVICE=sberdevices

ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/$SERVICE

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git

COPY go.mod go.sum /go/src/$SERVICE/
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . /go/src/$SERVICE/
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 go build -gcflags "all=-N -l" -o $SERVICE cmd/sberhack/main.go

FROM scratch

ENV SERVICE=sberdevices

COPY --from=builder /go/bin/dlv /dlv
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/$SERVICE/config/config.yaml /config.yaml
COPY --from=builder /go/src/$SERVICE/$SERVICE /$SERVICE

EXPOSE 2345 2345
ENTRYPOINT ["/dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./sberdevices"]

CMD []
