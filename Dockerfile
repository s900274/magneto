FROM golang:1.10.4-alpine AS build-env
RUN apk add --no-cache --update git openssl bzr build-base \
    && go get -u github.com/kardianos/govendor && go get github.com/s900274/magneto \
    && cd github.com/s900274/magneto \
    && govendor sync && go build -o bin/magneto github.com/s900274/magneto/cmd/magneto




MAINTAINER s900274

