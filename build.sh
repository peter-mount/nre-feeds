#!/bin/sh

clear
docker build -t test . &&\
docker run -it --rm \
  --name test \
  -v /home/peter/Downloads:/data:ro \
  test \
  loaddarwinref -f /data/20180103020732_ref_v3.xml &&\
docker run -it --rm \
  --name test \
  -v /home/peter/Downloads:/data:ro \
  test \
  loaddarwintimetable -f /data/20180103020732_v8.xml
