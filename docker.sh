#!/bin/bash

export GOPATH=$(pwd)/../../../../

export CGO_ENABLED=0

export GOOS=linux

export GOARCH=amd64

go build -o bin/magneto ./cmd/magneto

docker build -t weiwen/magneto -f ./build/Dockerfile .

docker-compose -f ./deployments/docker-compose.yml up -d
