package env

import "github.com/caarlos0/env/v6"

var Config struct {
	Address string `env:"ADDRESS"`
}

func Init() error {
	if err := env.Parse(&Config); err != nil {
		return err
	}

	return nil
}
