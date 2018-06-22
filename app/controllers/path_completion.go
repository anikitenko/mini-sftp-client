package controllers

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// GetRemotePathCompletion acts like double tab for remote and returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// items: list of found files/folders
func (c App) GetRemotePathCompletion() revel.Result {
	path := c.Params.Get("path")
	data := make(map[string]interface{})

	if connection := c.EstablishSSHConnection(); connection != nil {
		return connection
	}

	defer SSHsession.Close()
	defer SSHclient.Close()

	dataCompletion, errString, err := RemotePathCompletion(SSHsession, path)
	if err != nil {
		logger.Warnf("%s: %v", errString, err)
		response := CompileJSONResult(false, errString)
		return c.RenderJSON(response)
	}

	data["items"] = dataCompletion
	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}

func RemotePathCompletion(session *ssh.Session, path string) ([]map[string]interface{}, string, error) {
	var dataCompletion []map[string]interface{}
	var errorMessage string

	dataResult, err := session.Output(`compgen -o default ` + path)
	if err != nil {
		errorMessage = "Problem with getting path"
		return dataCompletion, errorMessage, err
	}

	listOfCompletionFound := strings.Split(string(dataResult), "\n")
	for _, data := range listOfCompletionFound {
		if strings.TrimSpace(data) == "" {
			continue
		}
		dirFilePath := map[string]interface{}{}
		dirFilePath["id"] = data
		dirFilePath["text"] = data
		dataCompletion = append(dataCompletion, dirFilePath)
	}

	return dataCompletion, "", nil
}

// GetLocalPathCompletion acts like double tab for local and returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// items: list of found files/folders
func (c App) GetLocalPathCompletion() revel.Result {
	path := c.Params.Get("path")
	data := make(map[string]interface{})

	dataCompletion, errString, err := LocalPathCompletion(path)
	if err != nil {
		logger.Warnf("%s: %v", errString, err)
		response := CompileJSONResult(false, errString)
		return c.RenderJSON(response)
	}

	data["items"] = dataCompletion
	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}

func LocalPathCompletion(path string) ([]map[string]interface{}, string, error) {
	var dataCompletion []map[string]interface{}
	var errorMessage string

	directory, file := filepath.Split(path)

	listDirectory, err := ioutil.ReadDir(directory)
	if err != nil {
		errorMessage = "Problem with listing directory"
		return dataCompletion, errorMessage, err
	}

	r := regexp.MustCompile(`^` + file)
	for _, file := range listDirectory {
		if r.MatchString(file.Name()) {
			dirFilePath := map[string]interface{}{}
			dirFilePath["id"] = directory + file.Name()
			dirFilePath["text"] = directory + file.Name()
			dataCompletion = append(dataCompletion, dirFilePath)
		}
	}

	return dataCompletion, "", nil
}
