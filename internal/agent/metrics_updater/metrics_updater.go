package metricsupdater

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"log"
	"strconv"
	"time"
)

var (
	ErrorInvalidMetricValueType = errors.New("invalid metric value type")
)

type Updater struct {
	client *resty.Client
	store  *metrics.RuntimeMetrics
}

func New(client *resty.Client, store *metrics.RuntimeMetrics) *Updater {
	return &Updater{
		client: client,
		store:  store,
	}
}

func (u Updater) UpdateMetrics() {
	time.Sleep(time.Second * time.Duration(config.Data.ReportInterval))

	for k, v := range u.store.Runtime {
		if err := u.updateMetric(k, v); err != nil {
			log.Printf("[Warning] %s - %v", k, err)
		}
	}

	if err := u.updateMetric("PollCount", u.store.PollCount); err != nil {
		log.Printf("[Warning] PollCount - %v", err)
	}

	if err := u.updateMetric("RandomValue", u.store.RandomValue); err != nil {
		log.Printf("[Warning] RandomValue - %v", err)
	}
}

func (u Updater) updateMetric(name string, metric metrics.Metric) error {
	value, err := u.parseMetric(metric)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/update/%s/%s/%v", config.Data.Address, metric.Type, name, value)
	_, err = u.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (u Updater) parseMetric(metric metrics.Metric) (interface{}, error) {
	var value interface{}

	switch metric.Value.(type) {
	case float64:
		if metric.Type != metrics.GaugeType {
			return nil, ErrorInvalidMetricValueType
		}

		valueFloat64 := metric.Value.(float64)
		value = strconv.FormatFloat(valueFloat64, 'f', 1, 64)
	case int64:
		if metric.Type != metrics.CounterType {
			return nil, ErrorInvalidMetricValueType
		}

		value = metric.Value.(int64)
	default:
		return nil, ErrorInvalidMetricValueType
	}

	return value, nil
}

type MetricUpdater interface {
	UpdateMetrics()
	updateMetric(name string, metric metrics.Metric) error
	parseMetric(metric metrics.Metric) (interface{}, error)
}