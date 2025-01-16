package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SystemInfo struct {
	RunTimeOS string
	HostName  string
	VmTotal   uint64
	VmUsed    uint64
}

type CpuInfo struct {
	CpuType  string
	CpuCores uint64
}

type DiskInfo struct {
	DiscTotal uint64
	DiscUsed  uint64
}

// Get system information
func GetSystem() (SystemInfo, error) {
	//sys := new(SystemInfo) //New pointer to and instance of SystemInfo
	sys := SystemInfo{} // New instance of SystemInfo

	caser := cases.Title(language.AmericanEnglish)
	sys.RunTimeOS = caser.String(runtime.GOOS)

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return SystemInfo{}, err
	}
	sys.VmTotal = vmStat.Total

	hostStat, err := host.Info()
	if err != nil {
		return SystemInfo{}, err
	}
	sys.HostName = hostStat.Hostname

	//output := fmt.Sprintf("Hostname: %s\nTotal Memory: %d\nUsed Memory: %d\nOS: %s", sys.hostName, sys.vmTotal, sys.vmTotal, sys.runTimeOS)
	output := sys

	fmt.Printf("Output: %v", output)

	return output, nil
}

// Get CPU data
func GetCPU() (string, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("CPU: %s\nCores %d", cpuStat[0].ModelName, len(cpuStat))
	return output, nil
}

// Get disk data
func GetDisk() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Total Disk Space: %d\nFree Disk Space: %d", diskStat.Total, diskStat.Free)
	return output, nil
}
