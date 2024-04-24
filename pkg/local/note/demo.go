package note

import (
	_ "embed"
	"encoding/json"
	"time"

	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/notebox/nb-crdt-go/block"
)

//go:embed demo.note.json
var demo []byte

const EMBED_SEPARATOR = string(filepath.Separator)

func addDemoNote(path string) error {
	notePath := filepath.Join(path, "demo.note")
	_, err := os.Lstat(notePath)
	if !os.IsNotExist(err) {
		return err
	}

	blocksPath := filepath.Join(notePath, "blocks")
	err = os.MkdirAll(blocksPath, 0777)
	if err != nil {
		return err
	}

	var note struct {
		ID     *uuid.UUID        `json:"id"`
		Blocks []json.RawMessage `json:"blocks"`
	}
	err = json.Unmarshal(demo, &note)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(notePath, "meta.json"), []byte(fmt.Sprintf(`{"id":"%s"}`, note.ID)), 0777)
	if err != nil {
		return err
	}

	fname := time.Now().Format("060102150405") + ".json"
	for _, bl := range note.Blocks {
		var b block.Block
		err := json.Unmarshal(bl, &b)
		if err != nil {
			return err
		}

		fpath := filepath.Join(blocksPath, b.BlockID.String(), fname)
		err = os.MkdirAll(filepath.Dir(fpath), 0777)
		if err != nil {
			return err
		}
		err = os.WriteFile(fpath, bl, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
