package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name string

		metricName  string
		metricType  models.MetricType
		metricValue interface{}

		wantedBody       string
		wantedStatusCode int
	}{
		{
			name: "Positive (gauge)",

			metricName:  "TestGauge",
			metricType:  models.GaugeType,
			metricValue: 10.5,

			wantedBody:       "10.5",
			wantedStatusCode: http.StatusOK,
		},
		{
			name: "Positive (counter)",

			metricName:  "TestCounter",
			metricType:  models.CounterType,
			metricValue: int64(10500),

			wantedBody:       "10500",
			wantedStatusCode: http.StatusOK,
		},
		{
			name: "Negative (gauge)",

			metricType: models.GaugeType,

			wantedBody:       "",
			wantedStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative (counter)",

			metricType: models.CounterType,

			wantedBody:       "",
			wantedStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative (invalid metric type)",

			metricType: models.MetricType("invalid"),

			wantedBody:       "",
			wantedStatusCode: http.StatusNotFound,
		},
	}

	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.metricName != "" {
				switch tt.metricType {
				case models.CounterType:
					_ = storage.AddCounter(tt.metricName, tt.metricValue.(int64))
				case models.GaugeType:
					_ = storage.SetGauge(tt.metricName, tt.metricValue.(float64))
				}
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/value/%s/%s/", tt.metricType, tt.metricName), nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.wantedStatusCode)
			assert.Equal(t, string(body), tt.wantedBody)
		})
	}
}

func TestValueByBody(t *testing.T) {
	tests := []struct {
		name string

		metricName  string
		metricType  models.MetricType
		metricValue interface{}

		wantedBody       string
		wantedStatusCode int
	}{
		{
			name: "Positive (gauge)",

			metricName:  "TestGauge",
			metricType:  models.GaugeType,
			metricValue: 10.5,

			wantedBody:       "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":10.5}",
			wantedStatusCode: http.StatusOK,
		},
		{
			name: "Positive (counter)",

			metricName:  "TestCounter",
			metricType:  models.CounterType,
			metricValue: int64(10500),

			wantedBody:       "{\"id\":\"TestCounter\",\"type\":\"counter\",\"delta\":10500}",
			wantedStatusCode: http.StatusOK,
		},
		{
			name: "Negative (gauge)",

			metricType: models.GaugeType,

			wantedBody:       "{\"error\":\"Field validation for \\\"ID\\\" failed on the 'required' tag.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative (counter)",

			metricType: models.CounterType,

			wantedBody:       "{\"error\":\"Field validation for \\\"ID\\\" failed on the 'required' tag.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
	}

	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.metricName != "" {
				switch tt.metricType {
				case models.CounterType:
					_ = storage.AddCounter(tt.metricName, tt.metricValue.(int64))
				case models.GaugeType:
					_ = storage.SetGauge(tt.metricName, tt.metricValue.(float64))
				}
			}

			jsonBytes, err := json.Marshal(&models.MetricsValue{
				ID:    tt.metricName,
				MType: string(tt.metricType),
			})
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewReader(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.wantedStatusCode)
			assert.Equal(t, string(body), tt.wantedBody)
		})
	}
}
