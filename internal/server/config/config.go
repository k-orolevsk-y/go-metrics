package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var Config struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func Load() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
	flag.Int64Var(&Config.StoreInterval, "i", 0, "store interval in seconds")
	flag.StringVar(&Config.FileStoragePath, "f", "tmp/metrics-db.json", "json file mem_storage path")
	flag.BoolVar(&Config.Restore, "r", true, "whether to load old values from a file")
	flag.StringVar(&Config.DatabaseDSN, "d", "", "postgresql dsn")
}

func Parse() error {
	flag.Parse()

	return env.Parse(&Config)
}
