#!/bin/sh

# Module to build & run
#MODULE=darwinref
#MODULE=darwintimetable
MODULE=darwind3
#MODULE=ldb
#MODULE=darwinkb

# db directory
DB=/home/peter/tmp/nre

./build.sh test amd64 latest ${MODULE} &&\
docker run \
  -it \
  --rm \
  --name ${MODULE} \
  -v ${DB}:/database \
  -v $(pwd)/config.yaml:/config.yaml:ro \
  test:${MODULE}-amd64-latest
