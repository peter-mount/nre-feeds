#!/bin/sh

# Module to build & run
#MODULE=darwinref
MODULE=darwintimetable

# db directory
DB=/home/peter/tmp/nre

./scripts/build.sh test ${MODULE} amd64 latest &&\
docker run \
  -it \
  --rm \
  --name ${MODULE} \
  -v ${DB}:/database \
  -v $(pwd)/config.yaml:/config.yaml:ro \
  test:${MODULE}-amd64-latest
