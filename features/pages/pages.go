package pages

import (
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/felixge/felixge.de/commons/appcontext"
)

type templateData struct {
	Title string
	HTML  template.HTML
}

func Handler(d *appcontext.AppContext, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := d.Markdown.Render(path)
		if err != nil {
			return
		}
		// string(data).pri
		tmplData := &templateData{
			Title: result.Title,
			HTML:  template.HTML(result.HTML),
		}
		d.Templates.Render(w, "page.html", tmplData)
	}
}

const pagesDir = "content/pages"

func Routes(a *appcontext.AppContext) ([]*Route, error) {
	var routes []*Route
	fs.WalkDir(a.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(path, pagesDir) && strings.Contains(d.Name(), ".md") {
			name := strings.TrimSuffix(d.Name(), ".md")
			route := &Route{
				Route: "/" + name,
				Path:  path,
			}
			if name == "home" {
				route.Route = "/"
			}
			routes = append(routes, route)
		}
		return nil
	})
	return routes, nil
}

type Route struct {
	Route string
	Path  string
}
