package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateContentType(t *testing.T) {
	tests := []struct {
		name               string
		contentTypeRequest string
		contentTypeNeed    string
		withoutContentType bool
		want               bool
	}{
		{
			name:               "Positive text/plain",
			contentTypeRequest: "text/plain",
			contentTypeNeed:    "text/plain",
			withoutContentType: true,
			want:               true,
		},
		{
			name:               "Positive json/application",
			contentTypeRequest: "json/application",
			contentTypeNeed:    "json/application",
			withoutContentType: true,
			want:               true,
		},
		{
			name:               "Positive without request content-type",
			contentTypeRequest: "",
			contentTypeNeed:    "text/plain",
			withoutContentType: true,
			want:               true,
		},
		{
			name:               "Negative without request content-type",
			contentTypeRequest: "",
			contentTypeNeed:    "json/application",
			withoutContentType: false,
			want:               false,
		},
		{
			name:               "Negative json/application & text/plain",
			contentTypeRequest: "json/application",
			contentTypeNeed:    "text/plain",
			withoutContentType: true,
			want:               false,
		},
	}

	gin.SetMode(gin.ReleaseMode)
	bh := baseHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.contentTypeRequest != "" {
				ctx.Request.Header.Set("Content-Type", tt.contentTypeRequest)
			}

			assert.Equal(t, tt.want, bh.validateContentType(ctx, tt.contentTypeNeed, tt.withoutContentType))
		})
	}
}

func getRandomFloat64() *float64 {
	v := rand.Float64()
	return &v
}

func getRandomInt64() *int64 {
	v := rand.Int63()
	return &v
}

func TestValidateAndShouldBindJSON(t *testing.T) {
	tests := []struct {
		name string

		obj         models.MetricsUpdate
		body        models.MetricsUpdate
		withoutBody bool

		wantedErrorResponse *models.ErrorResponse
		wantedStatusCode    int
		wantedErr           bool
	}{
		{
			name: "Positive gauge",

			obj: models.MetricsUpdate{},
			body: models.MetricsUpdate{
				ID:    "TestGauge",
				MType: string(models.GaugeType),
				Value: getRandomFloat64(),
			},

			wantedErrorResponse: nil,
			wantedStatusCode:    0,
			wantedErr:           false,
		},
		{
			name: "Positive counter",

			obj: models.MetricsUpdate{},
			body: models.MetricsUpdate{
				ID:    "TestCounter",
				MType: string(models.CounterType),
				Delta: getRandomInt64(),
			},

			wantedErrorResponse: nil,
			wantedStatusCode:    0,
			wantedErr:           false,
		},
		{
			name: "Negative without params",

			obj:  models.MetricsUpdate{},
			body: models.MetricsUpdate{},

			wantedErrorResponse: &models.ErrorResponse{Error: "Field validation for \"ID\" failed on the 'required' tag."},
			wantedStatusCode:    http.StatusBadRequest,
			wantedErr:           true,
		},
		{
			name: "Negative invalid type",

			obj: models.MetricsUpdate{},
			body: models.MetricsUpdate{
				ID:    "NTest",
				MType: "heh",
			},

			wantedErrorResponse: &models.ErrorResponse{Error: "Field validation for \"MType\" failed on the 'oneof=counter gauge' tag."},
			wantedStatusCode:    http.StatusBadRequest,
			wantedErr:           true,
		},
		{
			name: "Negative without value (gauge)",

			obj: models.MetricsUpdate{},
			body: models.MetricsUpdate{
				ID:    "NTest",
				MType: string(models.GaugeType),
				Delta: getRandomInt64(),
			},

			wantedErrorResponse: &models.ErrorResponse{Error: "Field validation for \"Value\" failed on the 'required_if=MType gauge' tag."},
			wantedStatusCode:    http.StatusBadRequest,
			wantedErr:           true,
		},
		{
			name: "Negative without value (counter)",

			obj: models.MetricsUpdate{},
			body: models.MetricsUpdate{
				ID:    "NTest",
				MType: string(models.CounterType),
				Value: getRandomFloat64(),
			},

			wantedErrorResponse: &models.ErrorResponse{Error: "Field validation for \"Delta\" failed on the 'required_if=MType counter' tag."},
			wantedStatusCode:    http.StatusBadRequest,
			wantedErr:           true,
		},
		{
			name: "Negative without request body",

			withoutBody: true,

			wantedErrorResponse: &models.ErrorResponse{Error: "Request body not provided."},
			wantedStatusCode:    http.StatusBadRequest,
			wantedErr:           true,
		},
	}

	gin.SetMode(gin.ReleaseMode)
	bh := baseHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if !tt.withoutBody {
				jsonBytes, err := json.Marshal(&tt.body)
				require.NoError(t, err)

				body = bytes.NewReader(jsonBytes)
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", body)

			response, statusCode, err := bh.validateAndShouldBindJSON(ctx, &tt.obj)

			assert.Equal(t, response, tt.wantedErrorResponse)
			assert.Equal(t, statusCode, tt.wantedStatusCode)
			if tt.wantedErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.obj, tt.body)
			}
		})
	}
}
