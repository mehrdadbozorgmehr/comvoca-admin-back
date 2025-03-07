package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

var Logger zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Info(msg string, keysAndValues ...interface{}) {
	Logger.Info().Msg(joinMessage(msg, keysAndValues))
}

func Error(msg string, keysAndValues ...interface{}) {
	Logger.Error().Msg(joinMessage(msg, keysAndValues))
}

func Debug(msg string, keysAndValues ...interface{}) {
	Logger.Debug().Msg(joinMessage(msg, keysAndValues))
}

func Warn(msg string, keysAndValues ...interface{}) {
	Logger.Warn().Msg(joinMessage(msg, keysAndValues))
}

func joinMessage(msg string, keysAndValues []interface{}) string {
	s := []string{msg}
	for _, v := range keysAndValues {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return strings.Join(s, " ")
}
