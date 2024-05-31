package appcontext

import (
	"io/fs"
	"log/slog"

	"github.com/felixge/felixge.de/commons/config"
	"github.com/felixge/felixge.de/commons/markdown"
	"github.com/felixge/felixge.de/commons/templates"
)

type AppContext struct {
	Config    *config.Config
	Log       *slog.Logger
	Templates *templates.Templates
	Markdown  *markdown.Markdown
	FS        fs.FS
}
