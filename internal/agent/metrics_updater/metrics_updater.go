package metricsupdater

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"strconv"
)

var (
	ErrorInvalidMetricValueType = errors.New("invalid metric value type")
)

type (
	Updater struct {
		client *resty.Client
		store  store
		log    logger
	}

	logger interface {
		Errorf(template string, args ...interface{})
	}

	store interface {
		GetRuntime() map[string]metrics.Metric
		GetPollCount() metrics.Metric
		GetRandomValue() metrics.Metric
	}
)

func New(client *resty.Client, store store, log logger) *Updater {
	return &Updater{
		client: client,
		store:  store,
		log:    log,
	}
}

func (u Updater) UpdateMetrics() {
	for k, v := range u.store.GetRuntime() {
		if err := u.updateMetric(k, v); err != nil {
			u.log.Errorf("Failed to update metric \"%s\": %s", k, err)
		}
	}

	if err := u.updateMetric("PollCount", u.store.GetPollCount()); err != nil {
		u.log.Errorf("Failed to update metric \"PollCount\": %s", err)
	}

	if err := u.updateMetric("RandomValue", u.store.GetRandomValue()); err != nil {
		u.log.Errorf("Failed to update metric \"PollCount\": %s", err)
	}
}

func (u Updater) updateMetric(name string, metric metrics.Metric) error {
	value, err := u.parseMetric(metric)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/update/%s/%s/%v", config.Config.Address, metric.Type, name, value)
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
