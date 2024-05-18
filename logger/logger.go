package logger

import (
	"os"

	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/logger/scoped"
	"github.com/sirupsen/logrus"
)

func SetupLogger(scope string) *logrus.Logger {
	parentFormatter := []logrus.Formatter{}
	if cmd.Params.JsonOutput {
		parentFormatter = append(parentFormatter, &logrus.JSONFormatter{
			// PrettyPrint: true,
		})
	}
	log := logrus.New()
	log.SetLevel(cmd.Params.LogLevel)
	log.SetFormatter(scoped.New(scope, parentFormatter...))
	log.SetOutput(os.Stdout)
	return log
}
