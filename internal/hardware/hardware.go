package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func GetSystem() (string, error) {
	runTimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Hostname: %s\nTotal Memory: %d\nUsed Memory: %d\nOS: %s", hostStat.Hostname, vmStat.Total, vmStat.Used, runTimeOS)

	return output, nil
}

func GetCPU() (string, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("CPU: %s\nCores %d", cpuStat[0].ModelName, len(cpuStat))
	return output, nil
}

func GetDisk() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Total Disk Space: %d\nFree Disk Space: %d", diskStat.Total, diskStat.Free)
	return output, nil
}
