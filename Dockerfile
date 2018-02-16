# ============================================================
# Dockerfile used to build the various microservices
#
# To build run something like:
#
# docker build -t mytag --build-arg service=darwinref .
#
# where the value of service is:
#   darwinref     The Darwin Reference API
# ============================================================

# ============================================================
# gcc container required to build the docker-entrypoint.
# See bin/docker/main.c for why we need this
FROM alpine as gcc
RUN apk add --no-cache gcc musl-dev

# ============================================================
# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
FROM golang:alpine as build

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      git

# We want to build our final image under /dest
# A copy of /etc/ssl is required if we want to use https datasources
RUN mkdir -p /dest/etc &&\
    cp -rp /etc/ssl /dest/etc/

# Ensure we have the libraries - docker will cache these between builds
RUN go get -v \
      flag \
      github.com/coreos/bbolt/... \
      github.com/gorilla/mux \
      github.com/jlaffaye/ftp \
      github.com/muesli/cache2go \
      github.com/peter-mount/golib/codec \
      github.com/peter-mount/golib/rabbitmq \
      github.com/peter-mount/golib/rest \
      github.com/peter-mount/golib/statistics \
      github.com/peter-mount/golib/util \
      gopkg.in/robfig/cron.v2 \
      gopkg.in/yaml.v2 \
      io/ioutil \
      log \
      net/http \
      path/filepath \
      time

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /go/src
ADD . .

# ============================================================
# Run all tests in a new container so any output won't affect
# the final build
FROM source as test
RUN go test -v util

# ============================================================
# Now test have past build the docker-entrypoint
# see bin/docker/main.c for why we need this
FROM gcc as wrapper
ARG service
WORKDIR /work
ADD bin/docker .
RUN sed -i "s/@@service@@/${service}/g" main.c &&\
    gcc -o main -static main.c &&\
    strip main

# ============================================================
# Compile the source.
FROM source as compiler
ARG service

# Static compile
ENV CGO_ENABLED=0
ENV GOOS=linux

# Microservice version is the branch & commit hash from git
RUN branch="$(git branch --no-color|cut -f2- -d' ')" &&\
    version="$(git rev-parse --short "$branch")" &&\
    sed -i "s/@@version@@/${version}(${branch})/g" bin/version.go

# Build the microservice
RUN go build -o /dest/${service} bin/${service}

# Install the docker-entrypoint
COPY --from=wrapper /work/main /dest/docker-entrypoint

# ============================================================
# Finally build the final runtime container for the specific
# microservice
FROM scratch

# The default database directory
Volume /database

# Install our built image
COPY --from=compiler /dest/ /

ENTRYPOINT ["/docker-entrypoint"]
CMD [ "-c", "/config.yaml"]
