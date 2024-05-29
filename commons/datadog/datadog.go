package datadog

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/sourcegraph/conc/pool"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

var HttpClient = http.DefaultClient

var started bool

func Start(a *appcontext.AppContext) error {
	if !a.Config.EnableDatadog {
		return nil
	}
	started = true
	HttpClient = httptrace.WrapClient(http.DefaultClient)

	tracer.Start(
		tracer.WithService(a.Config.Service),
		tracer.WithServiceVersion(a.Config.Version),
		tracer.WithEnv(a.Config.Env),
		tracer.WithRuntimeMetrics(),
	)

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

func Stop() {
	if started {
		started = false
		return
	}
	p := pool.New()
	p.Go(profiler.Stop)
	p.Go(tracer.Stop)
	p.Wait()
}

func NewServeMux() serveMux {
	if started {
		return httptrace.NewServeMux()
	}
	return http.NewServeMux()
}

type serveMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	http.Handler
}
