FROM golang:1.21
RUN mkdir /go/app
WORKDIR /go/app

COPY src/ .

RUN go mod download
RUN go build -o /app/bin/qpusdt ./main.go

FROM alpine:latest
COPY --from=0 /app/bin/qpusdt /app/bin/qpusdt