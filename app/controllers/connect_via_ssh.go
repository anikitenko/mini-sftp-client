package controllers

import (
	"os/user"
	"path/filepath"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

// ConnectViaSSH returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// remote_path: remote home directory
// local_path: local home directory
// local_path_separator: local separator which is needed for navigation compatible with Windows/*nix based systems
// errors: array of messages received during the whole process of connecting via SSH
func (c App) ConnectViaSSH() revel.Result {
	var resultMessage []string
	data := make(map[string]interface{})

	if connection := c.EstablishSSHConnection(); connection != nil {
		return connection
	}

	defer SSHsession.Close()
	defer SSHclient.Close()

	if currentUserPathBytes, err := SSHsession.Output(`echo -n "$PWD"`); err == nil {
		data["remote_path"] = string(currentUserPathBytes)
	} else {
		data["remote_path"] = ""
		logger.Warnf("Unable to determine remote path: %v", err)
		resultMessage = append(resultMessage, "Unable to determine remote path")
	}

	var homeDirectory string
	if username, err := user.Current(); err != nil {
		if currentAbsPath, err := filepath.Abs("./"); err == nil {
			homeDirectory = currentAbsPath
		} else {
			homeDirectory = ""
			logger.Warnf("Problem with getting absolute path: %v", err)
		}
		logger.Warnf("Problem with getting current user: %v", err)
	} else {
		homeDirectory = username.HomeDir
	}

	data["local_path"] = homeDirectory

	data["errors"] = resultMessage

	data["local_path_separator"] = string(filepath.Separator)

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
