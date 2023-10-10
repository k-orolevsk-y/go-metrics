package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
)

func New() (*sql.DB, error) {
	return sql.Open("pgx", config.Config.DatabaseDSN)
}
