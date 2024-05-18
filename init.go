package main

import (
	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/fallback"
	"github.com/fmotalleb/watch2do/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	cmd.Execute()
	log = logger.SetupLogger("Watch2Do")
	defer func() {
		fallback.CaptureError(log, recover())
	}()
	ValidateCliParams()
}
