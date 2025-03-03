package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(level int, nonConsoleWriter io.Writer) (*Logger, error) {
	var writer io.Writer = os.Stdout
	if nonConsoleWriter != nil {
		writer = zerolog.MultiLevelWriter(os.Stdout, nonConsoleWriter)
	}

	zerolog.MessageFieldName = "_msg"

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(zerolog.Level(level))
	return &Logger{logger}, nil
}
