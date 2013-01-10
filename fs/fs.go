package fs

import (
	"net/http"
	"path"
	"runtime"
)

var (
	_, filename, _, _ = runtime.Caller(0)
	root              = path.Join(path.Dir(filename), "..")
)

func New() http.FileSystem {
	pages := newPages(http.Dir(root))
	public := newLessProcessor(http.Dir(root + "/public"))

	fs := newChain(public, pages)
	return fs
}
