package dbstorage

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
	"time"
)

const schema = `
	CREATE TABLE IF NOT EXISTS metrics (
	    "id" SERIAL,
		"name" TEXT NOT NULL,
		"mtype" VARCHAR(12) NOT NULL DEFAULT 'gauge',
		"delta" INTEGER NOT NULL DEFAULT 0,
		"value" DOUBLE PRECISION NOT NULL DEFAULT 0.0,
		CONSTRAINT unique_id_mtype UNIQUE (name, mtype),
		PRIMARY KEY (id)
	)
`

type (
	databaseStorage struct {
		db  *sqlx.DB
		log logger.Logger

		prepares prepares
	}

	prepares struct {
		getGaugeMetric    *sqlx.NamedStmt
		getCounterMetric  *sqlx.NamedStmt
		setOrUpdateMetric *sqlx.NamedStmt
	}
)

func New(db *sqlx.DB, log logger.Logger) (*databaseStorage, error) {
	dbStorage := &databaseStorage{
		db:  db,
		log: log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := dbStorage.Ping(ctx); err != nil {
		return dbStorage, nil
	}

	if _, err := dbStorage.db.ExecContext(ctx, schema); err != nil {
		return nil, err
	}

	if err := dbStorage.buildPrepares(ctx); err != nil {
		return nil, err
	}

	return dbStorage, nil
}

func (dbStorage *databaseStorage) buildPrepares(ctx context.Context) error {
	preparesData := map[string]string{
		"getGaugeMetric":   "SELECT value FROM metrics WHERE name = :name AND mtype = 'gauge'",
		"getCounterMetric": "SELECT delta FROM metrics WHERE name = :name AND mtype = 'counter'",
		"setOrUpdateMetric": `
			INSERT INTO metrics (name, mtype, delta, value) 
				VALUES (:name, :mtype, :delta, :value)
			ON CONFLICT (name, mtype) DO 
			    UPDATE SET delta = metrics.delta + excluded.delta, value = excluded.value 
		`,
	}

	for key, sql := range preparesData {
		p, err := dbStorage.db.PrepareNamedContext(ctx, sql)
		if err != nil {
			return err
		}

		switch key {
		case "getGaugeMetric":
			dbStorage.prepares.getGaugeMetric = p
		case "getCounterMetric":
			dbStorage.prepares.getCounterMetric = p
		case "setOrUpdateMetric":
			dbStorage.prepares.setOrUpdateMetric = p
		}
	}

	return nil
}

func (dbStorage *databaseStorage) Close() error {
	return dbStorage.db.Close()
}

func (dbStorage *databaseStorage) SetGauge(name string, value float64) (err error) {
	_, err = dbStorage.prepares.setOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "gauge", "delta": 0, "value": value})
	return err
}

func (dbStorage *databaseStorage) AddCounter(name string, value int64) (err error) {
	_, err = dbStorage.prepares.setOrUpdateMetric.ExecContext(context.Background(), map[string]interface{}{"name": name, "mtype": "counter", "delta": value, "value": 0.0})
	return
}

func (dbStorage *databaseStorage) GetGauge(name string) (value float64, err error) {
	err = dbStorage.prepares.getGaugeMetric.GetContext(context.Background(), &value, map[string]interface{}{"name": name})
	return
}

func (dbStorage *databaseStorage) GetCounter(name string) (value int64, err error) {
	err = dbStorage.prepares.getCounterMetric.GetContext(context.Background(), &value, map[string]interface{}{"name": name})
	return
}

func (dbStorage *databaseStorage) GetAll() (metrics []models.MetricsValue, err error) {
	err = dbStorage.db.SelectContext(context.Background(), &metrics, "SELECT name, mtype, delta, value FROM metrics")
	return
}

func (dbStorage *databaseStorage) Ping(ctx context.Context) error {
	return dbStorage.db.PingContext(ctx)
}

func (dbStorage *databaseStorage) GetMiddleware() gin.HandlerFunc {
	return func(_ *gin.Context) {}
}
