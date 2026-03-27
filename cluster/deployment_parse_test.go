package cluster

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCluster_deployment_parse(t *testing.T) {
	assert := assert.New(t)

	t.Run("file_syntax_validate_test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		configDir := "testdata/deployments"
		files, err := os.ReadDir(configDir)
		assert.Nil(err)

		for _, file := range files {
			configFile := filepath.Join(workingDir, configDir, file.Name())

			err := parseDeploymentFileSyntax(configFile)
			if strings.Contains(file.Name(), "malformed") {
				assert.Error(err)
			} else {
				assert.Nil(err)
			}
		}
	})

	t.Run("parse_deployment", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		configDir := "testdata/deployments"
		files, err := os.ReadDir(configDir)
		assert.Nil(err)

		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		for _, file := range files {
			exceptions := []string{"malformed", "error_container_image"}
			if !slices.Contains(exceptions, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))) {
				configFile := filepath.Join(workingDir, configDir, file.Name())
				data, err := os.ReadFile(configFile)
				require.NoError(t, err)

				_, err = cluster.parseDeployment(data)
				if strings.Contains(file.Name(), "success") {
					assert.Nil(err, file.Name())
				} else {
					assert.Error(err, file.Name())
				}
			}
		}
	})
}
