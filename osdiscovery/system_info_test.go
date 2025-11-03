package osdiscovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsDiscovery_system_info(t *testing.T) {
	assert := assert.New(t)

	t.Run("test_data", func(t *testing.T) {
		s := NewSystemInfo()
		assert.Greater(s.CPU.CPU, uint32(0))
	})
}
