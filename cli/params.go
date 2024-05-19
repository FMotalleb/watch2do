package cli

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// Params parsed using viper
type Params struct {
	// Shell that commands should run on default to cmd in windows and sh in linux
	Shell string
	// List of directory to watch
	WatchList []string
	// List of directory to ignore
	ExcludeWatchList []string
	// Will Match event operand against given match list
	MatchList []string
	// Commands's first value is always a program and the rest values are its params
	Commands []string
	// Debounce is the amount time to wait before calling the `Command`
	//   and if the program receive any other events from the given `WatchList` current action will be discarded and rescheduled
	Debounce time.Duration
	// LogLevel of the application
	LogLevel logrus.Level
	// Accepted operations to act on
	Operations []fsnotify.Op
	// Use Json logger instead of text/ansi logger
	JsonOutput bool
	// Recursively watch all subdirectories
	Recursive bool
	// Kill old processes from last invocations
	KillBeforeExecute bool
}
