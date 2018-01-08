# Dockerfile used to build the application

# Build container containing our pre-pulled libraries
FROM golang:latest as build

# Static compile
ENV CGO_ENABLED=0
ENV GOOS=linux

# We want to build our final image under /dest
# A copy of /etc/ssl is required if we want to use https datasources
RUN mkdir -p /dest/etc &&\
    cp -rp /etc/ssl /dest/etc/

# Ensure we have the libraries - docker will cache these between builds
RUN go get -v \
      flag \
      github.com/gorilla/mux \
      github.com/peter-mount/golib/rabbitmq \
      github.com/peter-mount/golib/statistics \
      gopkg.in/robfig/cron.v2 \
      gopkg.in/yaml.v2 \
      io/ioutil \
      log \
      net/http \
      path/filepath \
      time

# Import the source and compile
WORKDIR /usr/local/go/src
ADD . .

# Now each binary
RUN go build -v -x \
      -o /dest/bin/loaddarwinref bin/loaddarwinref

# Finally build the final runtime container will all required files
FROM scratch
COPY --from=build /dest/ /
#CMD ["td"]
