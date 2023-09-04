package flags

import "flag"

var Data struct {
	ServerHost     string
	ReportInterval int
	PollInterval   int
}

func Init() {
	flag.StringVar(&Data.ServerHost, "a", "localhost:8080", "server address")
	flag.IntVar(&Data.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Data.PollInterval, "p", 2, "poll interval")

	flag.Parse()
}
