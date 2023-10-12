package dbstorage

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type databaseStorage struct {
	db  *sql.DB
	log logger.Logger
}

func New(db *sql.DB, log logger.Logger) (*databaseStorage, error) {
	dbStorage := &databaseStorage{
		db:  db,
		log: log,
	}

	if err := dbStorage.buildTable(context.Background()); err != nil {
		return nil, err
	}

	return dbStorage, nil
}

func (dbStorage *databaseStorage) Close() error {
	//TODO implement me
	panic("implement me")
}

func (dbStorage *databaseStorage) buildTable(ctx context.Context) error {
	_, err := dbStorage.db.ExecContext(ctx, `SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME='metrics'`)
	if err != nil {
		return err
	}

	_, err = dbStorage.db.ExecContext(ctx, `CREATE TABLE metrics (
    	"id" TEXT PRIMARY KEY,
    	"type" VARCHAR(12) NOT NULL DEFAULT 'gauge',
    	"delta" DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    	"value" INTEGER NOT NULL DEFAULT 0
	)`)

	return err
}

func (dbStorage *databaseStorage) SetGauge(name string, value float64) {
	//TODO implement me
	panic("implement me")
}

func (dbStorage *databaseStorage) AddCounter(name string, value int64) {
	//TODO implement me
	panic("implement me")
}

func (dbStorage *databaseStorage) GetGauge(name string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (dbStorage *databaseStorage) GetCounter(name string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (dbStorage *databaseStorage) GetAll() []models.MetricsValue {
	//TODO implement me
	panic("implement me")
}

func (dbStorage *databaseStorage) Ping(ctx context.Context) error {
	return dbStorage.db.PingContext(ctx)
}

func (dbStorage *databaseStorage) GetMiddleware() gin.HandlerFunc {
	return func(_ *gin.Context) {}
}
