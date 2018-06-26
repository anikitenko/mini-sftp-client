package controllers

import (
	"github.com/revel/revel"
)

type ApiV1 struct {
	*revel.Controller
}

// @title Mini sFTP client API
// @version 1.0.0
// @description This API was created to automate your work with sFTP client

// @license.name MIT
// @license.url https://github.com/anikitenko/mini-sftp-client/blob/staging/LICENSE

// @host 127.0.0.1:9000
// @BasePath /api/v1

func (c ApiV1) Help() revel.Result {
	return c.Redirect("/api/v1/index.html")
}