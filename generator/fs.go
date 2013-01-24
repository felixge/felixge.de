package generator

import (
	"encoding/json"
	"net/http"
	"path"
	"runtime"
	"github.com/felixge/makefs"
	"html/template"
	"io"
	"io/ioutil"
)

var (
	_, filename, _, _ = runtime.Caller(0)
	__dirname         = path.Dir(filename)
	root              = path.Join(__dirname, "..")
)

func NewFs() http.FileSystem {
	fs := makefs.NewFs(http.Dir(root))

	fs.ExecMake("%.html", "%.md", __dirname+"/processors/bin/markdown.js")
	fs.ExecMake("%.css", "%.less", __dirname+"/processors/bin/less.js")

	fs.Make(
		"/public/index.html",
		[]string{"/pages/index.html", "/layouts/default.html", "/talks.json"},
		func(t *makefs.Task) error {
			sources := t.Sources()

			talksJson, err := ioutil.ReadAll(sources["/talks.json"])
			if err != nil {
				return err
			}

			talks := make([]*talk, 0)
			if err := json.Unmarshal(talksJson, &talks); err != nil {
				return err
			}

			viewVars := map[string]interface{}{
				"Talks": talks,
			}

			return render(
				t.Target(),
				sources["/pages/index.html"],
				sources["/layouts/default.html"],
				viewVars,
			)
		},
	)

	return fs.SubFs("/public")
}

func render(w io.Writer, page, layout io.Reader, viewVars interface{}) error {
	pageHtml, err := ioutil.ReadAll(page)
	if err != nil {
		return err
	}

	tmpl, err := template.New("page").Parse(string(pageHtml))
	if err != nil {
		return err
	}

	layoutHtml, err := ioutil.ReadAll(layout)
	if err != nil {
		return err
	}

	tmpl, err = tmpl.New("layout").Parse(string(layoutHtml))
	if err != nil {
		return err
	}

	return tmpl.Execute(w, viewVars)
}

type talk struct {
	Title    string
	Location string
	Date     string
	Url      string
	EventUrl string
	VideoUrl string
	PdfUrl   string
	CodeUrl  string
}
