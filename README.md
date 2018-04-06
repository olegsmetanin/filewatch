# Filewatch utility in docker

A small utility (556Kb) for watching file changes. Exits with a timeout after the last file modification.

## Using

```
filewatch -t 5 -verbose -filenames ./test/text.txt

Options:
  -filenames string
    	files to watch separated by commas
  -t int
    	debounce interval
  -verbose
    	verbose mode
```

With docker
```
docker pull olegsmetanin/filewatch:latest-alpine3.7
```

Look at test.sh for example
```
#!/bin/sh

SVC=./test/text.txt

rm $SVC

(echo "$(date +"%T"): create file"; echo "create file" > $SVC; sleep 4; echo "$(date +"%T"): some change"; echo "some change" > $SVC; sleep 4; echo "$(date +"%T"): delete file"; rm $SVC; sleep 4; echo "$(date +"%T"): last change"; echo "last change" > $SVC) &

while true; do
  while [ ! -f $SVC ]; do sleep 1; done
  filewatch -t 5 -verbose -filenames $SVC
  echo "$(date +"%T"): file change finished"
done
```

## Test

```
docker-compose up
```

## Test builded image

```
docker-compose -f docker-compose.test.yaml up
```
