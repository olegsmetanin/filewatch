#!/bin/sh

SVC=./text.txt

(echo "first change"; echo "first change" > $SVC; sleep 5; echo "second change"; echo "second change" > $SVC; sleep 5; echo "stop changing";) &

while true; do
  while [ ! -f $SVC ]; do sleep 1; done
  filewatch -t 5 -verbose -filenames $SVC
  echo "File change finished"
done
