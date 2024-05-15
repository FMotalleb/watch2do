package main

import (
	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/prefix"
	"github.com/sirupsen/logrus"
)

func init() {
	cmd.Execute()
	logrus.SetLevel(cmd.Params.LogLevel)
	logrus.SetFormatter(prefix.Set("Watch2Do"))
}
