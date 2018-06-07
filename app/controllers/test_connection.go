package controllers

import (
	"github.com/revel/revel"
)

// TestSSHConnection returns if connection via SSH was made successfully or not
// it just executed function EstablishSSHConnection
func (c App) TestSSHConnection() revel.Result {
	if connect := c.EstablishSSHConnection(); connect != nil {
		return connect
	}

	SSHsession.Close()
	SSHclient.Close()
	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
