package controllers

import (
	"path/filepath"
	"strconv"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

// Download returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// If success function downloads file or directory locally
func (c App) Download() revel.Result {
	sourcePath, localPath := c.Params.Get("source_path"), c.Params.Get("local_path")
	isDir, fileNamePost := c.Params.Get("is_dir"), c.Params.Get("file_name")
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

	if dirOrNot {
		if errString, err := downloadDirectory(sourcePath, localPath, fileName); err != nil {
			logger.Warnf("%s: %v", errString, err)
			response := CompileJSONResult(false, errString)
			return c.RenderJSON(response)
		}
	} else {
		if errString, err := downloadFile(localPath, fileName, fileNamePost); err != nil {
			logger.Warnf("%s: %v", errString, err)
			response := CompileJSONResult(false, errString)
			return c.RenderJSON(response)
		}
	}
	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
