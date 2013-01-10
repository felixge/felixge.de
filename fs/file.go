package fs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// BUG(felixge) should be named "file", but getting compiler error
type mFile struct {
	http.File
	name string
	data []byte
	buf  *bytes.Buffer
}

type mStat struct {
	os.FileInfo
	size int64
	name string
}

func (stat *mStat) Name() string {
	return stat.name
}

func (stat *mStat) Size() int64 {
	return stat.size
}

func (file *mFile) Read(buf []byte) (int, error) {
	if file.buf == nil {
		file.buf = bytes.NewBuffer(file.data)
	}
	return file.buf.Read(buf)
}

func (file *mFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case os.SEEK_SET:
		file.buf = bytes.NewBuffer(file.data)
		n, err := io.CopyN(ioutil.Discard, file.buf, offset)
		return int64(n), err
	}

	return 0, fmt.Errorf("not implemented")
}

func (file *mFile) Stat() (os.FileInfo, error) {
	stat, err := file.File.Stat()
	if err != nil {
		return stat, err
	}

	size := int64(len(file.data))
	return &mStat{
		FileInfo: stat,
		size: size,
		name: file.name,
	}, nil
}

func newFile(file http.File, path string, data []byte) http.File {
	return &mFile{
		File: file,
		name: path,
		data: data,
	}
}
