package cluster

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)

	t.Run("validate_test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		configDir := "testdata/config"
		files, err := os.ReadDir(configDir)
		assert.Nil(err)

		for _, file := range files {
			configFile := filepath.Join(workingDir, configDir, file.Name())
			config := ClusterInitialConfig{
				ConfigFile: configFile,
				Test:       true,
			}

			_, err = NewCluster(config)
			if strings.Contains(file.Name(), "success") {
				assert.Nil(err)
			} else {
				assert.Error(err)
			}
		}
	})

	t.Run("err_no_such_file", func(t *testing.T) {
		fakeFile := filepath.Join(os.TempDir(), "no_such_file")
		config := ClusterInitialConfig{ConfigFile: fakeFile}
		_, err := NewCluster(config)
		assert.Error(err)
	})
}
