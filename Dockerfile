FROM golang:1.14.3-alpine3.11 AS builder

WORKDIR /go/src/github.com/sasalatart/api-demo/

COPY main.go pong.go ./
RUN GOOS=linux go build -o app .

###

FROM alpine:latest

LABEL maintainer="Sebastian Salata R-T <sa.salatart@gmail.com>"

WORKDIR /root/

COPY --from=builder /go/src/github.com/sasalatart/api-demo/app .
CMD ["./app"]
