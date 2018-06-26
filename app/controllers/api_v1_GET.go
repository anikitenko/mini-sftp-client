package controllers

import (
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"path/filepath"
)

// @Summary Get connections
// @Description get list of stored connections
// @ID get-connections
// @Produce  json
// @Success 200 {object} controllers.GetConnectionsStruct	"Success"
// @Failure 403 {string} string "Not authorized!"
// @Router /getConnections [get]
func (c ApiV1) GetConnections(id string) revel.Result {
	data := make(map[string]interface{})

	if id == "" {
		data["connections"] = ApiConnections
		response := CompileJSONResult(true, "", data)
		return c.RenderJSON(response)
	}

	if connection, ok := ApiConnections[id]; !ok {
		response := CompileJSONResult(false, "Connection does not exist!")
		return c.RenderJSON(response)
	} else {
		connections := make(map[string]ApiConnectionStruct)
		connections[id] = connection
		data["connections"] = connections
	}

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}

// @Summary Get remote home directory
// @Description get remote home directory path
// @ID get-remote-home-directory-path
// @Produce  json
// @Param   id      	path   	string     	true  	"Connection ID"
// @Success 200 {object} controllers.GeneralResponse	"Success with home directory path in message"
// @Failure 403 {string} string "Not authorized!"
// @Router /getRemoteHomeDirectory/{id} [get]
func (c ApiV1) GetRemoteHomeDirectory(id string) revel.Result {
	var sshSessionConnect SSHSessionStruct
	var homePath string

	if connection, ok := ApiConnections[id]; !ok {
		response := CompileJSONResult(false, "Connection does not exist!")
		return c.RenderJSON(response)
	} else {
		sshSessionConnect = createSession(connection.Ip, connection.User, connection.Password, connection.Port)

		if sshSessionConnect.ErrorErr != nil {
			logger.Warnf("%s: %v", sshSessionConnect.ErrorStr, sshSessionConnect.ErrorErr)
			response := CompileJSONResult(false, sshSessionConnect.ErrorStr)
			return c.RenderJSON(response)
		}
	}

	sshClient := sshSessionConnect.Client
	sshSession := sshSessionConnect.Session

	defer sshSession.Close()
	defer sshClient.Close()

	if currentUserPathBytes, err := sshSession.Output(`echo -n "$HOME"`); err == nil {
		homePath = string(currentUserPathBytes)
	} else {
		logger.Warnf("Unable to determine remote path: %v", err)
		response := CompileJSONResult(false, "Unable to determine remote path")
		return c.RenderJSON(response)
	}

	response := CompileJSONResult(true, homePath)
	return c.RenderJSON(response)
}

// @Summary Get remote path completion
// @Description list remote directories like "double tab"
// @ID get-remote-path-completion
// @Produce  json
// @Param   id      	path   	string     	true  	"Connection ID"
// @Param   path      	query   	string     	true  	"Remote path"
// @Success 200 {object} controllers.GetPathCompletionStruct	"Success with list of items"
// @Failure 403 {string} string "Not authorized!"
// @Router /getRemotePathCompletion/{id} [get]
func (c ApiV1) GetRemotePathCompletion(id string) revel.Result {
	var sshSessionConnect SSHSessionStruct
	var completionItems []string
	path := c.Params.Query.Get("path")
	data := make(map[string]interface{})

	if connection, ok := ApiConnections[id]; !ok {
		response := CompileJSONResult(false, "Connection does not exist!")
		return c.RenderJSON(response)
	} else {
		sshSessionConnect = createSession(connection.Ip, connection.User, connection.Password, connection.Port)

		if sshSessionConnect.ErrorErr != nil {
			logger.Warnf("%s: %v", sshSessionConnect.ErrorStr, sshSessionConnect.ErrorErr)
			response := CompileJSONResult(false, sshSessionConnect.ErrorStr)
			return c.RenderJSON(response)
		}
	}

	sshClient := sshSessionConnect.Client
	sshSession := sshSessionConnect.Session

	defer sshSession.Close()
	defer sshClient.Close()

	dataCompletion, errString, err := RemotePathCompletion(sshSession, path)
	if err != nil {
		logger.Warnf("%s: %v", errString, err)
		response := CompileJSONResult(false, errString)
		return c.RenderJSON(response)
	}

	for _, name := range dataCompletion {
		if text, ok := name["text"].(string); !ok {
			continue
		} else {
			completionItems = append(completionItems, text)
		}
	}

	data["items"] = completionItems

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}

// @Summary Get local home directory
// @ID get-local-home-directory
// @Produce  json
// @Success 200 {object} controllers.GeneralResponse	"Success with local path in message"
// @Failure 403 {string} string "Not authorized!"
// @Router /getLocalHomeDirectory [get]
func (c ApiV1) GetLocalHomeDirectory(id string) revel.Result {
	response := CompileJSONResult(true, LocalHomeDirectory())
	return c.RenderJSON(response)
}

// @Summary Get local path completion
// @Description list local directories like "double tab"
// @ID get-local-path-completion
// @Produce  json
// @Param   path      	query   	string     	true  	"Local path"
// @Success 200 {object} controllers.GetPathCompletionStruct	"Success with list of items"
// @Failure 403 {string} string "Not authorized!"
// @Router /getLocalPathCompletion [get]
func (c ApiV1) GetLocalPathCompletion(id string) revel.Result {
	var completionItems []string
	path := c.Params.Query.Get("path")
	data := make(map[string]interface{})

	dataCompletion, errString, err := LocalPathCompletion(path)
	if err != nil {
		logger.Warnf("%s: %v", errString, err)
		response := CompileJSONResult(false, errString)
		return c.RenderJSON(response)
	}

	for _, name := range dataCompletion {
		if text, ok := name["text"].(string); !ok {
			continue
		} else {
			completionItems = append(completionItems, text)
		}
	}

	data["items"] = completionItems

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}

// @Summary Download
// @Description download file or directory
// @ID download
// @Produce  json
// @Param   id      	path   	string     	true  	"Connection ID"
// @Param   path      	query   string     	true  	"Path to file OR directory"
// @Param	save_to		query	string		false	"If not empty than download to.. otherwise file/directory will be downloaded to local home directory"
// @Success 200 {object} controllers.GeneralResponse	"Success"
// @Failure 403 {string} string "Not authorized!"
// @Router /download/{id} [get]
func (c ApiV1) Download(id string) revel.Result {
	var sshSessionConnect SSHSessionStruct
	fileToDownload := c.Params.Query.Get("path")
	saveTo := c.Params.Query.Get("save_to")

	if saveTo == "" {
		saveTo = LocalHomeDirectory()
	}

	if connection, ok := ApiConnections[id]; !ok {
		response := CompileJSONResult(false, "Connection does not exist!")
		return c.RenderJSON(response)
	} else {
		sshSessionConnect = createSession(connection.Ip, connection.User, connection.Password, connection.Port)

		if sshSessionConnect.ErrorErr != nil {
			logger.Warnf("%s: %v", sshSessionConnect.ErrorStr, sshSessionConnect.ErrorErr)
			response := CompileJSONResult(false, sshSessionConnect.ErrorStr)
			return c.RenderJSON(response)
		}
	}

	sshClient := sshSessionConnect.Client
	sshSession := sshSessionConnect.Session

	defer sshSession.Close()
	defer sshClient.Close()

	commandToTestDir := "test -r '" + fileToDownload + "' -a -d '" + fileToDownload + "'; echo -n $?"
	commandToTestFile := "test -r '" + fileToDownload + "' -a -f '" + fileToDownload + "'; echo -n $?"

	testFileBytes, err := sshSession.Output(commandToTestFile)
	if err != nil {
		logger.Warnf("Unable to check if path is a file: %v", err)
		response := CompileJSONResult(false, "Unable to check if path is a file")
		return c.RenderJSON(response)
	}
	sshSession.Close()

	sshSessionCheckDir, err := sshClient.NewSession()
	if err != nil {
		logger.Warnf("Unable to create new SSH session to check if path is a directory: %v", err)
		response := CompileJSONResult(false, "Unable to create new SSH session to check if path is a directory")
		return c.RenderJSON(response)
	}

	testDirBytes, err := sshSessionCheckDir.Output(commandToTestDir)
	if err != nil {
		logger.Warnf("Unable to check if path is a directory: %v", err)
		response := CompileJSONResult(false, "Unable to check if path is a directory")
		return c.RenderJSON(response)
	}
	sshSessionCheckDir.Close()

	sshSessionDownload, err := sshClient.NewSession()
	if err != nil {
		logger.Warnf("Unable to create new SSH session to start downloading: %v", err)
		response := CompileJSONResult(false, "Unable to create new SSH session to start downloading")
		return c.RenderJSON(response)
	}

	fileName := filepath.Base(fileToDownload)
	filePath, _ := filepath.Split(fileToDownload)

	switch {
	case string(testFileBytes) == "0":
		if errString, err := DownloadFile(saveTo, fileName, fileToDownload, sshSessionDownload); err != nil {
			logger.Warnf("%s: %v", errString, err)
			response := CompileJSONResult(false, errString)
			return c.RenderJSON(response)
		}
	case string(testDirBytes) == "0":
		if errString, err := DownloadDirectory(filePath, saveTo, fileName, sshSessionDownload); err != nil {
			logger.Warnf("%s: %v", errString, err)
			response := CompileJSONResult(false, errString)
			return c.RenderJSON(response)
		}
	default:
		response := CompileJSONResult(false, "Unknown file type to download!")
		return c.RenderJSON(response)
	}

	response := CompileJSONResult(true, "File was downloaded and saved to '"+saveTo+"'")
	return c.RenderJSON(response)
}
