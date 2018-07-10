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

ARG arch=amd64
ARG goos=linux

# ============================================================
# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
FROM golang:alpine as build

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      git \
      tzdata

# Our build scripts
ADD scripts/ /usr/local/bin/

# Ensure we have the libraries - docker will cache these between builds
RUN get.sh

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /go/src/github.com/peter-mount/nre-feeds
ADD . .

# ============================================================
# Run all tests in a new container so any output won't affect
# the final build.
FROM source as test
ARG skipTest
RUN if [ -z "$skipTest" ] ;then test.sh; fi

# ============================================================
# Compile the source.
FROM source as compiler
ARG service
ARG arch
ARG goos
ARG goarch
ARG goarm

# Microservice version is the commit hash from git
#RUN version="$(git rev-parse --short HEAD)" &&\
#    sed -i "s/@@version@@/${version} ${goos}(${arch})/g" bin/version.go

# Build the microservice.
# NB: CGO_ENABLED=0 forces a static build
RUN CGO_ENABLED=0 \
    GOOS=${goos} \
    GOARCH=${goarch} \
    GOARM=${goarm} \
    compile.sh /dest/${service} ${service}

# ============================================================
# Finally build the final runtime container for the specific
# microservice
FROM alpine
RUN apk add --no-cache \
      curl \
      tzdata

# The default database directory
Volume /database

# Install our built image
COPY --from=compiler /dest/ /usr/bin/

#ENTRYPOINT ["/docker-entrypoint"]
ENTRYPOINT ["@@entrypoint@@"]
CMD [ "-c", "/config.yaml"]
