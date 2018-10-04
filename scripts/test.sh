#!/bin/sh
#
# Runs all tests
#
for bin in \
  util \
  darwinref \
  darwind3 \
  ldb \
  issues
do
  echo "Testing ${bin}"
  CGO_ENABLED=0 go test -v github.com/peter-mount/nre-feeds/${bin}
done
