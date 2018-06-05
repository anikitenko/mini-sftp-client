package controllers

import (
	"os/user"
	"path/filepath"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

func (c App) ConnectViaSSH() revel.Result {
	var resultMessage []string
	data := make(map[string]interface{})

	if connection := c.EstablishSSHConnection(); connection != nil {
		return connection
	}

	defer SSHsession.Close()
	defer SSHclient.Close()

	if currentUserPathBytes, err := SSHsession.Output("echo -n $PWD"); err == nil {
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
		}
	} else {
		homeDirectory = username.HomeDir
	}

	data["local_path"] = homeDirectory

	data["errors"] = resultMessage

	data["local_path_separator"] = string(filepath.Separator)

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
