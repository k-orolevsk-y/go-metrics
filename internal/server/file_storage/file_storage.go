package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"io"
	"os"
	"strings"
	"time"
)

type (
	fileStorage struct {
		*memstorage.MemStorage

		file *os.File
		log  logger

		encoder *json.Encoder
		decoder *json.Decoder
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}
)

func New(log logger) (*fileStorage, error) {
	file, err := os.OpenFile(config.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	store := memstorage.NewMem()

	return &fileStorage{
		MemStorage: store,

		file: file,
		log:  log,

		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (fStorage *fileStorage) Close() error {
	return fStorage.file.Close()
}

func (fStorage *fileStorage) Restore() error {
	if err := fStorage.file.Sync(); err != nil {
		return err
	}

	var metrics []models.MetricsValue
	if err := fStorage.decoder.Decode(&metrics); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	var errorsCount int
	for _, metric := range metrics {
		switch metric.MType {
		case string(models.GaugeType):
			_ = fStorage.SetGauge(metric.ID, *metric.Value)
		case string(models.CounterType):
			_ = fStorage.AddCounter(metric.ID, *metric.Delta)
		default:
			errorsCount++
			fStorage.log.Errorf("The metric couldn't be restored, it has an unknown type: %+v", metrics)
		}
	}

	fStorage.log.Infof("Successfully retrieved metrics (%d) from the file.", len(metrics)-errorsCount)
	return nil
}

func (fStorage *fileStorage) Start() {
	storeInterval := config.Config.StoreInterval
	if storeInterval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(storeInterval))
		for range ticker.C {
			if count, err := fStorage.update(); err != nil {
				fStorage.log.Errorf("Failed to save metrics to file: %s", err)
			} else {
				fStorage.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
			}
		}
	}()
}

func (fStorage *fileStorage) update() (int, error) {
	if err := fStorage.file.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := fStorage.file.Seek(0, 0); err != nil {
		return 0, err
	}

	metrics, err := fStorage.GetAll()
	if err != nil {
		return 0, err
	}

	if err = fStorage.encoder.Encode(&metrics); err != nil {
		return 0, err
	}

	return len(metrics), nil
}

func (fStorage *fileStorage) Ping(_ context.Context) error {
	_, err := fStorage.file.Stat()
	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func (fStorage *fileStorage) GetMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		storeInterval := config.Config.StoreInterval
		if storeInterval > 0 {
			return
		} else if !strings.Contains(ctx.FullPath(), "/update") {
			return
		}

		ctx.Next()

		if count, err := fStorage.update(); err != nil {
			fStorage.log.Errorf("Failed to save metrics to file: %s", err)
		} else {
			fStorage.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
		}
	}
}

func (fStorage *fileStorage) String() string {
	return fmt.Sprintf("FileStorage - %s", fStorage.file.Name())
}
