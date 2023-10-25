package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config"
)

func (bm baseMiddleware) Hash(ctx *gin.Context) {
	secureKey := config.Config.Key
	if secureKey == "" {
		return
	}

	hexHashByClient := ctx.GetHeader("HashSHA256")
	if hexHashByClient == "" {
		return
	}

	hashByClient, err := hex.DecodeString(hexHashByClient)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()

		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		bm.log.Errorf("Error get body for hash check: %s (%T)", err, err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewReader(body)) // Необходимо вернуть body, тк handler-ы потом не смогут прочитать body...

	hash := hmac.New(sha256.New, []byte(secureKey))
	hash.Write(body)

	hashByServer := hash.Sum(nil)
	hexHashByServer := hex.EncodeToString(hashByServer)
	ctx.Header("HashSHA256", hexHashByServer)

	if !bytes.Equal(hashByServer, hashByClient) {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
	}
}
