package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
	"os"
)

func NewEventLogger() fxevent.Logger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return &fxevent.ConsoleLogger{W: log.Logger}
}
