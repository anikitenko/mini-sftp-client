package controllers

import (
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
)

func DownloadFile(localPath, fileName, fileNamePost string, session *ssh.Session) (string, error) {
	var errorMessage string
	stdoutNewPipe, err := session.StdoutPipe()
	if err != nil {
		errorMessage = "Cannot create a pipe to download file"
		return errorMessage, err
	}

	newLocalFile, err := os.OpenFile(localPath+string(filepath.Separator)+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errorMessage = "Cannot open file to write to"
		return errorMessage, err
	}

	if err := session.Start("cat '" + fileNamePost + "'"); err != nil {
		errorMessage = "Problem with running command via SSH"
		return errorMessage, err
	}
	stdoutNewPipe = &PassThru{Reader: stdoutNewPipe}
	numberTransferred, err := io.Copy(newLocalFile, stdoutNewPipe)
	if err != nil {
		errorMessage = "Problem with copying file via SSH"
		return errorMessage, err
	}

	logger.Infof("Transferred %s", FormatBytes(float64(numberTransferred)))

	session.Close()
	newLocalFile.Close()

	return "", nil
}
