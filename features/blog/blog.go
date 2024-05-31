package blog

import (
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/felixge/felixge.de/features/posts"
)

type templateData struct {
	Posts []*posts.Post
}

func Handler(d *appcontext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := posts.Posts(r.Context(), d.Log)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmplData := &templateData{Posts: posts}
		d.Templates.Render(w, "blog.html", tmplData)
	}
}
