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
  go test -v github.com/peter-mount/nre-feeds/${bin}
done
