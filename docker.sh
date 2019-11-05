#!/bin/bash

start(){
  export GO111MODULE=on
  export CGO_ENABLED=0
  export GOOS=linux
  export GOARCH=amd64

  go build -o bin/magneto ./cmd/magneto

  docker build -t weiwen/magneto -f ./build/Dockerfile .

  docker-compose -f ./deployments/docker-compose.yml up -d
}

stop(){
  docker-compose -f ./deployments/docker-compose.yml down
}

case C"$1" in
  Cstart)
    start
    echo "start Done!"
    ;;
  Cstop)
    stop
    echo "stop Done!"
    ;;
  Crestart)
    stop
    start
    echo "restart Done!"
    ;;
  C*)
    echo "Usage: $0 {start|stop|restart}"
    ;;
esac