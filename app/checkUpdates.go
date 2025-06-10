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

	//TODO: check downloaded app:
	// - compare md5 checksum
	// - size checks (not tiny)
	// - run with "--test" flag, then inspect output and exitcode

	err = copyFile(newFilePath, selfPath)
	if err != nil {
		log.Println("Error copying downloaded version over ourselves!", err)
		return
	}

	log.Println("Restarting app...")
	// Filler, to increase odds of log messages getting flushed to stdout before we go away.
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

func copyFile(sourceFilePath, destinationFilePath string) error {
	log.Printf("Copying %s to %s.\n", sourceFilePath, destinationFilePath)

	// Windows fails with "Access Denied" if we overwrite our the currently running exe.
	// But it lets us do the replacement in two steps:
	//   1. Rename the in-use file to ".bak"
	//   2. Copy the new file to the current exe.
	// We can do the same on linux & mac for consistency. (Need to update permissions though.)
	// It's nice to have the ".bak" file too, for manual troubleshooting and rollback.

	// Detect old file's current permissions before we touch it.
	var perms os.FileMode
	perms = 0644
	if runtime.GOOS != "windows" {
		if fileInfo, err := os.Lstat(destinationFilePath); err != nil {
			log.Printf("Can't get current permissions for dest file %s: %v\n", destinationFilePath, err)
		} else {
			perms = fileInfo.Mode().Perm()
			log.Printf("Current permissions of dest file %s: %#o\n", destinationFilePath, perms)
		}
	}

	backupPath := destinationFilePath + ".bak"
	log.Printf("Backing up file %s to %s!\n", destinationFilePath, backupPath)
	if err := os.Rename(destinationFilePath, backupPath); err != nil {
		log.Printf("Can't back up file %s to %s!\n", destinationFilePath, backupPath)
		return err
	}

	inputFile, err := os.Open(sourceFilePath)
	if err != nil {
		log.Printf("Can't open source file %s: %e\n", sourceFilePath, err)
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destinationFilePath)
	if err != nil {
		log.Printf("Can't open destination file %s: %e\n", destinationFilePath, err)
		return err
	}
	defer outputFile.Close()

	n, err := io.Copy(outputFile, inputFile)
	if err != nil {
		log.Printf("Can't write destination file %s: %e\n", destinationFilePath, err)
		return err
	}
	log.Printf("Copied %s to %s. wrote %d bytes.\n", sourceFilePath, destinationFilePath, n)

	if runtime.GOOS != "windows" {
		log.Printf("Restoring original permissions (%v) to file %s.\n", perms, destinationFilePath)
		if err = os.Chmod(destinationFilePath, perms); err != nil {
			log.Printf("Can't set permissions to %v on %s: %s\n", perms, destinationFilePath, err)
		}
	}

	return nil
}
