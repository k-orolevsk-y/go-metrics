package memstorage

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"sync"
)

type tx struct {
	storage *MemStorage

	mx   sync.Mutex
	rows []models.MetricsUpdate
}

func (t *tx) SetGauge(name string, value float64) error {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.rows = append(t.rows, models.MetricsUpdate{
		ID:    name,
		MType: "gauge",
		Value: &value,
	})

	return nil
}

func (t *tx) AddCounter(name string, value int64) error {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.rows = append(t.rows, models.MetricsUpdate{
		ID:    name,
		MType: "counter",
		Delta: &value,
	})

	return nil
}

func (t *tx) Commit() error {
	t.mx.Lock()
	defer t.mx.Unlock()

	for _, row := range t.rows {
		switch row.MType {
		case "gauge":
			_ = t.storage.SetGauge(row.ID, *row.Value)
		case "counter":
			_ = t.storage.AddCounter(row.ID, *row.Delta)
		}
	}

	t.rows = []models.MetricsUpdate{}
	return nil
}

func (t *tx) RollBack() error {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.rows = []models.MetricsUpdate{}
	return nil
}
