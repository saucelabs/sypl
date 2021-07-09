###
# Provides an env for CI to run tests, and lint.
###

FROM golang:1.15-alpine

RUN apk update && apk add curl \
     git \
     bash \
     py-pip \
     musl-dev \
     gcc \
     make && \
     rm -rf /var/cache/apk/*

RUN pip install bumpversion

# Go test deps
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1
RUN mv bin/golangci-lint /usr/local/bin/
