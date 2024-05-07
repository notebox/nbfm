package pool

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/notebox/nb-crdt-go/block"
	"github.com/notebox/nb-crdt-go/common"
	local "github.com/notebox/nbfm/pkg/local/note"
	"github.com/notebox/nbfm/pkg/nav"
)

type SynchronizedBlock struct {
	sync.Mutex
	blockPath string
	fname     string
	block     *block.Block
	changed   bool
	toFile    bool
	fromFile  bool
}

func NewEmptySynchronizedBlock(blockPath string) (*SynchronizedBlock, error) {
	f, err := os.Lstat(blockPath)
	if os.IsNotExist(err) || !f.IsDir() {
		return nil, errors.New("invalid block path")
	}
	return &SynchronizedBlock{blockPath: blockPath}, nil
}

func NewUnSynchronizedBlock(blockPath string, block *block.Block) (*SynchronizedBlock, error) {
	err := os.MkdirAll(blockPath, 0777)
	if err != nil {
		return nil, err
	}
	sb, err := NewEmptySynchronizedBlock(blockPath)
	if err != nil {
		return nil, err
	}
	sb.block = block
	sb.toFile = true
	return sb, nil
}

func (sb *SynchronizedBlock) Sync(db *sql.DB, note *Note) (*block.Block, error) {
	sb.Lock()
	defer sb.Unlock()

	if !sb.fromFile {
		err := sb.write()
		if err != nil {
			return nil, err
		}
		sb.toFile = false
		return nil, nil
	}

	err := local.InsertBlock(db, &note.ID, sb.block)
	if err != nil {
		return nil, err
	}

	fname, err := latestFileName(sb.blockPath)
	if err != nil {
		return nil, err
	}

	if fname == sb.fname || fname == "" {
		if sb.changed {
			err := sb.write()
			if err != nil {
				return nil, err
			}
		}
		sb.fromFile = false
		sb.toFile = false
		sb.changed = false
		return nil, nil
	}

	block, err := blockFromFilename(sb.blockPath, fname)
	if err != nil {
		return nil, err
	}
	sb.block = block
	sb.fname = fname
	updated, err := sb.update(db, note)
	if err != nil {
		return nil, err
	}

	if sb.changed || updated {
		err := sb.write()
		if err != nil {
			return nil, err
		}
		sb.fromFile = false
		sb.toFile = false
		sb.changed = false
	}

	return sb.block, nil
}

func (sb *SynchronizedBlock) Apply(ctrb block.Contribution) error {
	sb.Lock()
	defer sb.Unlock()
	if err := sb.block.Apply(ctrb); err != nil {
		return err
	}
	sb.changed = true
	return nil
}

func (sb *SynchronizedBlock) update(db *sql.DB, note *Note) (bool, error) {
	blockNonce, textNonce := nonceOf(sb.block, note.replicaID)
	ctrbs, err := local.SelectAllAfter(db, note.replicaID, &note.ID, &sb.block.BlockID, blockNonce, textNonce)
	if err != nil {
		return false, err
	}
	for _, ctrb := range ctrbs {
		err := sb.block.Apply(*ctrb)
		if err != nil {
			return true, err
		}
	}
	return len(ctrbs) > 0, nil
}

func (sb *SynchronizedBlock) write() error {
	bytes, err := json.Marshal(sb.block)
	if err != nil {
		return err
	}
	fname := nav.NewBlockFileName()
	err = os.WriteFile(filepath.Join(sb.blockPath, fname), bytes, 0777)
	if err != nil {
		return err
	}
	oldfname := sb.fname
	sb.fname = fname
	if oldfname != "" && oldfname != fname {
		_, err := os.Lstat(filepath.Join(sb.blockPath, fname))
		if err != nil {
			return err
		}
		err = os.Remove(filepath.Join(sb.blockPath, oldfname))
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func nonceOf(block *block.Block, replicaID common.ReplicaID) (blockNonce, textNonce common.Nonce) {
	nonces, ok := block.Version[replicaID]
	if ok {
		blockNonce = nonces[0]
		textNonce = nonces[1]
	}
	return
}

func latest(blockPath string) (fname string, block *block.Block, err error) {
	fname, err = latestFileName(blockPath)
	if err != nil {
		return
	}
	block, err = blockFromFilename(blockPath, fname)
	return
}

func latestFileName(blockPath string) (fname string, err error) {
	var names []string

	names, err = nav.ReadDirFilenames(blockPath)
	if err != nil {
		return
	}
	l := len(names)
	if l == 0 {
		return
	}
	for _, name := range names {
		if name > fname {
			fname = name
		}
	}
	return
}

func blockFromFilename(blockPath, fname string) (*block.Block, error) {
	bytes, err := os.ReadFile(filepath.Join(blockPath, fname))
	if err != nil {
		return nil, err
	}

	var block *block.Block
	return block, json.Unmarshal(bytes, &block)
}
