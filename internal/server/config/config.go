package config

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config/env"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/config/flags"
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
