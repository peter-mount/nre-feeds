#!/bin/sh

# Microservice to build & run
SERVICE=darwinref

# docker image:tag to build
IMAGE=test:${SERVICE}

# Port to run against
PORT=8081

# Local paths to where to store the DB's and local config.yaml
DBPATH=/home/peter/tmp/
CONFIG=$(pwd)/config.yaml

# End of customisations

clear

docker build -t ${IMAGE} --build-arg service=${SERVICE} . || exit 1

docker run -it --rm \
  --name test \
  -v ${DBPATH}:/database \
  -v ${CONFIG}:/config.yaml:ro \
  -p ${PORT}:${PORT} \
  ${IMAGE}
