package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var Config struct {
	Address string `env:"ADDRESS"`
}

func Load() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
}

func Parse() error {
	flag.Parse()

	return env.Parse(&Config)
}
