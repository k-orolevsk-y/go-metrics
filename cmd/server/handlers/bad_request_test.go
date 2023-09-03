package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequest(t *testing.T) {
	tests := []struct {
		name           string
		argHttpMethod  string
		wantStatusCode int
	}{
		{
			name:           "Positive GET",
			argHttpMethod:  http.MethodGet,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Positive POST",
			argHttpMethod:  http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Positive DELETE",
			argHttpMethod:  http.MethodDelete,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.argHttpMethod, "/", nil)

			w := httptest.NewRecorder()
			BadRequest(w, request)

			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}
