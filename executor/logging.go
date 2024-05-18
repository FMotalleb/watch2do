package executor

import (
	"github.com/sirupsen/logrus"

	"github.com/fmotalleb/watch2do/logger"
)

var log *logrus.Logger

func setupLog() {
	log = logger.SetupLogger("Executor")
}
