#!/bin/sh

# Module to build & run
#MODULE=darwinref
#MODULE=darwintimetable
#MODULE=darwind3
MODULE=ldb
#MODULE=darwinkb
#MODULE=darwindb
#MODULE=darwinkb

# db directory
DB=/home/peter/tmp/nre

ARGS="$ARGS -it --rm --name ${MODULE} --hostname ${MODULE}"
ARGS="$ARGS -v ${DB}:/database"
ARGS="$ARGS -v $(pwd)/config.yaml:/config.yaml:ro"

if [ "$MODULE" = "darwindb" ]
then
  ARGS="$ARGS --link postgres"
  ARGS="$ARGS -e POSTGRESDB='postgres://postgres:temppass@postgres/postgres?sslmode=disable&connect_timeout=3'"
fi

ARGS="$ARGS -e CACHEDIR=/database/${MODULE}"

if [ ! -z "$IPP" ]
then
    ARGS="$ARGS -e IPP=$IPP"
fi

ARGS="$ARGS test:${MODULE}-amd64-latest"

echo $ARGS

./build.sh test amd64 latest ${MODULE}

exit 0

# &&\
exec docker run $ARGS

