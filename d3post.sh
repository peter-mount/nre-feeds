#!/bin/bash
#
# Simple bash file to read a file line by line and post it to our test
# rest service
#
# Simply run with the filename of captured logs from the live d3 feed as an
# argument and it will submit each line to the service
#
while IFS= read -r var
do
  #echo $var
  curl -s \
    -X POST \
    -H "Content-Type: text/xml" \
    --data-binary "$var" \
    http://localhost:8081/live/test
  #echo
done <$1
