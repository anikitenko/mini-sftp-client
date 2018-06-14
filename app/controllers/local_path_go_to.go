package controllers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

// LocalPathGoTo returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// local_files: list of local files
// local_path_separator: local separator which is needed for navigation compatible with Windows/*nix based systems
func (c App) LocalPathGoTo() revel.Result {
	var localFilesList []FileStructureStruct

	path := c.Params.Get("path")
	data := make(map[string]interface{})

	fileDirectoryStat, err := os.Stat(path)
	if err != nil {
		logger.Warnf("Problem with accessing local files. Lack of permissions OR directory doesn't exist: %v", err)
		response := CompileJSONResult(false, "Problem with listing local files. Lack of permissions OR directory doesn't exist")
		return c.RenderJSON(response)
	}

	switch mode := fileDirectoryStat.Mode(); {
	case mode.IsDir():
		if err := os.Chdir(path); err != nil {
			logger.Warnf("Problem with changing directory. Lack of permissions OR directory doesn't exist: %v", err)
			response := CompileJSONResult(false, "Problem with changing directory. Lack of permissions OR directory doesn't exist")
			return c.RenderJSON(response)
		}
		localFiles, err := ioutil.ReadDir(".")
		if err != nil {
			logger.Warnf("Problem with listing local files. Lack of permissions OR directory doesn't exist: %v", err)
			response := CompileJSONResult(false, "Problem with listing local files. Lack of permissions OR directory doesn't exist")
			return c.RenderJSON(response)
		} else {
			for _, line := range localFiles {
				var file FileStructureStruct

				file.Path = line.Name()

				if line.IsDir() {
					file.Directory = true
				}

				if symPath, err := filepath.EvalSymlinks(file.Path); err == nil {
					if symPath != file.Path {
						file.Path += " -> " + symPath
						file.Symlink = true

						if osStat, err := os.Stat(symPath); err == nil {
							if osStat.IsDir() {
								file.Directory = true
							}
						}
					}
				} else {
					logger.Warnf("Could not resolve symlink: %v", err)
				}
				localFilesList = append(localFilesList, file)
			}
		}
		data["local_files"] = localFilesList
		data["local_path_separator"] = string(filepath.Separator)
	default:
		logger.Warnf("Specified path is not a directory: %v", err)
		response := CompileJSONResult(false, "Specified path is not a directory")
		return c.RenderJSON(response)
	}

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
