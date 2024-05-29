package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/felixge/felixge.de/commons/config"
	"github.com/felixge/felixge.de/commons/datadog"
	"github.com/felixge/felixge.de/commons/logger"
	"github.com/felixge/felixge.de/commons/router"
	"github.com/felixge/felixge.de/commons/templates"
	"github.com/felixge/felixge.de/features/posts"
)

//go:embed features/*/*.html commons/*/*.html
var html embed.FS

func main() {
	log := logger.NewLogger()
	if err := run(log); err != nil {
		log.Error("run error", "err", err)
		os.Exit(1)
	}
	return
}

func run(log *slog.Logger) error {
	ctx := context.Background()
	a := &appcontext.AppContext{
		Config:    config.EnvConfig(),
		Log:       log,
		Templates: templates.NewTemplates(html),
	}

	a.Log.Info("starting", "config", a.Config)

	if err := datadog.StartProfiling(a); err != nil {
		return err
	}

	flag.Parse()
	switch flag.Arg(0) {
	default:
		router := router.NewRouter(a)
		a.Log.Info("listening", "url", "http://"+a.Config.Addr)
		if err := http.ListenAndServe(a.Config.Addr, router); err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
	case "generate":
		return posts.WriteFeeds(ctx, a)
	}
	return nil
}
