package main

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/collectors"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
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

	collector := collectors.NewCollector(sugarLogger)
	go collector.Run(context.Background())

	client := resty.New()
	updater := metricsupdater.New(client, collector, sugarLogger)

	sugarLogger.Debugf("Metrics updater successfully initialized.")
	for {
		time.Sleep(time.Second * time.Duration(config.Config.ReportInterval))
		updater.UpdateMetrics()
	}
}
