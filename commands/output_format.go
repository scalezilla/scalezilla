package commands

import (
	"slices"

	"github.com/scalezilla/scalezilla/cluster"
)

// outputFormat returns nil if n error found
func outputFormat(o string) error {
	formats := []string{"table", "json"}
	if !slices.Contains(formats, o) {
		return cluster.ErrWrongFormat
	}
	return nil
}
