package fs

import (
	"net/http"
	"github.com/transloadit/magicfs"
)

func New(baseFs http.FileSystem) http.FileSystem {
	return magicfs.New(baseFs)
}
