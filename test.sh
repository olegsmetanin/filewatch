#!/bin/sh

SVC=./test/test/text.txt

echo "$(date +"%T"): change file"
echo "change file" > $SVC
sleep 4
echo "$(date +"%T"): delete folder, create and write file"
rm -rf ./test
mkdir -p ./test/test
echo "delete, create and write file" > $SVC
