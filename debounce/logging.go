package debounce

import (
	"os"

	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/prefix"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLog() {
	log = logrus.New()
	log.SetLevel(cmd.Params.LogLevel)
	log.Formatter = prefix.Set("Debounce")
	log.Out = os.Stdout
}
