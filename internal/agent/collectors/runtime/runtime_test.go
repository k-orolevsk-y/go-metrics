package runtime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeCollector(t *testing.T) {
	collector := NewRuntimeCollector()
	require.NoError(t, collector.Collect())
}
