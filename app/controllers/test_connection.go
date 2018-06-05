package controllers

import (
	"github.com/revel/revel"
)

func (c App) TestSSHConnection() revel.Result {
	if connect := c.EstablishSSHConnection(); connect != nil {
		return connect
	}

	SSHsession.Close()
	SSHclient.Close()
	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
