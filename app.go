package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/notebox/nbfm/pkg/config"
	"github.com/notebox/nbfm/pkg/editor"
	"github.com/notebox/nbfm/pkg/identifier"
	"github.com/notebox/nbfm/pkg/local"
	"github.com/notebox/nbfm/pkg/local/note"
	"github.com/notebox/nbfm/pkg/logger"
	"github.com/notebox/nbfm/pkg/nav"
)

// App struct
type App struct {
	Editor *editor.Editor
	Logger *logger.NBLogger

	ctx       context.Context
	replicaID uint32
	db        *sql.DB
	wd        *nav.WorkingDir
	config    *config.Config
}

func NewApp(config *config.Config) *App {
	app := &App{config: config}
	logger, err := logger.New(filepath.Join(app.WD().Path, "nbfm.log"))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	app.Logger = logger

	wd := app.WD()
	db, err := local.ConnectDB(filepath.Join(wd.HomeDir, ".nbfm", "local", "sqlite.db"))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	err = note.Prepare(db, wd.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	app.db = db

	macAddr, err := identifier.MacAddrHex()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	hash := sha256.Sum256(append([]byte(macAddr), 0))
	app.replicaID = binary.BigEndian.Uint32(hash[:4])
	app.Editor = editor.NewEditor(config)

	return app
}

func (a *App) Startup(ctx context.Context) {
	ctx = context.WithValue(ctx, identifier.DB, a.db)
	ctx = context.WithValue(ctx, identifier.ReplicaID, a.replicaID)
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	a.Editor.Wait()
	a.ctx.Done()
	a.ctx.Value(identifier.DB).(*sql.DB).Close()
	a.Logger.Close()
}

func (a *App) WD() *nav.WorkingDir {
	if a.wd != nil {
		return a.wd
	}

	wd, err := nav.NewWorkingDir()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	a.wd = wd
	return wd
}

func (a *App) ReadDirFiles(path string, hidden bool) ([]*nav.FileInfo, error) {
	return nav.ReadDirFiles(path, hidden)
}

func (a *App) Preview(path string, hidden bool) (*nav.PreviewInfo, error) {
	a.Editor.Close()
	return nav.Preview(path, hidden)
}

func (a *App) NBError(path string, err string) {
	a.Editor.Error(path, err)
}

func (a *App) NBConnected(path string, connected bool) {
	a.Editor.Connected(a.ctx, path, connected)
}

func (a *App) NBContribute(path string, ctrbs string) {
	a.Editor.Contribute(a.ctx, path, ctrbs)
}

/** @category node manipulation */
func (a *App) DeleteFile(path string) error {
	return nav.DeleteFile(path)
}

func (a *App) AddFile(path string) error {
	return nav.AddFile(path)
}

func (a *App) MoveFile(src, dst string) error {
	return nav.MoveFile(src, dst)
}

func (a *App) CopyFile(src, dst string) error {
	return nav.CopyFile(src, dst)
}
