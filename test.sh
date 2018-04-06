#!/bin/sh

SVC=./text.txt

rm $SVC

(echo "$(date +"%T"): create file"; echo "create file" > $SVC; sleep 4; echo "$(date +"%T"): some change"; echo "some change" > $SVC; sleep 4; echo "$(date +"%T"): delete file"; rm $SVC; sleep 4; echo "$(date +"%T"): last change"; echo "last change" > $SVC) &

while true; do
  while [ ! -f $SVC ]; do sleep 1; done
  filewatch -t 5 -verbose -filenames $SVC
  echo "$(date +"%T"): file change finished"
done
