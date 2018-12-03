#!/bin/bash
#
# Script to generate a multi-architecture docker image from the
# individual images
#
# SYNTAX
#
# multiarch.sh imagename version architectures...
#
# Where arch is one of the following: amd64 arm32v6 arm32v7 arm64v8
#
# image should be the full name, e.g. area51/nre-feeds:latest or area51/nre-feeds:0.2
# This script will append -{microservice}-{arch} to that name
#

IMAGE=$1
shift

VERSION=$1
shift

MODULE=$1
shift

# The final multiarch image
if [ "$MODULE" = "Build" ]
then
  MULTIIMAGE=${IMAGE}:${VERSION}
else
  MULTIIMAGE=${IMAGE}:${MODULE}-${VERSION}
fi

. functions.sh

CMD="docker manifest create -a ${MULTIIMAGE}"
for ARCH in $@
do
  if [ "$MODULE" = "Build" ]
  then
    TAG="${IMAGE}:${ARCH}-${VERSION}"
  else
    TAG="${IMAGE}:${MODULE}-${ARCH}-${VERSION}"
  fi

  CMD="$CMD $TAG"
done
execute $CMD

for ARCH in $@
do
  if [ "$MODULE" = "Build" ]
  then
    TAG="${IMAGE}:${ARCH}-${VERSION}"
  else
    TAG="${IMAGE}:${MODULE}-${ARCH}-${VERSION}"
  fi

  # ensure this node has the latest image for this architecture
  execute "docker pull ${TAG}"

  CMD="docker manifest annotate"
  CMD="$CMD --os linux"
  CMD="$CMD --arch $(goarch $ARCH)"

  if [ "$(goarch $ARCH)" = "arm" ]
  then
    CMD="$CMD --variant v$(goarm $ARCH)"
  fi

  CMD="$CMD $MULTIIMAGE $TAG"

  execute $CMD
done

execute docker manifest push -p $MULTIIMAGE
