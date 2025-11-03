package osdiscovery

// NewSystemInfo initialize os infos
func NewSystemInfo() *SystemInfo {
	s := &SystemInfo{}
	s.systemInfo()
	return s
}

// systemInfo return parsed data from OS
func (s *SystemInfo) systemInfo() {
	osi := newOsInfo()
	osi.osInfo()
	s.OS = osi

	cpu := newCPU()
	cpu.cpu()
	s.CPU = cpu

	mem := newMemory()
	mem.memory()
	s.Memory = mem

	cgroups := newCgroups()
	cgroups.cgroups()
	s.Cgroups = cgroups
}
