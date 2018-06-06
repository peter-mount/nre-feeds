#!/bin/sh
#
# Builds all binaries
#
DEST=$1
BIN=$2

PACKAGE=$BIN
if [ "$PACKAGE" = "darwintt" ]
then
  PACKAGE="darwintimetable"
fi

echo "Building ${BIN}"
go build -o ${DEST} github.com/peter-mount/nre-feeds/${PACKAGE}/bin
#go build -o ${DEST} github.com/peter-mount/nre-feeds/bin/${BIN}
