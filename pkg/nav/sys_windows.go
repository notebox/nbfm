//go:build windows

package nav

import (
	"os"
	"os/exec"
)

func sysNames(_ os.FileInfo) (username, groupname string) {
	return "", ""
}

func sysCopyFile(src, dst string) error {
	return exec.Command("xcopy", src, dst).Run()
}

func sysMoveFile(src, dst string) error {
	return exec.Command("move", src, dst).Run()
}
