#!/bin/sh

clear
docker build -t test . || exit 1

#rm -fv /home/peter/tmp/darwin.db

docker run -it --rm \
  --name test \
  -v /home/peter/tmp/:/database \
  -v /home/peter/Downloads:/data:ro \
  -p 8081:8081 \
  test \
  darwin \
    -p 8081 \
    -ref /database/dwref.db

exit

docker run -it --rm \
  --name test \
  -v /home/peter/tmp/:/database \
  -v /home/peter/Downloads:/data:ro \
  test \
  loaddarwinref \
    -d /database/darwin.db \
    -f /data/20180103020732_ref_v3.xml

docker run -it --rm \
  --name test \
  -v /home/peter/tmp/:/database \
  -v /home/peter/Downloads:/data:ro \
  test \
  loaddarwintimetable \
    -d /database/darwin.db \
    -f /data/20180103020732_v8.xml
