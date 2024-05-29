package router

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/felixge/felixge.de/features/home"
)

func NewRouter(d *appcontext.AppContext) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home.HomeHandler(d))
	return mux
}
