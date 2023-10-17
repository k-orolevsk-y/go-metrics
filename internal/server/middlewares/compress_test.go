package middlewares

import (
	memstorage "github.com/k-orolevsk-y/go-metricts-tpl/internal/server/mem_storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareCompress(t *testing.T) {
	storage := memstorage.NewMem()
	r := setupRouter(storage, zaptest.NewLogger(t).Sugar())

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
}
