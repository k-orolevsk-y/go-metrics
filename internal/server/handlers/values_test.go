package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
)

func TestValues(t *testing.T) {
	tests := []struct {
		name   string
		method string

		wantedStatusCode  int
		wantedContentType string
	}{
		{
			name:   "Positive",
			method: http.MethodGet,

			wantedStatusCode:  http.StatusOK,
			wantedContentType: "text/html; charset=utf-8",
		},
		{
			name:   "Negative (POST)",
			method: http.MethodPost,

			wantedStatusCode: http.StatusBadRequest,
		},
	}

	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/", nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantedStatusCode, res.StatusCode)
			assert.Equal(t, tt.wantedContentType, res.Header.Get("Content-Type"))
		})
	}
}
