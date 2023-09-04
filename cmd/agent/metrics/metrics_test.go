package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntimeMetrics(t *testing.T) {
	var runtimeMetrics RuntimeMetrics
	runtimeMetrics.Init()

	assert.NoError(t, runtimeMetrics.Update())
}
