package file

import (
	"net/http"
	"os"
)

type AbsFileSystem struct{}

var last http.File

func (fs *AbsFileSystem) Open(path string) (http.File, error) {
	if last != nil {
		last.Close()
	}
	file, err := os.Open(path)
	last = file
	return file, err
}

func NewFileServer() http.Handler {
	return http.FileServer(&AbsFileSystem{})
}
