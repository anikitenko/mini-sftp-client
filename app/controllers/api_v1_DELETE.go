package controllers

import "github.com/revel/revel"

func (c ApiV1) Disconnect(id string) revel.Result {
	if _, ok := ApiConnections[id]; !ok {
		response := CompileJSONResult(false, "Connection does not exist!")
		return c.RenderJSON(response)
	}

	delete(ApiConnections, id)

	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
