package controllers

import "github.com/revel/revel"

func (c App) GetStoredConnections() revel.Result {
	data := make(map[string]interface{})

	data["connections"] = StoredConnection

	response := CompileJSONResult(true, "", data)
	return c.RenderJSON(response)
}