package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func router() {
	log.Println("port is", serverPort)

	router := http.NewServeMux()

	router.HandleFunc("GET /api/releases/{os}/{version}/info", func(w http.ResponseWriter, r *http.Request) {
		osParam := r.PathValue("os")
		versionParam := r.PathValue("version")
		if matched, err := regexp.MatchString("^linux|darwin|windows$", osParam); err != nil || !matched {
			fmt.Fprintf(w, "Validation failed for os=%s. Must be a supported os.\n", osParam)
			return
		}
		if matched, err := regexp.MatchString("^([0-9]+[.][0-9]+)|latest$", versionParam); err != nil || !matched {
			fmt.Fprintf(w, "Validation failed for version=%s. Must be x.x or 'latest'.\n", versionParam)
			return
		}
		log.Printf("api/releases/.../info: ! os=%s version=%s\n", osParam, versionParam)
		release, err := getRelease(osParam, versionParam)
		if err != nil {
			log.Println("getRelease did not return a valid release:", err)
			notFoundPage(w)
			return
		}
		log.Printf("found %s for %s\n", release.Version, release.Platform)
		jsonString, err := json.Marshal(release)
		if err != nil {
			log.Println("error:", err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonString))
	})

	router.HandleFunc("GET /api/releases/darwin/1.1", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("api/testdownload\n")

		data, err := os.ReadFile("./server")
		if err != nil {
			log.Printf("Error reading data file!")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Oh no! An error occurred on our end. Maybe try again later?")
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// homepage plus catch-all for all other paths and methods:
	router.Handle("/", http.HandlerFunc(infoPage))

	log.Println("listening on port", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatal("Couldn't start server! ", err)
	}
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

func infoPage(w http.ResponseWriter, req *http.Request) {
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
