package cli

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Params parsed using viper
type Params struct {
	// Shell that commands should run on default to cmd in windows and sh in linux
	Shell string
	// WatchList contains list of glob's to watch over
	WatchList []string
	// Commands's first value is always a program and the rest values are its params
	Commands []string
	// Debounce is the amount time to wait before calling the `Command`
	//   and if the program receive any other events from the given `WatchList` current action will be discarded and rescheduled
	Debounce time.Duration
	// LogLevel of the application
	LogLevel logrus.Level
}
