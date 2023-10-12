package handlers

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type (
	baseHandler struct {
		storage models.Storage
		log     logger.Logger
	}
)

func NewBase(storage models.Storage, log logger.Logger) *baseHandler {
	return &baseHandler{storage: storage, log: log}
}
