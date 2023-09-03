package main

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlerServer(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Test_updateMetric(t *testing.T) {
	type args struct {
		name   string
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Positive Gauge Test",
			args:    args{"TestGauge", Metric{Type: GaugeType, Value: float64(0.5)}},
			wantErr: false,
		},
		{
			name:    "Positive Counter Test",
			args:    args{"TestCounter", Metric{Type: CounterType, Value: int64(1)}},
			wantErr: false,
		},
		{
			name:    "Negative Gauge Test",
			args:    args{"TestGauge", Metric{Type: GaugeType, Value: int64(10)}},
			wantErr: true,
		},
		{
			name:    "Positive Counter Test",
			args:    args{"TestCounter", Metric{Type: CounterType, Value: float64(5.10)}},
			wantErr: true,
		},
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(handlerServer))
	if err = server.Listener.Close(); err != nil {
		t.Fatal("failed to close default listener:", err)
	}
	server.Listener = l

	server.Start()
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = updateMetric(tt.args.name, tt.args.metric)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrorInvalidMetricValueType)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
