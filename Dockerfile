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

# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
FROM golang:latest as build

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
# compiler container used just for this build. Changing the

FROM build as compiler
ARG service

# Static compile
ENV CGO_ENABLED=0
ENV GOOS=linux

# Import the source and compile
WORKDIR /go/src
ADD . .

# Run any tests
RUN go test -v util

# Build the microservice
#RUN go build -v -x -o /dest/bin/${service} bin/${service}
RUN go build -o /dest/bin/${service} bin/${service}

# The docker entrypoint
RUN cd /dest && \
    ln -s bin/${service} docker-entrypoint

# ============================================================
# Finally build the final runtime container for the specific
# microservice
FROM scratch
ARG service
COPY --from=compiler /dest/ /
CMD ["/docker-entrypoint", "-c", "/config.yaml"]
