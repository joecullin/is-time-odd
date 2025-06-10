package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

type Release struct {
	Checksum string
	Platform string
	Tags     []string
	Version  string
	Md5      string
}

func checkForUpdates() {
	log.Printf("This is version %s for %s.\n", currentVersion, runtime.GOOS)

	latest, err := getReleaseInfo("latest")
	if err != nil {
		log.Println("Error getting release info!", err)
		return
	}
	log.Printf("Latest version=%v, md5=%s\n", latest.Version, latest.Md5)
	if latest.Version == currentVersion {
		// We're up to date! Nothing more to do.
		return
	}

	newFilePath := selfPath + ".new_release"
	err = downloadRelease(newFilePath, latest)
	if err != nil {
		log.Println("Error downloading new release!", err)
		return
	}
	log.Println("Restarting app...")
	// increase chances of above log message getting flushed to stdout before we go away
	for range 5 {
		log.Println(".")
	}
	restartApp()
}

func getReleaseInfo(version string) (Release, error) {
	var releaseInfo Release
	infoUrl := fmt.Sprintf("%s/api/releases/%s/%s/info", serverUrl, runtime.GOOS, version)
	log.Println("Checking for updates at", infoUrl)

	res, err := http.Get(infoUrl)
	if err != nil {
		return releaseInfo, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return releaseInfo, fmt.Errorf("unexpected response status: %s", res.Status)
	}
	json.NewDecoder(res.Body).Decode(&releaseInfo)

	return releaseInfo, nil
}

func downloadRelease(filePath string, release Release) error {
	url := fmt.Sprintf("%s/api/releases/%s/%s", serverUrl, runtime.GOOS, release.Version)
	log.Printf("Downloading %s to %s\n", url, filePath)
	out, err := os.Create(filePath)
	if err != nil {
		log.Print("Can't create output file!", err)
		return err
	}
	defer out.Close()

	res, err := http.Get(url)
	if err != nil {
		log.Println("Error downloading app!", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status when downloading app: %s", res.Status)
		return err
	}

	n, err := io.Copy(out, res.Body)
	if err != nil {
		log.Println("Can't write downloaded file", err)
		return err
	}
	log.Printf("wrote %d bytes to %s\n", n, filePath)
	return nil
}

func restartApp() error {

	restartParams := []string{`--server="http://localhost:5005"`} //TODO - get these from argv

	if runtime.GOOS != "windows" {
		// Exec is nicer when it's available.
		// (Replaces current process, so things like job control in shell still work.)
		params := append([]string{"_"}, restartParams...)
		env := os.Environ()
		if err := syscall.Exec(selfPath, params, env); err != nil {
			log.Printf("Error re-starting %s with params %v: %e\n", selfPath, params, err)
			return err
		}
		log.Println("done re-starting!")
	} else {
		// Windows: no syscall.exec.
		// So we'll start a new process then exit this one.
		cmd := exec.Command(selfPath, strings.Join(restartParams, " "))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Println("Error starting new process!", err)
			return err
		} else {
			log.Println("Started new process! Exiting this one")
			os.Exit(0)
		}
	}
	return nil //TODO - it's an error if we reach this, right?
}
