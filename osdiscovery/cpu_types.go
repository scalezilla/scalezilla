package osdiscovery

import "regexp"

var (
	// procCPUPath is hold content of cpu informations
	procCPUPath = "/proc/cpuinfo"

	// reProcessor match processor field
	reProcessor = regexp.MustCompile(`^processor\s+:\s+(\d+)`)

	// reCores match cpu cores field
	reCores = regexp.MustCompile(`^cpu cores\s+:\s+(\d+)`)

	// reFrequency match cpu MHz field
	reFrequency = regexp.MustCompile(`^cpu MHz\s+:\s+(.*)`)

	// reCapabilities match cpu flags field
	reCapabilities = regexp.MustCompile(`^flags\s+:\s+(.*)`)

	// reVendorID match vendor_id field
	reVendorID = regexp.MustCompile(`^vendor_id\s+:\s+(.*)`)

	// reModel match model name field
	reModel = regexp.MustCompile(`^model name\s+:\s+(.*)`)
)

// CPU holds cpu informations
type CPU struct {
	// procCPUPath holds content of cpu informations
	procCPUPath string

	// CPU is the total number of cpu
	CPU uint32 `json:"cpu"`

	// Cores is the number of cores per cpu
	Cores uint32 `json:"cores"`

	// Frequency is the cpu frequency in Mhz
	Frequency float32 `json:"frequency"`

	// CumulativeFrequency is the sum of all cpu frequencies in Mhz
	// All cpu does not have the same frequencies
	CumulativeFrequency float64 `json:"cumulative_frequency"`

	// Capabilitites holds the features of the processor
	Capabilitites []string `json:"capabilitites"`

	// Vendor holds the vendor name
	Vendor string `json:"vendor"`

	// Model holds the vendor's model name
	Model string `json:"model"`
}
