package prefix

import (
	"github.com/sirupsen/logrus"
)

var termFormatter *logrus.TextFormatter = &logrus.TextFormatter{
	DisableColors: false,
	FullTimestamp: true,
	ForceColors:   true,
}

func Set(prefix string) *PrefixFormatter {
	return &PrefixFormatter{
		prefix: prefix,
	}
}

type PrefixFormatter struct {
	prefix string
}

func (f *PrefixFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = "[" + f.prefix + "] " + entry.Message
	return termFormatter.Format(entry)
}
