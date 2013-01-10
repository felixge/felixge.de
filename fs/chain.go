package fs

import (
	"net/http"
)

type chain struct {
	list []http.FileSystem
}

func newChain(list ...http.FileSystem) http.FileSystem {
	return &chain{list}
}

func (chain *chain) Open(path string) (http.File, error) {
	var err error
	for _, fs := range chain.list {
		file, openErr := fs.Open(path)
		if openErr == nil {
			return file, nil
		}
		err = openErr
	}
	return nil, err
}
