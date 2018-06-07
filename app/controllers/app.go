package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

// Index returns just index page
func (c App) Index() revel.Result {
	return c.Render()
}
