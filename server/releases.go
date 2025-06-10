package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type ReleaseData struct {
	File     string
	Platform string
	Md5      string
	Tags     []string
	Version  string
}
type AppData struct {
	Releases []ReleaseData
}

var appData AppData

// Get a release that matches a given os + version id.
// You can supply "latest" as the version id.
// (If there are multiple matches, this only gets the first. Label releases with unique versions/tags to avoid that ambiguity.)
func getRelease(os, version string) (ReleaseData, error) {
	i := slices.IndexFunc(appData.Releases, func(r ReleaseData) bool {
		if r.Platform == os {
			if version == r.Version || (version == "latest" && slices.Contains(r.Tags, "latest")) {
				return true
			}
		}
		return false
	})
	if i > 0 {
		return appData.Releases[i], nil
	} else {
		return ReleaseData{}, errors.New("no matching releases found")
	}
}

// Read json file with details about all releases
func loadAppData() {
	dataFilePath := filepath.Join(appDataPath, "appData.json")
	dataFileReader, err := os.Open(dataFilePath)
	if err != nil {
		log.Printf("Error opening data file %s: %e", dataFilePath, err)
	}
	defer dataFileReader.Close()

	jsonData, err := io.ReadAll(dataFileReader)
	if err != nil {
		log.Printf("Error reading data file %s: %e", dataFilePath, err)
	}
	json.Unmarshal(jsonData, &appData)
	log.Printf("Loaded info for %d releases", len(appData.Releases))
	swapLatestTagOddEven()
}

// Dump current appData value to log
func logAppData() {
	jsonString, err := json.MarshalIndent(appData, "", "  ")
	if err != nil {
		log.Println("Error trying to log app data:", err)
	}
	log.Printf("%s\n", jsonString)
}

func swapLatestTagOddEven() {
	var oldPattern, newPattern string
	if time.Now().Minute()%2 != 0 {
		oldPattern, newPattern = "even", "odd"
	} else {
		oldPattern, newPattern = "odd", "even"
	}
	log.Printf("Adding latest tag to %s", newPattern)
	removeTag("latest", oldPattern)
	addTag("latest", newPattern)
	logAppData()
}

// Add newTag to all releases that have searchTag
// This only modifies our copy. It doesn't write the change to disk.
func addTag(newTag, searchTag string) {
	updateCount := 0
	for i := 0; i < len(appData.Releases); i++ {
		if slices.Contains(appData.Releases[i].Tags, searchTag) && !slices.Contains(appData.Releases[i].Tags, newTag) {
			log.Printf("addTag: adding '%s' to %s", newTag, appData.Releases[i].File)
			appData.Releases[i].Tags = append(appData.Releases[i].Tags, newTag)
			updateCount++
		}
	}
	log.Printf("addTag: updated %d releases", updateCount)
}

// Remove deleteTag from all releases that have searchTag
// This only modifies our copy. It doesn't write the change to disk.
func removeTag(deleteTag, searchTag string) {
	updateCount := 0
	for i := 0; i < len(appData.Releases); i++ {
		if slices.Contains(appData.Releases[i].Tags, searchTag) &&
			slices.Contains(appData.Releases[i].Tags, deleteTag) {
			log.Printf("removeTag: removing '%s' from %s", deleteTag, appData.Releases[i].File)
			appData.Releases[i].Tags = slices.DeleteFunc(appData.Releases[i].Tags, func(tag string) bool {
				return tag == deleteTag
			})
			updateCount++
		}
	}
	log.Printf("removeTag: updated %d releases", updateCount)
}
