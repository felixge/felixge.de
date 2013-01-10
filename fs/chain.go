package fs

import (
	"net/http"
	"os"
)

type chain struct {
	list []http.FileSystem
}

type chainFile struct {
	http.File
	others []http.FileSystem
	path string
}

func (file *chainFile) Readdir(count int) ([]os.FileInfo, error) {
	dirs, err := file.File.Readdir(count)
	if err != nil {
		return dirs, err
	}

	for _, fs := range file.others {
		otherFile, err := fs.Open(file.path)
		if err != nil {
			continue
		}

		remaining := 0
		if count > 0 {
			remaining = count - len(dirs)
		}

		otherDirs, err := otherFile.Readdir(remaining)
		if err != nil {
			continue
		}

		dirs = append(dirs, otherDirs...)
	}

	return removeDuplicates(dirs), nil
}

func removeDuplicates(stats []os.FileInfo) []os.FileInfo {
	names := make(map[string]bool)
	results := make([]os.FileInfo, 0, len(stats))

	for _, stat := range stats {
		name := stat.Name()
		if _, ok := names[name]; ok {
			continue
		}

		names[name] = true
		results = append(results, stat)
	}
	return results
}

func newChain(list ...http.FileSystem) http.FileSystem {
	return &chain{list}
}

func (chain *chain) Open(path string) (http.File, error) {
	var err error
	for i, fs := range chain.list {
		file, openErr := fs.Open(path)
		if openErr == nil {
			return &chainFile{
				File: file,
				path: path,
				others: chain.list[i+1:],
			}, nil
		}
		err = openErr
	}
	return nil, err
}
