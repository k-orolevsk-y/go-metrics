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

	var store metrics.RuntimeMetrics
	store.Init()

	client := resty.New()
	updater := metricsupdater.New(client, &store)

	for {
		go func() {
			err := store.Update()
			if err != nil {
				panic(err)
			}
		}()

		updater.UpdateMetrics()
	}
}
