package controllers

import (
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (c App) RemoveLocalItem() revel.Result {
	name := c.Params.Get("name")
	path := c.Params.Get("path")
	newPath := path + string(filepath.Separator) + name

	if name == "" {
		response := CompileJSONResult(false, "Name is empty")
		return c.RenderJSON(response)
	}

	if err := os.RemoveAll(newPath); err != nil {
		logger.Warnf("Problem with removing "+newPath+": %v", err)
		response := CompileJSONResult(false, "Problem with removing "+newPath)
		return c.RenderJSON(response)
	}

	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}