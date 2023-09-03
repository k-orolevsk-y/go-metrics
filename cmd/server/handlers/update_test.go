package handlers

import (
	stor "github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateGauge(t *testing.T) {
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
			name: "Positive",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/gauge/test/123.0",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Positive (with big value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/gauge/test/123.123456789123456789",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Negative (without params)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/gauge/",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative (invalid http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "/update/gauge/test/123.0",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative (invalid value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/gauge/test/invalid_value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	storage := stor.NewMem()
	handlerUpdateGauge := UpdateGauge(&storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, nil)

			w := httptest.NewRecorder()
			handlerUpdateGauge(w, request)

			res := w.Result()

			require.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestUpdateCounter(t *testing.T) {
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
			name: "Positive",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/123",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Negative (without params)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative (invalid http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "/update/counter/test/123",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative (big value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/9223372036854775808",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative (text value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/text_value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative (float64 value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/12345.6789",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	storage := stor.NewMem()
	handlerUpdateCounter := UpdateCounter(&storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, nil)

			w := httptest.NewRecorder()
			handlerUpdateCounter(w, request)

			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}
