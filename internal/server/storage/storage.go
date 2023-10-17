package storage

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/database_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/file_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/database"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func Setup(log logger.Logger) (models.Storage, error) {
	if config.Config.DatabaseDSN != "" {
		db, err := database.New()
		if err != nil {
			return nil, err
		}

		return dbstorage.New(db, log)
	} else if config.Config.FileStoragePath != "" {
		fs, err := filestorage.New(log)
		if err != nil {
			return nil, err
		}

		if err = fs.Restore(); err != nil {
			return nil, err
		}
		fs.Start()

		return fs, nil
	} else {
		return memstorage.NewMem(), nil
	}
}
