package entity

import (
	"strconv"
)

type ComputerInfoModel struct {
	CpuInfo         []CpuInfoModel  `json:"cpu_info"`
	MemInfo         MemInfoModel    `json:"mem_info"`
	HostInfo        HostInfoModel   `json:"host_info"`
	DiskInfo        []DiskInfoModel `json:"disk_info"`
	NetInfo         []NetInfoModel  `json:"net_info"`
	PoolActivateNum int             `json:"pool_activate_num"`
	PoolJobNum      int             `json:"pool_job_num"`
}

type CpuInfoModel struct {
	CpuName       string `json:"cpu_name"`
	CpuCoreNum    string `json:"cpu_core_num"`
	CpuFrequency  string `json:"cpu_frequency"`
	CpuVendorId   string `json:"cpu_vendor_id"`
	CpuPhysicalId string `json:"cpu_physical_id"`
	CpuUsage      string `json:"cpu_usage"`
}

type MemInfoModel struct {
	MemTotal string `json:"mem_total"`
	MemUsed  string `json:"mem_used"`
	MemAvi   string `json:"mem_avi"`
	MemUsage string `json:"mem_usage"`
}

type HostInfoModel struct {
	HostName          string `json:"host_name"`
	HostOs            string `json:"host_os"`
	HostOsVersion     string `json:"host_os_version"`
	HostKernelVersion string `json:"host_kernel_version"`
	HostKernelArch    string `json:"host_kernel_arch"`
	HostId            string `json:"host_id"`
}

type DiskInfoModel struct {
	DiskName  string `json:"disk_name"`
	DiskSize  string `json:"disk_size"`
	DiskUsed  string `json:"disk_used"`
	DiskAvi   string `json:"disk_avi"`
	DiskUsage string `json:"disk_usage"`
	FsType    string `json:"fs_type"`
	Opts      string `json:"opts"`
}

type NetInfoModel struct {
	NetName         string `json:"net_name"`
	BytesSent       string `json:"bytes_sent"`
	BytesReceived   string `json:"bytes_received"`
	PacketsSent     string `json:"packets_sent"`
	PacketsReceived string `json:"packets_received"`
}

func NewComputerInfoModel(c []CpuInfoModel, m MemInfoModel,
	h HostInfoModel, d []DiskInfoModel, n []NetInfoModel, activateNum int, jobNum int) *ComputerInfoModel {
	return &ComputerInfoModel{
		CpuInfo:         c,
		MemInfo:         m,
		HostInfo:        h,
		DiskInfo:        d,
		NetInfo:         n,
		PoolActivateNum: activateNum,
		PoolJobNum:      jobNum,
	}
}

func NewCpuInfoModel(cpuName string, coreNum int32, freq float64, VendorId string, physicalId string, usage float64) *CpuInfoModel {
	frequency := strconv.FormatFloat(freq/1000, 'f', 4, 64) + " GHz"
	coreNumber := strconv.Itoa(int(coreNum))
	usagePercent := strconv.FormatFloat(usage, 'f', 4, 64) + " %"
	return &CpuInfoModel{
		CpuName:       cpuName,
		CpuCoreNum:    coreNumber,
		CpuFrequency:  frequency,
		CpuVendorId:   VendorId,
		CpuPhysicalId: physicalId,
		CpuUsage:      usagePercent,
	}
}

func NewMemInfoModel(memTotal uint64, memAvi uint64, memUsed uint64, memUsage float64) *MemInfoModel {
	MemTotal := strconv.FormatFloat(float64(memTotal)/1024/1024/1024, 'f', 4, 64) + " GB"
	MemAvi := strconv.FormatFloat(float64(memAvi)/1024/1024/1024, 'f', 4, 64) + " GB"
	MemUsed := strconv.FormatFloat(float64(memUsed)/1024/1024/1024, 'f', 4, 64) + " GB"
	MemUsage := strconv.FormatFloat(memUsage, 'f', 4, 64) + " %"
	return &MemInfoModel{
		MemTotal: MemTotal,
		MemAvi:   MemAvi,
		MemUsed:  MemUsed,
		MemUsage: MemUsage,
	}
}

func NewHostInfoModel(hostName string, os string, osVersion string, kernelVersion string, kernelArch string, hostId string) *HostInfoModel {
	return &HostInfoModel{
		HostName:          hostName,
		HostOs:            os,
		HostOsVersion:     osVersion,
		HostKernelVersion: kernelVersion,
		HostKernelArch:    kernelArch,
		HostId:            hostId,
	}
}

func NewDiskInfoModel(diskName string, diskSize uint64, diskUsed uint64, diskAvi uint64, diskUsage float64, fsType string, opts string) *DiskInfoModel {
	DiskSize := strconv.FormatFloat(float64(diskSize)/1024/1024/1024, 'f', 4, 64) + " GB"
	DiskUsed := strconv.FormatFloat(float64(diskUsed)/1024/1024/1024, 'f', 4, 64) + " GB"
	DiskAvi := strconv.FormatFloat(float64(diskAvi)/1024/1024/1024, 'f', 4, 64) + " GB"
	DiskUsage := strconv.FormatFloat(diskUsage, 'f', 4, 64) + " %"
	return &DiskInfoModel{
		DiskName:  diskName,
		DiskSize:  DiskSize,
		DiskUsed:  DiskUsed,
		DiskAvi:   DiskAvi,
		DiskUsage: DiskUsage,
		FsType:    fsType,
		Opts:      opts,
	}
}

func NewNetInfoModel(name string, bytesSent uint64, bytesReceived uint64, packetsSent uint64, packetsReceived uint64) *NetInfoModel {
	BytesSent := strconv.FormatFloat(float64(bytesSent)/1024/1024, 'f', 4, 64) + " MB"
	BytesReceived := strconv.FormatFloat(float64(bytesReceived)/1024/1024, 'f', 4, 64) + " MB"
	PacketsSent := strconv.FormatFloat(float64(packetsSent)/1024/1024, 'f', 4, 64) + " MB"
	PacketsReceived := strconv.FormatFloat(float64(packetsReceived)/1024/1024, 'f', 4, 64) + " MB"
	return &NetInfoModel{
		NetName:         name,
		BytesSent:       BytesSent,
		BytesReceived:   BytesReceived,
		PacketsSent:     PacketsSent,
		PacketsReceived: PacketsReceived,
	}
}
