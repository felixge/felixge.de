package main

import (
	"fmt"
	"github.com/felixge/felixge.de/fs"
	"net/http"
	"path/filepath"
)

func main() {
	fs := fs.New()

	if err := list(fs, "/"); err != nil {
		panic(err)
	}
}

func list(fs http.FileSystem, path string) error {
	file, err := fs.Open(path)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		fmt.Printf("file\t%s\t%d\n", path, stat.Size())
		return nil
	}

	fmt.Printf("dir\t%s\n", path)

	stats, err := file.Readdir(0)
	if err != nil {
		return err
	}

	for _, stat = range stats {
		if err := list(fs, filepath.Clean(path+"/"+stat.Name())); err != nil {
			return err
		}
	}
	return nil
}
