package osdiscovery

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsDiscovery_cpu(t *testing.T) {
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
				if !info.IsDir() && strings.Contains(file, "cpu") {
					cpu := newCPU()
					cpu.procCPUPath = file
					cpu.cpu()

					if strings.Contains(file, "empty") {
						assert.Equal(cpu.CPU, uint32(0))
					} else {
						assert.Greater(cpu.CPU, uint32(0))
						assert.Greater(cpu.Cores, uint32(0))
					}
				}
				return nil
			})
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		cpu := newCPU()
		cpu.procCPUPath = "testdata/os/absent/cpu"
		cpu.cpu()

		assert.Equal(cpu.CPU, uint32(0))
	})

	t.Run("fake_reader", func(t *testing.T) {
		cpu := newCPU()
		reader := &fakeReader{}
		cpu.cpuReader(reader)

		assert.Equal(cpu.CPU, uint32(0))
	})
}
