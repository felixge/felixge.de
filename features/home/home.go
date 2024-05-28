package home

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/deps"
	"github.com/felixge/felixge.de/features/posts"
)

type templateData struct {
	Posts []*posts.Post
}

func HomeHandler(d *deps.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := posts.Posts(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmplData := &templateData{Posts: posts}
		d.Templates.Render(w, "home.html", tmplData)
	}
}
