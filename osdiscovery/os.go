package osdiscovery

import (
	"bufio"
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
)

// newOsInfo initialize os infos
func newOsInfo() *OS {
	hostname, _ := os.Hostname()
	return &OS{
		osReleasePath: osReleasePath,
		Hostname:      hostname,
		Architecture:  runtime.GOARCH,
		OSType:        runtime.GOOS,
	}
}

// osInfo return parsed data from OS
func (osi *OS) osInfo() {
	file, err := os.Open(osi.osReleasePath)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	osi.osInfoReader(file)
}

// osInfoReader is used to read file content
func (osi *OS) osInfoReader(file io.Reader) {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}

		stripped := strings.ReplaceAll(line, `"`, "")
		if match := reOSID.FindStringSubmatch(stripped); match != nil {
			osi.Name = match[1]
			osi.Vendor = match[1]
		}

		if match := reOSVersionID.FindStringSubmatch(stripped); match != nil {
			osi.Version = match[1]
		}

		if match := reOSIDLike.FindStringSubmatch(stripped); match != nil {
			osi.Family = match[1]
		}
	}

	switch osi.Name {
	case "rhel":
		osi.Name = "redhat"
		osi.Vendor = "redhat"
		osi.Family = "redhat"

	case "fedora":
		osi.Vendor = "fedora"
		osi.Family = "redhat"
	}
}
