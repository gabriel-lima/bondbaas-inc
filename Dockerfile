FROM golang:1.12-alpine

RUN apk add --update --no-cache git bash && \
  rm -rf /tmp/* /var/cache/apk/*

ADD . /go/src/bondbaas
WORKDIR /go/src/app

RUN go get github.com/lib/pq
RUN go get github.com/pilu/fresh

RUN go install bondbaas

EXPOSE 3000