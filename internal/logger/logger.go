package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New() *Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{logger}
}

func (l *Logger) SetLevel(level int) {
	l.Logger = l.Logger.Level(zerolog.Level(level))
}
