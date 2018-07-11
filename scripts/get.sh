#!/bin/sh

# Get all external libraries we need
go get -v \
      github.com/coreos/bbolt/... \
      github.com/jlaffaye/ftp \
      github.com/muesli/cache2go \
      github.com/peter-mount/golib/codec \
      github.com/peter-mount/golib/rabbitmq \
      github.com/peter-mount/golib/kernel \
      github.com/peter-mount/golib/rest \
      github.com/peter-mount/golib/statistics \
      github.com/peter-mount/golib/util \
      github.com/peter-mount/goxml2json \
      github.com/peter-mount/sortfold \
      gopkg.in/yaml.v2

exit 0

github.com/gorilla/mux \
github.com/peter-mount/golib/codec \
github.com/peter-mount/golib/rabbitmq \
github.com/peter-mount/golib/rest \
github.com/peter-mount/golib/statistics \
github.com/peter-mount/golib/util \
gopkg.in/robfig/cron.v2 \
io/ioutil \
log \
net/http \
path/filepath \
time
