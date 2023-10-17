package filestorage

import (
	"bytes"
	"encoding/json"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"os"
	"testing"
)

func getPointerFloat64(v float64) *float64 {
	return &v
}

func getPointerInt64(v int64) *int64 {
	return &v
}

func TestSuccessFileStorage(t *testing.T) {
	metrics := []models.MetricsValue{
		{
			ID:    "TestGauge",
			MType: string(models.GaugeType),
			Value: getPointerFloat64(100.5),
		},
		{
			ID:    "TestCounter",
			MType: string(models.CounterType),
			Delta: getPointerInt64(321),
		},
		{
			ID:    "TestFloat64SimilarInt64",
			MType: string(models.GaugeType),
			Value: getPointerFloat64(300.0),
		},
	}

	file, err := os.CreateTemp(t.TempDir(), "tests-file-mem_storage-*.json")
	require.NoError(t, err)
	t.Setenv("FILE_STORAGE_PATH", file.Name())

	require.NoError(t, config.Parse())

	log := zaptest.NewLogger(t).Sugar()

	fStorage, err := New(log)
	require.NoError(t, err)

	for _, metric := range metrics {
		switch metric.MType {
		case string(models.GaugeType):
			_ = fStorage.SetGauge(metric.ID, metric.Value)
		case string(models.CounterType):
			_ = fStorage.AddCounter(metric.ID, metric.Delta)
		}
	}

	count, err := fStorage.update()

	require.NoError(t, err)
	require.Equal(t, count, len(metrics))

	fStorage, err = New(log)
	require.NoError(t, err)

	require.NoError(t, fStorage.Restore())

	metrics, _ = fStorage.GetAll()
	require.Equal(t, len(metrics), len(metrics))

	if err = os.Remove(file.Name()); err != nil {
		t.Logf("Не удалось удалить тестовый json-файл: %s", err)
	}
}

func TestNegativeFileStorage(t *testing.T) {
	metrics := []models.MetricsValue{
		{
			ID:    "TestGauge",
			MType: "InvalidType",
			Value: getPointerFloat64(100.5),
		},
	}

	file, err := os.CreateTemp(t.TempDir(), "tests-file-mem_storage-*.json")
	require.NoError(t, err)
	t.Setenv("FILE_STORAGE_PATH", file.Name())

	jsonBytes, err := json.Marshal(&metrics)
	require.NoError(t, err)

	_, err = file.Write(jsonBytes)
	require.NoError(t, err)

	require.NoError(t, config.Parse())

	buf := new(bytes.Buffer)
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.DebugLevel),
		zap.AddCaller(),
	)

	fStorage, err := New(log.Sugar())
	require.NoError(t, err)

	require.NoError(t, fStorage.Restore())
	require.Contains(t, buf.String(), "The metric couldn't be restored, it has an unknown type")

	if err = os.Remove(file.Name()); err != nil {
		t.Logf("Не удалось удалить тестовый json-файл: %s", err)
	}
}
