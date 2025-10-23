# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download || true
COPY . .
RUN go mod tidy && go build -o mock ./cmd/mock

FROM alpine:3.20
WORKDIR /app
COPY --from=base /app/mock /app/mock
EXPOSE 8090
ENTRYPOINT ["/app/mock"]
