/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/debounce"
	"github.com/fmotalleb/watch2do/executor"
	"github.com/fmotalleb/watch2do/watcher"
	"github.com/sirupsen/logrus"
)

func fallback() {
	if err := recover(); err != nil {
		logrus.Debugln(err)
		logrus.Warnln("Stopped with a panic")
	}
}

func main() {
	defer fallback()
	notifier := make(chan interface{})
	go watcher.New(notifier, cmd.Params.WatchList...)
	filtered := make(chan interface{})
	go debounce.Filter(notifier, filtered, cmd.Params.Debounce)
	// go func() {
	logrus.Infoln("Began watching the given directories, you can use flag `--verbose (-v)` if you want to debug the program")
	for range filtered {
		logrus.Infoln("received a filtered event from debouncer calling command")
		go executor.RunCommands()
	}
}
