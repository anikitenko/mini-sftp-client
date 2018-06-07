package controllers

import (
	"regexp"
	"strings"

	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

// RemotePathGoTo returns JSON which contains:
// result: true for success/false for any error
// message: empty if success, message if error
// remote_files: list of remote files
func (c App) RemotePathGoTo() revel.Result {
	var remoteFilesList []FileStructureStruct

	path := c.Params.Get("path")
	data := make(map[string]interface{})

	if connection := c.EstablishSSHConnection(); connection != nil {
		return connection
	}

	defer SSHclient.Close()

	checkDirectory, err := SSHsession.Output("test -r '" + path + "' -a -d '" + path + "'; echo -n $?")
	if err != nil {
		logger.Warnf("We were unable to check if this is a directory from the path specified: %v", err)
		response := CompileJSONResult(false, "We were unable to check if this is a directory from the path specified")
		return c.RenderJSON(response)
	}

	SSHsession.Close()

	if string(checkDirectory) == "0" {
		sessionListFiles, err := SSHclient.NewSession()
		if err != nil {
			logger.Warnf("Unable to create new SSH session to list files in directory: ", err)
			response := CompileJSONResult(false, "Unable to create new SSH session to list files in directory")
			return c.RenderJSON(response)
		}
		remoteFilesExecute, err := sessionListFiles.Output("cd '" + path + "' && ls -lA | grep -v 'total'")
		if err != nil {
			logger.Warnf("Problem with getting list of remote files: %v", err)
			response := CompileJSONResult(false, "Problem with getting list of remote files OR directory is empty")
			return c.RenderJSON(response)
		}
		sessionListFiles.Close()
		regexpDirectory := regexp.MustCompile(`^d`)
		regexpSymlink := regexp.MustCompile(`^l`)
		for _, line := range strings.Split(string(remoteFilesExecute), "\n") {
			var file FileStructureStruct

			if line == "" {
				continue
			}

			lineSplit := strings.Fields(line)
			if len(lineSplit) < 8 {
				continue
			}

			file.Path = strings.Join(lineSplit[8:], " ")

			if regexpDirectory.MatchString(line) {
				file.Directory = true
			}
			if regexpSymlink.MatchString(line) {
				file.Symlink = true
				sessionCheckSymlink, err := SSHclient.NewSession()
				if err != nil {
					logger.Warnf("Unable to create new SSH session to check if symlink is a directory: ", err)
					response := CompileJSONResult(false, "Unable to create new SSH session to check if symlink is a directory")
					return c.RenderJSON(response)
				}
				pathFromSymlinkSplit := strings.Split(line, "->")
				pathFromSymlink := strings.TrimSpace(pathFromSymlinkSplit[len(pathFromSymlinkSplit)-1])
				checkIfSymlinkPointsToDir, err := sessionCheckSymlink.Output("cd '" + path + "' && test -r '" + pathFromSymlink + "' -a -d '" + pathFromSymlink + "'; echo -n $?")
				sessionCheckSymlink.Close()
				if string(checkIfSymlinkPointsToDir) == "0" {
					file.Directory = true
				}
			}
			remoteFilesList = append(remoteFilesList, file)
		}
		data["remote_files"] = remoteFilesList
	} else {
		response := CompileJSONResult(false, "You specified wrong path OR you don't have permission to access it OR this is a file")
		return c.RenderJSON(response)
	}

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
