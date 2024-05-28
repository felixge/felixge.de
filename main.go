package main

import (
	"net/http"

	"github.com/felixge/felixge.de/common/config"
	"github.com/felixge/felixge.de/common/deps"
	"github.com/felixge/felixge.de/common/logger"
	"github.com/felixge/felixge.de/common/router"
)

func main() {
	d := &deps.Deps{
		Config: config.EnvConfig(),
		Log:    logger.NewLogger(),
	}
	router := router.NewRouter(d)
	d.Log.Info("listening", "url", "http://"+d.Config.Addr)
	if err := http.ListenAndServe(d.Config.Addr, router); err != nil {
		d.Log.Error("failed to listen", "err", err)
	}
	return
}
