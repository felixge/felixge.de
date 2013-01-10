package fs

import (
	"bytes"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
)

type lessProcessor struct {
	baseFs http.FileSystem
}

func newLessProcessor(baseFs http.FileSystem) http.FileSystem {
	return &lessProcessor{baseFs: baseFs}
}

func (less *lessProcessor) Open(path string) (http.File, error) {
	ext := filepath.Ext(path)
	file, err := less.baseFs.Open(path)
	if err == nil || ext != ".css" {
		return file, err
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
