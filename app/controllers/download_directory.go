package controllers

import (
	logger "github.com/sirupsen/logrus"
	"archive/tar"
	"os"
	"path/filepath"
	"io"
	"fmt"
	"time"
	"compress/gzip"
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

	archiveOpen, err := os.Open(tempArchiveName)
	if err != nil {
		errorMessage = "Problem with opening temp archive"
		return errorMessage, err
	}
	defer func() {
		archiveOpen.Close()
		if err := os.Remove(tempArchiveName); err != nil {
			logger.Warnf("Unable to remove temp archive: %v", err)
		}
	}()

	gzipOpen, err := gzip.NewReader(archiveOpen)
	if err != nil {
		errorMessage = "Problem with creating stream from temp archive"
		return errorMessage, err
	}

	tarReader := tar.NewReader(gzipOpen)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			logger.Info("Archive was successfully extracted")
			break
		}
		if err != nil {
			errorMessage = "Problem with reading temp archive"
			return errorMessage, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(localPath+string(filepath.Separator)+header.Name, 0755); err != nil {
				errorMessage = "Cannot create a directory from archive"
				return errorMessage, err
			}
			logger.Infof("Creating new directory at %s", localPath+string(filepath.Separator)+header.Name)
		case tar.TypeReg:
			outFile, err := os.Create(localPath + string(filepath.Separator) + header.Name)
			if err != nil {
				errorMessage = "Cannot create a file from archive"
				return errorMessage, err
			}

			logger.Infof("Creating new file at %s", localPath+string(filepath.Separator)+header.Name)
			if _, err := io.Copy(outFile, tarReader); err != nil {
				errorMessage = "Failed to write to a file from archive"
				return errorMessage, err
			}
			outFile.Close()
		default:
			errorMessage = "General archive problem, refer to logs"
			return errorMessage, err
		}
	}
	return "", nil
}