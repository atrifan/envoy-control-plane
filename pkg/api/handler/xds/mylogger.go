package xds

import (
	log "github.com/sirupsen/logrus"
)

type Logger struct {}

func (handler Logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
func (handler Logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

