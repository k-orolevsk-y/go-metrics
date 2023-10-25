package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	memstorage "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
)

func TestUpdates(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		body             any
		wantedBody       string
		wantedStatusCode int
	}{
		{
			name:   "Positive gauge",
			method: http.MethodPost,
			body: []models.MetricsUpdate{
				{
					ID:    "test",
					MType: string(models.GaugeType),
					Value: getPointerFloat64(123.0),
				},
			},
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Positive gauge (with big value)",
			method: http.MethodPost,
			body: []models.MetricsUpdate{
				{
					ID:    "test",
					MType: string(models.GaugeType),
					Value: getPointerFloat64(123.123456789123456789),
				},
			},
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Negative gauge (invalid value)",
			method: http.MethodPost,
			body: []map[string]any{
				{
					"id":    "test",
					"type":  "gauge",
					"value": "invalid_value",
				},
			},
			wantedBody:       "{\"error\":\"Field value \\\"value\\\" must be float64.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Positive counter",
			method: http.MethodPost,
			body: []models.MetricsUpdate{
				{
					ID:    "Test1",
					MType: "counter",
					Delta: getPointerInt64(123),
				},
			},
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Positive updating the counter several times",
			method: http.MethodPost,
			body: []models.MetricsUpdate{
				{
					ID:    "TestCounterN",
					MType: "counter",
					Delta: getPointerInt64(100),
				},
				{
					ID:    "TestCounterN",
					MType: "counter",
					Delta: getPointerInt64(200),
				},
				{
					ID:    "TestCounterN",
					MType: "counter",
					Delta: getPointerInt64(300),
				},
			},
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Negative counter (text value)",
			method: http.MethodPost,
			body: []map[string]any{
				{
					"id":    "Test3",
					"type":  "counter",
					"delta": "invalid_value",
				},
			},
			wantedBody:       "{\"error\":\"Field value \\\"delta\\\" must be int64.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Negative counter (float64 value)",
			method: http.MethodPost,
			body: []map[string]any{
				{
					"id":    "Test3",
					"type":  "counter",
					"delta": 12345.6789,
				},
			},
			wantedBody:       "{\"error\":\"Field value \\\"delta\\\" must be int64.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Negative (without params)",
			method:           http.MethodPost,
			body:             nil,
			wantedStatusCode: http.StatusOK,
		},
	}

	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(&tt.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.method, "/updates/", bytes.NewReader(jsonBytes))
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
