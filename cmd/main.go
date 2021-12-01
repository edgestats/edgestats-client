package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/edgestats/edgestats-client/handlers"
	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Initializing EdgeStats...")
	fmt.Printf("System version: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Golang version: %s\n", runtime.Version())

	fp, err := handlers.GetFilePath(runtime.GOOS)
	if err != nil {
		fmt.Println("Error initializing:", err)
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error initializing:", err)
		os.Exit(1)
	}
	defer watcher.Close()

	if err := watcher.Add(fp); err != nil {
		fmt.Println("Error initializing: no file", fp)
		os.Exit(1)
	}

	fmt.Println("EdgeStats watching file", fp)
	fmt.Println("EdgeStats is ready...")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	var offset int64

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// process event
				offset, err = handlers.ProcessEvent(watcher, event, fp, offset)
				if err != nil {
					continue // perhaps log to log file
				}
			case error, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error: ", error) // perhaps log to log file
			}
		}
	}()

	// goroutine to poke log file at interval
	// needed for windows & perhaps other OSs
	go func() {
		ticker := time.NewTicker(6000 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := handlers.PokeFilePath(fp); err != nil {
					continue // perhaps log to log file
				}
			}
		}
	}()

	sig := <-ch
	fmt.Printf("\nRecieved %s signal, shutting down...\n", sig)

	_, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
}
