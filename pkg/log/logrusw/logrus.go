package logrusw

import (
	"semesta-ban/pkg/log"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger so that the interface match with log.Logger
type Logger struct {
	*logrus.Logger
}

func (l Logger) Log(level log.Level, arg ...interface{}) {
	l.Logger.Log(logrus.Level(level), arg...)
}
