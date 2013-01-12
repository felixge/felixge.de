package fs

import (
	"net/http"
	"path"
	"runtime"
	"github.com/felixge/magicfs"
	"io"
	"os/exec"
)

var (
	_, filename, _, _ = runtime.Caller(0)
	__dirname         = path.Dir(filename)
	root              = path.Join(__dirname, "..")
)

func New() http.FileSystem {
	fs := magicfs.NewMagicFs(http.Dir(root + "/public"))
	fs.Exclude(".*")
	fs.Map(".less", ".css", func(less io.Reader) (io.Reader) {
		r, w := io.Pipe()
		cmd := exec.Command("lessc", "-", "--no-color")

		cmd.Stdin = less
		cmd.Stderr = w
		cmd.Stdout = w

		go func() {
			err := cmd.Run()
			if err != nil {
				w.Write([]byte("lessc: "+err.Error()))
			}
			w.CloseWithError(err)
		}()

		return r
	})
	fs.Or(newPages(http.Dir(root)))

	return fs
}
