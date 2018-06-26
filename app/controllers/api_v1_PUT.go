package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
)

// @Summary Connect via SSH
// @Description connect via SSH and store connection
// @ID connect
// @Accept  json
// @Produce  json
// @Param   JSON      	body   	controllers.ApiConnectionStruct     	true  	"Connection information"
// @Success 200 {object} controllers.GeneralResponse	"Success with connection ID in message"
// @Failure 403 {string} string "Not authorized!"
// @Router /connect [put]
func (c ApiV1) Connect() revel.Result {
	jsonReceive := c.Params.JSON
	connectionDetails := ApiConnectionStruct{}

	err := json.Unmarshal(jsonReceive, &connectionDetails)
	if err != nil {
		logger.Warnf("Invalid JSON: %v", err)
		response := CompileJSONResult(false, "Invalid JSON")
		return c.RenderJSON(response)
	}

	if errString, err := ConnectSSH(connectionDetails.Ip, connectionDetails.User, connectionDetails.Password, connectionDetails.Port); err != nil {
		logger.Warnf("%s: %v", errString, err)
		response := CompileJSONResult(false, errString)
		return c.RenderJSON(response)
	}

	defer SSHsession.Close()
	defer SSHclient.Close()

	newConnectionId := RandStringBytes(20)

	if connectionDetails.Ip == MockSSHHostString {
		connectionDetails.Ip = "127.0.0.1"
	}

	ApiConnections[newConnectionId] = connectionDetails

	response := CompileJSONResult(true, newConnectionId)
	return c.RenderJSON(response)
}
