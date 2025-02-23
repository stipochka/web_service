# syntax=docker/dockerfile:1.4
FROM golang:1.24.0 AS builder
#export DOCKER_BUILDKIT=1

WORKDIR /app

ENV GOPRIVATE=github.com/stipochka/*

COPY go.mod go.sum ./

RUN --mount=type=ssh \
    mkdir -p /root/.ssh \
    && ssh-keyscan github.com >> /root/.ssh/known_hosts \
    && git config --global url."git@github.com:".insteadOf "https://github.com" \
    && go mod download && go mod verify

COPY . .

WORKDIR /app/cmd/web_service

RUN go build -o web_service

FROM debian:bookworm-slim 

WORKDIR /app

COPY --from=builder /app/cmd/web_service .
COPY --from=builder /app/config .

# Add this to your Dockerfile to install wait-for-it
RUN apt-get update && apt-get install -y wait-for-it
ENV CONFIG_PATH=/app/config.yaml

CMD ["wait-for-it", "kafka:9092", "--", "./web_service"]

