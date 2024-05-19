package watcher

import (
	"os"
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
func New(notifier chan interface{}) {
	defer func() {
		fallback.CaptureError(log, recover())
	}()
	setupLog()
	log.WithFields(
		logrus.Fields{
			"is_recursive": cmd.Params.Recursive,
			"watch_list":   cmd.Params.WatchList,
			"exclude_list": cmd.Params.ExcludeWatchList,
			"matcher":      cmd.Params.MatchList,
		},
	).Infoln("initializing watcher")
	paths := cmd.Params.WatchList
	var finalWatchList []string

	if len(paths) < 1 {
		log.Panicln("No path given, falling back")
	}

	if cmd.Params.Recursive {
		for _, i := range paths {
			finalWatchList = append(finalWatchList, findSubDirs(i)...)
		}
		log.WithField("watch_list", finalWatchList).Infoln("flag recursive received, final list of directories to watch")
	} else {
		finalWatchList = paths
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithField("error", err).Panicln("failed to initialize watcher")
	}
	defer w.Close()
	go coreLoop(notifier, w, cmd.Params.MatchList)

	for _, p := range finalWatchList {
		// filepath.WalkDir(p,fs.WalkDirFunc)
		err = w.Add(p)
		if err != nil {
			log.WithFields(logrus.Fields{
				"path":  p,
				"error": err,
			}).Panicln("Cannot Watch given path")
		}
	}
	log.WithField("paths", finalWatchList).
		WithField("match_list", cmd.Params.MatchList).
		Debugln("Started Watching")
	<-make(chan struct{})
}

func findSubDirs(root string) []string {
	for _, i := range cmd.Params.ExcludeWatchList {
		if glob.Glob(i, root) {
			log.WithField("directory", root).Debugf("ignoring %s, it was listed on exclude list\n", root)
			return []string{}
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		log.WithField("directory", err).Warningln("Error reading directory")
		return []string{}
	}
	log.WithField("directory", root).Debugf("adding %s to watch list\n", root)
	result := []string{root}
	for _, entry := range entries {
		if entry.IsDir() {
			subDirPath := filepath.Join(root, entry.Name())
			result = append(result, findSubDirs(subDirPath)...)
		}
	}
	return result
}
func coreLoop(notifier chan interface{}, w *fsnotify.Watcher, paths []string) {
	for {
		select {
		case err, ok := <-w.Errors:
			log.WithField("event", err).Debugln("event received")
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
				log.WithField("event_data", e.Name).Debugln("unmatched")
			}
		}
	}
}
