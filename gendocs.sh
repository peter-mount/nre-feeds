#!/bin/sh
#
# Utility script to generate documentation.
#
# This will generate text & html docs under the docs directory and update the
# README.md file within the library directyories.
#
# To run simply be in the base directory of the repository and then run as:
#
# ./gendocks.sh
#
# It should then generate documentation under docs and any README.md files
#
if [ -z "$1" ]
then
  # This runs locally
  SRC="$(pwd)"
  PROJECT="$(basename $SRC)"
  IMPORT="github.com/peter-mount/$PROJECT"
  WORK="/go/src/$IMPORT"
  docker run \
    -it --rm \
    --name build-$PROJECT-$$ \
    -v $SRC:$WORK \
    -w $WORK \
    -e UID=$(id -u) \
    -e IMPORT=$IMPORT \
    golang:latest \
    ./gendocs.sh $$
else
  # This runs within the docker container

  # Needed for github README.md generation
  go get -v \
        github.com/robertkrimen/godocdown/godocdown

  echo "Generating documentation"

  # The docs directory
  mkdir -p docs
  chown -R $UID docs

  for LIB in bin darwind3 darwinkb darwinref darwinrest darwintimetable darwinupdate ldb util
  do
    PACKAGE=${IMPORT}/${LIB}

    echo "package ${PACKAGE}"

    godoc ${PACKAGE} >docs/${LIB}.txt
    godoc -html ${PACKAGE} >docs/${LIB}.html

    # cif markdown in the cif directory as that does get committed
    godocdown -output=docs/${LIB}.md ${PACKAGE}

    # Ensure permissions are correct
    chown $UID docs/${LIB}.*

    # Install the markdown docs
    cp -pv docs/${LIB}.md ${LIB}/README.md
  done
fi
