package controllers

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

func (c App) GetPath() revel.Result {
	path := c.Params.Get("path")
	remote := c.Params.Get("remote")
	data := make(map[string]interface{})

	isRemote, err := strconv.ParseBool(remote)
	if err != nil {
		logger.Warnf("Are we dealing with remote path or local?: %v", err)
		response := CompileJSONResult(false, "Are we dealing with remote path or local?")
		return c.RenderJSON(response)
	}

	directory, _ := filepath.Split(path)
	if isRemote {
		if directory != "/" {
			directory = strings.TrimRight(directory, "/")
		}
	} else {
		directory = filepath.Clean(directory)
	}

	data["path"] = directory
	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
