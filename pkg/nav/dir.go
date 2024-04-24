package nav

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func ReadDirFiles(path string, hidden bool) ([]*FileInfo, error) {
	files, err := readFiles(path, hidden)
	if err != nil {
		if os.IsPermission(err) {
			return nil, os.ErrPermission
		}
		return nil, err
	}
	return files, nil
}

func ReadDirFilenames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	return names, err
}

func readFiles(path string, hidden bool) ([]*FileInfo, error) {
	filenames, err := ReadDirFilenames(path)

	files := make([]*FileInfo, 0, len(filenames))
	for _, fname := range filenames {
		fpath := filepath.Join(path, fname)

		f, err := os.Lstat(fpath)

		if os.IsNotExist(err) {
			continue
		}
		if !hidden && strings.HasPrefix(fname, ".") {
			continue
		}
		if err != nil {
			files = append(files, &FileInfo{Name: fname})
			continue
		}

		files = append(files, NewFileInfo(fname, f))
	}

	slices.SortFunc(files, func(a, b *FileInfo) int {
		if a.IsDir == b.IsDir {
			return strings.Compare(a.Name, b.Name)
		}
		if a.IsDir {
			return -1
		}
		return 1
	})

	return files, err
}
