package metricsupdater

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

var (
	retries = []int{1, 3, 5}

	ErrorNotNeedHash       = errors.New("not need hash")
	ErrorInvalidStatusCode = errors.New("invalid status code")
)

type (
	Updater struct {
		client *resty.Client
		col    collector
		log    logger.Logger
	}

	collector interface {
		GetMetrics() []metrics.Metric
	}
)

func New(client *resty.Client, col collector, log logger.Logger) *Updater {
	return &Updater{
		client: client,
		col:    col,
		log:    log,
	}
}

func (u Updater) UpdateMetrics() {
	currentMetrics := u.col.GetMetrics()
	if err := u.updateMetrics(currentMetrics); err != nil {
		u.log.Errorf("Failed to update collectors: %s (%T)", err, err)
	}
}

func (u Updater) updateMetrics(metricForUpdate []metrics.Metric) error {
	url := fmt.Sprintf("http://%s/updates", config.Config.Address)

	req, err := u.compileRequest(metricForUpdate)
	if err != nil {
		return err
	}

	for _, timeSleep := range retries {
		resp, err := req.Post(url)
		if err != nil {
			u.log.Errorf("Failed to send collectors to server: %s. Retrying after %ds...", err, timeSleep)
			time.Sleep(time.Duration(timeSleep) * time.Second)

			continue
		}

		if resp.StatusCode() != http.StatusOK {
			return ErrorInvalidStatusCode
		} else {
			return nil
		}
	}

	return err
}

func (u Updater) compileRequest(metricsForRequest []metrics.Metric) (*resty.Request, error) {
	req := u.client.R().
		SetBody(metricsForRequest)

	hash, err := u.hashMetrics(metricsForRequest)
	if err != nil {
		if !errors.Is(err, ErrorNotNeedHash) {
			return nil, err
		}
	} else {
		req.SetHeader("HashSHA256", hash)
	}

	return req, nil
}

func (u Updater) hashMetrics(metricsForHash []metrics.Metric) (string, error) {
	secureKey := config.Config.Key
	if secureKey == "" {
		return "", ErrorNotNeedHash
	}

	bodyBytes, err := json.Marshal(metricsForHash)
	if err != nil {
		return "", fmt.Errorf("hashMetrics: %w", err)
	}

	hash := hmac.New(sha256.New, []byte(secureKey))
	hash.Write(bodyBytes)

	hashed := hash.Sum(nil)
	hexHashed := hex.EncodeToString(hashed)

	return hexHashed, nil
}
