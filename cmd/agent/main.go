package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics_updater"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
	"time"
)

func main() {
	sugarLogger, err := logger.New()
	if err != nil {
		panic(err)
	}

	defer func(log logger.Logger) {
		if err = log.Sync(); err != nil {
			panic(err)
		}
	}(sugarLogger)

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
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
