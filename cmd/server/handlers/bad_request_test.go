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
		argHTTPMethod  string
		wantStatusCode int
	}{
		{
			name:           "Positive GET",
			argHTTPMethod:  http.MethodGet,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Positive POST",
			argHTTPMethod:  http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Positive DELETE",
			argHTTPMethod:  http.MethodDelete,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.argHTTPMethod, "/", nil)

			w := httptest.NewRecorder()
			BadRequest(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}
