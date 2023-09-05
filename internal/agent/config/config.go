package config

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config/env"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config/flags"
)

func Init() error {
	if err := env.Init(); err != nil {
		return err
	}
	flags.Init()

	return nil
}

func GetAddress() string {
	if env.Config.Address == "" {
		return flags.Config.Address
	}

	return env.Config.Address
}

func GetReportInterval() int {
	if env.Config.ReportInterval == 0 {
		return flags.Config.ReportInterval
	}

	return env.Config.ReportInterval
}

func GetPollInterval() int {
	if env.Config.PollInterval == 0 {
		return flags.Config.PollInterval
	}

	return env.Config.PollInterval
}
