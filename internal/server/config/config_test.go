package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]map[string]string
		wantedAddress string
	}{
		{
			name: "Positive (without env & flags)",
			args: map[string]map[string]string{
				"env":  {},
				"flag": {},
			},

			wantedAddress: "localhost:8080",
		},
		{
			name: "Positive (with env & flags)",
			args: map[string]map[string]string{
				"env": {
					"ADDRESS": "localhost:8081",
				},
				"flag": {
					"a": "localhost:8082",
				},
			},
			wantedAddress: "localhost:8081",
		},
		{
			name: "Positive (with flags)",
			args: map[string]map[string]string{
				"flag": {
					"a": "localhost:8082",
				},
			},
			wantedAddress: "localhost:8082",
		},
		{
			name: "Positive (with env)",
			args: map[string]map[string]string{
				"env": {
					"ADDRESS": "localhost:9090",
				},
			},
			wantedAddress: "localhost:9090",
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
		})
	}
}
