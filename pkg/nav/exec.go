package nav

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/notebox/nbfm/pkg/config"
	local "github.com/notebox/nbfm/pkg/local/note"
)

func DeleteFile(path string) error {
	return os.RemoveAll(path)
}

func AddFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	if strings.HasSuffix(path, ".note") {
		noteID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		noteIDStr := noteID.String()
		err = NewNoteBlockIfNeeded(nil, path, noteIDStr)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(path, "meta.json"), []byte(fmt.Sprintf(`{"id":"%s"}`, noteIDStr)), 0777)
		if err != nil {
			return err
		}
		return nil
	}

	if strings.HasSuffix(path, config.Separator) {
		return os.Mkdir(path, 0777)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func CopyFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0777)
	if err != nil {
		return err
	}
	return sysCopyFile(src, dst)
}

func MoveFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0777)
	if err != nil {
		return err
	}
	return sysMoveFile(src, dst)
}

func NewNoteBlockIfNeeded(db *sql.DB, path, noteIDStr string) error {
	noteBlockPath := filepath.Join(path, "blocks", noteIDStr, NewBlockFileName())
	dirPath := filepath.Dir(noteBlockPath)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(noteBlockPath), 0777)
	if err != nil {
		return err
	}

	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		return err
	}
	var data []byte
	if db != nil {
		data, err = local.SelectBlockData(db, &noteID, &noteID)
		if err != nil {
			return err
		}
	}
	if data == nil {
		data = []byte(fmt.Sprintf(`["%s",{},[[0,0,1]],{"TYPE":[null,"NOTE"]},false,[]]`, noteIDStr))
	}

	return os.WriteFile(noteBlockPath, data, 0777)
}

func NewBlockFileName() string {
	return time.Now().Format("060102150405") + ".json"
}
