package gopsutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGopsutilCollector(t *testing.T) {
	collector := NewGopsutilCollector()
	require.NoError(t, collector.Collect())
}
