package home

import (
	"net/http"

	"github.com/felixge/felixge.de/common/deps"
)

func HomeHandler(d *deps.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Nice World!"))
	}
}
