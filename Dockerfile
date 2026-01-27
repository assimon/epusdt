FROM golang:alpine AS builder

RUN apk add --no-cache --update git build-base
ENV CGO_ENABLED=1
WORKDIR /app
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .
RUN  go mod tidy \
	&& go build -ldflags "-s -w"  \
	-o epusdt .

FROM alpine:latest AS runner
ENV TZ=Asia/Shanghai
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/static /app/static
COPY --from=builder /app/static /static
COPY --from=builder /app/epusdt .
VOLUME /app/conf

ENTRYPOINT ["./epusdt" ,"http","start"]
