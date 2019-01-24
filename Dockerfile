ARG arch=amd64
ARG goos=linux

# ============================================================
# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
FROM golang as build

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      git \
      tzdata \
      zip

WORKDIR /work
COPY go.mod .
RUN go mod download

#RUN go get -v \
#      github.com/etcd-io/bbolt \
#      github.com/jlaffaye/ftp \
#      github.com/muesli/cache2go \
#      github.com/peter-mount/golib/codec \
#      github.com/peter-mount/golib/rabbitmq \
#      github.com/peter-mount/golib/kernel \
#      github.com/peter-mount/golib/rest \
#      github.com/peter-mount/golib/statistics \
#      github.com/peter-mount/golib/util \
#      github.com/peter-mount/goxml2json \
#      github.com/peter-mount/sortfold \
#      gopkg.in/yaml.v2

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /work
ADD . .

# ============================================================
# Run all tests in a new container so any output won't affect
# the final build.
FROM source as test
ARG skipTest
RUN if [ -z "$skipTest" ] ;then \
      for bin in \
        util \
        darwinref \
        darwind3 \
        ldb \
        issues ;\
      do\
        echo "Testing ${bin}" ;\
        CGO_ENABLED=0 go test -v ./${bin} ;\
      done ;\
    fi

# ============================================================
# Compile the source.
FROM source as compiler
ARG module=
ARG arch
ARG goos
ARG goarch
ARG goarm
WORKDIR /work

# NB: CGO_ENABLED=0 forces a static build
RUN PACKAGE=${module};\
    if [ "$PACKAGE" = "darwintt" ];\
    then\
      PACKAGE="darwintimetable";\
    fi;\
    echo "Building ${module} as ${PACKAGE}";\
    CGO_ENABLED=0 \
    GOOS=${goos} \
    GOARCH=${goarch} \
    GOARM=${goarm} \
    go build \
      -o /dest/${module} \
      ./${PACKAGE}/bin

# ============================================================
# Optional stage, upload the binaries as a tar file
FROM compiler AS upload
ARG uploadPath=
ARG uploadCred=
ARG uploadName=
RUN if [ -n "${uploadCred}" -a -n "${uploadPath}" -a -n "${uploadName}" ] ;\
    then \
      cd /dest; \
      tar cvzpf /tmp/${uploadName}.tgz * && \
      zip /tmp/${uploadName}.zip * && \
      curl -u ${uploadCred} --upload-file /tmp/${uploadName}.tgz ${uploadPath}/ && \
      curl -u ${uploadCred} --upload-file /tmp/${uploadName}.zip ${uploadPath}/; \
    fi

# ============================================================
# Finally build the final runtime container for the specific
# microservice
FROM alpine
RUN apk add --no-cache \
      curl \
      tzdata

COPY --from=compiler /dest/ /usr/bin/

ENTRYPOINT ["@@entrypoint@@"]
CMD [ "-c", "/config.yaml"]
