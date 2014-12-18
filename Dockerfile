FROM golang

MAINTAINER Abhi Yerra <abhi@berkeley.edu>

ADD . /go/src/github.com/workmachine/workmachine

RUN cd /go/src/github.com/workmachine/workmachine && go get ./...
RUN go install github.com/workmachine/workmachine

WORKDIR /go/src/github.com/workmachine/workmachine

ENTRYPOINT /go/bin/workmachine

EXPOSE 3001
