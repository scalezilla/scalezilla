package osdiscovery

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsDiscovery_memory(t *testing.T) {
	assert := assert.New(t)

	t.Run("test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		testdata := "testdata/os"
		err = filepath.Walk(filepath.Join(workingDir, testdata),
			func(file string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.Contains(file, "meminfo") {
					mem := newMemory()
					mem.procMemoryPath = file
					mem.memory()

					if strings.Contains(file, "empty") {
						assert.Equal(mem.Total, uint(0))
					} else {
						assert.Greater(mem.Total, uint(0))
						assert.Greater(mem.Available, uint(0))
					}
				}
				return nil
			})
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		mem := newMemory()
		mem.procMemoryPath = "testdata/os/absent/memory"
		mem.memory()

		assert.Equal(mem.Total, uint(0))
	})

	t.Run("fake_reader", func(t *testing.T) {
		mem := newMemory()
		reader := &fakeReader{}
		mem.memoryReader(reader)

		assert.Equal(mem.Total, uint(0))
	})
}
