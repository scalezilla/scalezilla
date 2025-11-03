package osdiscovery

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// newCgroups initialize cgroups infos
func newCgroups() *Cgroups {
	return &Cgroups{
		procCgroupsPath: procCgroupsPath,
	}
}

// cgroups return parsed data from OS
func (c *Cgroups) cgroups() {
	file, err := os.Open(c.procCgroupsPath)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	c.cgroupsReader(file)
}

// cgroupsReader is used to read file content
func (c *Cgroups) cgroupsReader(file io.Reader) {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}

		if match := reCgroupsVersion1.FindStringSubmatch(line); match != nil {
			c.Version = 1
		}

		if match := reCgroupsVersion2.FindStringSubmatch(line); match != nil {
			c.Version = 2
			break
		}
	}
}
