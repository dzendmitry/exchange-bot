FROM golang:alpine

RUN mkdir -p /go/src/github.com/dzendmitry/exchange-bot
ADD . /go/src/github.com/dzendmitry/exchange-bot
RUN \
 cd /go/src/github.com/dzendmitry/exchange-bot && \
 go build -o /srv/exchange-bot && \
 rm -rf /go/src/*

EXPOSE 8080
WORKDIR /srv
CMD ["/srv/exchange-bot"]