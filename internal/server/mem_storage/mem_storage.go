package memstorage

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/errs"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"strings"
	"sync"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64

	mx sync.Mutex
}

func NewMem() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (mStorage *MemStorage) Close() error {
	*mStorage = MemStorage{}
	return nil
}

func (mStorage *MemStorage) NewTx() (models.StorageTx, error) {
	return &tx{
		storage: mStorage,
	}, nil
}

func (mStorage *MemStorage) GetGauge(name string) (float64, error) {
	mStorage.mx.Lock()
	defer mStorage.mx.Unlock()

	name = mStorage.normalizeName(name)

	value, ok := mStorage.gauge[name]
	if !ok {
		return 0, errs.ErrStorageInvalidGaugeName
	}

	return value, nil
}

func (mStorage *MemStorage) SetGauge(name string, value float64) error {
	mStorage.mx.Lock()
	defer mStorage.mx.Unlock()

	name = mStorage.normalizeName(name)
	mStorage.gauge[name] = value

	return nil
}

func (mStorage *MemStorage) GetCounter(name string) (int64, error) {
	mStorage.mx.Lock()
	defer mStorage.mx.Unlock()

	name = mStorage.normalizeName(name)
	value, ok := mStorage.counter[name]

	if !ok {
		return 0, errs.ErrStorageInvalidCounterName
	}

	return value, nil
}

func (mStorage *MemStorage) AddCounter(name string, value int64) error {
	name = mStorage.normalizeName(name)
	_, err := mStorage.GetCounter(name)

	mStorage.mx.Lock()
	defer mStorage.mx.Unlock()

	if err != nil {
		mStorage.counter[name] = value
	} else {
		mStorage.counter[name] += value
	}

	return nil
}

func (mStorage *MemStorage) GetAll() ([]models.MetricsValue, error) {
	var values []models.MetricsValue

	for k, v := range mStorage.gauge {
		values = append(values, models.MetricsValue{
			ID:    k,
			MType: string(models.GaugeType),
			Value: &v,
		})
	}

	for k, v := range mStorage.counter {
		values = append(values, models.MetricsValue{
			ID:    k,
			MType: string(models.CounterType),
			Delta: &v,
		})
	}

	return values, nil
}

func (mStorage *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (mStorage *MemStorage) GetMiddleware() gin.HandlerFunc {
	return func(_ *gin.Context) {}
}

func (mStorage *MemStorage) String() string {
	return fmt.Sprintf("MemStorage - Pointer(%+v)", &mStorage)
}

func (mStorage *MemStorage) normalizeName(name string) string {
	return strings.TrimSpace(name)
}
