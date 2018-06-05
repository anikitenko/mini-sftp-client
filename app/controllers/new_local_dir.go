package controllers

import (
	"os"
	"path/filepath"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

func (c App) CreateNewLocalDirectory() revel.Result {
	name := c.Params.Get("name")
	path := c.Params.Get("path")

	data := make(map[string]interface{})
	newPath := path + string(filepath.Separator) + name

	if name == "" {
		logger.Warn("Directory name is empty")
		response := CompileJSONResult(false, "Directory name is empty")
		return c.RenderJSON(response)
	}

	if err := os.Mkdir(newPath, 0755); err != nil {
		logger.Warnf("Problem with creating new directory: %v", err)
		response := CompileJSONResult(false, "Problem with creating new directory")
		return c.RenderJSON(response)
	}

	data["new_path"] = newPath

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
