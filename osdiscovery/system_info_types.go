package osdiscovery

// SystemInfo holds cpu, memory, os informations
type SystemInfo struct {
	// OS return os contents
	OS *OS `json:"os"`

	// CPU return cpu contents
	CPU *CPU `json:"cpu"`

	// Memory return os contents
	Memory *Memory `json:"memory"`

	// Cgroups return os contents
	Cgroups *Cgroups `json:"cgroups"`
}
