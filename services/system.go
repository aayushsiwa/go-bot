package services

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetSystemStats() string {
	// CPU
	cpuPercent, _ := cpu.Percent(0, false)

	// RAM
	vmStat, _ := mem.VirtualMemory()

	// Disk
	diskStat, _ := disk.Usage("/")

	// Uptime
	hostStat, _ := host.Info()
	uptime := time.Duration(hostStat.Uptime) * time.Second

	return fmt.Sprintf(
		"🖥️ **System Stats**\nCPU: %.2f%%\nRAM: %.2f / %.2f GB\nDisk: %.2f / %.2f GB\nUptime: %s",
		cpuPercent[0],
		float64(vmStat.Used)/1e9,
		float64(vmStat.Total)/1e9,
		float64(diskStat.Used)/1e9,
		float64(diskStat.Total)/1e9,
		uptime.Round(time.Second),
	)
}
