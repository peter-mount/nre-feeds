#!/bin/sh

# Microservice to build & run
#SERVICE=darwinref
SERVICE=darwintt
#SERVICE=darwind3
#SERVICE=ldb

# docker image:tag to build
IMAGE=test:${SERVICE}

# Port to run against
PORT=8081

# Local paths to where to store the DB's and local config.yaml
DBPATH=/home/peter/tmp/
CONFIG=$(pwd)/config.yaml

# End of customisations

clear

for i in darwinref darwintt darwind3 ldb
do
  docker build -t test:$i --build-arg service=$i . || exit 1
done

exit

docker build -t ${IMAGE} --build-arg service=${SERVICE} . || exit 1

exit

docker run -it --rm \
  --name test \
  -v ${DBPATH}:/database \
  -v ${CONFIG}:/config.yaml:ro \
  -p ${PORT}:80 \
  ${IMAGE}
