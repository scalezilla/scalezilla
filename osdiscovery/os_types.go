package osdiscovery

import (
	"regexp"
)

var (
	// osReleasePath is hold content of OS informations
	osReleasePath = "/etc/os-release"

	// reOSID match os id binded as vendor
	reOSID = regexp.MustCompile(`^ID=(.*)`)

	// reOSName match os version id
	reOSVersionID = regexp.MustCompile(`^VERSION_ID=(.*)`)

	// reOSIDLike match os family
	reOSIDLike = regexp.MustCompile(`^ID_LIKE=(.*)`)
)

// OS holds os informations
type OS struct {
	// osReleasePath holds content of OS informations
	osReleasePath string

	// Name is the os name
	Name string `json:"name,omitempty"`

	// Vendor is the enterprise or organization that create
	// or manage this OS
	Vendor string `json:"vendor,omitempty"`

	// Version is the version of the OS
	Version string `json:"version,omitempty"`

	// Family tell us if this os is a derivation from
	// an another OS
	Family string `json:"family,omitempty"`

	// Hostname return the os hostname
	Hostname string `json:"hostname,omitempty"`

	// Architecture return the system architecture
	Architecture string `json:"architecture,omitempty"`

	// OSType return the if os is darwin, linux etc
	OSType string `json:"os_type,omitempty"`
}
