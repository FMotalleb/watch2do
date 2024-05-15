package watcher

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/ryanuber/go-glob"
)

// New directory watcher
func New(notifier chan interface{}, paths ...string) {
	setupLog()
	log.Debugf("Watching over %q\n", paths)
	if len(paths) < 1 {
		log.Panicf("No path given, falling back\n")
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicf("cannot init new watcher: %s\n", err)
	}
	defer w.Close()
	go coreLoop(notifier, w, paths)
	for _, p := range paths {
		err = w.Add(filepath.Dir(p))
		if err != nil {
			log.Warnf("Cannot watch %q: %s\n", p, err)
		}
	}
	<-make(chan struct{})
}

func coreLoop(notifier chan interface{}, w *fsnotify.Watcher, paths []string) {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				log.Debugf("received error at from fsnotify: %v", err)
				return
			}
		case e, ok := <-w.Events:
			log.Debugf("event received: %v\n", e)
			if !ok {
				log.Debugf("a not ok event received for: %v\n", e)
				return
			}
			for _, path := range paths {
				if glob.Glob(path, e.Name) {
					log.Debugf("matched %v using %v globe\n", e.Name, path)
					notifier <- 0
					break
				} else if strings.HasPrefix(e.Name, path) {
					log.Debugf("matched %v using %v by recursive rule\n", e.Name, path)
					notifier <- 0
					break
				}
			}
		}
	}
}
