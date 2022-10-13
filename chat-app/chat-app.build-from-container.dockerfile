FROM golang:alpine3.16 AS builder

COPY . /build
WORKDIR /build

RUN go mod download
RUN env CGO_ENABLED=0 go build -o './bin/chat-app' ./cmd/gigachat/*.go

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=0 /build/bin/chat-app .
COPY --from=0 /build/configs ./configs

CMD [ "/app/chat-app" ]