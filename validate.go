package main

import (
	"github.com/fmotalleb/watch2do/cmd"
)

func ValidateCliParams() {
	if len(cmd.Params.Commands) < 1 {
		log.Panicln("commands must have at least a single value")
	}
	if len(cmd.Params.WatchList) < 1 {
		log.Panicln("commands must have at least a single value")
	}
}
