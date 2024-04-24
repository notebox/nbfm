package editor

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/notebox/nbfm/pkg/config"
	"github.com/notebox/nbfm/pkg/editor/pool"
)

type Editor struct {
	pool *pool.Pool
}

func NewEditor(config *config.Config) *Editor {
	return &Editor{
		pool: pool.NewPool(config),
	}
}

func (editor *Editor) Error(path string, err string) {
	log.Error().Msg(err)
}

func (editor *Editor) Connected(ctx context.Context, path string, connected bool) {
	go func() {
		note, err := editor.pool.Open(ctx, path)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		bytes, err := note.Json()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		runtime.WindowExecJS(
			ctx,
			fmt.Sprintf("window.nbExecutor('%s', nb => nb.init(%s))", strings.ReplaceAll(path, "\\", "\\\\"), string(bytes)),
		)
	}()
}

func (editor *Editor) Contribute(ctx context.Context, path string, ctrbs string) {
	err := editor.pool.Contribute(ctx, path, []byte(ctrbs))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func (editor *Editor) Close() {
	err := editor.pool.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func (editor *Editor) Wait() {
	editor.pool.Wait()
}
