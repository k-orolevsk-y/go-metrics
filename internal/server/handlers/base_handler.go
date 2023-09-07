package handlers

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
)

type baseHandler struct {
	storage stor.Storage
}

func NewBase(storage stor.Storage) *baseHandler {
	return &baseHandler{storage: storage}
}
