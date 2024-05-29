package router

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/felixge/felixge.de/commons/datadog"
	"github.com/felixge/felixge.de/features/home"
)

func NewRouter(d *appcontext.AppContext) http.Handler {
	mux := datadog.NewServeMux()
	mux.HandleFunc("GET /", home.HomeHandler(d))
	return mux
}
