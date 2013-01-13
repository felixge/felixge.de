package fs

import (
	"net/http"
	"path"
	"runtime"
	"github.com/felixge/magicfs"
)

var (
	_, filename, _, _ = runtime.Caller(0)
	__dirname         = path.Dir(filename)
	root              = path.Join(__dirname, "..")
)

func New() http.FileSystem {
	pageFs := magicfs.
		NewMagicFs(http.Dir(root)).
		ExecMap(".md", ".html", __dirname+"/processors/bin/markdown.js");

	return magicfs.
		NewMagicFs(http.Dir(root + "/public")).
		Exclude(".*").
		ExecMap(".less", ".css", __dirname+"/processors/bin/less.js").
		Or(newPages(pageFs))
}
