package pool

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/notebox/nb-crdt-go/block"
	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nbfm/pkg/combine"
	"github.com/notebox/nbfm/pkg/identifier"
	local "github.com/notebox/nbfm/pkg/local/note"
	"github.com/notebox/nbfm/pkg/nav"
)

type Note struct {
	sync.Mutex

	ID common.BlockID

	replicaID common.ReplicaID
	path      string
	blocks    Blocks
	watcher   *watcher.Watcher
	wg        *sync.WaitGroup
	dt        *combine.DebouncingThrottle
}
type Blocks map[common.BlockID]*SynchronizedBlock

func NewNote(path string, replicaID common.ReplicaID, interval time.Duration) (*Note, error) {
	noteID, err := readNoteID(path)
	if err != nil {
		return nil, err
	}
	err = nav.NewNoteBlockIfNeeded(path, noteID.String())
	if err != nil {
		return nil, err
	}
	blocks, err := readBlocks(path)
	if err != nil {
		return nil, err
	}
	w := watcher.New()
	w.FilterOps(watcher.Write, watcher.Create)
	note := Note{
		ID:        *noteID,
		replicaID: replicaID,
		path:      path,
		blocks:    blocks,
		watcher:   w,
		dt:        combine.NewDebouncingThrottle(interval),
	}
	if note.ID == uuid.Nil {
		return nil, fmt.Errorf("invalid note")
	}
	return &note, nil
}

func (note *Note) Json() ([]byte, error) {
	blocks := make([]*block.Block, 0, len(note.blocks))

	for _, mb := range note.blocks {
		blocks = append(blocks, mb.block)
	}

	return json.Marshal(struct {
		ReplicaID uint32         `json:"replicaID"`
		Blocks    []*block.Block `json:"blocks"`
	}{
		ReplicaID: note.replicaID,
		Blocks:    blocks,
	})
}

func (note *Note) Contribute(ctx context.Context, bytes []byte) error {
	note.Lock()
	defer note.Unlock()

	var ctrbs []*block.Contribution
	err := json.Unmarshal(bytes, &ctrbs)
	if err != nil {
		return err
	}
	db := ctx.Value(identifier.DB).(*sql.DB)
	blockIDs, err := local.Insert(db, &note.ID, ctrbs)
	if err != nil {
		return err
	}
	for _, ctrb := range ctrbs {
		sb, ok := note.blocks[ctrb.BlockID]
		if !ok {
			blockPath := filepath.Join(note.path, "blocks", ctrb.BlockID.String())
			err := os.MkdirAll(blockPath, 0777)
			if err != nil {
				return err
			}
			sb, err = NewEmptySynchronizedBlock(blockPath)
			if err != nil {
				return err
			}
			sb.block = ctrb.Operations.BINS
			note.blocks[ctrb.BlockID] = sb
		}
		err := sb.Apply(*ctrb)
		if err != nil {
			return err
		}
	}
	return note.flagSync(ctx, blockIDs, false)
}

func (note *Note) FlagSyncFromFS(ctx context.Context, blockID *uuid.UUID) error {
	note.Lock()
	defer note.Unlock()

	return note.flagSync(ctx, []*uuid.UUID{blockID}, true)
}

func (note *Note) flagSync(ctx context.Context, blockIDs []*uuid.UUID, isFS bool) error {
	for _, blockID := range blockIDs {
		sb, ok := note.blocks[*blockID]
		if !ok {
			var err error
			sb, err = NewEmptySynchronizedBlock(filepath.Join(note.path, "blocks", blockID.String()))
			if err != nil {
				return err
			}
			note.blocks[*blockID] = sb
		}
		if isFS && !sb.toReadFile {
			sb.toReadFile = true
		}
		sb.toSyncFile = true
	}

	note.wg.Add(note.dt.Add(func() {
		go note.sync(ctx)
	}))

	return nil
}

func (note *Note) sync(ctx context.Context) {
	db := ctx.Value(identifier.DB).(*sql.DB)
	for blockID, sb := range note.blocks {
		if !sb.toSyncFile {
			continue
		}
		block, err := sb.Sync(db, note)
		if err != nil {
			log.Fatal().Err(err).Msg("")
			return
		}

		var script string
		if block != nil {
			bytes, err := json.Marshal(block)
			if err != nil {
				log.Fatal().Err(err).Msg("")
				return
			}
			script = fmt.Sprintf("window.nbExecutor('%s', nb => nb.applyRemoteBlock(JSON.parse(`%s`)))", strings.ReplaceAll(note.path, "\\", "\\\\"), string(bytes))
		} else {
			blockNonce, _ := nonceOf(sb.block, note.replicaID)
			script = fmt.Sprintf("window.nbExecutor('%s', nb => nb.trackMergedNonce('%s', %d))", strings.ReplaceAll(note.path, "\\", "\\\\"), blockID, blockNonce-1)
		}
		runtime.WindowExecJS(ctx, script)
	}
	note.wg.Done()
}

func readNoteID(path string) (*uuid.UUID, error) {
	var meta map[string]string
	bytes, err := os.ReadFile(filepath.Join(path, "meta.json"))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &meta)
	if err != nil {
		return nil, err
	}

	noteID, err := uuid.Parse(meta["id"])
	if err != nil {
		return nil, err
	}

	return &noteID, nil
}

func readBlocks(path string) (Blocks, error) {
	blocksPath := filepath.Join(path, "blocks")
	names, err := nav.ReadDirFilenames(blocksPath)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var mutex sync.Mutex
	blocks := make(Blocks)
	for _, bid := range names {
		wg.Add(1)
		go func(bid string) {
			defer wg.Done()
			blockPath := filepath.Join(blocksPath, bid)
			sb, err := NewEmptySynchronizedBlock(blockPath)
			if err != nil {
				return
			}
			fname, block, err := latest(sb.blockPath)
			if err != nil {
				return
			}
			sb.fname = fname
			sb.block = block
			mutex.Lock()
			blocks[sb.block.BlockID] = sb
			mutex.Unlock()
		}(bid)
	}
	wg.Wait()
	return blocks, nil
}
