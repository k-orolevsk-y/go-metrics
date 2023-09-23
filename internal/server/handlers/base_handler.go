package handlers

import "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"

type (
	baseHandler struct {
		storage stor
		log     logger
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}

	stor interface {
		GetGauge(name string) (float64, error)
		SetGauge(name string, value float64)
		GetCounter(name string) (int64, error)
		AddCounter(name string, value int64)
		GetAll() []storage.Value
	}
)

func NewBase(storage stor, log logger) *baseHandler {
	return &baseHandler{storage: storage, log: log}
}
