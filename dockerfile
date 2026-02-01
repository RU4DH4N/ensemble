FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git protobuf-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN protoc --proto_path=pkg/protocol/v1 \
    --go_out=pkg/protocol/v1 --go_opt=paths=source_relative \
    --go-grpc_out=pkg/protocol/v1 --go-grpc_opt=paths=source_relative \
    pkg/protocol/v1/ensemble.proto

ARG VERSION=dev
ARG GOOS=linux
ARG GOARCH=amd64

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build -ldflags "-X main.version=${VERSION} -w -s" \
    -o /app/bin/ensembled \
    ./cmd/ensembled

FROM alpine:latest

RUN apk add --no-cache ca-certificates

RUN addgroup -g 1000 ensemble && \
    adduser -D -u 1000 -G ensemble ensemble

WORKDIR /app

COPY --from=builder /app/bin/ensembled /app/ensembled
COPY --chown=ensemble:ensemble example.json /app/example.json

RUN mkdir -p /data/uploads && \
    chown -R ensemble:ensemble /data

USER ensemble

VOLUME ["/data/uploads"]

ENTRYPOINT ["/app/ensembled"]