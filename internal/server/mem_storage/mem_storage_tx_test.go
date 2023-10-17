package memstorage

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStorageTx(t *testing.T) {
	memStorage := NewMem()

	txx, err := memStorage.NewTx()
	require.NoError(t, err)

	require.NoError(t, txx.AddCounter("Test", 100))
	require.NoError(t, txx.AddCounter("Test", 123))

	require.NoError(t, txx.SetGauge("Wow", 13.5))
	require.NoError(t, txx.SetGauge("Go", 199.3492))
	require.NoError(t, txx.SetGauge("Wow", 20))

	require.NoError(t, txx.Commit())

	metrics, err := memStorage.GetAll()
	require.NoError(t, err)
	require.Equal(t, 3, len(metrics))

	metricsBe := []struct {
		name  string
		mType models.MetricType
		value interface{}
	}{
		{
			name:  "Test",
			mType: models.CounterType,
			value: int64(223),
		},
		{
			name:  "Wow",
			mType: models.GaugeType,
			value: float64(20),
		},
		{
			name:  "Go",
			mType: models.GaugeType,
			value: 199.3492,
		},
	}

	for _, metric := range metricsBe {
		var value interface{}

		switch metric.mType {
		case models.GaugeType:
			value, err = memStorage.GetGauge(metric.name)
		case models.CounterType:
			value, err = memStorage.GetCounter(metric.name)
		default:
			continue
		}

		require.NoError(t, err)
		require.Equal(t, metric.value, value)
	}
}
