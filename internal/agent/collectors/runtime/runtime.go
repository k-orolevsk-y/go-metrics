package runtime

import (
	"runtime"
	"sync"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
)

type RuntimeMetricsCollector struct {
	mx      sync.Mutex
	metrics []metrics.Metric
}

func NewRuntimeCollector() *RuntimeMetricsCollector {
	return &RuntimeMetricsCollector{}
}

func (c *RuntimeMetricsCollector) Collect() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.setMetric(metrics.NewMetric("Alloc", metrics.GaugeType, 0, float64(memStats.Alloc)))
	c.setMetric(metrics.NewMetric("BuckHashSys", metrics.GaugeType, 0, float64(memStats.BuckHashSys)))
	c.setMetric(metrics.NewMetric("Frees", metrics.GaugeType, 0, float64(memStats.Frees)))
	c.setMetric(metrics.NewMetric("GCCPUFraction", metrics.GaugeType, 0, memStats.GCCPUFraction))
	c.setMetric(metrics.NewMetric("GCSys", metrics.GaugeType, 0, float64(memStats.GCSys)))
	c.setMetric(metrics.NewMetric("HeapAlloc", metrics.GaugeType, 0, float64(memStats.HeapAlloc)))
	c.setMetric(metrics.NewMetric("HeapIdle", metrics.GaugeType, 0, float64(memStats.HeapIdle)))
	c.setMetric(metrics.NewMetric("HeapInuse", metrics.GaugeType, 0, float64(memStats.HeapInuse)))
	c.setMetric(metrics.NewMetric("HeapObjects", metrics.GaugeType, 0, float64(memStats.HeapObjects)))
	c.setMetric(metrics.NewMetric("HeapReleased", metrics.GaugeType, 0, float64(memStats.HeapReleased)))
	c.setMetric(metrics.NewMetric("HeapSys", metrics.GaugeType, 0, float64(memStats.HeapSys)))
	c.setMetric(metrics.NewMetric("LastGC", metrics.GaugeType, 0, float64(memStats.LastGC)))
	c.setMetric(metrics.NewMetric("Lookups", metrics.GaugeType, 0, float64(memStats.Lookups)))
	c.setMetric(metrics.NewMetric("MCacheInuse", metrics.GaugeType, 0, float64(memStats.MCacheInuse)))
	c.setMetric(metrics.NewMetric("MCacheSys", metrics.GaugeType, 0, float64(memStats.MCacheSys)))
	c.setMetric(metrics.NewMetric("MSpanInuse", metrics.GaugeType, 0, float64(memStats.MSpanInuse)))
	c.setMetric(metrics.NewMetric("MSpanSys", metrics.GaugeType, 0, float64(memStats.MSpanSys)))
	c.setMetric(metrics.NewMetric("Mallocs", metrics.GaugeType, 0, float64(memStats.Mallocs)))
	c.setMetric(metrics.NewMetric("NextGC", metrics.GaugeType, 0, float64(memStats.NextGC)))
	c.setMetric(metrics.NewMetric("NumForcedGC", metrics.GaugeType, 0, float64(memStats.NumForcedGC)))
	c.setMetric(metrics.NewMetric("NumGC", metrics.GaugeType, 0, float64(memStats.NumGC)))
	c.setMetric(metrics.NewMetric("OtherSys", metrics.GaugeType, 0, float64(memStats.OtherSys)))
	c.setMetric(metrics.NewMetric("PauseTotalNs", metrics.GaugeType, 0, float64(memStats.PauseTotalNs)))
	c.setMetric(metrics.NewMetric("StackInuse", metrics.GaugeType, 0, float64(memStats.StackInuse)))
	c.setMetric(metrics.NewMetric("StackSys", metrics.GaugeType, 0, float64(memStats.StackSys)))
	c.setMetric(metrics.NewMetric("Sys", metrics.GaugeType, 0, float64(memStats.Sys)))
	c.setMetric(metrics.NewMetric("TotalAlloc", metrics.GaugeType, 0, float64(memStats.TotalAlloc)))

	return nil
}

func (c *RuntimeMetricsCollector) setMetric(metric metrics.Metric) {
	c.mx.Lock()
	defer c.mx.Unlock()

	for key, m := range c.metrics {
		if m.ID == metric.ID {
			c.metrics[key] = metric
			return
		}
	}

	c.metrics = append(c.metrics, metric)
}

func (c *RuntimeMetricsCollector) GetResults() (results []metrics.Metric) {
	c.mx.Lock()
	defer c.mx.Unlock()

	return c.metrics
}
