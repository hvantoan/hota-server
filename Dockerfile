# Builder
FROM golang:1.11.4-alpine3.8 as builder

RUN apk update && apk upgrade && \
    apk --update add git gcc make && \
    go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/hota-server

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8888

COPY --from=builder /go/src/hota-server/engine /app

CMD /app/engine
