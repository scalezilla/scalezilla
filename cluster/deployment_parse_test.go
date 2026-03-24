package cluster

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_parseDeployment_file_syntax(t *testing.T) {
	assert := assert.New(t)

	t.Run("validate_test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		configDir := "testdata/deployments"
		files, err := os.ReadDir(configDir)
		assert.Nil(err)

		for _, file := range files {
			configFile := filepath.Join(workingDir, configDir, file.Name())

			err := parseDeploymentFileSyntax(configFile)
			if strings.Contains(file.Name(), "success") {
				assert.Nil(err)
			} else {
				assert.Error(err)
			}
		}
	})
}
