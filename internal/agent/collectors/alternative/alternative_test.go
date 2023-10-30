package alternative

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlternativeMetricsCollector(t *testing.T) {
	collector := NewAlternativeCollector()
	require.NoError(t, collector.Collect())
}
