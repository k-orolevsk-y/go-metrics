package collectors

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func TestCollector(t *testing.T) {
	config.Config.RateLimit = 3
	config.Config.PollInterval = 2

	log, err := logger.New()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	c := NewCollector(log)
	go c.Run(ctx)

	var countMetrics int

outerLoop:
	for {
		select {
		case <-ctx.Done():
			break outerLoop
		default:
			countMetrics = len(c.GetMetrics())
			if countMetrics > 0 {
				break outerLoop
			}
			time.Sleep(time.Second)
		}
	}

	assert.Greater(t, countMetrics, 0)
}
