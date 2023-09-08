package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var Data struct {
	Address string `env:"ADDRESS"`
}

func Init() {
	flag.StringVar(&Data.Address, "a", "localhost:8080", "server address")
}

func Parse() error {
	flag.Parse()

	return env.Parse(&Data)
}
