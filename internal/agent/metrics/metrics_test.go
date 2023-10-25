package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeMetrics(t *testing.T) {
	runtimeMetrics := NewRuntimeMetrics()

	assert.NoError(t, runtimeMetrics.Update())
}
