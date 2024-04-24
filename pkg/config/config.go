package config

import (
	"path/filepath"
	"time"
)

type Config struct {
	SyncInterval time.Duration
}

var Separator = string(filepath.Separator)
