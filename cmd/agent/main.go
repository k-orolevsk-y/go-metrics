package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/agent/metrics"
	"log"
	"strconv"
	"time"
)

var (
	ErrorInvalidMetricValueType = errors.New("invalid metric value type")
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	var metricsStore metrics.RuntimeMetrics
	metricsStore.Init()

	restyClient := resty.New()
	for {
		go func() {
			err := metricsStore.Update()
			if err != nil {
				panic(err)
			}
		}()

		updateMetrics(&metricsStore, restyClient)
	}
}

func updateMetrics(metricsStore *metrics.RuntimeMetrics, restyClient *resty.Client) {
	time.Sleep(time.Second * time.Duration(config.GetReportInterval()))

	for k, v := range metricsStore.Runtime {
		if err := updateMetric(k, v, restyClient); err != nil {
			log.Printf("[Warning] %s - %v", k, err)
		}
	}

	if err := updateMetric("PollCount", metricsStore.PollCount, restyClient); err != nil {
		log.Printf("[Warning] PollCount - %v", err)
	}

	if err := updateMetric("RandomValue", metricsStore.RandomValue, restyClient); err != nil {
		log.Printf("[Warning] RandomValue - %v", err)
	}
}

func updateMetric(name string, metric metrics.Metric, restyClient *resty.Client) error {
	var value interface{}

	switch metric.Value.(type) {
	case float64:
		if metric.Type != metrics.GaugeType {
			return ErrorInvalidMetricValueType
		}

		valueFloat64 := metric.Value.(float64)
		value = strconv.FormatFloat(valueFloat64, 'f', 1, 64)
	case int64:
		if metric.Type != metrics.CounterType {
			return ErrorInvalidMetricValueType
		}

		value = metric.Value.(int64)
	default:
		return ErrorInvalidMetricValueType
	}

	url := fmt.Sprintf("http://%s/update/%s/%s/%v", config.GetAddress(), metric.Type, name, value)

	_, err := restyClient.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		return err
	}

	return nil
}
