package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics_updater"
)

func main() {
	config.Load()
	if err := config.Parse(); err != nil {
		panic(err)
	}

	client := resty.New()
	store := metrics.NewRuntimeMetrics()

	updater := metricsupdater.New(client, store)
	go func() {
		for {
			err := store.Update()
			if err != nil {
				panic(err)
			}
		}
	}()

	for {
		updater.UpdateMetrics()
	}
}
