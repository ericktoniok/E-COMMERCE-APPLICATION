# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache git build-base
COPY go.mod go.sum ./
RUN go mod download || true
COPY . .
RUN go mod tidy && go build -o app ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN adduser -D appuser
USER appuser
COPY --from=base /app/app /app/app
COPY --from=base /app/migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["/app/app"]
