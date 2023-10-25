package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		url              string
		wantedStatusCode int
	}{
		{
			name:             "Positive gauge",
			method:           http.MethodPost,
			url:              "gauge/test/123.0",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:             "Positive gauge (with big value)",
			method:           http.MethodPost,
			url:              "gauge/test/123.123456789123456789",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:             "Negative gauge (without id)",
			method:           http.MethodPost,
			url:              "gauge",
			wantedStatusCode: http.StatusNotFound,
		},
		{
			name:             "Negative gauge (without params)",
			method:           http.MethodPost,
			url:              "gauge//123.45",
			wantedStatusCode: http.StatusNotFound,
		},
		{
			name:             "Negative gauge (invalid http method)",
			method:           http.MethodGet,
			url:              "gauge/test/123.0",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Negative gauge (invalid value)",
			method:           http.MethodPost,
			url:              "gauge/test/invalid_value",
			wantedStatusCode: http.StatusBadRequest,
		},

		{
			name:             "Positive counter",
			method:           http.MethodPost,
			url:              "counter/test/123",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:             "Negative counter (without id)",
			method:           http.MethodPost,
			url:              "counter",
			wantedStatusCode: http.StatusNotFound,
		},
		{
			name:             "Negative counter (without params)",
			method:           http.MethodPost,
			url:              "counter//10",
			wantedStatusCode: http.StatusNotFound,
		},
		{
			name:             "Negative counter (invalid http method)",
			method:           http.MethodGet,
			url:              "counter/test/123",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Negative counter (big value)",
			method:           http.MethodPost,
			url:              "counter/test/9223372036854775808",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Negative counter (text value)",
			method:           http.MethodPost,
			url:              "counter/test/text_value",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Negative counter (float64 value)",
			method:           http.MethodPost,
			url:              "counter/test/12345.6789",
			wantedStatusCode: http.StatusBadRequest,
		},
	}

	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, fmt.Sprintf("/update/%s", tt.url), nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantedStatusCode)
		})
	}
}

func TestUpdateByBody(t *testing.T) {
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
			body: models.MetricsUpdate{
				ID:    "test",
				MType: string(models.GaugeType),
				Value: getPointerFloat64(123.0),
			},
			wantedBody:       "{\"id\":\"test\",\"type\":\"gauge\",\"value\":123}",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Positive gauge (with big value)",
			method: http.MethodPost,
			body: models.MetricsUpdate{
				ID:    "test",
				MType: string(models.GaugeType),
				Value: getPointerFloat64(123.123456789123456789),
			},
			wantedBody:       "{\"id\":\"test\",\"type\":\"gauge\",\"value\":123.12345678912345}",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Negative gauge (invalid http method)",
			method: http.MethodDelete,
			body: models.MetricsUpdate{
				ID:    "test",
				MType: string(models.GaugeType),
				Value: getPointerFloat64(123.0),
			},
			wantedBody:       "",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Negative gauge (invalid value)",
			method: http.MethodPost,
			body: map[string]any{
				"id":    "test",
				"type":  "gauge",
				"value": "invalid_value",
			},
			wantedBody:       "{\"error\":\"Field value \\\"value\\\" must be float64.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Positive counter",
			method: http.MethodPost,
			body: models.MetricsUpdate{
				ID:    "Test1",
				MType: "counter",
				Delta: getPointerInt64(123),
			},
			wantedBody:       "{\"id\":\"Test1\",\"type\":\"counter\",\"delta\":123}",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:   "Negative counter (invalid http method)",
			method: http.MethodDelete,
			body: models.MetricsUpdate{
				ID:    "Test2",
				MType: "counter",
				Delta: getPointerInt64(123),
			},
			wantedBody:       "",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Negative counter (text value)",
			method: http.MethodPost,
			body: map[string]any{
				"id":    "Test3",
				"type":  "counter",
				"delta": "invalid_value",
			},
			wantedBody:       "{\"error\":\"Field value \\\"delta\\\" must be int64.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Negative counter (float64 value)",
			method: http.MethodPost,
			body: map[string]any{
				"id":    "Test3",
				"type":  "counter",
				"delta": 12345.6789,
			},
			wantedBody:       "{\"error\":\"Field value \\\"delta\\\" must be int64.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			name:             "Negative (without params)",
			method:           http.MethodPost,
			body:             models.MetricsUpdate{},
			wantedBody:       "{\"error\":\"Field validation for \\\"ID\\\" failed on the 'required' tag.\"}",
			wantedStatusCode: http.StatusBadRequest,
		},
	}

	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(&tt.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.method, "/update/", bytes.NewReader(jsonBytes))
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
