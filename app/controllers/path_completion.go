package controllers

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

// GetRemotePathCompletion acts like double tab for remote and returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// items: list of found files/folders
func (c App) GetRemotePathCompletion() revel.Result {
	var dataCompletion []map[string]interface{}

	path := c.Params.Get("path")
	data := make(map[string]interface{})

	if connection := c.EstablishSSHConnection(); connection != nil {
		return connection
	}

	defer SSHsession.Close()
	defer SSHclient.Close()

	dataResult, err := SSHsession.Output(`compgen -o default ` + path)
	if err != nil {
		logger.Warnf("Problem with getting path: %v", err)
		response := CompileJSONResult(false, "Problem with getting path")
		return c.RenderJSON(response)
	}

	listOfCompletionFound := strings.Split(string(dataResult), "\n")
	for _, data := range listOfCompletionFound {
		dirFilePath := map[string]interface{}{}
		dirFilePath["id"] = data
		dirFilePath["text"] = data
		dataCompletion = append(dataCompletion, dirFilePath)
	}

	data["items"] = dataCompletion
	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}

// GetLocalPathCompletion acts like double tab for local and returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// items: list of found files/folders
func (c App) GetLocalPathCompletion() revel.Result {
	var (
		dataCompletion []map[string]interface{}
	)

	path := c.Params.Get("path")
	data := make(map[string]interface{})

	directory, file := filepath.Split(path)

	listDirectory, err := ioutil.ReadDir(directory)
	if err != nil {
		logger.Warnf("Problem with listing directory: %v", err)
		response := CompileJSONResult(false, "Problem with listing directory")
		return c.RenderJSON(response)
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
	data["items"] = dataCompletion
	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
