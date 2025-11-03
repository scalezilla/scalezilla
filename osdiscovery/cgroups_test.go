package osdiscovery

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsDiscovery_cgroups(t *testing.T) {
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
				if !info.IsDir() && strings.Contains(file, "filesystems") {
					c := newCgroups()
					c.procCgroupsPath = file
					c.cgroups()

					if strings.Contains(file, "empty") {
						assert.Equal(c.Version, uint(0))
					} else {
						assert.Greater(c.Version, uint(0))
					}
				}
				return nil
			})
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		c := newCgroups()
		c.procCgroupsPath = "testdata/os/absent/filesystems"
		c.cgroups()

		assert.Equal(c.Version, uint(0))
	})

	t.Run("fake_reader", func(t *testing.T) {
		c := newCgroups()
		reader := &fakeReader{}
		c.cgroupsReader(reader)

		assert.Equal(c.Version, uint(0))
	})
}
