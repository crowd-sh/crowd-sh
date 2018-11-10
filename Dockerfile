FROM golang:latest

WORKDIR /go/src/workmachine

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...