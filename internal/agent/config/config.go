package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func Load() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
	flag.IntVar(&Config.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Config.PollInterval, "p", 2, "poll interval")
	flag.StringVar(&Config.Key, "k", "", "key for hash")
	flag.IntVar(&Config.RateLimit, "l", 2, "rate limit for worker pool")
}

func Parse() error {
	flag.Parse()

	return env.Parse(&Config)
}
