package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	zglob "github.com/mattn/go-zglob"
)

var fileNames = flag.String("filenames", "", "files to watch separated by commas")
var debounceInterval = flag.Int("t", 0, "debounce interval")
var verbose = flag.Bool("verbose", false, "verbose mode")
var command = flag.String("command", "", "command to execute")
var initial = flag.Bool("initial", false, "run command before any change happens")

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

func debounceThen(events <-chan fsnotify.Event, cb func()) {
	event := <-events
	if *verbose {
		log.Printf("event: %s, wait for next\n", event)
	}

LOOP:
	for {
		select {
		case event := <-events:
			if *verbose {
				log.Printf("event: %s, wait for next\n", event)
			}
		case <-time.After(time.Duration(*debounceInterval) * time.Second):
			break LOOP
		}
	}
	cb()
}

func watchForChanges(files []string) chan fsnotify.Event {
	fileMap := make(map[string]struct{})
	for _, f := range files {
		fileMap[f] = struct{}{}
	}

	events := make(chan fsnotify.Event)

	go func() {
		for {
			select {
			case event := <-watch.Events:
				if _, ok := fileMap[event.Name]; ok {
					if event.Op == fsnotify.Chmod {
						continue
					}
					events <- event
				}
			case err := <-watch.Errors:
				if err != nil {
					log.Fatalf("watch error: %s", err)
				} else {
					log.Fatalf("unexpected watch error")
				}
			}
		}
	}()

	return events
}

func runCommand(ctx context.Context, command string) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("can't get stdout for command: %s %s", command, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("can't get stderr for command: %s %s", command, err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("can't start command: %s %s", command, err)
	}

	go func() {
		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			log.Printf("[STDERR] %s", errScanner.Text())
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		log.Printf("%s", scanner.Text())
	}

	if err = cmd.Wait(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			log.Printf("%s", e.ProcessState)
		} else {
			log.Printf("can't wait for process: %s %s", command, err)
		}

	}
}

func main() {
	flag.Parse()

	if *verbose {
		log.Printf("filewatch version 0.0.4\n")
	}

	var err error
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()

	files := make([]string, 0)

	parts := strings.Split(*fileNames, ",")
	for _, part := range parts {
		matches, err := zglob.Glob(part)
		if err != nil {
			log.Fatalf("invalid pattern: %s %s", part, err)
		}
		for _, f := range matches {
			if *verbose {
				log.Printf("watching for: %s", f)
			}
			f, err := filepath.Abs(f)
			if err != nil {
				log.Fatalf("can't get absolute path for file: %s", f)
			}
			files = append(files, f)
		}
	}

	if len(files) == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	if err := addFilesToWatch(files); err != nil {
		log.Fatal(err)
	}

	events := watchForChanges(files)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if *initial {
		go runCommand(ctx, *command)
	}

	for {
		debounceThen(events, func() {
			if *command == "" {
				os.Exit(0)
				return
			}

			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			go runCommand(ctx, *command)
		})
	}

}
