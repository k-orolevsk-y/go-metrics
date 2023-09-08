package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var Data struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func Init() {
	flag.StringVar(&Data.Address, "a", "localhost:8080", "server address")
	flag.IntVar(&Data.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Data.PollInterval, "p", 2, "poll interval")
}

func Parse() error {
	flag.Parse()

	return env.Parse(&Data)
}
