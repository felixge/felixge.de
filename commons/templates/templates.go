package templates

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

func NewTemplates(resources embed.FS) *Templates {
	var paths []string
	fs.WalkDir(resources, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(d.Name(), ".html") {
			paths = append(paths, path)
		}
		return nil
	})
	return &Templates{t: template.Must(template.ParseFS(resources, paths...))}
}

type Templates struct {
	t *template.Template
}

func (t *Templates) Render(w http.ResponseWriter, name string, data any) {
	var buffer bytes.Buffer

	err := t.t.ExecuteTemplate(&buffer, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	buffer.WriteTo(w)
}
