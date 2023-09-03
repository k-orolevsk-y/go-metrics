package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	poolInterval   = time.Duration(2)
	reportInterval = time.Duration(10)
)

const APIUrl = "http://localhost:8080"

func main() {
	var metrics RuntimeMetrics
	metrics.New()

	for {
		go func() {
			err := metrics.Renew()
			if err != nil {
				panic(err)
			}
		}()

		updateMetrics(&metrics)
	}
}

func updateMetrics(m *RuntimeMetrics) {
	time.Sleep(time.Second * reportInterval)

	for k, v := range m.Runtime {
		if err := updateMetric(k, v); err != nil {
			panic(err)
		}
	}

	if err := updateMetric("PollCount", m.PollCount); err != nil {
		panic(err)
	}

	if err := updateMetric("RandomValue", m.RandomValue); err != nil {
		panic(err)
	}
}

func updateMetric(name string, metric Metric) error {
	var value interface{}

	if metric.Type == GaugeType {
		valueFloat64 := metric.Value.(float64)
		value = strconv.FormatFloat(valueFloat64, 'f', 1, 64)
	} else {
		value = metric.Value.(int64)
	}

	url := fmt.Sprintf("%s/update/%s/%s/%v", APIUrl, metric.Type, name, value)

	_, err := http.Post(url, "text/plain", nil)
	return err
}
