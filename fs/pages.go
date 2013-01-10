package fs

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type pages struct {
	baseFs http.FileSystem
}

func newPages(baseFs http.FileSystem) http.FileSystem {
	return &pages{baseFs: baseFs}
}

func (pages *pages) Open(path string) (http.File, error) {
	file, err := pages.baseFs.Open("/pages" + path)
	if err != nil {
		return file, err
	}

	ext := filepath.Ext(path)
	if ext != ".html" {
		return file, nil
	}

	layout, err := pages.layout()
	if err != nil {
		return nil, err
	}

	pageHtml, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if _, err := layout.New("page").Parse(string(pageHtml)); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := layout.Execute(buf, nil); err != nil {
		return nil, err
	}

	result := newFile(file, path, buf.Bytes())
	return result, err
}

func (pages *pages) layout() (*template.Template, error) {
	layoutFile, err := pages.baseFs.Open("/layouts/default.html")
	if err != nil {
		return nil, err
	}
	defer layoutFile.Close()

	layoutHtml, err := ioutil.ReadAll(layoutFile)
	if err != nil {
		return nil, err
	}

	return template.New("layout").Parse(string(layoutHtml))
}
