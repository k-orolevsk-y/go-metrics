package gopsutil

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
)

type GopsutilMetricsCollector struct {
	total                 metrics.Metric
	free                  metrics.Metric
	cpuUtilizationMetrics []metrics.Metric
}

func NewGopsutilCollector() *GopsutilMetricsCollector {
	return &GopsutilMetricsCollector{}
}

func (c *GopsutilMetricsCollector) Collect() error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	cpuUtilizationMetrics, err := cpu.Percent(time.Millisecond*100, true)
	if err != nil {
		return err
	}

	c.total = metrics.NewMetric("TotalMemory", metrics.GaugeType, 0, float64(memory.Total))
	c.free = metrics.NewMetric("FreeMemory", metrics.GaugeType, 0, float64(memory.Total))

	c.cpuUtilizationMetrics = make([]metrics.Metric, 0)
	for i, cpuUtilizationMetric := range cpuUtilizationMetrics {
		c.cpuUtilizationMetrics = append(
			c.cpuUtilizationMetrics,
			metrics.NewMetric(fmt.Sprintf("CPUUtilization%d", i+1), metrics.GaugeType, 0, cpuUtilizationMetric),
		)
	}

	return nil
}

func (c *GopsutilMetricsCollector) GetResults() (results []metrics.Metric) {
	results = append(results, c.total)
	results = append(results, c.free)
	results = append(results, c.cpuUtilizationMetrics...)

	return
}
