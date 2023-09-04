package flags

import "flag"

var Data struct {
	Host string
}

func Init() {
	flag.StringVar(&Data.Host, "a", "localhost:8080", "server address")
	flag.Parse()
}
