FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD cmd /app/cmd
ADD internal /app/internal
ADD pkg /app/pkg
ADD main.go /app/

RUN go build -o app

FROM alpine:latest

RUN apk update
RUN apk add bash ncurses
ENV TERM=xterm-256color

RUN apk --no-cache add --no-check-certificate ca-certificates \
    && update-ca-certificates

RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /home/appuser/

COPY --from=builder /app/app .

RUN chown appuser:appuser app
USER appuser

ENTRYPOINT ["./app"]
