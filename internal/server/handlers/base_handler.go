package handlers

import "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"

type baseHandler struct {
	storage stor
}

func NewBase(storage stor) *baseHandler {
	return &baseHandler{storage: storage}
}

type stor interface {
	GetGauge(name string) (float64, error)
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, error)
	AddCounter(name string, value int64)
	GetAll() []storage.Value
}
