FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git upx

WORKDIR /builder
COPY go.mod .
COPY go.sum .
COPY . .

RUN go mod download \
    && go build -ldflags "-s -w" -o cli \
    && upx -9 cli

FROM alpine:latest
ARG CONFIG
ENV config=$CONFIG
WORKDIR /app
COPY --from=builder /builder/cli .
ENTRYPOINT ["./cli"]