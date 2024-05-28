package deps

import (
	"log/slog"

	"github.com/felixge/felixge.de/commons/config"
	"github.com/felixge/felixge.de/commons/templates"
)

type Deps struct {
	Config    *config.Config
	Log       *slog.Logger
	Templates *templates.Templates
}
