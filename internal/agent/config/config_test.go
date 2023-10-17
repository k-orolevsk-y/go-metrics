package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name string

		args map[string]map[string]string

		wantedAddress        string
		wantedPollInterval   int
		wantedReportInterval int
	}{
		{
			name: "Positive (without env & flags)",

			args: map[string]map[string]string{
				"env":  {},
				"flag": {},
			},

			wantedAddress:        "localhost:8080",
			wantedPollInterval:   2,
			wantedReportInterval: 10,
		},
		{
			name: "Positive (with env & flags)",

			args: map[string]map[string]string{
				"env": {
					"ADDRESS":         "localhost:8081",
					"REPORT_INTERVAL": "1",
					"POLL_INTERVAL":   "2",
				},
				"flag": {
					"a": "localhost:8082",
					"r": "3",
					"p": "4",
				},
			},

			wantedAddress:        "localhost:8081",
			wantedPollInterval:   2,
			wantedReportInterval: 1,
		},
		{
			name: "Positive (with flags)",

			args: map[string]map[string]string{
				"flag": {
					"a": "localhost:8082",
					"r": "3",
					"p": "4",
				},
			},

			wantedAddress:        "localhost:8082",
			wantedPollInterval:   4,
			wantedReportInterval: 3,
		},
		{
			name: "Positive (with env)",

			args: map[string]map[string]string{
				"env": {
					"ADDRESS":         "localhost:9090",
					"REPORT_INTERVAL": "100",
					"POLL_INTERVAL":   "200",
				},
			},

			wantedAddress:        "localhost:9090",
			wantedPollInterval:   200,
			wantedReportInterval: 100,
		},
	}

	Load()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			os.Args = []string{"main"}

			if tt.args["env"] != nil && len(tt.args["env"]) > 0 {
				for k, v := range tt.args["env"] {
					require.NoError(t, os.Setenv(k, v))
				}
			}

			if tt.args["flag"] != nil && len(tt.args["flag"]) > 0 {
				for k, v := range tt.args["flag"] {
					require.NoError(t, flag.Set(k, v))
				}
			}

			require.NoError(t, Parse())

			assert.Equal(t, tt.wantedAddress, Config.Address)
			assert.Equal(t, tt.wantedPollInterval, Config.PollInterval)
			assert.Equal(t, tt.wantedReportInterval, Config.ReportInterval)
		})
	}
}
