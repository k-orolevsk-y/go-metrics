package filestorage

import (
	"context"
	"encoding/json"
	"errors"
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
	Storage struct {
		*memstorage.Mem

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

func New(log logger) (*Storage, error) {
	file, err := os.OpenFile(config.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	store := memstorage.NewMem()

	return &Storage{
		Mem: store,

		file: file,
		log:  log,

		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (s *Storage) Close() error {
	return s.file.Close()
}

func (s *Storage) Restore() error {
	if err := s.file.Sync(); err != nil {
		return err
	}

	var metrics []models.MetricsValue
	if err := s.decoder.Decode(&metrics); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	var errorsCount int
	for _, metric := range metrics {
		switch metric.MType {
		case string(models.GaugeType):
			_ = s.SetGauge(metric.ID, *metric.Value)
		case string(models.CounterType):
			_ = s.AddCounter(metric.ID, *metric.Delta)
		default:
			errorsCount++
			s.log.Errorf("The metric couldn't be restored, it has an unknown type: %+v", metrics)
		}
	}

	s.log.Infof("Successfully retrieved metrics (%d) from the file.", len(metrics)-errorsCount)
	return nil
}

func (s *Storage) Start() {
	storeInterval := config.Config.StoreInterval
	if storeInterval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(storeInterval))
		for range ticker.C {
			if count, err := s.update(); err != nil {
				s.log.Errorf("Failed to save metrics to file: %s", err)
			} else {
				s.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
			}
		}
	}()
}

func (s *Storage) update() (int, error) {
	if err := s.file.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := s.file.Seek(0, 0); err != nil {
		return 0, err
	}

	metrics, err := s.GetAll()
	if err != nil {
		return 0, err
	}

	if err = s.encoder.Encode(&metrics); err != nil {
		return 0, err
	}

	return len(metrics), nil
}

func (s *Storage) Ping(_ context.Context) error {
	_, err := s.file.Stat()
	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func (s *Storage) GetMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		storeInterval := config.Config.StoreInterval
		if storeInterval > 0 {
			return
		} else if !strings.Contains(ctx.FullPath(), "/update") {
			return
		}

		ctx.Next()

		if count, err := s.update(); err != nil {
			s.log.Errorf("Failed to save metrics to file: %s", err)
		} else {
			s.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
		}
	}
}
