//go:build linux || darwin

package nav

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"syscall"
)

func sysNames(f os.FileInfo) (username, groupname string) {
	if stat, ok := f.Sys().(*syscall.Stat_t); ok {
		uid := fmt.Sprint(stat.Uid)
		if u, err := user.LookupId(uid); err == nil {
			username = u.Username
		} else {
			username = uid
		}
		gid := fmt.Sprint(stat.Gid)
		if g, err := user.LookupGroupId(gid); err == nil {
			groupname = g.Name
		} else {
			groupname = gid
		}
	}
	return
}

func sysCopyFile(src, dst string) error {
	return exec.Command("cp", "-r", src, dst).Run()
}

func sysMoveFile(src, dst string) error {
	return exec.Command("mv", src, dst).Run()
}
