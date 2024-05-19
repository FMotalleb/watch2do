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
	validateCliParams()
}

func validateCliParams() {
	log := log.WithField("params", cmd.Params)
	if len(cmd.Params.Commands) < 1 {
		log.Panicln("commands must have at least a single value.")
	}
	if len(cmd.Params.WatchList) < 1 {
		log.Panicln("commands must have at least a single value.")
	}
	if len(cmd.Params.Operations) < 1 {
		log.Panicln("no operation to act on, stopping.")
	}
	if len(cmd.Params.MatchList) < 1 {
		log.Panicln("no matchList to match against, stopping.")
	}
}
