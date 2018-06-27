package controllers

import (
	"path/filepath"
	"strconv"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"os"
	"fmt"
	"time"
)

// Download returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// If success function downloads file or directory locally
func (c App) Download() revel.Result {
	sourcePath, localPath := c.Params.Get("source_path"), c.Params.Get("local_path")
	isDir, fileNamePost := c.Params.Get("is_dir"), c.Params.Get("file_name")
	needToBackup := c.Params.Get("backup")
	fileName := filepath.Base(fileNamePost)

	if connection := c.EstablishSSHConnection(); connection != nil {
		return connection
	}

	defer SSHclient.Close()

	dirOrNot, err := strconv.ParseBool(isDir)
	if err != nil {
		logger.Warnf("We don't understand if it's a file or directory: %v", err)
		response := CompileJSONResult(false, "We don't understand if it's a file or directory")
		return c.RenderJSON(response)
	}

	backup, err := strconv.ParseBool(needToBackup)
	if err != nil {
		logger.Warnf("We don't understand if we need to backup or not: %v", err)
		response := CompileJSONResult(false, "We don't understand if we need to backup or not")
		return c.RenderJSON(response)
	}

	if dirOrNot {
		if errString, err := DownloadDirectory(sourcePath, localPath, fileName, SSHsession); err != nil {
			logger.Warnf("%s: %v", errString, err)
			response := CompileJSONResult(false, errString)
			return c.RenderJSON(response)
		}
	} else {
		if backup {
			newName := fmt.Sprintf("%s_%v", localPath+string(filepath.Separator)+fileName, time.Now().UnixNano())
			os.Rename(localPath+string(filepath.Separator)+fileName, newName)
		}

		if errString, err := DownloadFile(localPath, fileName, fileNamePost, SSHsession); err != nil {
			logger.Warnf("%s: %v", errString, err)
			response := CompileJSONResult(false, errString)
			return c.RenderJSON(response)
		}
	}
	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
