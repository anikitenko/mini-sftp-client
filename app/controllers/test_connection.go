package controllers

import (
	"github.com/revel/revel"
)

// TestSSHConnection returns if connection via SSH was made successfully or not
// it just executed function EstablishSSHConnection
func (c App) TestSSHConnection() revel.Result {
	data := make(map[string]interface{})
	if connect := c.EstablishSSHConnection(); connect != nil {
		return connect
	}

	SSHsession.Close()
	SSHclient.Close()

	data["pin_code"] = PinCode

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}
