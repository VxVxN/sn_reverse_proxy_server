FROM golang:1.16

ENV GOFLAGS="-mod=vendor" GO111MODULE=on

ADD . /app
WORKDIR /app

RUN go build -o app/reverse_proxy_server ./cmd/reverse_proxy_server

RUN cd app

CMD ["./app/reverse_proxy_server"]