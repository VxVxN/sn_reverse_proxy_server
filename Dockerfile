FROM golang:1.16

ENV GOFLAGS="-mod=vendor" GO111MODULE=on

ADD . /app
WORKDIR /app

RUN go build -o app/reverse_proxy_server ./app

RUN cd app

CMD ["./app/reverse_proxy_server"]