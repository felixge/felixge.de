package appcontext

import (
	"log/slog"

	"github.com/felixge/felixge.de/commons/config"
	"github.com/felixge/felixge.de/commons/templates"
)

type AppContext struct {
	Config    *config.Config
	Log       *slog.Logger
	Templates *templates.Templates
}
