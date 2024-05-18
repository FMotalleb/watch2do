package debounce

import (
	"time"

	"github.com/fmotalleb/watch2do/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLog() {
	log = logger.SetupLogger("Debounce")
}

// Filter the input channel with given delay and signals output channel
func Filter(input <-chan interface{}, output chan<- interface{}, delay time.Duration) {
	setupLog()
	log.Debugf("spawning a debouncer with %v wait time", delay)
	var (
		timer *time.Timer
	)
	timerChan := make(chan time.Time)
	go func() {
		for i := range timerChan {
			log.Debugf("debounce duration passed sending a signal into output channel")
			output <- i
			timer = nil
		}
	}()
	for range input {

		log.Debugf("received event")
		if timer == nil {
			log.Debugf("no timer found creating a new one")
			timer = time.NewTimer(delay)
			drain(timer.C, timerChan)
		} else {
			log.Debugf("found an old timer resetting the timer")
			if !timer.Stop() {
				<-timerChan
			}
			timer.Reset(delay)
		}
	}
}

func drain[T any](input <-chan T, output chan<- T) {
	go func() {
		for i := range input {
			output <- i
		}
	}()
}
