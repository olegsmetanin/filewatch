package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var fileNames = flag.String("filenames", "", "files to watch separated by commas")
var debounceInterval = flag.Int("t", 0, "debounce interval")
var verbose = flag.Bool("verbose", false, "verbose mode")

var watch *fsnotify.Watcher

func addFilesToWatch(files []string) error {
	for _, f := range files {
		stat, err := os.Stat(f)
		if err != nil {
			return fmt.Errorf("can't get stat for file: %s, %s", f, err)
		}

		if err := watch.Add(f); err != nil {
			return fmt.Errorf("can't add file to watch: %s, %s", f, err)
		}
		if !stat.IsDir() {
			if err := watch.Add(path.Dir(f)); err != nil {
				return fmt.Errorf("can't add file to watch: %s, %s", f, err)
			}
		}
	}
	return nil
}

func main() {
	flag.Parse()

	if *fileNames == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	var err error
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()

	files := strings.Split(*fileNames, ",")
	if err := addFilesToWatch(files); err != nil {
		log.Fatal(err)
	}

	fileMap := make(map[string]struct{})
	for _, f := range files {
		fileMap[f] = struct{}{}
	}

	events := make(chan fsnotify.Event)
	go func() {
		for {
			select {
			case event := <-watch.Events:
				if event.Op&fsnotify.Chmod != 0 {
					continue
				}

				if _, ok := fileMap[event.Name]; ok {
					events <- event
					addFilesToWatch(files)
				}
			case err := <-watch.Errors:
				log.Fatalf("watch error: %s", err)
			}
		}
	}()

	<-events
	for {
		select {
		case event := <-events:
			if *verbose {
				log.Printf("event: %s, wait for next\n", event)
			}
		case <-time.After(time.Duration(*debounceInterval) * time.Second):
			return
		}
	}

}
