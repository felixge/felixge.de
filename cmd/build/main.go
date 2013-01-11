package main

import (
	"fmt"
	"github.com/felixge/felixge.de/fs"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	outputDir := cwd + "/build"

	fs := fs.New()
	if err := list(fs, "/", outputDir); err != nil {
		panic(err)
	}
}

func list(fs http.FileSystem, path string, outputDir string) error {
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
		outputFile, err := os.OpenFile(outputDir +path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		if _, err := io.Copy(outputFile, file); err != nil {
			return err
		}

		return nil
	}

	if err := os.MkdirAll(outputDir + path, 0777); err != nil {
		return err
	}

	fmt.Printf("dir\t%s\n", path)

	stats, err := file.Readdir(0)
	if err != nil {
		return err
	}

	for _, stat = range stats {
		if err := list(fs, filepath.Clean(path+"/"+stat.Name()), outputDir); err != nil {
			return err
		}
	}
	return nil
}
