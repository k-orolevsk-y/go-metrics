package alternative

import (
	"math/rand"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
)

type AlternativeMetricsCollector struct {
	pollCount   metrics.Metric
	randomValue metrics.Metric
}

func NewAlternativeCollector() *AlternativeMetricsCollector {
	return &AlternativeMetricsCollector{}
}

func (c *AlternativeMetricsCollector) Collect() error {
	var currentPollCount int64
	if !c.pollCount.IsNil() {
		currentPollCount = *c.pollCount.Delta
	}

	c.pollCount = metrics.NewMetric("PollCount", metrics.CounterType, currentPollCount+1, 0)
	c.randomValue = metrics.NewMetric("RandomValue", metrics.GaugeType, 0, rand.Float64())

	return nil
}

func (c *AlternativeMetricsCollector) GetResults() []metrics.Metric {
	return []metrics.Metric{
		c.pollCount,
		c.randomValue,
	}
}
