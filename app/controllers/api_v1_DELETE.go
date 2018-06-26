package controllers

import "github.com/revel/revel"

// @Summary Disconnect
// @Description remove stored connection
// @ID disconnect
// @Accept  json
// @Produce  json
// @Param   id      	path   	string     	true  	"Connection ID"
// @Success 200 {object} controllers.GeneralResponse	"Success"
// @Failure 403 {string} string "Not authorized!"
// @Router /disconnect/{id} [delete]
func (c ApiV1) Disconnect(id string) revel.Result {
	if _, ok := ApiConnections[id]; !ok {
		response := CompileJSONResult(false, "Connection does not exist!")
		return c.RenderJSON(response)
	}

	delete(ApiConnections, id)

	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
