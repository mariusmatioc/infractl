package global

import (
	"fmt"
	"math"
)

// bestMemCpuMatch returns best match based on desired values and acceptable values
// mem values are in MiB, cpu values are such that 1024 == 1vCPU
func bestMemCpuMatch(mem, cpu int) (memMatch, cpuMatch int, err error) {
	if mem == 0 && cpu == 0 {
		// Defaults
		memMatch = mibs(8)
		cpuMatch = 1024
		return
	}
	if mem != 0 {
		if mem <= mibs(2) {
			cpuMatch = 256
			if mem < mibs(1) {
				memMatch = 512
				return
			}
		} else if mem <= mibs(4) {
			cpuMatch = 512
		} else if mem <= mibs(8) {
			cpuMatch = 1024
		} else if mem <= mibs(16) {
			cpuMatch = 2048
		} else if mem <= mibs(30) {
			cpuMatch = 4096
		} else if mem <= mibs(60) {
			cpuMatch = 8192
			mem -= mibs(16)
			g := mem / KiB
			// in 4G increments
			g /= 4
			g *= 4
			memMatch = mibs(16 + g)
			return
		} else if mem <= mibs(120) {
			cpuMatch = 16384
			mem -= mibs(32)
			g := mem / KiB
			// in 8G increments
			g /= 8
			g *= 8
			memMatch = mibs(32 + g)
			return
		} else {
			err = fmt.Errorf("max allowed memory is 120G")
			return
		}
		memMatch = floor(mem)
		return
	}

	// cpu based
	if cpu <= 256 {
		memMatch = mibs(2)
		cpuMatch = 256
	} else if cpu <= 512 {
		memMatch = mibs(3)
		cpuMatch = 512
	} else if cpu <= 1024 {
		memMatch = mibs(5)
		cpuMatch = 1024
	} else if cpu <= 2048 {
		memMatch = mibs(10)
		cpuMatch = 2048
	} else if cpu <= 4096 {
		memMatch = mibs(12)
		cpuMatch = 4096
	} else if cpu <= 8192 {
		memMatch = mibs(24)
		cpuMatch = 8192
	} else if cpu <= 16384 {
		memMatch = mibs(48)
		cpuMatch = 16384
	} else {
		err = fmt.Errorf("max allowed cpu is 16384")
	}
	return
}

// Converts to MiB
func mibs(gigs int) int {
	return gigs * KiB
}

func floor(mibs int) int {
	return int(math.Floor(float64(mibs)/KiB)) * KiB
}
