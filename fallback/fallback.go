package fallback

import (
	"os"

	"github.com/sirupsen/logrus"
)

func CaptureError(log *logrus.Logger, recovered any) {
	if recovered != nil {
		log.Debugf("%#v\n", recovered)
		log.Warnln("Stopped with a panic")
		os.Exit(1)
	}
}
