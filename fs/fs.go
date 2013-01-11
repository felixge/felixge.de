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
	pages := newPages(http.Dir(root))
	public := newLessProcessor(http.Dir(root + "/public"))

	chain := magicfs.NewChainFs(public, pages)
	fs := newExclude(chain, ".*")
	return fs
}
