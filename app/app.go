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
	} else {
		go keepShowingOddTime()
	}

	// Also check for updates every second.
	ticker := time.Tick(1 * time.Second)
	for range ticker {
		log.Println("") // blank line space to separate chunks log lines
		checkForUpdates()
	}
}

// The "core function" of the app, aside from its self-update ability.
// Display a message saying whether the current minute is odd or not.
func showOddTime() {
	// ANSI colors, to make this message stand out from all the update-checker log messages.
	var Highlight = "\033[30;46m" // white letters on cyan background
	var Reset = "\033[0m"

	// A value like "odd" or "not odd" gets compiled into each version of the app. (see Makefile)
	log.Println(Highlight + "************ Time is PLACEHOLDER_TIME_VALUE! ************" + Reset)
}

// Keep showing odd time forever, once a second
func keepShowingOddTime() {
	ticker := time.Tick(1 * time.Second)
	for range ticker {
		time.Sleep(1 * time.Second)
		showOddTime()
	}
}

func startupInit() error {
	// "make build" overwrites the below placeholder with a version string, like 1.5 or 0.9.
	currentVersion = "PLACEHOLDER_VERSION_STRING"
	log.Printf("Starting version %s.", currentVersion)

	// Process command-line params
	flag.BoolVar(&selfTestMode, "test", false, "run self-test: print short output then exit")
	flag.StringVar(&serverUrl, "server", "http://localhost:3008", "base url for app updates API")
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
