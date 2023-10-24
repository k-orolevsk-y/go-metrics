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
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
	"io"
	"os"
	"strings"
	"time"
)

var retries = []int{1, 3, 5}

type (
	fileStorage struct {
		*memstorage.MemStorage

		file *os.File
		log  logger.Logger

		encoder *json.Encoder
		decoder *json.Decoder
	}
)

func New(log logger.Logger) (*fileStorage, error) {
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
	var (
		errs    []error
		metrics []models.MetricsValue
	)

	for _, timeSleep := range retries {
		if err := fStorage.file.Sync(); err != nil {
			errs = append(errs, err)
		}

		if err := fStorage.decoder.Decode(&metrics); err != nil {
			if !errors.Is(err, io.EOF) {
				errs = append(errs, err)
			}
		}

		if len(errs) <= 0 {
			break
		}

		fStorage.log.Errorf("Failed to restore metrics in file: %s. Retrying after %ds...", errors.Join(errs...), timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	var errorsCount int
	for _, metric := range metrics {
		switch metric.MType {
		case string(models.GaugeType):
			_ = fStorage.SetGauge(metric.ID, metric.Value)
		case string(models.CounterType):
			_ = fStorage.AddCounter(metric.ID, metric.Delta)
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
		fStorage.log.Debugf("A ticker was created and launched to update the metrics in the file")

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
	metrics, err := fStorage.GetAll()
	if err != nil {
		return 0, err
	}

	var errs []error
	for _, timeSleep := range retries {
		if err = fStorage.file.Truncate(0); err != nil {
			errs = append(errs, err)
		}

		if _, err = fStorage.file.Seek(0, 0); err != nil {
			errs = append(errs, err)
		}

		if err = fStorage.encoder.Encode(&metrics); err != nil {
			errs = append(errs, err)
		}

		if len(errs) <= 0 {
			break
		}

		fStorage.log.Errorf("Failed to update metrics in file: %s. Retrying after %ds...", errors.Join(errs...), timeSleep)
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}

	if len(errs) > 0 {
		return 0, errors.Join(errs...)
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
	fStorage.log.Debugf("Created and received middleware to update metrics after a request.")
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
