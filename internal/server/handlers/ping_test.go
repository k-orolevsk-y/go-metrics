package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		wantedBody       string
		wantedStatusCode int
	}{
		{
			name:             "Connected",
			err:              nil,
			wantedBody:       "",
			wantedStatusCode: http.StatusOK,
		},
		{
			name:             "Not connected",
			err:              errors.New("test error"),
			wantedBody:       "{\"error\":\"test error\"}",
			wantedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStorage(ctrl)
			m.EXPECT().GetMiddleware().Return(func(_ *gin.Context) {})
			m.EXPECT().Ping(gomock.Any()).Return(tt.err)

			r := setupRouter(m, zaptest.NewLogger(t).Sugar())

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, res.StatusCode, tt.wantedStatusCode)
			if tt.wantedBody != "" {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				require.Equal(t, string(body), tt.wantedBody)
			}
		})
	}
}
