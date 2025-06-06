package logger

import (
	"log/slog"
	"marketfuck/internal/application/port"
	"os"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter() port.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	slogger := slog.New(handler)
	return &SlogAdapter{logger: slogger}
}

func (l *SlogAdapter) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args)
}

func (l *SlogAdapter) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args)
}

func (l *SlogAdapter) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args)
}

func (l *SlogAdapter) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args)
}
