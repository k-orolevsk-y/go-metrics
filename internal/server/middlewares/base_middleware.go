package middlewares

import "github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"

type (
	baseMiddleware struct {
		log logger.Logger
	}
)

func NewBase(log logger.Logger) *baseMiddleware {
	return &baseMiddleware{
		log: log,
	}
}
