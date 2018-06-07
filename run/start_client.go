package main

import (
	"bufio"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"runtime"
)

func StartClient() {
	extension := ""
	goOS := runtime.GOOS
	if goOS == "windows" {
		extension = ".exe"
	}
	fmt.Println("Starting client...")
	client := exec.Command("mini-sftp-client-"+goOS+extension, "-importPath", "mini-sftp-client", "-runMode", "prod")

	stdout, err := client.StdoutPipe()
	if nil != err {
		logger.Fatalf("Error obtaining stdout: %v", err)
	}

	stderr, err := client.StderrPipe()
	if nil != err {
		logger.Fatalf("Error obtaining stderr: %v", err)
	}

	readerOut := bufio.NewReader(stdout)
	readerErr := bufio.NewReader(stderr)

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}(readerOut)
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}(readerErr)

	if err := client.Start(); err != nil {
		logger.Fatalf("Problem with starting client: %v", err)
	}

	if err := client.Wait(); err != nil {
		logger.Fatalf("Problem with starting client: %v", err)
	}
}
