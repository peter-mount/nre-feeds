#!/bin/sh

# Module to build & run
#MODULE=darwinref
#MODULE=darwintimetable
#MODULE=darwind3
#MODULE=ldb
#MODULE=darwinkb
MODULE=darwindb

# db directory
DB=/home/peter/tmp/nre

./build.sh test amd64 latest ${MODULE} &&\
docker run \
  -it \
  --rm \
  --name ${MODULE} \
  --hostname darwin-db \
  -v ${DB}:/database \
  -v $(pwd)/config.yaml:/config.yaml:ro \
  --link postgres \
  -e POSTGRESDB='postgres://postgres:temppass@postgres/postgres?sslmode=disable&connect_timeout=3' \
  test:${MODULE}-amd64-latest
