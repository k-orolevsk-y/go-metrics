package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
)

func New() (*sqlx.DB, error) {
	return sqlx.Open("pgx", config.Config.DatabaseDSN)
}
