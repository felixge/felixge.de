package datadog

import (
	"github.com/felixge/felixge.de/commons/appcontext"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func StartProfiling(a *appcontext.AppContext) error {
	return profiler.Start(
		profiler.WithService(a.Config.Service),
		profiler.WithEnv(a.Config.Env),
		profiler.WithVersion(a.Config.Version),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.MutexProfile,
			profiler.GoroutineProfile,
		),
	)
}
