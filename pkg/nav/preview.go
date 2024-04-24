package nav

import (
	"os"
	"strings"

	"github.com/h2non/filetype"
)

type PreviewType string

const dirType PreviewType = "dir"
const textType PreviewType = "text"
const imageType PreviewType = "image"
const audioType PreviewType = "audio"
const videoType PreviewType = "video"
const embedType PreviewType = "embed"
const unknownType PreviewType = ""

type PreviewInfo struct {
	DirFiles []*FileInfo `json:"dirFiles,omitempty"`
	UTF8     *string     `json:"utf8,omitempty"`
	Type     PreviewType `json:"type"`
}

func Preview(path string, hidden bool) (*PreviewInfo, error) {
	f, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		files, err := ReadDirFiles(path, hidden)
		if err != nil {
			return nil, err
		}
		return &PreviewInfo{DirFiles: files, Type: dirType}, err
	}

	var pt PreviewType
	file, err := os.Open(path)
	if err == nil {
		head := make([]byte, 261)
		file.Read(head)
		file.Close()

		kind, err := filetype.Match(head)
		if err == nil {
			pt = previewType(kind.MIME.Value)
		}
	}

	if pt == unknownType {
		if strings.HasSuffix(path, ".csv") {
			pt = textType
		} else if strings.HasSuffix(path, ".svg") {
			pt = imageType
		}
	}

	if (pt == textType || pt == unknownType) && f.Size() < 1024*1024*1 {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		utf8 := string(bytes)
		return &PreviewInfo{UTF8: &utf8, Type: textType}, err
	}

	return &PreviewInfo{Type: pt}, err
}

func previewType(mime string) PreviewType {
	if strings.HasSuffix(mime, "pdf") {
		return embedType
	}
	if len(mime) < 5 {
		return unknownType
	}
	switch mime[:5] {
	case "text/":
		return textType
	case "image":
		return imageType
	case "audio":
		return audioType
	case "video":
		return videoType
	default:
		return unknownType
	}
}
