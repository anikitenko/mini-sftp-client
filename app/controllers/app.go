package controllers

import (
	"github.com/revel/revel"
	"time"
)

type App struct {
	*revel.Controller
}

// Index returns just index page
func (c App) Index() revel.Result {
	testParam := c.Params.Get("for_testing")
	userPinCode := c.Session["pin_code"]
	c.ViewArgs["noPinCode"] = false

	if (PinCode != userPinCode && c.ClientIP != "127.0.0.1") || testParam == "true" {
		c.ViewArgs["noPinCode"] = true
	}

	return c.Render()
}

func (c App) SetPinCode() revel.Result {
	userPinCode := c.Params.Get("pin_code")

	if PinCode != userPinCode {
		TimeToWaitInvalidPin = TimeToWaitInvalidPin + time.Second

		time.Sleep(TimeToWaitInvalidPin)

		response := CompileJSONResult(false, "Pin code is incorrect")
		return c.RenderJSON(response)
	}

	c.Session["pin_code"] = userPinCode

	response := CompileJSONResult(true, "")
	return c.RenderJSON(response)
}
