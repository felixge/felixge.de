package router

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/felixge/felixge.de/commons/datadog"
	"github.com/felixge/felixge.de/features/blog"
	"github.com/felixge/felixge.de/features/pages"
)

func NewRouter(a *appcontext.AppContext) (http.Handler, error) {
	pageRoutes, err := pages.Routes(a)
	if err != nil {
		return nil, err
	}

	mux := datadog.NewServeMux()
	mux.HandleFunc("GET /blog", blog.Handler(a))
	for _, pr := range pageRoutes {
		mux.HandleFunc("GET "+pr.Route, pages.Handler(a, pr.Path))
	}
	return mux, nil
}
