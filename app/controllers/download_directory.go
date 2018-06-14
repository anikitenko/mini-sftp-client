package controllers

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"io"
	"fmt"
	"time"
)

func downloadDirectory(sourcePath, localPath, fileName string) (string, error) {
	var errorMessage string
	stdoutNewPipe, err := SSHsession.StdoutPipe()
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

	if err := SSHsession.Start("tar cz --directory='" + sourcePath + "' '" + fileName + "'"); err != nil {
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

	SSHsession.Close()
	tempArchiveFile.Close()

	if errorMessage, err = UnTarArchive(tempArchiveName, localPath); err != nil {
		return errorMessage, err
	}

	return "", nil
}