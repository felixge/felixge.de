package router

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/felixge/felixge.de/commons/datadog"
	"github.com/felixge/felixge.de/features/pages"
)

func NewRouter(a *appcontext.AppContext) (http.Handler, error) {
	pageRoutes, err := pages.Routes(a)
	if err != nil {
		return nil, err
	}

	mux := datadog.NewServeMux()
	for _, pr := range pageRoutes {
		mux.HandleFunc("GET "+pr.Route, pages.PageHandler(a, pr.Path))
	}
	return mux, nil
}
