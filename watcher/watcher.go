package watcher

import (
	"path/filepath"

	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/fallback"
	"github.com/fmotalleb/watch2do/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/ryanuber/go-glob"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLog() {
	log = logger.SetupLogger("Watcher")
}

// New directory watcher
func New(notifier chan interface{}, paths ...string) {
	defer func() {
		fallback.CaptureError(log, recover())
	}()
	setupLog()
	log.WithField("paths", paths).Debugln("Started Watching")
	if len(paths) < 1 {
		log.Panicf("No path given, falling back\n")
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithField("error", err).Panicln("failed to initialize watcher")
	}
	defer w.Close()
	go coreLoop(notifier, w, paths)
	for _, p := range paths {
		err = w.Add(filepath.Dir(p))
		if err != nil {
			log.WithFields(logrus.Fields{
				"path":  p,
				"error": err,
			}).Panicln("Cannot Watch given path")
		}
	}
	<-make(chan struct{})
}

func coreLoop(notifier chan interface{}, w *fsnotify.Watcher, paths []string) {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				log.WithField("error", err).Debugln("received error at from fsnotify")
				return
			}
		case e, ok := <-w.Events:
			log.WithField("event", e).Debugln("event received")
			if !ok {
				log.WithField("event", e).Debugln("fail event received")
				return
			}
			for _, path := range paths {
				accepted := false
				for _, op := range cmd.Params.Operations {
					if op == e.Op {
						accepted = true
						break
					}
				}
				if accepted == false {
					break
				}
				if glob.Glob(path, e.Name) {
					log.WithFields(logrus.Fields{
						"matcher":    path,
						"event_data": e.Name,
					}).Debugln("matched found")
					notifier <- 0
					break
				}
			}
		}
	}
}
