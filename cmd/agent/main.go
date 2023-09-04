package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
	"time"
)

var (
	poolInterval   = time.Duration(2)
	reportInterval = time.Duration(10)

	ErrorInvalidMetricValueType = errors.New("invalid metric value type")
)

const APIUrl = "http://localhost:8080"

func main() {
	var metrics RuntimeMetrics
	metrics.New()

	restyClient := resty.New()
	for {
		go func() {
			err := metrics.Renew()
			if err != nil {
				panic(err)
			}
		}()

		updateMetrics(&metrics, restyClient)
	}
}

func updateMetrics(m *RuntimeMetrics, restyClient *resty.Client) {
	time.Sleep(time.Second * reportInterval)

	for k, v := range m.Runtime {
		if err := updateMetric(k, v, restyClient); err != nil {
			log.Printf("[Warning] %s - %v", k, err)
		}
	}

	if err := updateMetric("PollCount", m.PollCount, restyClient); err != nil {
		log.Printf("[Warning] PollCount - %v", err)
	}

	if err := updateMetric("RandomValue", m.RandomValue, restyClient); err != nil {
		log.Printf("[Warning] RandomValue - %v", err)
	}
}

func updateMetric(name string, metric Metric, restyClient *resty.Client) error {
	var value interface{}

	switch metric.Value.(type) {
	case float64:
		if metric.Type != GaugeType {
			return ErrorInvalidMetricValueType
		}

		valueFloat64 := metric.Value.(float64)
		value = strconv.FormatFloat(valueFloat64, 'f', 1, 64)
	case int64:
		if metric.Type != CounterType {
			return ErrorInvalidMetricValueType
		}

		value = metric.Value.(int64)
	default:
		return ErrorInvalidMetricValueType
	}

	url := fmt.Sprintf("%s/update/%s/%s/%v", APIUrl, metric.Type, name, value)

	_, err := restyClient.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		return err
	}

	return nil
}
