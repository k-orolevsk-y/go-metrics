package main

import (
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics_updater"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

func main() {
	sugarLogger, err := logger.New()
	if err != nil {
		panic(err)
	}
	sugarLogger.Debugf("The logger has been successfully initialized and configured.")

	defer func() {
		if err = sugarLogger.Sync(); err != nil {
			panic(err)
		}
	}()

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
	}
	sugarLogger.Debugf("The config was successfully received and configured.")

	client := resty.New()
	store := metrics.NewRuntimeMetrics()

	go func() {
		sugarLogger.Debugf("Metrics collector successfully initialized.")
		for {
			time.Sleep(time.Second * time.Duration(config.Config.PollInterval))

			if err = store.Update(); err != nil {
				sugarLogger.Panicf("Failed to update metrics: %s", err)
			}
		}
	}()

	updater := metricsupdater.New(client, store, sugarLogger)
	sugarLogger.Debugf("Metrics updater successfully initialized.")

	for {
		time.Sleep(time.Second * time.Duration(config.Config.ReportInterval))
		updater.UpdateMetrics()
	}
}
