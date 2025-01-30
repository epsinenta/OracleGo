package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	zerolog.Logger
}

func (cl *Logger) Log(args ...interface{}) error {
	var err interface{}
	if len(args) > 1 {
		err = args[1]
	} else {
		err = args[0]
	}
	cl.Error().Stack().Err(err.(error)).Msg("")
	return nil
}

func (cl *Logger) LogFatal(err error) {
	cl.Fatal().Stack().Err(err).Msg("")
}

func New() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := &Logger{Logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()}
	return logger
}
