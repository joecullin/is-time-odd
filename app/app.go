package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

var serverUrl string      // base url for the updates api
var selfPath string       // the currently running executable's path
var currentVersion string // version id of current executable
var selfTestMode bool     // print version and a single time result

func main() {
	if err := startupInit(); err != nil {
		log.Println("Error!", err)
		os.Exit(1)
	}

	if selfTestMode {
		showOddTime()
		os.Exit(0)
	}

	// Print result and check for updates every second.
	ticker := time.Tick(1 * time.Second)
	counter := 0
	for range ticker {
		// We don't do anything with counter, other than log it to help show context.
		counter++
		log.Println(counter)
		showOddTime()
		checkForUpdates()
	}
}

// The supposed "core function" of the app, other than the self-updating part.
// Display a line saying whether the current minute is odd or not.
func showOddTime() {
	// A value like "odd" or "not odd" gets compiled into different versions of the app.
	// (See Makefile.)
	log.Println("Time is PLACEHOLDER_TIME_VALUE!")
}

func startupInit() error {
	// "make build" overwrites the below placeholder with a version like 1.5 or 0.9.
	currentVersion = "PLACEHOLDER_VERSION_STRING"
	log.Printf("Starting version %s.", currentVersion)

	// Process command-line params
	flag.BoolVar(&selfTestMode, "test", false, "run self-test: print short output then exit")
	flag.StringVar(&serverUrl, "server", "http://localhost:3000", "base url for app updates API")
	flag.Parse()

	// Figure out our own path.
	// TODO: revisit this to handle symlinks
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	selfPath, err = filepath.Abs(ex)
	if err != nil {
		return err
	}
	log.Println("App running from:", selfPath)
	return nil
}
