package infrastructure

import (
	log "github.com/sirupsen/logrus"
)

type Log interface {
	Info(mess string)
	Error(mess string)
}

type Logger struct{}

func (logger Logger) Info(mess string) {
	log.Info(mess)
}

func (logger Logger) Error(mess string) {
	log.Error(mess)
}
