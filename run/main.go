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
	"strings"
	"path/filepath"
	"time"
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

type staticReleaseInfo struct {
	Size int
	URL  string
}

func main() {

	flag.Parse()

	if *NoVersionCheck {
		StartClient()
		return
	}

	fmt.Println("Checking for updates...")

	runDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal(err)
	}

	appVersion, err := ioutil.ReadFile(runDir + string(filepath.Separator) + ".version")
	if err != nil {
		logger.Warnf("Unable to get current version: %v", err)
		StartClient()
		return
	}

	releaseFound, staticReleaseFound, err := getReleaseInfo()
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

	cleanAppVersion := strings.TrimSuffix(string(appVersion), "\n")
	cv, err := semver.Make(cleanAppVersion)
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
	fmt.Println("Downloading binary...")
	bar := pb.New(releaseFound.Size).SetUnits(pb.U_BYTES)
	bar.Start()

	resp, err := http.Get(releaseFound.URL)
	if err != nil {
		logger.Warnf("Problem with downloading new version: %v", err)
		StartClient()
		return
	}

	newFileName := fmt.Sprintf("%s_%v", releaseFound.Name, time.Now().UnixNano())

	writer, err := os.Create(runDir + string(filepath.Separator) + newFileName)
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

	if err = os.Remove(runDir + string(filepath.Separator) + releaseFound.Name); err != nil {
		logger.Warnf("Problem with removing old file version: %v", err)
		StartClient()
		return
	}

	if err = os.Rename(runDir+string(filepath.Separator)+newFileName, runDir+string(filepath.Separator)+releaseFound.Name); err != nil {
		logger.Fatalf("Problem with renaming new file: %v", err)
	}

	if err = os.Chmod(runDir+string(filepath.Separator)+releaseFound.Name, 0755); err != nil {
		logger.Fatalf("Problem with setting up execute permissions permissions for new file version: %v", err)
	}

	fmt.Println("Downloading static files...")
	bar = pb.New(staticReleaseFound.Size).SetUnits(pb.U_BYTES)
	bar.Start()

	resp, err = http.Get(staticReleaseFound.URL)
	if err != nil {
		logger.Warnf("Problem with downloading static files: %v", err)
		StartClient()
		return
	}

	newFileName = fmt.Sprintf("%s_%v.zip", "update", time.Now().UnixNano())

	writer, err = os.Create(runDir + string(filepath.Separator) + newFileName)
	if err != nil {
		logger.Warnf("Problem with creating new file: %v", err)
		StartClient()
		return
	}

	multiWriter = io.MultiWriter(writer, bar)

	bytesWritten, err = io.Copy(multiWriter, resp.Body)
	if err != nil {
		logger.Warnf("Problem with downloading static files: %v", err)
		StartClient()
		return
	}

	if bytesWritten != int64(staticReleaseFound.Size) {
		logger.Warnf("Problem with saving new file (incorrect bytes count): %v", err)
		StartClient()
		return
	}

	bar.Finish()
	resp.Body.Close()
	writer.Close()

	publicFolderLocation := "src" + string(filepath.Separator) + "github.com" + string(filepath.Separator) +
		"anikitenko" + string(filepath.Separator) + "mini-sftp-client"
	if err = os.RemoveAll(runDir + string(filepath.Separator) + publicFolderLocation + string(filepath.Separator) + "app"); err != nil {
		logger.Warnf("Problem with removing old app folder: %v", err)
		StartClient()
	}

	if err = os.RemoveAll(runDir + string(filepath.Separator) + publicFolderLocation + string(filepath.Separator) + "conf"); err != nil {
		logger.Warnf("Problem with removing old conf folder: %v", err)
		StartClient()
	}

	if err = os.RemoveAll(runDir + string(filepath.Separator) + publicFolderLocation + string(filepath.Separator) + "public"); err != nil {
		logger.Warnf("Problem with removing old public folder: %v", err)
		StartClient()
	}

	if _, err = Unzip(runDir + string(filepath.Separator) + newFileName, runDir + string(filepath.Separator) + publicFolderLocation); err != nil {
		logger.Fatalf("Problem with extracting archive with static files: %v", err)
	}

	if err = os.Remove(runDir + string(filepath.Separator) + newFileName); err != nil {
		logger.Warnf("Problem with removing archive with static files: %v", err)
	}

	if err = ioutil.WriteFile(runDir+string(filepath.Separator)+".version", []byte(latestVersion), 0644); err != nil {
		logger.Warnf("Problem with updating new version")
	}

	StartClient()
}
