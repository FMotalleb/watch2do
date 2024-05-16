/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/debounce"
	"github.com/fmotalleb/watch2do/executor"
	"github.com/fmotalleb/watch2do/watcher"
)

func main() {
	notifier := make(chan interface{})
	go watcher.New(notifier, cmd.Params.WatchList...)
	filtered := make(chan interface{})
	go debounce.Filter(notifier, filtered, cmd.Params.Debounce)

	log.Infoln("Began watching the given directories, you can use flag `--verbose (-v)` if you want to debug the program")
	for range filtered {
		log.Infoln("received a filtered event from debouncer calling command")
		go executor.RunCommands()
	}
}
