FROM golang:1.10.4-alpine

ARG environment=testoce
ARG zookeeperlist

MAINTAINER s900274

ENV ENVIRONMENT=$environment
ENV PRJ_START_DELAY=5

COPY ./bin magneto/bin
COPY ./api magneto/api
COPY ./config magneto/config
COPY ./web magneto/web

ENTRYPOINT cd magneto \
    		&& mkdir -p status/magneto \
    		&& mkdir -p log \
    		&& sleep ${PRJ_START_DELAY}s \
    		&& ./bin/magneto -config ./config/app.toml