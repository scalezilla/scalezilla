package commands

import (
	"testing"

	"github.com/scalezilla/scalezilla/cluster"
	"github.com/stretchr/testify/assert"
)

func TestCommandsOutputFormat(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(outputFormat("json"))
	assert.ErrorIs(cluster.ErrWrongFormat, outputFormat("j"))
}
