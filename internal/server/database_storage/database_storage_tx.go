package dbstorage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type tx struct {
	txDB                     *sqlx.Tx
	prepareSetOrUpdateMetric *sqlx.NamedStmt

	log logger.Logger
}

func (t *tx) buildPrepares(ctx context.Context) (err error) {
	t.prepareSetOrUpdateMetric, err = t.txDB.PrepareNamedContext(
		ctx,
		`INSERT INTO metrics (name, mtype, delta, value) 
				VALUES (:name, :mtype, :delta, :value)
			ON CONFLICT (name, mtype) DO 
			    UPDATE SET delta = metrics.delta + excluded.delta, value = excluded.value`,
	)
	return
}

func (t *tx) SetGauge(name string, value *float64) (err error) {
	_, err = t.prepareSetOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "gauge", "delta": 0, "value": value})
	return
}

func (t *tx) AddCounter(name string, value *int64) (err error) {
	_, err = t.prepareSetOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "counter", "delta": value, "value": 0.0})
	return
}

func (t *tx) Commit() (err error) {
	err = t.txDB.Commit()
	return
}

func (t *tx) RollBack() (err error) {
	err = t.txDB.Rollback()
	return
}
