#!/bin/sh
#
# Script to run a build for a specific microservice and platform.
#
# SYNTAX
#
# docker.sh imagename microservice arch version
#
# Where arch is one of the following: amd64 arm32v6 arm32v7 arm64v8
#
# image should be the full name, e.g. area51/nre-feeds:latest or area51/nre-feeds:0.2
# This script will append -{microservice}-{arch} to that name
#

IMAGE=$1
SERVICE=$2
ARCH=$3
VERSION=$4

# Resolve the architecture
case $ARCH in
  amd64)
    GOARCH=amd64
    ;;
  arm32v6)
    GOARCH=arm
    GOARM=6
    ;;
  arm32v7)
    GOARCH=arm
    GOARM=7
    ;;
  arm64v8)
    GOARCH=arm64
    ;;
  *)
    echo "Unsupported architecture $ARCH"
    exit 1
    ;;
esac

# For now just support Linux
GOOS=linux

# The actual image being built
TAG=${IMAGE}:${SERVICE}-${ARCH}-${VERSION}

# Now customise Dockerfile
DOCKERFILE=Dockerfile.${SERVICE}.${GOOS}.${ARCH}
(
  echo "# GENERATED $DOCKERFILE"
  sed "s/@@entrypoint@@/${SERVICE}/g" Dockerfile
) >$DOCKERFILE

echo "Building $SERVICE image $TAG on $ARCH"

docker build \
  --force-rm=true \
  -t ${TAG} \
  -f ${DOCKERFILE} \
  --build-arg service=${SERVICE} \
  --build-arg arch=${ARCH} \
  --build-arg goos=${GOOS} \
  --build-arg goarch=${GOARCH} \
  --build-arg goarm=${GOARM} \
  .
