package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
)

func TestBadRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantStatusCode int
	}{
		{
			name:           "Positive POST",
			method:         http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Positive DELETE",
			method:         http.MethodDelete,
			wantStatusCode: http.StatusBadRequest,
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

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}
