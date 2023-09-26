package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics_updater"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Panic loading config: %s", err)
	}

	client := resty.New()
	store := metrics.NewRuntimeMetrics()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.Config.PollInterval))

			if err = store.Update(); err != nil {
				sugarLogger.Panicf("Failed to update metrics: %s", err)
			}
		}
	}()

	updater := metricsupdater.New(client, store, sugarLogger)
	for {
		time.Sleep(time.Second * time.Duration(config.Config.ReportInterval))
		updater.UpdateMetrics()
	}
}
