package nav

import (
	"os"
	"path/filepath"
	"strings"
)

func NewWorkingDir() (*WorkingDir, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if strings.Contains(homeDir, wd) {
		wd = filepath.Join(homeDir, "nbfm")
		if _, err = os.Stat(wd); os.IsNotExist(err) {
			err := os.Mkdir(wd, 0777)
			if err != nil {
				return nil, err
			}
		}
	}

	return &WorkingDir{
		Separator: string(filepath.Separator),
		Path:      wd,
		HomeDir:   homeDir,
	}, nil
}

type WorkingDir struct {
	Separator string
	Path      string
	HomeDir   string
}
