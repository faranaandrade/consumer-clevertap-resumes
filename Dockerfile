# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as builder
LABEL "com.occ.vendor"="Occ"
LABEL "version"="2023.11.07"

WORKDIR /go/src/consumer-clevertap-applies

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk update && \
    apk add --no-cache tzdata && \
    apk add ca-certificates && rm -rf /var/cache/apk/* && \
    GOOS=linux \
    GOARCH=amd64 \
    go build  -tags musl,appsec  -o consumer-clevertap-applies-app  ./cmd/consumer

FROM alpine:latest

ENV PROJECT_VERSION="1.0.0"
ENV TZ="America/Mexico_City"

COPY --from=builder /go/src/consumer-clevertap-applies/consumer-clevertap-applies-app .

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

CMD [ "/consumer-clevertap-applies-app"]