package metrics

import (
	"errors"
	"math/rand"
	"runtime"
	"sync"
)

type (
	MetricType string
	Metric     struct {
		Name  string
		Type  MetricType
		Value interface{}
	}
)

var (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"

	ErrorInvalidPoolCount = errors.New("invalid poll count")
)

type RuntimeMetrics struct {
	mx      sync.Mutex
	metrics []Metric
}

func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{}
}

func (m *RuntimeMetrics) Update() error {
	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)

	m.setMetric(Metric{Name: "Alloc", Type: GaugeType, Value: float64(runtimeMetrics.Alloc)})
	m.setMetric(Metric{Name: "BuckHashSys", Type: GaugeType, Value: float64(runtimeMetrics.BuckHashSys)})
	m.setMetric(Metric{Name: "Frees", Type: GaugeType, Value: float64(runtimeMetrics.Frees)})
	m.setMetric(Metric{Name: "GCCPUFraction", Type: GaugeType, Value: runtimeMetrics.GCCPUFraction})
	m.setMetric(Metric{Name: "GCSys", Type: GaugeType, Value: float64(runtimeMetrics.GCSys)})
	m.setMetric(Metric{Name: "HeapAlloc", Type: GaugeType, Value: float64(runtimeMetrics.HeapAlloc)})
	m.setMetric(Metric{Name: "HeapIdle", Type: GaugeType, Value: float64(runtimeMetrics.HeapIdle)})
	m.setMetric(Metric{Name: "HeapInuse", Type: GaugeType, Value: float64(runtimeMetrics.HeapInuse)})
	m.setMetric(Metric{Name: "HeapObjects", Type: GaugeType, Value: float64(runtimeMetrics.HeapObjects)})
	m.setMetric(Metric{Name: "HeapReleased", Type: GaugeType, Value: float64(runtimeMetrics.HeapReleased)})
	m.setMetric(Metric{Name: "HeapSys", Type: GaugeType, Value: float64(runtimeMetrics.HeapSys)})
	m.setMetric(Metric{Name: "LastGC", Type: GaugeType, Value: float64(runtimeMetrics.LastGC)})
	m.setMetric(Metric{Name: "Lookups", Type: GaugeType, Value: float64(runtimeMetrics.Lookups)})
	m.setMetric(Metric{Name: "MCacheInuse", Type: GaugeType, Value: float64(runtimeMetrics.MCacheInuse)})
	m.setMetric(Metric{Name: "MCacheSys", Type: GaugeType, Value: float64(runtimeMetrics.MCacheSys)})
	m.setMetric(Metric{Name: "MSpanInuse", Type: GaugeType, Value: float64(runtimeMetrics.MSpanInuse)})
	m.setMetric(Metric{Name: "MSpanSys", Type: GaugeType, Value: float64(runtimeMetrics.MSpanSys)})
	m.setMetric(Metric{Name: "Mallocs", Type: GaugeType, Value: float64(runtimeMetrics.Mallocs)})
	m.setMetric(Metric{Name: "NextGC", Type: GaugeType, Value: float64(runtimeMetrics.NextGC)})
	m.setMetric(Metric{Name: "NumForcedGC", Type: GaugeType, Value: float64(runtimeMetrics.NumForcedGC)})
	m.setMetric(Metric{Name: "NumGC", Type: GaugeType, Value: float64(runtimeMetrics.NumGC)})
	m.setMetric(Metric{Name: "OtherSys", Type: GaugeType, Value: float64(runtimeMetrics.OtherSys)})
	m.setMetric(Metric{Name: "PauseTotalNs", Type: GaugeType, Value: float64(runtimeMetrics.PauseTotalNs)})
	m.setMetric(Metric{Name: "StackInuse", Type: GaugeType, Value: float64(runtimeMetrics.StackInuse)})
	m.setMetric(Metric{Name: "StackSys", Type: GaugeType, Value: float64(runtimeMetrics.StackSys)})
	m.setMetric(Metric{Name: "Sys", Type: GaugeType, Value: float64(runtimeMetrics.Sys)})
	m.setMetric(Metric{Name: "TotalAlloc", Type: GaugeType, Value: float64(runtimeMetrics.TotalAlloc)})

	pollCount, ok := m.getMetric("PollCount")
	if !ok {
		m.setMetric(Metric{Name: "PollCount", Type: CounterType, Value: int64(0)})
	} else {
		pollCountValue, parsedOk := pollCount.Value.(int64)
		if !parsedOk {
			return ErrorInvalidPoolCount
		}
		m.setMetric(Metric{Name: "PollCount", Type: CounterType, Value: pollCountValue + 1})
	}

	m.setMetric(Metric{Name: "RandomValue", Type: GaugeType, Value: rand.Float64()})

	return nil
}

func (m *RuntimeMetrics) getMetric(name string) (Metric, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, metric := range m.metrics {
		if metric.Name == name {
			return metric, true
		}
	}

	return Metric{}, false
}

func (m *RuntimeMetrics) setMetric(metric Metric) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for key, mc := range m.metrics {
		if mc.Name == metric.Name {
			m.metrics[key] = metric
			return
		}
	}

	m.metrics = append(m.metrics, metric)
}

func (m *RuntimeMetrics) GetMetrics() []Metric {
	m.mx.Lock()
	defer m.mx.Unlock()

	return m.metrics
}
