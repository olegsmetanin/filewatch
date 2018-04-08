#!/bin/sh

echo "$(date +"%T"): Test begin"

LOCK=./lock.txt
echo "lock" > $LOCK

SVC=./test/test/text.txt

mkdir -p ./test/test
echo "initial" > $SVC

run()
{
sleep 1
echo "$(date +"%T"): change file"
echo "change file" > $SVC
sleep 4
echo "$(date +"%T"): change file"
echo "change file" > $SVC
sleep 6
echo "$(date +"%T"): delete folder, create and write file"
rm -rf ./test/test/*
echo "delete, create and write file" > $SVC
sleep 6
echo "$(date +"%T"): delete folder, create and write file"
rm -rf ./test/test/*
echo "delete, create and write file" > $SVC
sleep 6
echo "$(date +"%T"): delete folder, create and write file"
rm -rf ./test/*
sleep 1
mkdir -p ./test/test
echo "delete, create and write file" > $SVC

rm $LOCK
}

run &

while [ -f $LOCK ]; do
  filewatch -t 5 -verbose -filenames $SVC
  echo "$(date +"%T"): filewatch finished"
done

echo "// first comment" > $SVC

filewatch -t 1 --verbose -filenames='test/**/*.txt' --initial --command='cat ./test/test/text.txt' &

PID=$!

sleep 5
echo "// next comment" >> $SVC
sleep 5

kill $PID

echo "$(date +"%T"): Test done"
