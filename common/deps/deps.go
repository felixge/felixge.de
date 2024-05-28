package deps

import (
	"log/slog"

	"github.com/felixge/felixge.de/common/config"
)

type Deps struct {
	Config *config.Config
	Log    *slog.Logger
}
