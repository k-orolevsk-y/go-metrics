package logger

import "go.uber.org/zap"

func New() (Logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}

type Logger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Sync() error
}
