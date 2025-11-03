package osdiscovery

import "regexp"

var (
	// procMemoryPath is hold content of memory informations
	procMemoryPath = "/proc/meminfo"

	// reMemTotal match MemTotal field
	reMemTotal = regexp.MustCompile(`^MemTotal:\s+(\d+)(.*)`)

	// reMemAvailable match MemAvailable field
	reMemAvailable = regexp.MustCompile(`^MemAvailable:\s+(\d+)(.*)`)
)

// Memory holds memory informations
type Memory struct {
	// procMemoryPath holds content of memory informations
	procMemoryPath string

	// Total is the number return by MemTotal field
	Total uint `json:"total"`

	// Available is the number return by MemAvailable field
	Available uint `json:"available"`
}
