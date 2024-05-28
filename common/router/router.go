package router

import (
	"net/http"

	"github.com/felixge/felixge.de/common/deps"
	"github.com/felixge/felixge.de/features/home"
)

func NewRouter(d *deps.Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home.HomeHandler(d))
	return mux
}
