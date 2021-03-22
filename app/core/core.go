package core

import (
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/lib/websession"
)

// App -
type App struct {
	Router  *router.Mux
	Storage *datastorage.Storage
	Render  *htmltemplate.Engine
	Sess    *websession.Session
	Plugins *PluginSystem
}
