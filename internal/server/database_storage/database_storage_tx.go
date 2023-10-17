package dbstorage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type tx struct {
	txDB                     *sqlx.Tx
	prepareSetOrUpdateMetric *sqlx.NamedStmt
}

func (t *tx) buildPrepares(ctx context.Context) (err error) {
	t.prepareSetOrUpdateMetric, err = t.txDB.PrepareNamedContext(ctx, setOrUpdateMetricSQLRequest)
	return
}

func (t *tx) SetGauge(name string, value float64) (err error) {
	_, err = t.prepareSetOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "gauge", "delta": 0, "value": value})
	return
}

func (t *tx) AddCounter(name string, value int64) (err error) {
	_, err = t.prepareSetOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "counter", "delta": value, "value": 0.0})
	return
}

func (t *tx) Commit() error {
	return t.txDB.Commit()
}

func (t *tx) RollBack() error {
	return t.txDB.Rollback()
}
