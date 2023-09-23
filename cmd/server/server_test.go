package main

import (
	"fmt"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
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

	r := setupRouter(nil, zaptest.NewLogger(t).Sugar())

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
			name: "Negative gauge (without id)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge",
			},
			wantStatusCode: http.StatusNotFound,
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
			name: "Negative counter (without id)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter",
			},
			wantStatusCode: http.StatusNotFound,
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

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, zaptest.NewLogger(t).Sugar())

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

func TestValue(t *testing.T) {
	type args struct {
		name       string
		metricType storage.MetricType
		value      interface{}
	}
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive (gauge)",
			args: args{
				name:       "TestGauge",
				metricType: storage.GaugeType,
				value:      10.5,
			},
			want: want{
				body:       "10.5",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Positive (counter)",
			args: args{
				name:       "TestCounter",
				metricType: storage.CounterType,
				value:      int64(10500),
			},
			want: want{
				body:       "10500",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Negative (gauge)",
			args: args{
				metricType: storage.GaugeType,
			},
			want: want{
				body:       "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative (counter)",
			args: args{
				metricType: storage.CounterType,
			},
			want: want{
				body:       "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative (invalid metric type)",
			args: args{
				metricType: storage.MetricType("invalid"),
			},
			want: want{
				body:       "",
				statusCode: http.StatusNotFound,
			},
		},
	}

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, zaptest.NewLogger(t).Sugar())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.name != "" {
				switch tt.args.metricType {
				case storage.CounterType:
					memStorage.AddCounter(tt.args.name, tt.args.value.(int64))
				case storage.GaugeType:
					memStorage.SetGauge(tt.args.name, tt.args.value.(float64))
				}
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/value/%s/%s/", tt.args.metricType, tt.args.name), nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(body), tt.want.body)
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

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, zaptest.NewLogger(t).Sugar())

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
