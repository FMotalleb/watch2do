package watcher

import (
	"github.com/fmotalleb/watch2do/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLog() {
	log = logger.SetupLogger("Watcher ")
}
