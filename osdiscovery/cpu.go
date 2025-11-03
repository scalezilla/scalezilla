package osdiscovery

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// newCPU initialize cpu infos
func newCPU() *CPU {
	return &CPU{
		procCPUPath: procCPUPath,
	}
}

// cpu return parsed data from OS
func (cpu *CPU) cpu() {
	file, err := os.Open(cpu.procCPUPath)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	cpu.cpuReader(file)
}

// cpuReader is used to read file content
func (cpu *CPU) cpuReader(file io.Reader) {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}

		if match := reProcessor.FindStringSubmatch(line); match != nil {
			cpu.CPU++
		}

		if match := reCores.FindStringSubmatch(line); match != nil {
			if value, err := strconv.Atoi(match[1]); err == nil {
				cpu.Cores = uint32(value)
			}
		}

		if match := reFrequency.FindStringSubmatch(line); match != nil {
			if value, err := strconv.ParseFloat(match[1], 32); err == nil {
				// each cpu does not always have to same frequencies
				// so we keep the high one
				if float32(value) > cpu.Frequency {
					cpu.Frequency = float32(value)
				}
				cpu.CumulativeFrequency += value
			}
		}

		if match := reCapabilities.FindStringSubmatch(line); match != nil {
			cpu.Capabilitites = strings.Split(match[1], " ")
		}

		if match := reVendorID.FindStringSubmatch(line); match != nil {
			cpu.Vendor = strings.ReplaceAll(match[1], "Genuine", "")
		}

		if match := reModel.FindStringSubmatch(line); match != nil {
			cpu.Model = match[1]
		}
	}
}
