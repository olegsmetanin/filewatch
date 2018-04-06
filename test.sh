#!/bin/sh

SVC=./text.txt

(echo "$(date +"%T"): first change"; echo "first change" > $SVC; sleep 5; echo "$(date +"%T"): last change"; echo "last change" > $SVC) &

while true; do
  while [ ! -f $SVC ]; do sleep 1; done
  filewatch -t 5 -verbose -filenames $SVC
  echo "$(date +"%T"): file change finished"
done
