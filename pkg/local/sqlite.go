package local

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func ConnectDB(path string) (*sql.DB, error) {
	err := prepare(path)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	return db, err
}

func prepare(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}
