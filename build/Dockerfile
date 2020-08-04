FROM golang:1.14 as builder


WORKDIR /kbc
COPY . /kbc

ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN go build -mod=vendor -v -o chatbot 

# RUN go mod tidy && go mod download
# RUN go build -mod=readonly -v -o chatbot

FROM alpine:3

RUN apk add --no-cache ca-certificates

RUN mkdir -p /kbc
WORKDIR /kbc

COPY --from=builder /kbc/chatbot /kbc/chatbot
COPY --from=builder /kbc/config /kbc/config
COPY --from=builder /kbc/pubsub_credentials.json /kbc/pubsub_credentials.json
COPY --from=builder /kbc/build kdc/kbc_bin

ENTRYPOINT ./chatbot