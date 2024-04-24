package nav

import (
	"fmt"
	"io/fs"
)

type FileInfo struct {
	Name    string `json:"name"`
	RawSize int64  `json:"rawSize"`
	IsDir   bool   `json:"isDir"`

	Mode      string `json:"mode"`
	Username  string `json:"username"`
	GroupName string `json:"groupName"`
	Size      string `json:"size"`
	ModTime   string `json:"modTime"`
}

func NewFileInfo(name string, f fs.FileInfo) *FileInfo {
	username, groupName := sysNames(f)
	return &FileInfo{
		Name:    name,
		RawSize: f.Size(),
		IsDir:   f.IsDir(),

		Mode:      f.Mode().String(),
		Username:  username,
		GroupName: groupName,
		Size:      readableSize(f.Size()),
		ModTime:   f.ModTime().Format("2006-01-02 15:04:05Z07:00"),
	}
}

func readableSize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / 1000
	for _, s := range suffix {
		if curr < 10 {
			return fmt.Sprintf("%.1f%s", curr-0.0499, s)
		} else if curr < 1000 {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= 1000
	}

	return ""
}
