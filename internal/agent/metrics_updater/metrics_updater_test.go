package metricsupdater

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
)

func handlerServer(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestUpdater_updateMetric(t *testing.T) {
	tests := []struct {
		name      string
		metrics   []metrics.Metric
		wantedErr bool
	}{
		{
			name:      "Positive Gauge Test",
			metrics:   []metrics.Metric{metrics.NewMetric("TestGaugePositive", metrics.GaugeType, 0, float64(0.5))},
			wantedErr: false,
		},
		{
			name:      "Positive Counter Test",
			metrics:   []metrics.Metric{metrics.NewMetric("TestCounterPositive", metrics.CounterType, 1, 0)},
			wantedErr: false,
		},
	}

	require.NoError(t, os.Setenv("POLL_INTERVAL", "1"))
	require.NoError(t, os.Setenv("REPORT_INTERVAL", "1"))

	config.Load()
	require.NoError(t, config.Parse())

	l, err := net.Listen("tcp", ":8080")
	require.NoError(t, err)

	server := httptest.NewUnstartedServer(http.HandlerFunc(handlerServer))
	if err = server.Listener.Close(); err != nil {
		t.Fatal("failed to close default listener:", err)
	}
	server.Listener = l

	server.Start()
	defer server.Close()

	buf := new(bytes.Buffer)
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.DebugLevel),
		zap.AddCaller(),
	)

	client := resty.New()
	updater := New(client, nil, log.Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, updater.updateMetrics(tt.metrics))

			if tt.wantedErr {
				require.Contains(t, buf.String(), "Invalid metric type:")
			}

			buf.Reset()
		})
	}
}
