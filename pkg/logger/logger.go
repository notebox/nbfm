package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type NBLogger struct {
	file *os.File
}

func New(path string) (*NBLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(io.MultiWriter(os.Stdout, file))
	return &NBLogger{file}, nil
}

func (l *NBLogger) Close() {
	l.file.Close()
}

func (l *NBLogger) Print(message string) {
	log.Print(message)
}

func (l *NBLogger) Trace(message string) {
	log.Trace().Msg(message)
}

// ignore
func (l *NBLogger) Debug(message string) {}

func (l *NBLogger) Info(message string) {
	log.Info().Msg(message)
}

func (l *NBLogger) Warning(message string) {
	log.Warn().Msg(message)
}

func (l *NBLogger) Error(message string) {
	log.Error().Msg(message)
}

func (l *NBLogger) Fatal(message string) {
	log.Fatal().Msg(message)
}
