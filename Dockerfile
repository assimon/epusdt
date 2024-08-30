FROM golang:1.21 as builder
RUN mkdir /go/app
WORKDIR /go/app

COPY src/ .

RUN go mod download
RUN go build -o /app/bin/epusdt ./main.go

FROM alpine:latest
COPY --from=builder /app/bin/epusdt /app/bin/epusdt

RUN apk update && apk upgrade && apk add ca-certificates && update-ca-certificates
# Change TimeZone
RUN apk add --update tzdata
# Clean APK cache
RUN rm -rf /var/cache/apk/*

ENTRYPOINT [ "/app/bin/epusdt" ]