package fs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type lessProcessor struct {
	baseFs http.FileSystem
}

func newLessProcessor(baseFs http.FileSystem) http.FileSystem {
	return &lessProcessor{baseFs: baseFs}
}

type lessFile struct {
	http.File
}

type lessStat struct {
	os.FileInfo
	name string
}

func (stat *lessStat) Name() string  {
	return stat.name
}

func (file *lessFile) Readdir(count int) ([]os.FileInfo, error) {
	stats, err := file.File.Readdir(count)
	if err != nil {
		return stats, err
	}

	for i, stat := range stats {
		name := stat.Name()
		ext := filepath.Ext(name)
		if ext == ".less" {
			name = name[0:len(name)-len(ext)] + ".css"
			stats[i] = &lessStat{FileInfo: stat, name: name}
		}
	}
	return stats, nil
}

func (less *lessProcessor) Open(path string) (http.File, error) {
	ext := filepath.Ext(path)
	if ext == ".less" {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	file, err := less.baseFs.Open(path)
	if ext != ".css" {
		return &lessFile{File: file}, err
	}

	lessPath := path[0:len(path)-len(ext)] + ".less"
	lessFile, lessErr := less.baseFs.Open(lessPath)
	if lessErr != nil {
		return file, err
	}

	buf := &bytes.Buffer{}
	cmd := exec.Command("lessc", "-", "--no-color")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	io.Copy(stdin, lessFile)

	if err := stdin.Close(); err != nil {
		return nil, err
	}

	io.Copy(buf, stdout)
	io.Copy(buf, stderr)

	cmd.Wait()

	result := newFile(lessFile, path, buf.Bytes())
	return result, nil
}
