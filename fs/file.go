package fs

import (
	"bytes"
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

func (file *mFile) Stat() (os.FileInfo, error) {
	stat, err := file.File.Stat()
	if err != nil {
		return stat, err
	}

	size := int64(len(file.data))
	return &mStat{FileInfo: stat, size: size}, nil
}

func newFile(file http.File, path string, data []byte) http.File {
	return &mFile{
		File: file,
		name: path,
		data: data,
	}
}
