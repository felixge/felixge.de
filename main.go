package main

import (
	"embed"
	"net/http"

	"github.com/felixge/felixge.de/commons/config"
	"github.com/felixge/felixge.de/commons/deps"
	"github.com/felixge/felixge.de/commons/logger"
	"github.com/felixge/felixge.de/commons/router"
	"github.com/felixge/felixge.de/commons/templates"
)

//go:embed features/*/*.html commons/*/*.html
var html embed.FS

func main() {
	d := &deps.Deps{
		Config:    config.EnvConfig(),
		Log:       logger.NewLogger(),
		Templates: templates.NewTemplates(html),
	}
	router := router.NewRouter(d)
	d.Log.Info("listening", "url", "http://"+d.Config.Addr)
	if err := http.ListenAndServe(d.Config.Addr, router); err != nil {
		d.Log.Error("failed to listen", "err", err)
	}
	return
}
