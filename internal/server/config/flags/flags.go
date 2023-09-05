package flags

import "flag"

var Config struct {
	Address string
}

func Init() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
	flag.Parse()
}
