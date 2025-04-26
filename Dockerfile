FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN make build

FROM alpine:3.13 AS release

LABEL maintainer="oka4shi"
WORKDIR /app

COPY --from=builder /app/bin/kusamochi .

RUN chmod +x ./kusamochi

ENV KUSAMOCHI_GITHUB_TOKEN=${KUSAMOCHI_GITHUB_TOKEN}
ENV KUSAMOCHI_WEBHOOK_URL=${KUSAMOCHI_WEBHOOK_URL}
ENTRYPOINT /app/kusamochi
