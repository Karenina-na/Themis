package util

import (
	"Themis/src/entity"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

// F is a function to handle error
var F func(err error)

//
// SetStatusErrorHandle
// @Description: Set the error handler for the status module
// @param        f : the error handler
//
func SetStatusErrorHandle(f func(err error)) {
	F = f
}

//
// GetCpuInfo
// @Description: Get the CPU information
// @return       []entity.CpuInfoModel : the CPU information
//
func GetCpuInfo() []entity.CpuInfoModel {
	cpuInfos, err1 := cpu.Info()
	if err1 != nil {
		F(err1)
	}
	percent, err2 := cpu.Percent(time.Second, false)
	if err2 != nil {
		F(err2)
	}
	Infos := make([]entity.CpuInfoModel, 0, 1)
	for i, cpuInfo := range cpuInfos {
		info := entity.NewCpuInfoModel(cpuInfo.ModelName,
			cpuInfo.Cores, cpuInfo.Mhz,
			cpuInfo.VendorID, cpuInfo.PhysicalID, percent[i])
		Infos = append(Infos, *info)
	}
	return Infos
}

//
// GetMemInfo
// @Description: Get the memory information
// @return       *entity.MemInfoModel : the memory information
//
func GetMemInfo() *entity.MemInfoModel {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		F(err)
	}
	MemInfo := entity.NewMemInfoModel(memInfo.Total, memInfo.Available, memInfo.Used, memInfo.UsedPercent)
	return MemInfo
}

//
// GetHostInfo
// @Description: Get the host information
// @return       *entity.HostInfoModel : the host information
//
func GetHostInfo() *entity.HostInfoModel {
	hostInfo, err := host.Info()
	if err != nil {
		F(err)
	}
	HostInfo := entity.NewHostInfoModel(hostInfo.Hostname, hostInfo.Platform,
		hostInfo.PlatformVersion, hostInfo.KernelVersion,
		hostInfo.KernelArch, hostInfo.HostID)
	return HostInfo
}

//
// GetDiskInfo
// @Description: Get the disk information
// @return       []entity.DiskInfoModel : the disk information
//
func GetDiskInfo() []entity.DiskInfoModel {
	parts, err := disk.Partitions(true)
	if err != nil {
		F(err)
	}
	Infos := make([]entity.DiskInfoModel, 0, 1)
	for _, part := range parts {
		partInfo, err := disk.Usage(part.Mountpoint)
		if err != nil {
			F(err)
		}
		info := entity.NewDiskInfoModel(part.Device, partInfo.Total, partInfo.Used, partInfo.Free, partInfo.UsedPercent, part.Fstype, part.Opts)
		Infos = append(Infos, *info)
	}
	return Infos
}

//
// GetNetInfo
// @Description: Get the network information
// @return       []entity.NetInfoModel : the network information
//
func GetNetInfo() []entity.NetInfoModel {
	netIOs, err := net.IOCounters(true)
	if err != nil {
		F(err)
	}
	Infos := make([]entity.NetInfoModel, 0, 1)
	for _, netIO := range netIOs {
		info := entity.NewNetInfoModel(netIO.Name, netIO.BytesSent, netIO.BytesRecv, netIO.PacketsSent, netIO.PacketsRecv)
		Infos = append(Infos, *info)
	}
	return Infos
}
