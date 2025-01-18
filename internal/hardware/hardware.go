package hardware

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gldebug/gpuinfo"
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
	CpuCores int
}

type DiskInfo struct {
	DiscTotal uint64
	DiscUsed  uint64
	DiskFree  uint64
}

type GpuInfo struct {
	GpuType string
	GpuFree int
	GpuUsed int
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
func GetCPU() (CpuInfo, error) {

	cpuInfo := CpuInfo{}

	cpuStat, err := cpu.Info()
	if err != nil {
		return CpuInfo{}, err
	}
	cpuInfo.CpuType = cpuStat[0].ModelName
	cpuInfo.CpuCores = len(cpuStat)

	output := cpuInfo

	fmt.Sprintf("CPU: %s\nCores %d", cpuInfo.CpuType, cpuInfo.CpuCores)
	return output, nil
}

// Get disk data
func GetDisk() (DiskInfo, error) {
	diskInf := DiskInfo{}
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskInfo{}, err
	}

	diskInf.DiscTotal = diskStat.Total
	diskInf.DiscUsed = diskStat.Used
	diskInf.DiscUsed = diskStat.Free

	fmt.Sprintf("Total Disk Space: %d\n Used Disk Space: %d\n Free Disk Space: %d", diskInf.DiscTotal, diskInf.DiscUsed, diskInf.DiscUsed)
	output := diskInf
	return output, nil
}

func GetGPU() GpuInfo {
	gpuInf := GpuInfo{}

	gpu := gpuinfo.NVGpu{}
	gpuInf.GpuFree = gpu.Free()

	println("GPU Free: %d", gpuInf.GpuFree)

	return gpuInf

}
