package middlewares

type (
	baseMiddleware struct {
		log logger
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}
)

func NewBase(log logger) *baseMiddleware {
	return &baseMiddleware{
		log: log,
	}
}
