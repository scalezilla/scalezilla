package osdiscovery

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsDiscovery_os_info(t *testing.T) {
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
				if !info.IsDir() {
					osi := newOsInfo()
					osi.osReleasePath = file
					osi.osInfo()

					if strings.Contains(file, "empty") {
						assert.Empty(osi.Name)
						assert.Empty(osi.Version)
					} else {
						assert.NotNil(osi.Name)
						assert.NotNil(osi.Version)
					}
				}
				return nil
			})
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		osi := newOsInfo()
		osi.osReleasePath = "testdata/os/absent/os-release"
		osi.osInfo()

		assert.Empty(osi.Name)
	})

	t.Run("fake_reader", func(t *testing.T) {
		osi := newOsInfo()
		reader := &fakeReader{}
		osi.osInfoReader(reader)

		assert.Equal("ubuntu", osi.Name)
		assert.Empty(osi.Version)
	})
}
