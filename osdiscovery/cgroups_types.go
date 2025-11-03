package osdiscovery

import "regexp"

var (
	// procCgroupsPath is hold content of cgroups informations
	procCgroupsPath = "/proc/filesystems"

	// reCgroupsVersion1 match nodev	cgroup field
	reCgroupsVersion1 = regexp.MustCompile(`^nodev\s+(cgroup)?\s+$`)

	// reCgroupsVersion2 match nodev	cgroup2 field
	reCgroupsVersion2 = regexp.MustCompile(`^nodev\s+(cgroup2)?\s+$`)
)

// Cgroups holds version found on linux OS
type Cgroups struct {
	// procCgroupsPath holds content of cgroups informations
	procCgroupsPath string

	// Version return the version of the cgroup
	Version uint `json:"version"`
}
