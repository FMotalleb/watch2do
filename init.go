package main

import (
	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/fallback"
	"github.com/fmotalleb/watch2do/prefix"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	cmd.Execute()
	log = logrus.New()
	log.SetLevel(cmd.Params.LogLevel)
	log.SetFormatter(prefix.Set("Watch2Do"))
	defer func() {
		fallback.CaptureError(log, recover())
	}()
	ValidateCliParams()
}
