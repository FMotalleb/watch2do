package logger

import (
	"os"

	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/logger/prefix"
	"github.com/sirupsen/logrus"
)

func SetupLogger(scope string) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(cmd.Params.LogLevel)
	log.SetFormatter(prefix.Set(scope))
	log.SetOutput(os.Stdout)
	return log
}
