FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN apk update && apk add make
RUN go mod download
RUN make build

FROM alpine:3.22 AS release

LABEL maintainer="oka4shi"
WORKDIR /app

COPY --from=builder /app/bin/kusamochi .

RUN chmod +x ./kusamochi

ENTRYPOINT /app/kusamochi
