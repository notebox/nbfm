package pool

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"

	"github.com/notebox/nb-crdt-go/common"
	"github.com/notebox/nbfm/pkg/config"
	"github.com/notebox/nbfm/pkg/identifier"
	local "github.com/notebox/nbfm/pkg/local/note"
)

type Path = string
type Pool struct {
	sync.Mutex

	opened       *Note
	notes        map[Path]*Note
	syncInterval time.Duration
}

func NewPool(config *config.Config) *Pool {
	return &Pool{
		syncInterval: config.SyncInterval,
		notes:        make(map[Path]*Note),
	}
}

func (pool *Pool) Open(ctx context.Context, path string) (*Note, error) {
	pool.Lock()
	defer pool.Unlock()

	var err error
	note, ok := pool.notes[path]
	if ok {
		note.wg.Add(1)
	} else {
		db := ctx.Value(identifier.DB).(*sql.DB)
		replicaID := ctx.Value(identifier.ReplicaID).(common.ReplicaID)
		note, err = NewNote(db, replicaID, path, pool.syncInterval)
		if err != nil {
			return nil, err
		}
		pool.notes[path] = note

		note.wg = new(sync.WaitGroup)
		note.wg.Add(1)

		cached, err := local.SelectBlocks(db, &note.ID)
		if err != nil {
			return nil, err
		}
		blocksPath := filepath.Join(path, "blocks")
		for _, b := range cached {
			sb, ok := note.blocks[b.BlockID]
			if ok {
				sb.toFile, err = sb.update(db, note)
				if err != nil {
					return nil, err
				}
				continue
			}
			sb, err = NewUnSynchronizedBlock(filepath.Join(blocksPath, b.BlockID.String()), b)
			if err != nil {
				return nil, err
			}
			note.blocks[b.BlockID] = sb
		}

		note.wg.Add(1)
		note.sync(ctx)

		go func() {
			note.wg.Wait()
			pool.Lock()
			note.watcher.Close()
			delete(pool.notes, note.path)
			defer pool.Unlock()
		}()

		go pool.watch(ctx, note)
		err = note.watcher.AddRecursive(note.path)
		if err != nil {
			return nil, err
		}

		go func() {
			if err := note.watcher.Start(time.Second); err != nil {
				log.Fatal().Err(err).Msg("")
			}
		}()
	}

	pool.opened = note
	return note, nil
}

func (pool *Pool) Close() error {
	pool.Lock()
	defer pool.Unlock()
	if pool.opened != nil {
		pool.opened.wg.Done()
		pool.opened = nil
	}
	return nil
}

func (pool *Pool) Wait() {
	pool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), pool.syncInterval*3/2)

	go func() {
		var wg sync.WaitGroup
		for _, note := range pool.notes {
			wg.Add(1)
			go func(n *Note) {
				n.dt.Clear()
				n.wg.Wait()
				wg.Done()
			}(note)
		}
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
}

func (pool *Pool) watch(ctx context.Context, note *Note) {
	for {
		select {
		case event := <-note.watcher.Event:
			blockID := updatedBlockIDFromPath(note.path, event.Path)
			if blockID != nil {
				note.FlagSyncFromFS(ctx, blockID)
			}
		case err := <-note.watcher.Error:
			if err == watcher.ErrWatchedFileDeleted {
				n := pool.notes[note.path]
				n.dt.Stop()
				n.wg.Done()
				return
			}
			log.Fatal().Err(err).Msg("")
		case <-note.watcher.Closed:
			return
		}
	}
}

func (pool *Pool) Contribute(ctx context.Context, notePath string, bytes []byte) error {
	pool.Lock()
	note, ok := pool.notes[notePath]
	pool.Unlock()
	if !ok {
		return errors.New("note not found")
	}
	return note.Contribute(ctx, bytes)
}

func updatedBlockIDFromPath(notePath, path string) *common.BlockID {
	sub, _ := strings.CutPrefix(path, notePath)
	frags := strings.Split(sub, config.Separator)
	l := len(frags)
	if l < 3 || l > 4 {
		return nil
	}
	if frags[1] != "blocks" {
		return nil
	}
	blockID, err := uuid.Parse(frags[2])
	if err != nil {
		return nil
	}
	return &blockID
}
