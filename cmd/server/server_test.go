package main

import (
	"fmt"
	stor "github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	r := setupRouter(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.argHTTPMethod, "/", nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		httpMethod string
		path       string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "Positive gauge",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge/test/123.0",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Positive gauge (with big value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge/test/123.123456789123456789",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Negative gauge (without params)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge//123.45",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative gauge (invalid http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "gauge/test/123.0",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative gauge (invalid value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge/test/invalid_value",
			},
			wantStatusCode: http.StatusBadRequest,
		},

		{
			name: "Positive counter",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/123",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Negative counter (without params)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter//10",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative counter (invalid http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "counter/test/123",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative counter (big value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/9223372036854775808",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative counter (text value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/text_value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative counter (float64 value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/12345.6789",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	storage := stor.NewMem()
	r := setupRouter(&storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.httpMethod, fmt.Sprintf("/update/%s", tt.args.path), nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestValues(t *testing.T) {
	type args struct {
		httpMethod string
	}
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive",
			args: args{
				httpMethod: http.MethodGet,
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Negative (POST)",
			args: args{
				httpMethod: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	storage := stor.NewMem()
	r := setupRouter(&storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.httpMethod, "/", nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
