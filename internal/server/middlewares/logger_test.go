package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	memstorage "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
)

func TestMiddlewareLogger(t *testing.T) {
	storage := memstorage.NewMem()

	buf := new(bytes.Buffer)
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.DebugLevel),
		zap.AddCaller(),
	)

	r := setupRouter(storage, log.Sugar())

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Contains(t, buf.String(), "GET Request")
}
