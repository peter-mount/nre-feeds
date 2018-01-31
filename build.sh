#!/bin/sh

clear
docker build -t test . || exit 1

#exit
#rm -rf /home/peter/tmp/dwlive.db
#rm -f /home/peter/tmp/dw*.db

# If first parameter is present then it's the National Rail ftp password
if [ -n "$1" ]
then
  FTP="-ftp $1"
fi

docker run -it --rm \
  --name test \
  -v /home/peter/tmp/:/database \
  -v $(pwd)/config.yaml:/config.yaml:ro \
  -p 8081:8081 \
  test

exit
docker run -it --rm \
  --name test \
  -v /home/peter/tmp/:/database \
  -v /home/peter/Downloads:/data:ro \
  -p 8081:8081 \
  test \
  darwin \
    -p 8081 $FTP \
    -ref /database/dwref.db \
    -timetable /database/dwtt.db

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
