package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntimeMetrics(t *testing.T) {
	var runtimeMetrics RuntimeMetrics
	runtimeMetrics.New()

	assert.NoError(t, runtimeMetrics.Renew())
}
