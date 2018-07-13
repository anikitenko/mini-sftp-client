package controllers

import (
	"github.com/revel/revel"
	"path/filepath"
	"os"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
)

func (c App) ShowFileContent() revel.Result {
	var fileContent []byte
	var response map[string]interface{}
	name := c.Params.Get("name")
	path := c.Params.Get("path")
	data := make(map[string]interface{})
	newPath := path + string(filepath.Separator) + name

	fileDirectoryStat, err := os.Stat(newPath)
	if err != nil {
		logger.Warnf("Problem with accessing local files. Lack of permissions OR directory doesn't exist: %v", err)
		response := CompileJSONResult(false, "Problem with listing local files. Lack of permissions OR directory doesn't exist")
		return c.RenderJSON(response)
	}

	switch mode := fileDirectoryStat.Mode(); {
	case mode.IsDir():
		response = CompileJSONResult(false, "This is a directory!")
		break
	case mode.IsRegular():
		fileContent, err = ioutil.ReadFile(newPath)
		data["contents"] = string(fileContent)
		response = CompileJSONResult(true, "", data)
		break
	default:
		response = CompileJSONResult(false, "Unknown file type!")
		break
	}

	return c.RenderJSON(response)
}