#!/bin/sh
#
# Builds all binaries
#
DEST=$1
BIN=$2

echo "Building ${BIN}"
go build -o ${DEST} github.com/peter-mount/nre-feeds/bin/${BIN}
