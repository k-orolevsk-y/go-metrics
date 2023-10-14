package models

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Storage interface {
	SetGauge(string, float64) error
	AddCounter(string, int64) error

	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)

	GetAll() ([]MetricsValue, error)

	GetMiddleware() gin.HandlerFunc
	Ping(context.Context) error
	Close() error
}
