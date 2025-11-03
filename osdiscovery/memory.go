package osdiscovery

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
)

// newMemory initialize memory infos
func newMemory() *Memory {
	return &Memory{
		procMemoryPath: procMemoryPath,
	}
}

// osInfo return os system name, vendor and version
func (mem *Memory) memory() {
	file, err := os.Open(mem.procMemoryPath)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	mem.memoryReader(file)
}

// memoryReader is used to read file content
func (mem *Memory) memoryReader(file io.Reader) {
	reader := bufio.NewReader(file)

	count := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}

		if match := reMemTotal.FindStringSubmatch(line); match != nil {
			if value, err := strconv.Atoi(match[1]); err == nil {
				mem.Total = uint(value)
			}
		}

		if match := reMemAvailable.FindStringSubmatch(line); match != nil {
			if value, err := strconv.Atoi(match[1]); err == nil {
				mem.Available = uint(value)
			}
		}

		if count == 2 {
			break
		}
		count++
	}
}
