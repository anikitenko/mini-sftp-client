package controllers

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"time"
)

func DownloadDirectory(sourcePath, localPath, fileName string, session *ssh.Session) (string, error) {
	var errorMessage string
	stdoutNewPipe, err := session.StdoutPipe()
	if err != nil {
		errorMessage = "Cannot create a pipe to download folder"
		return errorMessage, err
	}

	tempArchiveName := fmt.Sprintf("%sdownload_%v.tar.gz", localPath+string(filepath.Separator), time.Now().UnixNano())

	tempArchiveFile, err := os.OpenFile(tempArchiveName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errorMessage = "Cannot open temp archive to write to"
		return errorMessage, err
	}

	if err := session.Start("tar cz --directory='" + sourcePath + "' '" + fileName + "'"); err != nil {
		errorMessage = "Problem with running command via SSH"
		return errorMessage, err
	}

	stdoutNewPipe = &PassThru{Reader: stdoutNewPipe}

	numberTransferred, err := io.Copy(tempArchiveFile, stdoutNewPipe)
	if err != nil {
		errorMessage = "Problem with copying archive via SSH"
		return errorMessage, err
	}

	logger.Infof("Transferred %s", FormatBytes(float64(numberTransferred)))

	session.Close()
	tempArchiveFile.Close()

	if errorMessage, err = UnTarArchive(tempArchiveName, localPath); err != nil {
		return errorMessage, err
	}

	return "", nil
}
