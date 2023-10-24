package memstorage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemGetCounter(t *testing.T) {
	tests := []struct {
		name string

		metricName  string
		metricValue int64

		wantedErr bool
	}{
		{
			name: "Positive",

			metricName:  "Test",
			metricValue: 12345,

			wantedErr: false,
		},
		{
			name: "Negative",

			metricName:  "",
			metricValue: 0,

			wantedErr: true,
		},
	}

	storage := NewMem()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.metricName != "" && tt.metricValue != 0 {
				_ = storage.AddCounter(tt.metricName, &tt.metricValue)
			}

			got, err := storage.GetCounter(tt.metricName)
			if tt.wantedErr {
				assert.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.metricValue, *got)
			}
		})
	}
}

func TestMem_GetGauge(t *testing.T) {
	tests := []struct {
		name string

		metricName  string
		metricValue float64

		wantedErr bool
	}{
		{
			name: "Positive",

			metricName:  "Test",
			metricValue: 123.45,

			wantedErr: false,
		},
		{
			name: "Negative",

			metricName:  "",
			metricValue: 0,

			wantedErr: true,
		},
	}

	storage := NewMem()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.metricName != "" && tt.metricValue != 0 {
				_ = storage.SetGauge(tt.metricName, &tt.metricValue)
			}

			got, err := storage.GetGauge(tt.metricName)
			if tt.wantedErr {
				assert.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.metricValue, *got)
			}
		})
	}
}

func TestMem_normalizeName(t *testing.T) {
	tests := []struct {
		name       string
		metricName string
		want       string
	}{
		{
			name:       "Positive",
			metricName: " Test name ",
			want:       "Test name",
		},
		{
			name:       "Large name",
			metricName: " Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur ",
			want:       "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur",
		},
	}

	storage := NewMem()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, storage.normalizeName(tt.metricName))
		})
	}
}
