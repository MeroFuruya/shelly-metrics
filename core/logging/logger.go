package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetLogger(name string) zerolog.Logger {
	return log.With().Str("module", name).Logger()
}

func SetupLogger(format string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if format == "pretty" {
		consoleLogger := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(consoleLogger)
	}
}
