package main

import (
	"flag"
	"log"
	"time"
)

var serverPort string
var appDataPath string

func main() {
	startupInit()
	loadAppData()

	// serve http api:
	go router()

	// Each minute, swap the "latest" tags in app data.
	ticker := time.Tick(1 * time.Second)
	minute := time.Now().Minute()
	for range ticker {
		log.Println(".")
		if time.Now().Minute() != minute {
			minute = time.Now().Minute()
			swapLatestTagOddEven()
		}
	}
}

func startupInit() error {
	log.Printf("Starting server")

	// Process command-line params
	flag.StringVar(&serverPort, "port", "3008", "port to listen on")
	flag.StringVar(&appDataPath, "app-data", "./server_data", "data directory")
	flag.Parse()

	//TODO - stricter validation, like:
	// - make sure port is a valid TCP port (uint16)
	// - make sure appDataPath exists and contains expected files
	return nil
}
