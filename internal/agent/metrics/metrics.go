package metrics

import (
	"errors"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"math/rand"
	"runtime"
	"time"
)

type (
	MetricType string
	Metric     struct {
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
	Runtime     map[string]Metric
	PollCount   Metric
	RandomValue Metric
}

func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{
		Runtime:     make(map[string]Metric),
		PollCount:   Metric{Type: CounterType, Value: int64(0)},
		RandomValue: Metric{Type: GaugeType, Value: float64(0)},
	}
}

func (m *RuntimeMetrics) Update() error {
	time.Sleep(time.Second * time.Duration(config.Config.PollInterval))

	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)

	m.Runtime["Alloc"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Alloc)}
	m.Runtime["BuckHashSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.BuckHashSys)}
	m.Runtime["GCCPUFraction"] = Metric{Type: GaugeType, Value: runtimeMetrics.GCCPUFraction}
	m.Runtime["HeapAlloc"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapAlloc)}
	m.Runtime["HeapIdle"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapIdle)}
	m.Runtime["HeapInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapInuse)}
	m.Runtime["HeapObjects"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapObjects)}
	m.Runtime["HeapReleased"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapReleased)}
	m.Runtime["HeapSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapSys)}
	m.Runtime["LastGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.LastGC)}
	m.Runtime["Lookups"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Lookups)}
	m.Runtime["MCacheInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MCacheInuse)}
	m.Runtime["MCacheSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MCacheSys)}
	m.Runtime["MSpanInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MSpanInuse)}
	m.Runtime["Mallocs"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Mallocs)}
	m.Runtime["NumForcedGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.NumForcedGC)}
	m.Runtime["NumGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.NumGC)}
	m.Runtime["OtherSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.OtherSys)}
	m.Runtime["PauseTotalNs"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.PauseTotalNs)}
	m.Runtime["StackInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.StackInuse)}
	m.Runtime["StackSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.StackSys)}
	m.Runtime["Sys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Sys)}
	m.Runtime["TotalAlloc"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.TotalAlloc)}

	pollCount, ok := m.PollCount.Value.(int64)
	if !ok {
		return ErrorInvalidPoolCount
	}
	m.PollCount.Value = pollCount + 1

	m.RandomValue.Value = rand.Float64()
	return nil
}
