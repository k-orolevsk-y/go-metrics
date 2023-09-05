package flags

import "flag"

var Config struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func Init() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
	flag.IntVar(&Config.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Config.PollInterval, "p", 2, "poll interval")

	flag.Parse()
}
