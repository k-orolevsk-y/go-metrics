package metricsupdater

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

var (
	retries                = []int{1, 3, 5}
	ErrorInvalidStatusCode = errors.New("invalid status code")
)

type (
	Updater struct {
		client *resty.Client
		store  store
		log    logger.Logger
	}

	store interface {
		GetMetrics() []metrics.Metric
	}
)

func New(client *resty.Client, store store, log logger.Logger) *Updater {
	return &Updater{
		client: client,
		store:  store,
		log:    log,
	}
}

func (u Updater) UpdateMetrics() {
	currentMetrics := u.store.GetMetrics()
	if err := u.updateMetrics(currentMetrics); err != nil {
		u.log.Errorf("Failed to update metrics: %s (%T)", err, err)
	}
}

func (u Updater) updateMetrics(metricForUpdate []metrics.Metric) error {
	url := fmt.Sprintf("http://%s/updates", config.Config.Address)
	body := u.parseMetrics(metricForUpdate)

	var err error
	for _, timeSleep := range retries {
		resp, err := u.client.R().
			SetBody(body).
			Post(url)
		if err == nil {
			if resp.StatusCode() != http.StatusOK {
				return ErrorInvalidStatusCode
			} else {
				return nil
			}
		}

		u.log.Errorf("Failed to send metrics to server: %s. Retrying after %ds...", err, timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}

	return err
}

func (u Updater) parseMetrics(metricsForParse []metrics.Metric) *[]models.Metrics {
	var objects []models.Metrics

	for _, metric := range metricsForParse {
		var obj models.Metrics
		obj.ID = metric.Name

		switch metric.Value.(type) {
		case float64:
			if metric.Type != metrics.GaugeType {
				u.log.Errorf("Invalid metric type: %s - %s != %T", metric.Name, metric.Type, metric.Value)
				continue
			}
			value := metric.Value.(float64)

			obj.MType = string(metrics.GaugeType)
			obj.Value = &value
		case int64:
			if metric.Type != metrics.CounterType {
				u.log.Errorf("Invalid metric type: %s - %s != %T", metric.Name, metric.Type, metric.Value)
				continue
			}
			delta := metric.Value.(int64)

			obj.MType = string(metrics.CounterType)
			obj.Delta = &delta
		default:
			u.log.Errorf("Invalid metric type: %s - %s", metric.Name, metric.Type)
			continue
		}

		objects = append(objects, obj)
	}

	return &objects
}
