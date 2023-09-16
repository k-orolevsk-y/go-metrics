package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics_updater"
	"time"
)

func main() {
	config.Load()
	if err := config.Parse(); err != nil {
		panic(err)
	}

	client := resty.New()
	store := metrics.NewRuntimeMetrics()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.Config.PollInterval))

			err := store.Update()
			if err != nil {
				panic(err)
			}
		}
	}()

	updater := metricsupdater.New(client, store)
	for {
		updater.UpdateMetrics()
		time.Sleep(time.Second * time.Duration(config.Config.ReportInterval))
	}
}
