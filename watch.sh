#!/bin/sh

SVC=./test/test/text.txt

while true; do
  while [ ! -f $SVC ]; do sleep 1; done
  filewatch -t 5 -verbose -filenames $SVC
  echo "$(date +"%T"): file change finished"
done
