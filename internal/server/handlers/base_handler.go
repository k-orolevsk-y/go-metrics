package handlers

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type (
	baseHandler struct {
		storage stor
		log     logger.Logger
	}

	stor interface {
		GetGauge(name string) (float64, error)
		SetGauge(name string, value float64)
		GetCounter(name string) (int64, error)
		AddCounter(name string, value int64)
		GetAll() []models.MetricsValue
	}
)

func NewBase(storage stor, log logger.Logger) *baseHandler {
	return &baseHandler{storage: storage, log: log}
}
