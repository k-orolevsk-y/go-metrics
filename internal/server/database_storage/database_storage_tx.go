package dbstorage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
	"time"
)

type tx struct {
	txDB                     *sqlx.Tx
	prepareSetOrUpdateMetric *sqlx.NamedStmt

	log logger.Logger
}

func (t *tx) buildPrepares(ctx context.Context) (err error) {
	t.prepareSetOrUpdateMetric, err = t.txDB.PrepareNamedContext(ctx, setOrUpdateMetricSQLRequest)
	return
}

func (t *tx) SetGauge(name string, value *float64) (err error) {
	for _, timeSleep := range maximumNumberOfRetries {
		_, err = t.prepareSetOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "gauge", "delta": 0, "value": value})

		ok, parsedErr := parseRetriableError(err)
		if !ok {
			return
		}

		t.log.Errorf("Error set gauge metric (TX) %s: \"%s\". Retrying after %ds...", name, parsedErr, timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
	return
}

func (t *tx) AddCounter(name string, value *int64) (err error) {
	for _, timeSleep := range maximumNumberOfRetries {
		_, err = t.prepareSetOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "counter", "delta": value, "value": 0.0})

		ok, parsedErr := parseRetriableError(err)
		if !ok {
			return
		}

		t.log.Errorf("Error add counter metric (TX) %s: \"%s\". Retrying after %ds...", name, parsedErr, timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
	return
}

func (t *tx) Commit() (err error) {
	for _, timeSleep := range maximumNumberOfRetries {
		err = t.txDB.Commit()

		ok, parsedErr := parseRetriableError(err)
		if !ok {
			return
		}

		t.log.Errorf("Error commit transaction: \"%s\". Retrying after %ds...", parsedErr, timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}

	return
}

func (t *tx) RollBack() (err error) {
	for _, timeSleep := range maximumNumberOfRetries {
		err = t.txDB.Rollback()

		ok, parsedErr := parseRetriableError(err)
		if !ok {
			return
		}

		t.log.Errorf("Error rollback transaction: \"%s\". Retrying after %ds...", parsedErr, timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}

	return
}
