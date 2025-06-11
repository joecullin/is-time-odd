package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func testDelay(seconds int) {
	log.Printf("sleeping %d seconds to test how client app handles delays", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	log.Println("done sleeping")
}

func router() {
	log.Println("port is", serverPort)

	router := http.NewServeMux()

	// get release details
	router.HandleFunc("GET /api/releases/{os}/{version}/info", func(w http.ResponseWriter, r *http.Request) {
		release, err := validateAndGetRelease(r.PathValue("os"), r.PathValue("version"))
		if err != nil {
			notFoundPage(w)
			return
		}
		// Un-comment to manually verify that client app still prints the time when update check is slow:
		// testDelay(15)
		handleGetReleaseInfo(w, release)
	})

	// download a release
	router.HandleFunc("GET /api/releases/{os}/{version}", func(w http.ResponseWriter, r *http.Request) {
		release, err := validateAndGetRelease(r.PathValue("os"), r.PathValue("version"))
		if err != nil {
			notFoundPage(w)
			return
		}
		handleDownloadRelease(w, release)
	})

	// homepage, plus catch-all 404 for any other method+path
	router.Handle("/", http.HandlerFunc(apiInfoPage))

	log.Println("listening on port", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatal("Couldn't start server! ", err)
	}
}

// shared by a couple routes with similar paths: validate path params, and get release.
func validateAndGetRelease(osParam, versionParam string) (ReleaseData, error) {
	var release ReleaseData
	if matched, err := regexp.MatchString("^linux|darwin|windows$", osParam); err != nil || !matched {
		return release, fmt.Errorf("Validation failed for os=%s. Must be a supported os.\n", osParam)
	}
	if matched, err := regexp.MatchString("^([0-9]+[.][0-9]+)|latest$", versionParam); err != nil || !matched {
		return release, fmt.Errorf("Validation failed for version=%s. Must be x.x or 'latest'.\n", versionParam)
	}
	release, err := getRelease(osParam, versionParam)
	if err != nil {
		log.Println("validateAndGetRelease: getRelease did not return a valid release.", err)
		return release, err
	}
	return release, nil
}

func handleGetReleaseInfo(w http.ResponseWriter, release ReleaseData) {
	log.Printf("sending release info: os=%s version=%s", release.Platform, release.Version)
	jsonString, err := json.Marshal(release)
	if err != nil {
		log.Println("Error creating json of release info:", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonString))
}

func handleDownloadRelease(w http.ResponseWriter, release ReleaseData) {
	log.Printf("api/testdownload\n")

	appFilePath := getPathToExecutable(release)
	appFileData, err := os.ReadFile(appFilePath)
	if err != nil {
		log.Printf("Error reading app file %s: %e", appFilePath, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Oh no! An error occurred on our end. Maybe try again later?")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(appFileData)))
	w.WriteHeader(http.StatusOK)
	w.Write(appFileData)
}

func apiInfoPage(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		notFoundPage(w)
		return
	}

	fmt.Fprint(w, `
        <html>
        <head>Releases Server</head>
        <body>
			TODO: add some small "docs" here.
			- maybe - link to github readme
			- maybe - incorporate some dynamic data from request (in a template)
        </body>
        </html>`)
}

func notFoundPage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, `
        <html>
        <head>Not found!</head>
        <body>
        Not found! See <a href="/">docs</a> for help.
        </body>
        </html>`)
}
