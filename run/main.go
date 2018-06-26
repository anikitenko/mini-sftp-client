package main

import (
	"flag"
	"fmt"
	"github.com/blang/semver"
	"github.com/cheggaaa/pb"
	logger "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	releaseUrl     = "https://api.github.com/repos/anikitenko/mini-sftp-client/releases/latest"
	PortToListen   = flag.String("p", "", "Port to listen on")
	RunMode        = flag.String("m", "prod", "Run mode: dev OR prod")
	NoVersionCheck = flag.Bool("no-ver-check", false, "Skip version check?")
)

type releaseInfo struct {
	Name    string
	Size    int
	URL     string
	TagName string
}

func main() {

	flag.Parse()

	if *NoVersionCheck {
		StartClient()
		return
	}

	fmt.Println("Checking for updates...")

	appVersion, err := ioutil.ReadFile(".version")
	if err != nil {
		logger.Warnf("Unable to get current version: %v", err)
		StartClient()
		return
	}

	releaseFound, err := getReleaseInfo()
	if err != nil {
		switch err.Error() {
		case "url not found":
		case "invalid data":
		default:
			logger.Warnf("Problem with checking for updates: %v", err)
		}
		StartClient()
		return
	}

	cv, err := semver.Make(string(appVersion))
	if err != nil {
		logger.Warnf("Problem with parsing current version: %v", err)
		StartClient()
		return
	}

	latestVersion := releaseFound.TagName[1:]

	av, err := semver.Make(latestVersion)
	if err != nil {
		logger.Warnf("Problem with parsing new version: %v", err)
		StartClient()
		return
	}

	if av.LTE(cv) {
		StartClient()
		return
	}

	fmt.Println("New version is available: " + releaseFound.TagName)
	fmt.Println("Upgrading...")
	bar := pb.New(releaseFound.Size).SetUnits(pb.U_BYTES)
	bar.Start()

	resp, err := http.Get(releaseFound.URL)
	if err != nil {
		logger.Warnf("Problem with downloading new version: %v", err)
		StartClient()
		return
	}

	newFileName := releaseFound.Name + "-new"

	writer, err := os.Create(newFileName)
	if err != nil {
		logger.Warnf("Problem with creating new file: %v", err)
		StartClient()
		return
	}

	multiWriter := io.MultiWriter(writer, bar)

	bytesWritten, err := io.Copy(multiWriter, resp.Body)
	if err != nil {
		logger.Warnf("Problem with saving new file: %v", err)
		StartClient()
		return
	}

	if bytesWritten != int64(releaseFound.Size) {
		logger.Warnf("Problem with saving new file (incorrect bytes count): %v", err)
		StartClient()
		return
	}

	bar.Finish()
	resp.Body.Close()
	writer.Close()

	err = os.Remove(releaseFound.Name)
	if err != nil {
		logger.Warnf("Problem with removing old file version: %v", err)
		StartClient()
		return
	}

	err = os.Rename(newFileName, releaseFound.Name)
	if err != nil {
		logger.Fatalf("Problem with renaming new file: %v", err)
	}

	err = os.Chmod(releaseFound.Name, 0755)
	if err != nil {
		logger.Fatalf("Problem with setting up execute permissions permissions for new file version: %v", err)
	}
}
