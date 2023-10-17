package logger

import "go.uber.org/zap"

func New() (Logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	l.Core().Enabled(zap.DebugLevel)

	return l.Sugar(), nil
}

type Logger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Sync() error
}
