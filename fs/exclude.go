package fs

import (
	"fmt"
	"path/filepath"
	"net/http"
	"os"
)

type excludeFs struct {
	baseFs http.FileSystem
	pattern string
}

type excludeFile struct {
	http.File
	pattern string
}

func (file *excludeFile) Readdir(count int) ([]os.FileInfo, error) {
	stats, err := file.File.Readdir(count)
	if err != nil {
		return stats, err
	}

	var results []os.FileInfo

	for _, stat := range stats {
		name := stat.Name()
		if match, err := filepath.Match(file.pattern, name); err != nil {
			continue
		} else if match {
			continue
		}

		results = append(results, stat)
	}
	return results, nil
}

func newExclude(baseFs http.FileSystem, pattern string) http.FileSystem {
	return &excludeFs{baseFs: baseFs, pattern: pattern}
}

func (e *excludeFs) Open(path string) (http.File, error) {
	if match, err := filepath.Match(e.pattern, path); err != nil {
		return nil, err
	} else if match {
		return nil, fmt.Errorf("excluded file: %s", path)
	}

	file, err := e.baseFs.Open(path)
	if err != nil {
		return file, err
	}
	return &excludeFile{File: file, pattern: e.pattern}, nil
}
